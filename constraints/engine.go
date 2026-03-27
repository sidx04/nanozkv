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
	switch step.Op {

	case vm.PUSH:
		return verifyPush(step)

	case vm.ADD:
		return verifyAdd(step)

	case vm.SUB:
		return verifySub(step)

	case vm.MUL:
		return verifyMul(step)

	case vm.LOAD:
		return verifyLoad(step)

	case vm.STORE:
		return verifyStore(step)

	default:
		return fmt.Errorf("unknown opcode")
	}
}

func verifyAdd(step vm.TraceStep) error {
	if len(step.StackBefore) < 2 {
		return fmt.Errorf("not enough elements for ADD")
	}

	if len(step.StackAfter) < 1 {
		return fmt.Errorf("invalid stack after ADD")
	}

	a := step.StackBefore[len(step.StackBefore)-2]
	b := step.StackBefore[len(step.StackBefore)-1]
	c := step.StackAfter[len(step.StackAfter)-1]

	if a+b != c {
		return fmt.Errorf("ADD constraint failed: %d + %d != %d", a, b, c)
	}

	return nil
}

func verifySub(step vm.TraceStep) error {
	if len(step.StackBefore) < 2 {
		return fmt.Errorf("not enough elements for SUB")
	}

	if len(step.StackAfter) < 1 {
		return fmt.Errorf("invalid stack after SUB")
	}

	a := step.StackBefore[len(step.StackBefore)-2]
	b := step.StackBefore[len(step.StackBefore)-1]
	c := step.StackAfter[len(step.StackAfter)-1]

	if a-b != c {
		return fmt.Errorf("SUB constraint failed")
	}

	return nil
}

func verifyMul(step vm.TraceStep) error {
	if len(step.StackBefore) < 2 {
		return fmt.Errorf("not enough elements for MUL")
	}

	if len(step.StackAfter) < 1 {
		return fmt.Errorf("invalid stack after MUL")
	}

	a := step.StackBefore[len(step.StackBefore)-2]
	b := step.StackBefore[len(step.StackBefore)-1]
	c := step.StackAfter[len(step.StackAfter)-1]

	if a*b != c {
		return fmt.Errorf("MUL constraint failed")
	}

	return nil
}

func verifyPush(step vm.TraceStep) error {
	if len(step.StackAfter) != len(step.StackBefore)+1 {
		return fmt.Errorf("PUSH should increase stack size by 1")
	}
	return nil
}

func verifyStore(step vm.TraceStep) error {
	if len(step.StackAfter) != len(step.StackBefore)-1 {
		return fmt.Errorf("STORE should pop stack")
	}
	return nil
}

func verifyLoad(step vm.TraceStep) error {
	if len(step.StackAfter) != len(step.StackBefore)+1 {
		return fmt.Errorf("LOAD should push to stack")
	}
	return nil
}
