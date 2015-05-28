package ui

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/gizak/termui"
)

type podList struct {
	pods []pod
}

func (pl podList) Len() int           { return len(pl.pods) }
func (pl podList) Less(i, j int) bool { return pl.pods[i].p.Name < pl.pods[j].p.Name }
func (pl podList) Swap(i, j int)      { pl.pods[i], pl.pods[j] = pl.pods[j], pl.pods[i] }

func (pl podList) toRows() []*termui.Row {
	rows := make([]*termui.Row, 0, len(pl.pods)+1)

	// header
	lname := label("name")
	lname.TextFgColor = termui.ColorWhite | termui.AttrBold
	lstatus := label("status")
	lstatus.TextFgColor = termui.ColorWhite | termui.AttrBold
	lhost := label("host")
	lhost.TextFgColor = termui.ColorWhite | termui.AttrBold

	rows = append(rows, termui.NewRow(
		termui.NewCol(2, 0, lname),
		termui.NewCol(1, 0, lstatus),
		termui.NewCol(1, 0, lhost),
	))
	for _, p := range pl.pods {
		rows = append(rows, p.toRow())
	}
	return rows
}

type pod struct {
	p *api.Pod
}

func (pr pod) toRow() *termui.Row {
	lname := label(pr.p.Name)
	lstatus := label(string(pr.p.Status.Phase))
	lhost := label(pr.p.Spec.Host)

	return termui.NewRow(
		termui.NewCol(2, 0, lname),
		termui.NewCol(1, 0, lstatus),
		termui.NewCol(1, 0, lhost),
	)
}
