package vm

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	code := []Instruction{
		{Op: PUSH, Arg: 3},
		{Op: PUSH, Arg: 4},
		{Op: ADD},
		{Op: HALT},
	}

	vm := NewVM(code)
	err := vm.Run()

	if err != nil {
		t.Fatal(err)
	}

	if len(vm.Stack) != 1 || vm.Stack[0] != 7 {
		t.Fatalf("expected 7, got %v", vm.Stack)
	}

}

func TestMemory(t *testing.T) {
	code := []Instruction{
		{Op: PUSH, Arg: 42},
		{Op: STORE, Arg: 1},
		{Op: LOAD, Arg: 1},
		{Op: HALT},
	}

	vm := NewVM(code)
	err := vm.Run()
	if err != nil {
		t.Fatal(err)
	}

	if vm.Stack[0] != 42 {
		t.Fatalf("expected 42, got %v", vm.Stack[0])
	}
}

func TestStackUnderflow(t *testing.T) {
	code := []Instruction{
		{Op: ADD}, // invalid
		{Op: HALT},
	}

	vm := NewVM(code)
	err := vm.Run()

	if err == nil {
		t.Fatal("expected error for stack underflow")
	}
}

func TestTraceRecording(t *testing.T) {
	code := []Instruction{
		{Op: PUSH, Arg: 2},
		{Op: PUSH, Arg: 3},
		{Op: ADD},
		{Op: HALT},
	}

	vm := NewVM(code)
	err := vm.Run()
	if err != nil {
		t.Fatal(err)
	}

	trace := vm.Trace.Steps

	// Trace, without considering HALT, should look something like:
	// [{0 PUSH} {1 PUSH} {2 ADD}]
	if len(trace) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(trace))
	}

	fmt.Println(trace)

	// Check ADD step
	addStep := trace[2]

	if addStep.Op != ADD {
		t.Fatalf("expected ADD, got %v", addStep.Op)
	}

	if len(addStep.StackBefore) != 2 {
		t.Fatalf("expected stack size 2 before ADD")
	}

	if len(addStep.StackAfter) != 1 || addStep.StackAfter[0] != 5 {
		t.Fatalf("expected result 5, got %v", addStep.StackAfter)
	}
}
