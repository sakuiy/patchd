package main

import (
	"fmt"
	"os"
)

const usage = `patchd: crash-to-verified-hot-patch pipeline

Usage:
  patchd <command> [flags]

Commands:
  ingest     receive crash reports (webhook)
  symbolize  map a stack trace to source lines
  repro      replay a crash in a sandbox
  verify     prove old/new function equivalence with Z3
  apply      hot-patch a running process at a function boundary
  report     open a PR with the fix and proof log
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "ingest", "symbolize", "repro", "verify", "apply", "report":
		fmt.Printf("patchd %s: not implemented yet (prototype skeleton)\n", os.Args[1])
	default:
		fmt.Print(usage)
		os.Exit(1)
	}
}
