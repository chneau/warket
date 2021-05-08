package display

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buger/goterm"
	"github.com/chneau/warket/client"
	"github.com/mattn/go-colorable"
	"github.com/schollz/progressbar/v3"
)

func floatsToStrings(floats []float64) (strings []string) {
	for i := range floats {
		strings = append(strings, strconv.FormatFloat(floats[i], 'f', 0, 64))
	}
	return
}

type TableContext struct {
	Username    string
	Sells       bool
	Buys        bool
	LinesOfLogs int
	Sorting     string
	TicksSecond int
	LastSells   *Table
	LastBuys    *Table
	Logs        []string
}

func (tc *TableContext) Run() {
	goterm.Output = bufio.NewWriter(colorable.NewColorable(os.Stdout))

	if tc.TicksSecond < 0 {
		goterm.Print(tc.PrepareTable())
		goterm.Flush()
		return
	}
	for {
		goterm.Clear()
		goterm.MoveCursor(1, 1)
		goterm.Print(tc.PrepareTable())
		goterm.Flush()
		bar := progressbar.New(tc.TicksSecond)
		_ = bar.RenderBlank()
		for i := 0; i < tc.TicksSecond; i++ {
			time.Sleep(time.Second)
			_ = bar.Add(1)
		}
		_ = bar.Clear()
	}
}

func (tc *TableContext) PrepareTable() string {
	buys, sells, err := client.FetchUserOrders(tc.Username)
	if err != nil {
		panic(err)
	}
	all := append(buys, sells...)
	result := ""
	if tc.Sells {
		lines := NewTable(tc.Username, all, "sell")
		tc.Logs = append(lines.Diff(tc.LastSells), tc.Logs...)
		tc.LastSells = &lines
		lines.Sort(tc.Sorting)
		result += lines.String()
	}
	if tc.Buys {
		lines := NewTable(tc.Username, all, "buy")
		tc.Logs = append(lines.Diff(tc.LastBuys), tc.Logs...)
		tc.LastBuys = &lines
		lines.Sort(tc.Sorting)
		result += lines.String()
	}
	if len(tc.Logs) > tc.LinesOfLogs {
		tc.Logs = tc.Logs[:tc.LinesOfLogs]
	}
	if tc.LinesOfLogs > 0 {
		result += "\n" + strings.Join(tc.Logs, "\n") + "\n"
	}
	return result
}
