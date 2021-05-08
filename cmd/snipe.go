package cmd

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/chneau/warket/client"
	"github.com/fatih/color"
	"github.com/gen2brain/beeep"
	"github.com/urfave/cli/v2"
)

func Snipe(c *cli.Context) error {
	minGain := c.Float64("gain")
	notficationEnabled := c.Bool("notification")
	copyEnabled := c.Bool("copy")
	log.Println("Sniping with a minimum gain of", minGain, "plat.")
	ordersChan := make(chan *client.Order)
	go func() {
		err := client.SubWS(ordersChan)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	for order := range ordersChan {
		if order.OrderType != "sell" {
			continue
		}
		orders, err := client.FetchItemOrders(order.Item.URLName)
		if err != nil {
			log.Fatalln(err)
		}
		all := []float64{}
		position := 1
		for _, o := range orders {
			if o.User.Status == "ingame" && o.OrderType == order.OrderType && o.ModRank == order.ModRank && o.Region == order.Region {
				all = append(all, o.Platinum)
				if order.OrderType == "buy" {
					if o.Platinum > order.Platinum {
						position++
					}
				} else {
					if order.Platinum > o.Platinum {
						position++
					}
				}
			}
		}
		sort.Float64s(all)
		if len(all) > 10 {
			all = all[:10]
		}
		if len(all) < 2 {
			continue
		}
		if all[0] > order.Platinum {
			continue
		}
		gain := all[1] - order.Platinum
		if gain < minGain {
			continue
		}

		fmt.Println(
			color.MagentaString(order.Item.Info.ItemName),
			color.GreenString(strconv.Itoa(position)),
			color.HiBlueString(fmt.Sprint(order.Platinum)),
			color.CyanString(fmt.Sprint(all)),
			color.HiRedString(fmt.Sprint("potential gain: ", gain)),
		)

		whisper := fmt.Sprint("/w ", order.User.IngameName, " Hi! I want to buy: ", order.Item.Info.ItemName, " for ", order.Platinum, " platinum. (warframe.market)")
		fmt.Println(whisper)

		if notficationEnabled {
			_ = beeep.Notify(order.Item.Info.ItemName, fmt.Sprint("At ", order.Platinum, "p for a gain of ", gain), "assets/information.png")
		}

		if copyEnabled {
			_ = clipboard.WriteAll(whisper)
		}
	}
	return nil
}
