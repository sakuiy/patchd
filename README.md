# Baby Care Robot
An AI agent built with Devin and Node.js for baby health monitoring and temperature tracking

**Self-healing patch pipeline: from production crash to verified hot patch.**

`patchd` takes a production crash, localizes it to a source line, lets an AI agent (Devin) write a fix, formally checks that the fix doesn't change any other behavior, and applies it as a hot patch at a function boundary, or falls back to a normal PR when it can't prove the fix is safe.

> Status: early prototype / architecture skeleton.

## Pipeline

```
crash report ──> ingest ──> symbolize ──> repro ──> fix (Devin) ──> verify ──> gate ──> hotpatch ──> report
```

1. **Ingest:** a crash reporter (for example a Sentry webhook) sends the stack trace, core dump, and build ID to `patchd ingest`.
2. **Symbolize:** `patchd symbolize` maps the trace to source lines using the build's DWARF debug info.
3. **Repro:** `patchd repro` replays the crash in a sandbox from the core dump and produces a failing test.
4. **Fix:** Devin receives the repro plus the source, writes a candidate fix, and adds a regression test.
5. **Verify:** `patchd verify` lowers the old and new function to LLVM IR and runs Alive2-style translation validation in Z3. The check covers every input except the crashing path. Concurrent code is checked with an axiomatic C11 memory model, in the style of Dartagnan.
6. **Gate:** if the proof fails or times out, `patchd` falls back to differential fuzzing and opens a normal PR without hot-patching.
7. **Hot-patch:** `patchd apply` waits for quiescence at the function boundary, swaps the function pointer atomically, and rolls back automatically if the error rate rises.
8. **Report:** `patchd report` opens a PR with the fix, the proof log, and the regression test.

## Safety rules

- Patches only happen at **function boundaries**. Signatures, types, and Rust lifetimes must be identical, and the new function must pass the borrow checker against the same signature.
- A function is never patched while a thread is inside it. `patchd` waits for quiescence first.
- Anything that can't be proven within the time budget is **never** hot-patched; it goes out as a normal deploy.

## CLI

```
patchd ingest    --webhook :8080
patchd symbolize --build-id <id> --trace crash.json
patchd repro     --core core.dump --out repro_test
patchd verify    --old old.bc --new new.bc --fn <symbol> [--memory-model c11]
patchd apply     --pid <pid> --fn <symbol> --patch patch.so [--rollback-on-error-rate 0.5%]
patchd report    --repo <owner/name>
```

## Agent workflow (Devin)

Devin runs steps 4–8 in a loop:

```
patchd repro ... && <write fix + test> && patchd verify ... && patchd apply ... || patchd report --fallback-pr
```

## Layout

```
cmd/patchd/          CLI entry point
internal/ingest/     crash webhook + storage
internal/symbolize/  DWARF symbolization
internal/repro/      sandboxed replay from core dumps
internal/verify/     LLVM IR -> SMT (Z3) equivalence, memory-model checks
internal/hotpatch/   quiescence detection + atomic function swap + rollback
internal/report/     PR + proof log generation
examples/            sample crash-to-patch walkthrough
docs/                design notes
```

## Supported targets (planned)

C and C++ (clang), Rust, and Go (via the gollvm toolchain).

## License

MIT
