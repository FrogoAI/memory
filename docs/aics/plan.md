# memory — Master Action Plan

Single source of truth for all remaining work items across sessions.
**Done items are removed** — git history preserves what was completed.

Last updated: 2026-03-27 <!-- S-055 completed -->

---

---

## Phase 2: Dependencies & API Consistency

### P2.1 — Naming Consistency

*All items complete.*

---

---

## Phase 4: Testing & Benchmarks

### P4.0 — Benchmarks

*All items complete.*

### P4.1 — Test Coverage Improvements

Current overall threshold: 35%. Target: **80% per package, 80% overall**.

| Story | Action | Package | Current | Target | Status |
|-------|--------|---------|:-------:|:------:|:------:|
| S-052 | Add edge-case tests to `lru/`: capacity 1, eviction under concurrent access, get-after-evict | lru | 95% | 95%+ | TODO |
| S-053 | Add edge-case tests to `btree/`: empty tree operations, single-element tree, duplicate keys | btree | 89% | 90%+ | TODO |
| **Theme: Coverage to 80%** | | | | | |
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
| **P4** Testing & Benchmarks | 11 | Medium–High |
| **P5** Tooling | 3 | Medium |
| **P6** Feature Gaps | 5 | Low |
| **TOTAL** | **19** | |

---

## Execution Order

1. ~~P5.0: S-061 (.golangci.yml)~~ Done
2. ~~P4.1 critical: S-054 (stack)~~ Done
3. ~~P4.1 critical: S-055 (comparator 24% — worst coverage gap)~~ Done
4. ~~P4.0: S-045..S-046 (benchmarks)~~ Done
5. P4.1 remaining: S-050..S-053, S-056..S-083 (coverage to 80% across all packages)
6. P5.1: S-065..S-067 (lint cleanup)
7. P6: S-070..S-075 (feature gaps)
