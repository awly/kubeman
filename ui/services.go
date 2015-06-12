package ui

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/alytvynov/termui"
)

func servicesTab(ui *UI) tab {
	return &listTab{
		ui:       ui,
		mu:       &sync.Mutex{},
		headers:  serviceHeaders,
		itemType: reflect.TypeOf(serviceItem{}),
	}
}

var serviceHeaders = []header{
	{"name", 4},
	{"type", 2},
	{"ip", 3},
	{"port", 3},
}

type serviceItem struct {
	service api.Service
	ui      *UI
}

func (s serviceItem) toRows() []*termui.Row {
	lname := label(s.service.Name)
	ltype := label(string(s.service.Spec.Type))
	lip := label(s.service.Spec.ClusterIP)

	if len(s.service.Spec.DeprecatedPublicIPs) > 0 {
		lip.Text = s.service.Spec.DeprecatedPublicIPs[0]
	}

	rows := make([]*termui.Row, 0, len(s.service.Spec.Ports))
	for i, p := range s.service.Spec.Ports {
		if i > 0 {
			lname = label("")
			ltype = label("")
			lip = label("")
		}
		lport := label(fmt.Sprintf("%d -> %s", p.Port, p.TargetPort.String()))
		rows = append(rows, termui.NewRow(
			termui.NewCol(4, 0, lname),
			termui.NewCol(2, 0, ltype),
			termui.NewCol(3, 0, lip),
			termui.NewCol(3, 0, lport),
		))
	}
	return rows
}

func (p *serviceItem) init(ui *UI)                { p.ui = ui }
func (p *serviceItem) setData(d interface{})      { p.service = *d.(*api.Service) }
func (p serviceItem) sameData(d interface{}) bool { return p.service.Name == (*d.(*api.Service)).Name }
func (p serviceItem) less(i listItem) bool        { return p.service.Name < i.(*serviceItem).service.Name }

func (p serviceItem) handleEvent(e termui.Event) {}
