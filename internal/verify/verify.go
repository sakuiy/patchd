package verify

// Package verify checks that a patched function is equivalent to the
// original on every input except the crashing path.
//
// Plan:
//   1. Lower old and new function to LLVM IR (bitcode).
//   2. Encode both as SMT formulas over shared symbolic inputs
//      (bitvectors for integers, array theory for memory).
//   3. Ask Z3 whether any input exists where outputs or memory effects
//      differ. unsat means equivalent; sat gives a counterexample, which
//      becomes a regression test.
//   4. Loops are unrolled to a fixed bound; timeouts fall back to
//      differential fuzzing.
//   5. Concurrent code: encode atomic accesses as events and solve
//      happens-before / reads-from / modification order under the C11
//      memory model.

// Result is the outcome of an equivalence check.
type Result int

const (
	Equivalent Result = iota
	Counterexample
	Timeout
)
