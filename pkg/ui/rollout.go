package ui

import (
	"github.com/anarcher/kroller/pkg/resource"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func RolloutList(rl resource.RolloutList) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Namespace", "Kind", "Name")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, r := range rl {
		tbl.AddRow(r.Namespace(), r.Kind(), r.Name())
	}

	tbl.Print()
}
