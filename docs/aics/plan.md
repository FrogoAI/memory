# memory — Master Action Plan

Single source of truth for all remaining work items across sessions.
**Done items are removed** — git history preserves what was completed.

Last updated: 2026-03-27 <!-- S-016 completed -->

---

---

## Phase 2: Dependencies & API Consistency

### P2.1 — Naming Consistency

| Story | Action | Status |
|-------|--------|:------:|
| S-017 | Standardize `Remove()` everywhere (replace any `Delete()` usage) | TODO |
| S-018 | Add missing `Len()` to `lru`, `bloom`, `hll` | TODO |
| S-019 | Add missing `Clear()` to `lru`, `stack`, `linkedlist`, `hll` | TODO |

---

## Phase 3: Thread-Safety & Documentation

### P3.0 — Thread-Safety Audit

| Story | Action | Status |
|-------|--------|:------:|
| S-020 | Document thread-safety status in package-level godoc for ALL packages. Currently undocumented everywhere | TODO |
| S-021 | Add `sync.RWMutex` to `bloom/` — shared counter mutations are unsafe | TODO |
| S-022 | Add `sync.RWMutex` to `btree/` — tree rebalancing is unsafe | TODO |
| S-023 | Add `sync.RWMutex` to `linkedlist/` — complex pointer manipulation is unsafe | TODO |
| S-024 | Add `sync.RWMutex` to `stack/` — bare slice modifications are unsafe | TODO |
| S-025 | Add `sync.Mutex` to `hll/` — wraps external lib without synchronization | TODO |
| S-026 | Upgrade `lru/` from `sync.Mutex` to `sync.RWMutex` for read-heavy workloads | TODO |

### P3.1 — Package Documentation

| Story | Action | Status |
|-------|--------|:------:|
| S-030 | Add package-level godoc comments to ALL packages (bloom, btree, comparator, fuzzysearch, hll, linkedlist, lru, orderedmap, registry, sortedset, stack, utils) | TODO |
| S-031 | Add godoc comments to all exported types and constructors in `lru/`, `stack/`, `linkedlist/` (currently zero comments) | TODO |
| S-032 | Add `Example*` test functions for `bloom`, `btree`, `lru`, `sortedset`, `fuzzysearch` — these appear in godoc | Enhancement |

---

## Phase 4: Testing & Benchmarks

### P4.0 — Benchmarks

| Story | Action | Status |
|-------|--------|:------:|
| S-040 | Add benchmarks to `bloom/`: `BenchmarkAdd`, `BenchmarkTest`, `BenchmarkRemove` | TODO |
| S-041 | Add benchmarks to `btree/`: `BenchmarkPut`, `BenchmarkGet`, `BenchmarkRemove`, `BenchmarkIterate` | TODO |
| S-042 | Add benchmarks to `lru/`: `BenchmarkPut`, `BenchmarkGet`, `BenchmarkEviction` | TODO |
| S-043 | Add benchmarks to `sortedset/`: `BenchmarkAdd`, `BenchmarkGetByRank`, `BenchmarkRange` | TODO |
| S-044 | Add benchmarks to `fuzzysearch/`: `BenchmarkMatch`, `BenchmarkRankFind` | TODO |
| S-045 | Add benchmarks to `orderedmap/`: `BenchmarkAdd`, `BenchmarkGet`, `BenchmarkIterate` | TODO |
| S-046 | Add benchmarks to `linkedlist/`, `stack/`, `registry/`, `hll/` | Enhancement |

### P4.1 — Test Coverage Improvements

Current overall threshold: 35%. Target: **80% per package, 80% overall**.

| Story | Action | Package | Current | Target | Status |
|-------|--------|---------|:-------:|:------:|:------:|
| S-050 | Add concurrent test cases for all thread-safe packages (`lru`, `orderedmap`, `sortedset`, `registry`) using `-race` | multiple | — | — | TODO |
| S-051 | Add edge-case tests to `bloom/`: capacity 0, remove non-existent, serialization round-trip | bloom | 81% | 85%+ | TODO |
| S-052 | Add edge-case tests to `lru/`: capacity 1, eviction under concurrent access, get-after-evict | lru | 95% | 95%+ | TODO |
| S-053 | Add edge-case tests to `btree/`: empty tree operations, single-element tree, duplicate keys | btree | 89% | 90%+ | TODO |
| **Theme: Coverage to 80%** | | | | | |
| S-054 | Add tests to `stack/`: Push, Pop, Peek, Len on empty stack, large stack, Pop-empty | stack | **0%** | 80%+ | CRITICAL |
| S-055 | Add tests to `comparator/`: cover all 18 comparator functions (int, float, string, time, byte, rune, uint, diff) with table-driven tests | comparator | **24%** | 80%+ | CRITICAL |
| S-056 | Add tests to `orderedmap/`: Add/Get/Remove, iteration order, Truncate, Copy, index-based access, concurrent operations | orderedmap | **39%** | 80%+ | TODO |
| S-057 | Add tests to `registry/`: group create/destroy, item add/remove, Iterator sync+async, GetGroup not found, lifecycle hooks | registry | **47%** | 80%+ | TODO |
| S-058 | Add tests to `sortedset/`: Add/Remove/GetByRank, Range, GetByKeyRange, Dump/Restore, score updates, edge ranks | sortedset | **49%** | 80%+ | TODO |
| S-059 | Add tests to `utils/`: SafeMap concurrent ops, SafeList all methods, CRC32/CRC16 known vectors, string normalization edge cases | utils | **53%** | 80%+ | TODO |
| S-080 | Add tests to `linkedlist/`: PushFront/PushBack, Remove by ID, Append, List() ordering, ByID not found, empty list ops | linkedlist | **63%** | 80%+ | TODO |
| S-081 | Add tests to `fuzzysearch/`: MatchNormalized, MatchNormalizedFold, RankFind with ties, empty input, unicode edge cases | fuzzysearch | **76%** | 80%+ | TODO |
| S-082 | Add tests to `hll/`: Union, intersection estimation, large cardinality, serialization round-trip | hll | **79%** | 80%+ | TODO |
| **Theme: Coverage threshold** | | | | | |
| S-083 | Update `.testcoverage.yml`: raise `total` from 35 to 80, set `package` to 80, set `file` to 50 | config | 35% | 80% | TODO |

---

## Phase 5: Tooling & Infrastructure

### P5.0 — Build Tooling

| Story | Action | Status |
|-------|--------|:------:|
| S-061 | Create `.golangci.yml` with strict linting config (errcheck, govet, staticcheck, gosec, revive, misspell, lll, forbidigo, mnd, wsl_v5, gofmt, goimports, gci) | TODO |
| S-062 | Add `benchmarks/` directory with `.gitkeep` for baseline tracking | Enhancement |

### P5.1 — Linting Cleanup

| Story | Action | Status |
|-------|--------|:------:|
| S-065 | Audit and fix all `nolint` suppressions — remove unnecessary ones, add justification comments to remaining | TODO |
| S-066 | Extract magic numbers to named constants in `bloom/bloom.go`, `sortedset/sortedset.go`, `utils/strings.go` | TODO |
| S-067 | Fix `sortedset/sortedset.go` `GetByKeyRange()` cyclomatic complexity — refactor into smaller functions | Enhancement |

---

## Phase 6: Feature Gaps

### P6.0 — Missing Standard Operations

| Story | Action | Status |
|-------|--------|:------:|
| S-070 | Add `Copy()` / `Clone()` to `bloom`, `btree`, `lru`, `linkedlist`, `stack`, `sortedset` | Enhancement |
| S-071 | Add `Iterator()` to `lru/` — iterate over cached entries (most-recent to least-recent) | Enhancement |
| S-072 | Add `Keys()` and `Values()` helpers to `orderedmap`, `btree` | Enhancement |
| S-073 | Add `MarshalBinary()` / `UnmarshalBinary()` to packages that support serialization (`bloom`, `sortedset`, `hll`) for stdlib compatibility | Enhancement |

### P6.1 — Error Handling Improvements

| Story | Action | Status |
|-------|--------|:------:|
| S-075 | Add `(value, bool)` return pattern to `orderedmap.Get()` — currently returns zero-value on miss with no way to distinguish | Enhancement |
| S-076 | ~~Removed — simdict moved to github.com/FrogoAI/lsh~~ | — |

---

## Summary — Remaining Work

| Phase | Items | Priority |
|-------|:-----:|----------|
| **P2** Dependencies & API | 5 | High |
| **P3** Thread-Safety & Docs | 9 | High |
| **P4** Testing & Benchmarks | 21 | Medium–High |
| **P5** Tooling | 4 | Medium |
| **P6** Feature Gaps | 6 | Low |
| **TOTAL** | **45** | |

---

## Execution Order

1. P5.0: S-060..S-061 (Makefile + linter config — enables CI for all subsequent work)
4. P4.1 critical: S-054..S-055 (stack 0%, comparator 24% — worst coverage gaps)
5. P2.1: S-015..S-019 (naming consistency)
6. P3: S-020..S-032 (thread-safety + docs)
7. P4.0: S-040..S-046 (benchmarks)
8. P4.1 remaining: S-050..S-053, S-056..S-083 (coverage to 80% across all packages)
9. P5.1: S-065..S-067 (lint cleanup)
10. P6: S-070..S-076 (feature gaps)
