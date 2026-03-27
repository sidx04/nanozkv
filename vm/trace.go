package vm

type TraceStep struct {
	PC int
	Op Opcode

	// Argument to certain opcodes like PUSH, MSTORE, etc.
	Arg int

	// Stack snapshots
	StackBefore []int
	StackAfter  []int

	// Memory access info
	MemAddr  int
	MemValue int
}

type ExecutionTrace struct {
	Steps []TraceStep
}

func (t *ExecutionTrace) Add(step TraceStep) {
	t.Steps = append(t.Steps, step)
}
