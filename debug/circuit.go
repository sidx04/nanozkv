package debug

import (
	"fmt"
	"os"

	"nanozkv/zk"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintCircuitInputTable(c zk.TraceCircuit) {
	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)

	tw.SetTitle("ZKVM CIRCUIT TRACE TABLE")

	tw.AppendHeader(table.Row{
		"Step",
		"PC",
		"TOP",
		"Opcode",
		"Read1 (addr,val)",
		"Read2 (addr,val)",
		"Write (addr,val)",
		"Selectors",
	})

	for i, s := range c.Steps {

		// skip padding rows
		if s.IsNoop == 1 {
			continue
		}

		read1 := fmt.Sprintf("[%v, %v]", s.ReadAddr1, s.ReadVal1)
		read2 := fmt.Sprintf("[%v, %v]", s.ReadAddr2, s.ReadVal2)
		write := fmt.Sprintf("[%v, %v]", s.WriteAddr, s.WriteVal)

		selectors := fmt.Sprintf(
			"A:%v M:%v S:%v P:%v N:%v",
			s.IsAdd, s.IsMul, s.IsSub, s.IsPush, s.IsNoop,
		)

		tw.AppendRow(table.Row{
			i,
			s.PCBefore,
			s.StackPointerBefore,
			s.Opcode,
			read1,
			read2,
			write,
			selectors,
		})
	}

	tw.Render()
}
