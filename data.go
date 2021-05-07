package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buger/goterm"
	"github.com/chneau/warket/pkg/client"
	"github.com/mattn/go-colorable"
	"github.com/schollz/progressbar/v3"
)

func floatsToStrings(floats []float64) (strings []string) {
	for i := range floats {
		strings = append(strings, strconv.FormatFloat(floats[i], 'f', 0, 64))
	}
	return
}

type data struct {
	username  string
	sells     bool
	buys      bool
	logging   int
	sorting   string
	t         int
	lastSells *lines
	lastBuys  *lines
	logs      []string
}

func (d *data) run() {
	goterm.Output = bufio.NewWriter(colorable.NewColorable(os.Stdout))

	if d.t < 0 {
		goterm.Print(d.prepare())
		goterm.Flush()
		return
	}
	for {
		goterm.Clear()
		goterm.MoveCursor(1, 1)
		goterm.Print(d.prepare())
		goterm.Flush()
		bar := progressbar.New(d.t)
		_ = bar.RenderBlank()
		for i := 0; i < d.t; i++ {
			time.Sleep(time.Second)
			_ = bar.Add(1)
		}
		_ = bar.Clear()
	}
}

func (d *data) prepare() string {
	buys, sells, err := client.FetchUserOrders(d.username)
	if err != nil {
		panic(err)
	}
	all := append(buys, sells...)
	result := ""
	if d.sells {
		lines := newLines(d.username, all, "sell")
		d.logs = append(lines.diff(d.lastSells), d.logs...)
		d.lastSells = &lines
		lines.sort(d.sorting)
		result += lines.String()
	}
	if d.buys {
		lines := newLines(d.username, all, "buy")
		d.logs = append(lines.diff(d.lastBuys), d.logs...)
		d.lastBuys = &lines
		lines.sort(d.sorting)
		result += lines.String()
	}
	if len(d.logs) > d.logging {
		d.logs = d.logs[:d.logging]
	}
	if d.logging > 0 {
		result += "\n" + strings.Join(d.logs, "\n") + "\n"
	}
	return result
}

func (d *data) String() string {
	return ""
}
