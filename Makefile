.PHONY: test lint coverage bench bench-baseline bench-compare ci

COVER_PROFILE  := coverage.out
BENCH_DIR      := benchmarks
BENCH_FILE     := $(BENCH_DIR)/current.txt
BENCH_BASELINE := $(BENCH_DIR)/baseline.txt
NPROC          := $(shell nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)

## test: run all tests with race detection
test:
	GOMAXPROCS=$(NPROC) go test -race -p $(NPROC) -parallel=$(NPROC) ./...

## lint: run golangci-lint
lint:
	golangci-lint run -v ./...

## coverage: run tests with coverage and check thresholds
coverage:
	GOMAXPROCS=$(NPROC) go test -coverprofile=$(COVER_PROFILE) -cover -race -p $(NPROC) -parallel=$(NPROC) ./...
	go-test-coverage --config=./.testcoverage.yml

## bench: run all benchmarks
bench:
	go test -bench=. -benchmem -run="^$$" ./...

## bench-baseline: save current benchmark output as baseline
bench-baseline:
	@mkdir -p $(BENCH_DIR)
	go test -bench=. -benchmem -count=6 -run="^$$" ./... | tee $(BENCH_BASELINE)

## bench-compare: run benchmarks and compare against baseline
bench-compare:
	@mkdir -p $(BENCH_DIR)
	@test -f $(BENCH_BASELINE) || { echo "No baseline found. Run 'make bench-baseline' first."; exit 1; }
	go test -bench=. -benchmem -count=6 -run="^$$" ./... | tee $(BENCH_FILE)
	benchstat $(BENCH_BASELINE) $(BENCH_FILE)

## ci: full CI pipeline — lint, test with coverage, bench smoke run
ci: lint coverage bench
