package ui

import (
	"log"
	"sort"
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
	"github.com/gizak/termui"
)

type nodesTab struct {
	log *log.Logger

	mu    *sync.Mutex
	nodes []node
}

func (rt *nodesTab) Len() int           { return len(rt.nodes) }
func (rt *nodesTab) Less(i, j int) bool { return rt.nodes[i].node.Name < rt.nodes[j].node.Name }
func (rt *nodesTab) Swap(i, j int)      { rt.nodes[i], rt.nodes[j] = rt.nodes[j], rt.nodes[i] }

func (rt *nodesTab) update(e Event) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	p := *e.Data.(*api.Node)
	switch e.Type {
	case watch.Added:
		rt.nodes = append(rt.nodes, node{node: p})
	case watch.Modified:
		found := false
		for i, up := range rt.nodes {
			if up.node.Name == p.Name {
				found = true
				rt.nodes[i].node = p
				break
			}
		}
		if !found {
			rt.nodes = append(rt.nodes, node{node: p})
		}
	case watch.Deleted:
		for i, up := range rt.nodes {
			if up.node.Name == p.Name {
				rt.Swap(i, rt.Len()-1)
				rt.nodes = rt.nodes[:rt.Len()-1]
				break
			}
		}
	}
	sort.Sort(rt)
}

func (rt *nodesTab) toRows() []*termui.Row {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rows := make([]*termui.Row, 0)

	// header
	lname := header("name")
	lstate := header("state")
	laddr := header("ip")
	lcpu := header("cpu")
	lmem := header("mem")

	rows = append(rows, termui.NewRow(
		termui.NewCol(3, 0, lname),
		termui.NewCol(1, 0, lstate),
		termui.NewCol(1, 0, laddr),
		termui.NewCol(1, 0, lcpu),
		termui.NewCol(1, 0, lmem),
	))
	for _, r := range rt.nodes {
		rows = append(rows, r.toRow())
	}
	return rows
}

type node struct {
	node api.Node
}

func (r node) toRow() *termui.Row {
	lname := label(r.node.Name)
	lcpu := label(r.node.Status.Capacity.Cpu().String())
	lmem := label(r.node.Status.Capacity.Memory().String())

	var state string
	for _, s := range r.node.Status.Conditions {
		if s.Status == api.ConditionTrue {
			state = string(s.Type)
		}
	}
	lstate := label(state)

	var addr string
	for _, a := range r.node.Status.Addresses {
		if a.Type == api.NodeExternalIP {
			addr = a.Address
		}
	}
	laddr := label(addr)

	return termui.NewRow(
		termui.NewCol(3, 0, lname),
		termui.NewCol(1, 0, lstate),
		termui.NewCol(1, 0, laddr),
		termui.NewCol(1, 0, lcpu),
		termui.NewCol(1, 0, lmem),
	)
}
