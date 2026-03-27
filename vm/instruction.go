package vm

type Opcode int

const (
	PUSH Opcode = iota
	ADD
	SUB
	MUL
	LOAD
	STORE
	HALT
)

type Instruction struct {
	// Codes to instructions
	Op Opcode
	// Argument to an opcode, for instructions like PUSH, STORE, LOAD
	Arg int
}
