package main

import (
	"fmt"
	"sort"

	ui "github.com/gizak/termui"
)

type Grid struct {
	cursorPos  uint
	containers map[string]*Container
}

func (g *Grid) AddContainer(id string, names []string) {
	g.containers[id] = NewContainer(id, names)
}

// Return number of containers/rows
func (g *Grid) Len() uint {
	return uint(len(g.containers))
}

// Return sorted list of active container IDs
func (g *Grid) CIDs() []string {
	var ids []string
	for id, _ := range g.containers {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Redraw the cursor with the currently selected row
func (g *Grid) Cursor() {
	for n, id := range g.CIDs() {
		c := g.containers[id]
		if uint(n) == g.cursorPos {
			c.widgets.cid.TextFgColor = ui.ColorDefault
			c.widgets.cid.TextBgColor = ui.ColorWhite
		} else {
			c.widgets.cid.TextFgColor = ui.ColorWhite
			c.widgets.cid.TextBgColor = ui.ColorDefault
		}
	}
	ui.Render(ui.Body)
}

func (g *Grid) Rows() (rows []*ui.Row) {
	for _, cid := range g.CIDs() {
		c := g.containers[cid]
		rows = append(rows, c.widgets.MakeRow())
	}
	return rows
}

func header() *ui.Row {
	return ui.NewRow(
		ui.NewCol(1, 0, headerPar("CID")),
		ui.NewCol(2, 0, headerPar("CPU")),
		ui.NewCol(2, 0, headerPar("MEM")),
		ui.NewCol(2, 0, headerPar("NET RX/TX")),
		ui.NewCol(2, 0, headerPar("NAMES")),
	)
}

func headerPar(s string) *ui.Par {
	p := ui.NewPar(fmt.Sprintf(" %s", s))
	p.Border = false
	p.Height = 2
	p.Width = 20
	p.TextFgColor = ui.ColorWhite
	return p
}

func Display(g *Grid) {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// build layout
	ui.Body.AddRows(header())

	for _, row := range g.Rows() {
		ui.Body.AddRows(row)
	}

	// calculate layout
	ui.Body.Align()
	g.Cursor()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		if g.cursorPos > 0 {
			g.cursorPos -= 1
			g.Cursor()
		}
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		if g.cursorPos < (g.Len() - 1) {
			g.cursorPos += 1
			g.Cursor()
		}
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()
}