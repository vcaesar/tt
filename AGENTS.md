# AGENTS.md — github.com/vcaesar/tt

Small, colorful Go test-assertion helpers. Stdlib only (`go 1.17`, no deps).

## Layout

- `tt.go` — assertions (`Equal/Expect/Nil/Error/...`), `BM`, `Pprof`, `Version`.
- `assert.go` — `Assertions` method forms (`tt.New(t)`).
- `fn.go` — ANSI colors, `CallerInfo()`. `dbg.go` — `Log/Err/Drop/UnUsed`.
- `cov.go`, `sql.go`, `http.go` — empty placeholders. `mock/` — mock objects.
- `http/` — in-process `http.Handler` testing (`New(h).Get(p).Run()` →
  `*Response`; `Status/Body/JSON/...` funcs and `r.Expect(t)` chain). Failures
  go through `tt.Fmt`; default caller depth is `depth` (4) and `chain()` adds
  one for `Expect` methods — `TestTrace` pins these.
- `example/` — `add` fixture used by `tt_test.go` and README. Never import `tt`
  from `example/add.go` (only `add_test.go` may).

## Commands

```sh
go test -v ./...              # CI (also runs example/)
go test -bench=. .
gofmt -l . && go vet ./...    # no linter config; keep gofmt clean
```

## Conventions

- Apache-2.0 header on new `.go` files (copy from `tt.go`). Doc comments on exports.
- Assertions return `bool` and use `t.Errorf`; never `FailNow`.
- Trailing `args ...interface{}`: `[0]` message, `[1]` caller depth `int` or trace
  `string`, `[2]` trace `string` (see `argsFn` in `tt.go`).
- Tests are in-package; negative cases use `mockT := new(testing.T)`.
- Explicit traces/depths in tests (`"tt_test.go:95"`, `4`) pin `Error Trace:`
  output — update them when moving lines.
- Pointer receivers for new `Assertions` methods.

## Pitfalls

- **Caller depth is hand-tuned** (`callSub`/`callAdd`/`typeCall`, `call+1`).
  Adding a wrapper silently mis-points traces; eyeball `go test -v .` output.
- **Global `Type`** (`tt.go`) makes `Equal` also compare types. `TestType` sets it
  and never resets, so tests are order-dependent; reset it if you depend on it.
- `Pprof` binds `127.0.0.1:6060`; `TestTT` calls `Pprof(1)` (port must be free, ≥1s).
- `Err()` calls `log.Fatalln` — never use in tests.
- `Error()` only matches `errors.New` (`*errors.errorString`), not wrapped/custom errors.
