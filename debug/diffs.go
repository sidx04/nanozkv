package debug

import (
	"fmt"
	"nanozkv/vm"
	"nanozkv/zk"
)

func CompareTraceAndCircuit(trace vm.ExecutionTrace, circuit zk.TraceCircuit) {
	fmt.Println("==== TRACE vs CIRCUIT DIFF ====")

	for i := 0; i < len(trace.Steps); i++ {
		t := trace.Steps[i]
		s := circuit.Steps[i]

		// Check opcode
		if int(s.Opcode.(int)) != int(t.Op) {
			fmt.Printf("❌ Step %d Opcode mismatch: trace=%v circuit=%v\n",
				i, t.Op, s.Opcode)
		}

		// Check SP
		if int(s.StackPointerBefore.(int)) != len(t.StackBefore) {
			fmt.Printf("❌ Step %d SPBefore mismatch\n", i)
		}

		if len(t.StackAfter) > 0 {
			expected := t.StackAfter[len(t.StackAfter)-1]

			if int(s.WriteVal.(int)) != expected {
				fmt.Printf("❌ Step %d WriteVal mismatch: expected=%d got=%v\n",
					i, expected, s.WriteVal)
			}
		}
	}
}
