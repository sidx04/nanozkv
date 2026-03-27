package vm

import "fmt"

type VM struct {
	PC     int
	Stack  []int
	Memory map[int]int
	Code   []Instruction
	Halted bool

	Trace *ExecutionTrace
}

func NewVM(code []Instruction) *VM {
	return &VM{
		PC:     0,
		Stack:  []int{},
		Memory: make(map[int]int),
		Code:   code,
		Halted: false,
		Trace:  &ExecutionTrace{},
	}
}

func (vm *VM) push(val int) {
	vm.Stack = append(vm.Stack, val)
}

func (vm *VM) pop() (int, error) {
	if len(vm.Stack) == 0 {
		return 0, fmt.Errorf("stack underflow")
	}
	val := vm.Stack[len(vm.Stack)-1]
	vm.Stack = vm.Stack[:len(vm.Stack)-1]
	return val, nil
}

func (vm *VM) Run() error {
	for !vm.Halted {
		if err := vm.Step(); err != nil {
			return err
		}
	}
	return nil
}

func (vm *VM) Step() error {
	if vm.Halted {
		return nil
	}

	if vm.PC >= len(vm.Code) {
		return fmt.Errorf("pc out of bounds")
	}

	inst := vm.Code[vm.PC]

	if err := vm.executeInstruction(inst); err != nil {
		return err
	}

	vm.PC++

	return nil
}

func (vm *VM) executeInstruction(inst Instruction) error {
	beforeStack := copyStack(vm.Stack)

	memAddr := 0
	memValue := 0

	switch inst.Op {

	case PUSH:
		vm.push(inst.Arg)

	case ADD:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(a + b)

	case SUB:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(a - b)

	case MUL:
		b, err := vm.pop()
		if err != nil {
			return err
		}
		a, err := vm.pop()
		if err != nil {
			return err
		}
		vm.push(a * b)

	case LOAD:
		memAddr = inst.Arg
		memValue = vm.Memory[memAddr]
		vm.push(memValue)

	case STORE:
		val, err := vm.pop()
		if err != nil {
			return err
		}
		memAddr = inst.Arg
		memValue = val
		vm.Memory[memAddr] = val

	case HALT:
		vm.Halted = true
		return nil

	default:
		return fmt.Errorf("unknown opcode")
	}

	afterStack := copyStack(vm.Stack)

	vm.Trace.Add(TraceStep{
		PC:  vm.PC,
		Op:  inst.Op,
		Arg: inst.Arg,

		StackBefore: beforeStack,
		StackAfter:  afterStack,

		MemAddr:  memAddr,
		MemValue: memValue,
	})

	return nil
}

func copyStack(stack []int) []int {
	cp := make([]int, len(stack))
	copy(cp, stack)
	return cp
}
