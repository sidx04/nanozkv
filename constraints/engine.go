package constraints

import (
	"fmt"
	"nanozkv/vm"
)

type ConstraintEngine struct{}

func NewConstraintEngine() *ConstraintEngine {
	return &ConstraintEngine{}
}

func (ce *ConstraintEngine) Verify(trace vm.ExecutionTrace) error {
	for i, step := range trace.Steps {
		if err := ce.verifyStep(step); err != nil {
			return fmt.Errorf("step %d failed %w", i, err)
		}
	}
	return nil
}

func (ce *ConstraintEngine) verifyStep(step vm.TraceStep) error {
	// skip non-arithmetic ops for now
	if len(step.StackBefore) < 2 {
		return nil
	}

	a := step.StackBefore[len(step.StackBefore)-2]
	b := step.StackBefore[len(step.StackBefore)-1]

	if len(step.StackAfter) < 1 {
		return fmt.Errorf("invalid stack after")
	}

	c := step.StackAfter[len(step.StackAfter)-1]

	isAdd, isMul, isSub := getSelectors(step.Op)

	val := isAdd*(a+b-c) + isMul*(a*b-c) + isSub*(a-b-c)

	if val != 0 {
		return fmt.Errorf("constraint failed: got %d", val)
	}

	return nil
}

func getSelectors(op vm.Opcode) (isAdd, isMul, isSub int) {
	switch op {
	case vm.ADD:
		return 1, 0, 0
	case vm.MUL:
		return 0, 1, 0
	case vm.SUB:
		return 0, 0, 1
	default:
		return 0, 0, 0
	}
}
