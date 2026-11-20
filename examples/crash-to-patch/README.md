# Example: null deref in a request parser

1. Sentry reports `SIGSEGV in parse_header()` on build `a1b2c3`.
2. `patchd symbolize` points at `src/http/header.c:142`.
3. `patchd repro` rebuilds the input from the core dump and gives a failing test.
4. Devin adds a bounds/null check and a regression test.
5. `patchd verify` gives unsat on all inputs outside the crashing path, so the fix is equivalent.
6. `patchd apply` waits for quiescence, swaps `parse_header`, and watches error rate for 10 minutes.
7. `patchd report` opens a PR with the diff, the Z3 proof log and the test.
