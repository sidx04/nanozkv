package zk

import (
	"nanozkv/vm"

	"github.com/consensys/gnark/frontend"
)

func TraceToCircuit(trace vm.ExecutionTrace) TraceCircuit {
	var steps [MAX_STEPS]Step

	i := 0

	for _, t := range trace.Steps {
		spBefore := len(t.StackBefore)
		spAfter := len(t.StackAfter)
		pcBefore := t.PC
		pcAfter := t.PC + 1

		op := t.Op

		var readAddr1, readAddr2, writeAddr int
		var readVal1, readVal2, writeVal int

		if len(t.StackBefore) >= 2 {
			readAddr1 = spBefore - 2
			readVal1 = t.StackBefore[spBefore-2]

			readAddr2 = spBefore - 1
			readVal2 = t.StackBefore[spBefore-1]
		}

		if len(t.StackAfter) > 0 {
			writeAddr = spAfter - 1
			writeVal = t.StackAfter[spAfter-1]
		}

		step := Step{
			PCBefore: pcBefore,
			PCAfter:  pcAfter,

			StackPointerBefore: spBefore,
			StackPointerAfer:   spAfter,

			ReadAddr1: readAddr1,
			ReadVal1:  readVal1,

			ReadAddr2: readAddr2,
			ReadVal2:  readVal2,

			WriteAddr: writeAddr,
			WriteVal:  writeVal,

			Val:    0,
			Opcode: op,

			IsAdd:  0,
			IsMul:  0,
			IsSub:  0,
			IsPush: 0,
			IsNoop: 0,
		}

		// selectors
		switch t.Op {
		case vm.ADD:
			step.IsAdd = 1
			step.Opcode = OP_ADD
		case vm.MUL:
			step.IsMul = 1
			step.Opcode = OP_MUL
		case vm.SUB:
			step.IsSub = 1
			step.Opcode = OP_SUB
		case vm.PUSH:
			step.IsPush = 1
			step.Val = t.Arg
			step.Opcode = OP_PUSH
		default:
			step.IsNoop = 1
			step.Opcode = OP_HALT
		}

		steps[i] = step
		i++
	}

	var lastSP frontend.Variable
	var lastPC frontend.Variable

	if i > 0 {
		lastSP = steps[i-1].StackPointerAfer
		lastPC = steps[i-1].PCAfter
	}
	// initialise remaining steps to 0 otherwise gnark cannot
	// deal with null values.
	for ; i < MAX_STEPS; i++ {
		steps[i] = Step{
			PCBefore: lastPC,
			PCAfter:  lastPC,

			StackPointerBefore: lastSP,
			StackPointerAfer:   lastSP,

			ReadAddr1: 0,
			ReadVal1:  0,
			ReadAddr2: 0,
			ReadVal2:  0,

			WriteAddr: 0,
			WriteVal:  0,

			Val:    0,
			Opcode: OP_HALT,

			IsAdd:  0,
			IsMul:  0,
			IsSub:  0,
			IsPush: 0,
			IsNoop: 1,
		}
	}

	return TraceCircuit{
		Steps: steps,
	}
}
