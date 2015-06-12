package ui

import (
	"reflect"
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/alytvynov/termui"
)

func nodesTab(ui *UI) tab {
	return &listTab{
		ui:       ui,
		mu:       &sync.Mutex{},
		headers:  nodeHeaders,
		itemType: reflect.TypeOf(nodeItem{}),
	}
}

var nodeHeaders = []header{
	{"name", 4},
	{"state", 2},
	{"ip", 2},
	{"cpu", 2},
	{"mem", 2},
}

type nodeItem struct {
	node api.Node
	ui   *UI
}

func (r nodeItem) toRows() []*termui.Row {
	lname := label(r.node.Name)
	lcpu := label(r.node.Status.Capacity.Cpu().String())
	lmem := label(r.node.Status.Capacity.Memory().String())

	lstate := label(string(r.node.Status.Phase))
	switch r.node.Status.Phase {
	case api.NodeRunning:
		lstate.TextFgColor = termui.ColorGreen
	case api.NodePending:
		lstate.TextFgColor = termui.ColorYellow
	case api.NodeTerminated:
		lstate.TextFgColor = termui.ColorRed
	}
	if lstate.Text == "" {
		for _, c := range r.node.Status.Conditions {
			if c.Status == api.ConditionTrue {
				lstate.Text = string(c.Type)
				if c.Type == api.NodeReady {
					lstate.TextFgColor = termui.ColorGreen
				}
			}
		}
	}

	var addr string
	for _, a := range r.node.Status.Addresses {
		if a.Type == api.NodeExternalIP {
			addr = a.Address
		}
	}
	laddr := label(addr)

	return []*termui.Row{termui.NewRow(
		termui.NewCol(4, 0, lname),
		termui.NewCol(2, 0, lstate),
		termui.NewCol(2, 0, laddr),
		termui.NewCol(2, 0, lcpu),
		termui.NewCol(2, 0, lmem),
	)}
}

func (p *nodeItem) init(ui *UI)                { p.ui = ui }
func (p *nodeItem) setData(d interface{})      { p.node = *d.(*api.Node) }
func (p nodeItem) sameData(d interface{}) bool { return p.node.Name == (*d.(*api.Node)).Name }
func (p nodeItem) less(i listItem) bool        { return p.node.Name < i.(*nodeItem).node.Name }

func (p nodeItem) handleEvent(e termui.Event) {}
