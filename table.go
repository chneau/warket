package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chneau/warket/client"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

type Row struct {
	ItemName string
	Quantity int
	Place    int
	Price    float64
	Diff     float64
	Bests    []float64
}

type Table struct {
	OrderType string
	Rows      []Row
}

// Suppose we have ll as new, other as old
func (t Table) Diff(otherTable *Table) (result []string) {
	if otherTable == nil {
		return
	}
	rowsByItemName := map[string]Row{}
	for _, row := range t.Rows {
		rowsByItemName[row.ItemName] = row
	}
	for _, otherRow := range otherTable.Rows {
		if row, exists := rowsByItemName[otherRow.ItemName]; exists { // o old, n new
			// it exists, needs to diff
			quantity := row.Quantity != otherRow.Quantity
			place := row.Place != otherRow.Place
			price := row.Price != otherRow.Price
			diff := row.Diff != otherRow.Diff
			if quantity || place || price || diff {
				res := fmt.Sprintf("%s is now ", color.MagentaString(row.ItemName))
				logs := []string{}
				if quantity {
					logs = append(logs, color.WhiteString("# %v -> %v", otherRow.Quantity, row.Quantity))
				}
				if place {
					logs = append(logs, color.GreenString("N° %v -> %v", otherRow.Place, row.Place))
				}
				if price {
					logs = append(logs, color.HiBlueString("$ %v -> %v", otherRow.Price, row.Price))
				}
				if diff {
					logs = append(logs, color.HiRedString("± %v -> %v", otherRow.Diff, row.Diff))
				}
				res += strings.Join(logs, ", ")
				res += "."
				result = append(result, res)
			}
			delete(rowsByItemName, otherRow.ItemName)
		} else {
			// it doesnt exists anymore
			result = append(result, fmt.Sprintf("%s is no more.", otherRow.ItemName))
		}
	}
	for _, n := range rowsByItemName { // new stuff
		result = append(result, fmt.Sprintf("%s is new.", color.MagentaString(n.ItemName)))
	}
	now := time.Now()
	for i := range result {
		result[i] = now.Format("15:04:05") + " " + t.OrderType + " " + result[i]
	}
	return result
}
func (t Table) String() string {
	if len(t.Rows) == 0 {
		return ""
	}
	data := [][]string{}
	for _, row := range t.Rows {
		diff := strconv.FormatFloat(row.Diff, 'f', 0, 64)
		if row.Diff > 0 {
			diff = "+" + diff
		}
		data = append(data, []string{
			row.ItemName,
			strconv.Itoa(row.Quantity),
			strconv.Itoa(row.Place),
			strconv.FormatFloat(row.Price, 'f', 0, 64),
			diff,
			strings.Join(floatsToStrings(row.Bests), " "),
		})
	}
	reader, writer := io.Pipe()
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Item", "#", "N°", "$", "±", "best " + t.OrderType})
	table.SetColWidth(math.MaxInt64)
	table.SetBorder(false)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.FgHiRedColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
	)
	table.AppendBulk(data)
	result := make(chan string)
	go func() {
		b, err := ioutil.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		result <- string(b)
	}()
	table.Render()
	writer.Close()
	return <-result
}

func (t Table) Sort(sorting string) {
	switch sorting {
	case "qtt":
		sort.Slice(t.Rows, func(i, j int) bool {
			if t.Rows[i].Quantity == t.Rows[j].Quantity {
				return t.Rows[i].Diff > t.Rows[j].Diff
			}
			return t.Rows[i].Quantity > t.Rows[j].Quantity
		})
	case "place":
		sort.Slice(t.Rows, func(i, j int) bool {
			if t.Rows[i].Place == t.Rows[j].Place {
				return t.Rows[i].Diff > t.Rows[j].Diff
			}
			return t.Rows[i].Place > t.Rows[j].Place
		})
	case "price":
		fallthrough
	case "plat":
		sort.Slice(t.Rows, func(i, j int) bool {
			if t.Rows[i].Price == t.Rows[j].Price {
				return t.Rows[i].Diff > t.Rows[j].Diff
			}
			return t.Rows[i].Price > t.Rows[j].Price
		})
	case "diff":
		sort.Slice(t.Rows, func(i, j int) bool {
			return t.Rows[i].Diff > t.Rows[j].Diff
		})
	case "name":
		fallthrough
	default:
		sort.Slice(t.Rows, func(i, j int) bool {
			if t.Rows[i].ItemName == t.Rows[j].ItemName {
				return t.Rows[i].Diff > t.Rows[j].Diff
			}
			return t.Rows[i].ItemName < t.Rows[j].ItemName
		})
	}
}

func NewTable(username string, orders []client.Order, orderType string) Table {
	t := Table{
		OrderType: orderType,
	}
	if len(orders) > 0 {
		for _, i := range orders {
			if i.OrderType != orderType {
				continue
			}
			orders, err := client.FetchItemOrders(i.Item.URLName)
			if err != nil {
				panic(err)
			}
			all := []float64{}
			position := 1
			for _, o := range orders {
				if o.User.Status == "ingame" && o.OrderType == i.OrderType && o.ModRank == i.ModRank && o.User.IngameName != username && o.Region == i.Region {
					all = append(all, o.Platinum)
					if i.OrderType == "buy" {
						if o.Platinum > i.Platinum {
							position++
						}
					} else {
						if i.Platinum > o.Platinum {
							position++
						}
					}
				}
			}
			sort.Float64s(all)
			if i.OrderType == "buy" {
				for i := len(all)/2 - 1; i >= 0; i-- { // reverse
					opp := len(all) - 1 - i
					all[i], all[opp] = all[opp], all[i]
				}
			}
			diff := 0.
			if len(all) > 0 {
				diff = i.Platinum - all[0]
			}
			itemName := i.Item.Info.ItemName
			if i.ModRank != 0 {
				itemName = itemName + " " + strconv.Itoa(i.ModRank)
			}
			five := int(math.Min(float64(len(all)), 5.))
			l := Row{
				ItemName: itemName,
				Quantity: i.Quantity,
				Place:    position,
				Price:    i.Platinum,
				Diff:     diff,
				Bests:    all[:five],
			}
			t.Rows = append(t.Rows, l)
		}
	}
	return t
}
