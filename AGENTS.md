# memory — AI Agent Guide

This file provides guidance to AI coding assistants (Claude Code, Copilot, Cursor, etc.) when working with this repository.

## 1. Project Overview

**Go vendor library** — a self-contained module of generic data structures and algorithms, imported by other projects via `go get github.com/FrogoAI/memory`. This is **not** an application — there is no `main()`, no `cmd/`, no HTTP server, no CLI, no deployment.

- **Module**: `github.com/FrogoAI/memory`
- **Go**: 1.25+ with generics
- **Type**: Reusable library (`go get`-able)
- **License**: MIT

## 2. Directory Structure

```
memory/
  go.mod
  go.sum

  # --- Packages (each is an independent data structure) ---
  bloom/              # Counting Bloom Filter — probabilistic membership with removal
  btree/              # B-tree — ordered key-value storage with configurable order
  comparator/         # Type-aware comparison functions (int, string, float, time, etc.)
  fuzzysearch/        # Fuzzy string matching with Levenshtein distance ranking
  hll/                # HyperLogLog — efficient cardinality estimation
  linkedlist/         # Generic doubly-linked list with ID-based indexing
  lru/                # Thread-safe LRU cache with generic values
  orderedmap/         # Insertion-ordered map with thread-safe operations
  registry/           # Thread-safe registry for grouping items by category and ID
  sortedset/          # Redis-like ZSET using skip list with scoring
  stack/              # Generic LIFO stack with slice-based backing
  utils/              # Thread-safe collections, hashing, sorting, string helpers

  # --- Tooling ---
  .testcoverage.yml
  .gitignore
  .github/
    workflows/go.yml
    dependabot.yml

  # --- Metadata ---
  README.md
  AGENTS.md
  CODEOWNERS
  LICENSE

  # --- Documentation ---
  docs/
    aics/plan.md      # Master action plan
    source/            # Architecture & engineering reference (gitignored except index)
```

### What is NOT here

- No `cmd/` — this is not an application.
- No `main.go` — there is nothing to run.
- No `internal/` — everything is a sub-package. The consumer decides what to import.
- No `pkg/` — each top-level directory IS a package.
- No `docker-compose.yml` — the consuming app owns infrastructure.
- No `services/` — this is a library, not a monorepo.

## 3. Architecture

Each package is an independent, self-contained data structure. No package depends on another (except `utils/` and `comparator/` which are shared infrastructure).

```
Consumer App
  |
  v
[bloom]  [btree]  [lru]  [orderedmap]  [sortedset]  ...
  |         |        |         |              |
  v         v        v         v              v
[utils]  [comparator]                    [utils]
```

### Dependency rules

**Allowed:**
```
any package  -->  comparator, utils
bloom        -->  packer (external, for serialization)
btree        -->  comparator
sortedset    -->  comparator
```

**Forbidden:**
```
Package A    -/->  Package B (no cross-dependencies between data structure packages)
ANY package  -/->  framework, HTTP, CLI, DI container, or application-level packages
```

### Key constraint: no infrastructure leakage

Packages must never import web frameworks, database drivers, message queues, observability SDKs, or any application-level infrastructure. Heavy external dependencies should be avoided — prefer stdlib or lightweight alternatives.

## 4. Package Design Principles

Each data structure package follows a consistent pattern:

| Component | Contains | Does NOT contain |
|---|---|---|
| Constructor (`New*()`) | Creates and returns the data structure | Business logic, config parsing |
| Core methods | CRUD operations, iteration, serialization | Infrastructure, I/O |
| `errors.go` (if needed) | `var Err... = errors.New(...)` sentinel errors | Logic, types |
| `*_test.go` | Table-driven tests, benchmarks | Infrastructure dependencies |
| `options.go` (if needed) | Configuration types, functional options | Business logic |

### API consistency across packages

| Operation | Preferred name | Avoid |
|---|---|---|
| Create/Insert | `Add()` | `Insert()`, `Put()` (unless semantic, like cache `Put`) |
| Read | `Get()` | `Fetch()`, `Retrieve()` |
| Delete | `Remove()` | `Delete()` |
| Count elements | `Len()` | `Size()`, `GetCount()`, `Count()` |
| Empty all | `Clear()` | `Reset()`, `Truncate()` |
| Check existence | `Has()` or `Test()` | `Contains()`, `Exists()` |
| Iterate | `Iterator()` or `Range()` | `ForEach()`, `List()` |

### Thread-safety convention

- Thread-safe packages MUST document it in the package comment: `// Package X provides a thread-safe ...`
- Non-thread-safe packages MUST document it: `// Package X is NOT safe for concurrent use.`
- When adding thread safety, use `sync.RWMutex` (read-heavy) or `sync.Mutex` (write-heavy).
- All public methods on a thread-safe type must be protected. No partial safety.

## 5. Coding Standards

### Formatting & Linting
- **Formatter**: `gofumpt` (strict superset of gofmt).
- **Linter**: `golangci-lint v2` (strict, includes `wsl_v5`, `revive`, `staticcheck`).
- **Import ordering** (enforced by gci): standard library → blank → third-party → project packages.
- **Line length**: 120 characters max.
- **WSL** (whitespace linter) enabled — follow its blank-line conventions.

### Naming
- **Full words**: `httpConfig` not `httpCfg`, `userID` not `uid`.
- **Acronyms all-caps**: `HTTPClient`, `userID`, `xmlParser` (not `HttpClient`, `userId`).
- **No Get prefix**: `user.Name()` not `user.GetName()`. Setters use `Set`: `user.SetName(n)`.
- **Doc comments start with name**: `// Tree represents...` not `// Represents a tree...`.
- **Package names**: lowercase, single-word, no underscores. No `util`, `common`, `helpers`, `misc`.

### Code Style
- **Forbidden**: `fmt.Print*`, `log.Print*`, bare `print*`. Use `log/slog` if logging is needed.
- **Never ignore errors**: always check and return. No `_ = fn()`.
- **No flag arguments**: A `bool` parameter means the function does two things. Split or use options struct.
- **Error strings**: Lowercase, no punctuation — `"open file"` not `"Open file."`. Wrap with `fmt.Errorf("context: %w", err)`.
- **Early return, no else**: Error cases return early; happy path stays left-aligned.
- **Function size**: Aim for 5–20 lines. Over 40 lines signals the function does too much.
- **Zero value should be useful**: Design types so the zero value is valid.
- **No `init()` I/O**: `init()` should only register things. Never read files or make network calls.
- **No commented-out code**: Delete it. Git remembers.
- **No panics in exported API**: Return errors. Let the caller decide. Panic only for truly unrecoverable programmer bugs.
- **Compile-time interface checks**: `var _ Interface = (*Impl)(nil)` where interfaces are used.
- **Avoid `fmt.Sprintf` in hot paths**: Uses reflect. Prefer `+` or `strings.Join`. Fine for errors/logs.

## 6. Testing

### Table-driven tests (mandatory, no exceptions)

```go
cases := []struct {
    name    string
    input   string
    want    string
    wantErr error
}{...}

for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) { ... })
}
```

### Test types

| Type | Build tag | CI runs it |
|---|---|---|
| Unit | None | Always |
| Benchmark | None | Always (smoke run) |
| Integration | `//go:build integration` | On demand |

### Rules

- Every exported function must have at least one test.
- Test boundary conditions: nil inputs, empty collections, max capacity, zero values.
- Thread-safe packages must have concurrent tests (use `-race` flag).
- Benchmarks go in `*_test.go` files using `b.Run` for sub-benchmarks.
- Coverage threshold: 35%+ overall (`.testcoverage.yml`), increasing over time.

### Benchmarks

- Sub-benchmarks via `b.Run` for different scenarios.
- Compare with `benchstat` between runs.
- Critical paths: insertion, lookup, deletion, iteration for each data structure.

## 7. Build & Run Commands

```bash
# Run all tests
go test ./...

# Run tests with race detection (CI mode)
GOMAXPROCS=$(nproc) go test -coverprofile=coverage.out -cover -race -p $(nproc) -parallel=$(nproc) ./...

# Run benchmarks
go test -bench=. -benchmem -run="^$" ./...

# Lint
golangci-lint run -v ./...

# Auto-fix formatting
golangci-lint run --fix
```

## 8. CI Pipeline

GitHub Actions (`.github/workflows/go.yml`):
1. **Lint** — golangci-lint v2
2. **Test** — race detection, coverage threshold via `go-test-coverage`
3. **Bench** — smoke run of all benchmarks

## 9. Design Principles (Priority Order)

Higher principles override lower ones when in conflict:

1. **Performance by Design** — performance is a core trait, not an afterthought. This is a data structures library; every allocation matters.
2. **KISS** — Keep It Simple and Smart. Simple = efficient, readable, maintainable. Minimize moving parts, maximize clarity.
3. **DRY** — Don't Repeat Yourself. Extract shared logic to `utils/` or `comparator/`. But don't abstract prematurely (KISS overrides DRY).
4. **Clean Code** — readable, testable, maintainable. See `docs/source/Clean Code.md` for full reference.
5. **Minimal Dependencies** — prefer stdlib. Every external dependency is a liability for a vendor library. Heavy deps (like mongo-driver for ID generation) should be replaced with lightweight alternatives.
6. **Consistent API Surface** — all packages should feel like they belong to the same library. Same naming conventions, same patterns, same error handling style.

## 10. Error Handling

All sentinel errors live in `errors.go` per package:

```go
var (
    ErrEmpty   = errors.New("empty")
    ErrNotFound = errors.New("not found")
)
```

| Rule | Detail |
|---|---|
| Never ignore errors | `_ = fn()` is forbidden |
| Critical errors | Return to caller |
| Logging | `log/slog` only. No `fmt.Print*`, no `log.Print*` |
| Wrapping | `fmt.Errorf("context: %w", err)` |
| No panics | Exported API returns errors, never panics |

The consumer checks errors with `errors.Is(err, bloom.ErrEmpty)`.

## 11. Documentation

### Package-level comments (mandatory)

Every package must have a doc comment in its primary `.go` file:

```go
// Package bloom provides a counting Bloom filter for probabilistic
// membership testing with element removal support.
//
// The implementation uses FNV-64 hashing with configurable false-positive
// rates. It is NOT safe for concurrent use; callers must synchronize access.
package bloom
```

### Exported symbol comments

Every exported type, function, and method must have a godoc comment starting with the symbol name:

```go
// CountingFilter is a probabilistic data structure that supports
// both membership testing and element removal.
type CountingFilter struct { ... }

// NewCounting creates a CountingFilter sized for n expected elements
// with a false-positive probability of p (0 < p < 1).
func NewCounting(n int, p float64) *CountingFilter { ... }
```

## 12. Working with plan.md

`docs/aics/plan.md` is the single source of truth for remaining work.

### Statuses

| Status | Agent action |
|--------|-------------|
| `TODO` | Implement |
| `PARTIAL` | Finish existing work |
| `NEEDS REVIEW` | Validate against spec, fix if needed |
| `CRITICAL` | Implement (high priority) |
| `Enhancement` | Implement (lower priority) |
| `Validate` | Check code against spec |
| `Blocked` | **SKIP** |
| `USER_INVOLVEMENT_NEEDED` | **SKIP** |

Done items are **removed** from plan.md (git history preserves them).

## 13. Agent Behavior Rules

1. **Ask, don't guess** — when requirements are unclear, mark as `Blocked` or `USER_INVOLVEMENT_NEEDED`. Never implement assumptions.
2. **Read AGENTS.md** before starting any work.
3. **Read existing code** before proposing changes. Understand the current implementation.
4. **One task at a time** — complete and verify before moving to next.
5. **Leave changes unstaged** — human reviews and commits.
6. **Run tests** — `go test ./...` after every change. Fix what you break.
7. **Run linter** — `golangci-lint run ./...` before marking a task done.
8. **Keep it KISS** — prefer simple implementations. Don't over-engineer data structures.
9. **Performance matters** — avoid unnecessary allocations. Benchmark before and after changes to hot paths.
10. **No new dependencies without justification** — every import in `go.mod` must earn its place.
11. **Consistent API** — follow the naming conventions in section 4 when adding or modifying public methods.
12. **Document thread-safety** — always state whether a type is safe for concurrent use.
