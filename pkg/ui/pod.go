package ui

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
	v1 "k8s.io/api/core/v1"
)

func PodList(pods []v1.Pod) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Namespace", "Name")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, p := range pods {
		tbl.AddRow(p.Namespace, p.Name)
	}

	tbl.Print()
}
