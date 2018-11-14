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

	"github.com/fatih/color"

	"github.com/bradfitz/slice"
	"github.com/chneau/warket/pkg/client"
	"github.com/olekukonko/tablewriter"
)

type line struct {
	item     string
	quantity int
	place    int
	price    int
	diff     int
	best     []int
}

type lines struct {
	orderType string
	ll        []line
}

// Suppose we have ll as new, other as old
func (ll lines) diff(other *lines) (result []string) {
	if other == nil {
		return
	}
	llmaps := map[string]line{}
	for _, l := range ll.ll {
		llmaps[l.item] = l
	}
	for _, o := range other.ll {
		if n, exist := llmaps[o.item]; exist { // o old, n new
			// it exists, needs to diff
			qtt := n.quantity != o.quantity
			place := n.place != o.place
			price := n.price != o.price
			diff := n.diff != o.diff
			if qtt || place || price || diff {
				res := fmt.Sprintf("%s is now ", color.MagentaString(n.item))
				things := []string{}
				if qtt {
					things = append(things, color.WhiteString("# %d -> %d", o.quantity, n.quantity))
				}
				if place {
					things = append(things, color.GreenString("N° %d -> %d", o.place, n.place))
				}
				if price {
					things = append(things, color.HiBlueString("$ %d -> %d", o.price, n.price))
				}
				if diff {
					things = append(things, color.HiRedString("± %d -> %d", o.diff, n.diff))
				}
				res += strings.Join(things, ", ")
				res += "."
				result = append(result, res)
			}
			delete(llmaps, o.item)
		} else {
			// it doesnt exists anymore
			result = append(result, fmt.Sprintf("%s is no more.", o.item))
		}
	}
	for _, n := range llmaps { // new stuff
		result = append(result, fmt.Sprintf("%s is new.", color.MagentaString(n.item)))
	}
	now := time.Now()
	for i := range result {
		result[i] = now.Format("15:04:05") + " " + ll.orderType + " " + result[i]
	}
	return result
}
func (ll lines) String() string {
	if len(ll.ll) == 0 {
		return ""
	}
	data := [][]string{}
	for _, i := range ll.ll {
		diff := strconv.Itoa(i.diff)
		if i.diff > 0 {
			diff = "+" + diff
		}
		data = append(data, []string{
			i.item,
			strconv.Itoa(i.quantity),
			strconv.Itoa(i.place),
			strconv.Itoa(i.price),
			diff,
			strings.Join(intsToStrings(i.best), " "),
		})
	}
	reader, writer := io.Pipe()
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Item", "#", "N°", "$", "±", "best " + ll.orderType})
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

func (ll lines) sort(sorting string) {
	switch sorting {
	case "place":
		slice.Sort(ll.ll, func(i, j int) bool {
			return ll.ll[i].place > ll.ll[j].place
		})
	case "name":
		fallthrough
	default:
		slice.Sort(ll.ll, func(i, j int) bool {
			return ll.ll[i].item < ll.ll[j].item
		})
	}
}

func newLines(username string, all []client.Order, orderType string) lines {
	ll := lines{
		orderType: orderType,
	}
	if len(all) > 0 {
		for _, i := range all {
			if i.OrderType != orderType {
				continue
			}
			orders, err := client.FetchItemOrders(i.Item.URLName)
			if err != nil {
				panic(err)
			}
			all := []int{}
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
			sort.Ints(all)
			if i.OrderType == "buy" {
				for i := len(all)/2 - 1; i >= 0; i-- { // reverse
					opp := len(all) - 1 - i
					all[i], all[opp] = all[opp], all[i]
				}
			}
			diff := 0
			if len(all) > 0 {
				diff = i.Platinum - all[0]
			}
			itemName := i.Item.Info.ItemName
			if i.ModRank != 0 {
				itemName = itemName + " " + strconv.Itoa(i.ModRank)
			}
			five := int(math.Min(float64(len(all)), 5.))
			l := line{
				item:     itemName,
				quantity: i.Quantity,
				place:    position,
				price:    i.Platinum,
				diff:     diff,
				best:     all[:five],
			}
			ll.ll = append(ll.ll, l)
		}
	}
	return ll
}
