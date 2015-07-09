package ui

import (
	"log"

	"github.com/alytvynov/termui"
)

type statusBar struct {
	ui      *UI
	content string
}

type statusUpdate struct {
	msg string
}

func (s *statusBar) dataUpdate(e Event) {
	d := e.Data.(statusUpdate)
	s.content = d.msg
	log.Println("status:", s.content)
	go s.ui.redrawStatus()
}

func (s *statusBar) toRows() []*termui.Row {
	l := label(s.content)
	l.TextFgColor = termui.ColorWhite | termui.AttrBold
	return []*termui.Row{termui.NewRow(
		termui.NewCol(12, 0, l),
	)}
}

func (s *statusBar) uiUpdate(termui.Event) {}
func (s *statusBar) clean()                {}
