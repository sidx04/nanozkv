package zk

import "github.com/consensys/gnark/frontend"

type Step struct {
	// program counter validation
	PCBefore frontend.Variable
	PCAfter  frontend.Variable

	// stack pointer
	StackPointerBefore frontend.Variable
	StackPointerAfer   frontend.Variable

	// memory reads
	ReadAddr1 frontend.Variable
	ReadVal1  frontend.Variable

	ReadAddr2 frontend.Variable
	ReadVal2  frontend.Variable

	// memory write
	WriteAddr frontend.Variable
	WriteVal  frontend.Variable

	// push value
	Val frontend.Variable

	// check correct opcodes
	Opcode frontend.Variable

	// selectors
	IsAdd  frontend.Variable
	IsMul  frontend.Variable
	IsSub  frontend.Variable
	IsPush frontend.Variable
	IsNoop frontend.Variable
}

const (
	OP_PUSH  = 0
	OP_ADD   = 1
	OP_SUB   = 2
	OP_MUL   = 3
	OP_LOAD  = 4
	OP_STORE = 5
	OP_HALT  = 6
)

type TraceCircuit struct {
	Steps [MAX_STEPS]Step
}

func (tc *TraceCircuit) Define(api frontend.API) error {

	for _, s := range tc.Steps {
		// PC transition: pc_after = pc_before + 1
		api.AssertIsEqual(
			api.Mul(
				api.Add(s.IsAdd, s.IsMul, s.IsSub, s.IsPush),
				api.Sub(s.PCAfter, api.Add(s.PCBefore, 1)),
			),
			0,
		)

		// Opcode constraints
		api.AssertIsEqual(
			s.Opcode,
			api.Add(
				api.Add(
					api.Mul(s.IsAdd, OP_ADD),
					api.Mul(s.IsMul, OP_MUL),
				),
				api.Add(
					api.Mul(s.IsSub, OP_SUB),
					api.Mul(s.IsPush, OP_PUSH),
				),
				api.Mul(s.IsNoop, OP_HALT),
			),
		)

		// Arithmetic Constraints

		// ADD
		addConstraint := api.Sub(
			api.Add(s.ReadVal1, s.ReadVal2),
			s.WriteVal,
		)

		// MUL
		mulConstraint := api.Sub(
			api.Mul(s.ReadVal1, s.ReadVal2),
			s.WriteVal,
		)

		// SUB
		subConstraint := api.Sub(
			api.Sub(s.ReadVal1, s.ReadVal2),
			s.WriteVal,
		)

		// PUSH
		pushConstraint := api.Sub(s.WriteVal, s.Val)

		// unified constraint
		api.AssertIsEqual(
			api.Add(
				api.Add(
					api.Add(
						api.Mul(s.IsAdd, addConstraint),
						api.Mul(s.IsMul, mulConstraint),
					),
					api.Mul(s.IsSub, subConstraint),
				),
				api.Mul(s.IsPush, pushConstraint),
			),
			0,
		)

		// Stack Pointer Transitions

		// ADD/MUL/SUB: sp_after = sp_before - 1
		api.AssertIsEqual(
			api.Mul(
				api.Add(s.IsAdd, s.IsMul, s.IsSub),
				api.Sub(s.StackPointerAfer, api.Sub(s.StackPointerBefore, 1)),
			),
			0,
		)

		// PUSH: sp_after = sp_before + 1
		api.AssertIsEqual(
			api.Mul(
				s.IsPush,
				api.Sub(s.StackPointerAfer, api.Add(s.StackPointerBefore, 1)),
			),
			0,
		)

		// NOOP: sp unchanged
		api.AssertIsEqual(
			api.Mul(
				s.IsNoop,
				api.Sub(s.StackPointerAfer, s.StackPointerBefore),
			),
			0,
		)

		// ADDRESS RULES

		// ADD/MUL/SUB reads
		// ReadAddr1 = Top - 2
		// (isAdd + isMul + isSub) * (ReadAddr1 - (Top - 2)) = 0
		api.AssertIsEqual(
			api.Mul(
				api.Add(s.IsAdd, s.IsMul, s.IsSub),
				api.Sub(s.ReadAddr1, api.Sub(s.StackPointerBefore, 2)),
			),
			0,
		)

		// ReadAddr2 = Top - 1
		// (isAdd + isMul + isSub) * (ReadAddr2 - (Top - 1)) = 0
		api.AssertIsEqual(
			api.Mul(
				api.Add(s.IsAdd, s.IsMul, s.IsSub),
				api.Sub(s.ReadAddr2, api.Sub(s.StackPointerBefore, 1)),
			),
			0,
		)

		// write addr for binary ops
		// If opcode is ADD/MUL/SUB:
		// WriteAddr = Top - 2

		// say we have Stack = {0: 2, 1: 4, 2: ADD}
		// once ADD is encountered, we pop [0], [1] and write in [0]
		api.AssertIsEqual(
			api.Mul(
				api.Add(s.IsAdd, s.IsMul, s.IsSub),
				api.Sub(s.WriteAddr, api.Sub(s.StackPointerBefore, 2)),
			),
			0,
		)

		// PUSH write addr
		// If opcode is PUSH:
		// WriteAddr = SPBefore
		api.AssertIsEqual(
			api.Mul(
				s.IsPush,
				api.Sub(s.WriteAddr, s.StackPointerBefore),
			),
			0,
		)

		// SELECTOR CONSTRAINTS

		api.AssertIsEqual(
			api.Add(s.IsAdd, s.IsMul, s.IsSub, s.IsPush, s.IsNoop),
			1,
		)

		api.AssertIsEqual(api.Mul(s.IsAdd, api.Sub(s.IsAdd, 1)), 0)
		api.AssertIsEqual(api.Mul(s.IsMul, api.Sub(s.IsMul, 1)), 0)
		api.AssertIsEqual(api.Mul(s.IsSub, api.Sub(s.IsSub, 1)), 0)
		api.AssertIsEqual(api.Mul(s.IsPush, api.Sub(s.IsPush, 1)), 0)
		api.AssertIsEqual(api.Mul(s.IsNoop, api.Sub(s.IsNoop, 1)), 0)
	}

	for i := range len(tc.Steps) - 1 {
		curr := tc.Steps[i]
		next := tc.Steps[i+1]

		api.AssertIsEqual(curr.PCAfter, next.PCBefore)
	}

	return nil
}
