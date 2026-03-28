package main

import (
	"nanozkv/debug"
	"nanozkv/vm"
	"nanozkv/zk"
)

func main() {

	code := []vm.Instruction{
		{Op: vm.PUSH, Arg: 27},
		{Op: vm.PUSH, Arg: 3},
		{Op: vm.ADD},
		{Op: vm.PUSH, Arg: 81},
		{Op: vm.MUL},
		{Op: vm.PUSH, Arg: 1},
		{Op: vm.PUSH, Arg: 1},
		{Op: vm.ADD},
		{Op: vm.HALT},
	}

	vm := vm.NewVM(code)
	_ = vm.Run()

	trace := *vm.Trace
	circuitInput := zk.TraceToCircuit(trace)

	debug.PrintTrace(trace)
	debug.PrintCircuitInputTable(circuitInput)
	debug.CompareTraceAndCircuit(trace, circuitInput)
}
