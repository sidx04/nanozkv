package debug

import (
	"fmt"
	"nanozkv/vm"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintTrace(trace vm.ExecutionTrace) {
	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)

	tw.SetTitle("ZKVM VM TRACE")

	tw.AppendHeader(table.Row{
		"Step",
		"PC",
		"NextPC",
		"Opcode",
		"Arg",
		"StackBefore",
		"StackAfter",
	})

	for i, t := range trace.Steps {

		stackBefore := fmt.Sprintf("%v", t.StackBefore)
		stackAfter := fmt.Sprintf("%v", t.StackAfter)

		stackBefore = truncate(stackBefore, 30)
		stackAfter = truncate(stackAfter, 30)

		tw.AppendRow(table.Row{
			i,
			t.PC,
			t.PC + 1,
			t.Op,
			t.Arg,
			stackBefore,
			stackAfter,
		})
	}

	tw.Render()
}

// keeps table clean
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
