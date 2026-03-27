package constraints

import (
	"nanozkv/vm"
	"testing"
)

func TestConstraintEngineValid(t *testing.T) {
	code := []vm.Instruction{
		{Op: vm.PUSH, Arg: 2},
		{Op: vm.PUSH, Arg: 3},
		{Op: vm.ADD},
		{Op: vm.HALT},
	}

	m := vm.NewVM(code)
	err := m.Run()
	if err != nil {
		t.Fatal(err)
	}

	engine := NewConstraintEngine()
	err = engine.Verify(*m.Trace)

	if err != nil {
		t.Fatalf("expected valid trace, got error: %v", err)
	}
}

func TestConstraintEngineInvalid(t *testing.T) {
	trace := vm.ExecutionTrace{
		Steps: []vm.TraceStep{
			{
				Op:          vm.ADD,
				StackBefore: []int{2, 3},
				StackAfter:  []int{10}, // WRONG
			},
		},
	}

	engine := NewConstraintEngine()
	err := engine.Verify(trace)

	if err == nil {
		t.Fatal("expected constraint failure")
	}
}
