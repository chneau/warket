package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/chneau/warket/client"
	"github.com/fatih/color"
	"github.com/gen2brain/beeep"

	"github.com/urfave/cli/v2"
)

func init() {
	log.SetPrefix("[WARKET] ")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		println()
		os.Exit(0)
	}()
}

func main() {
	app := cli.NewApp()
	app.Name = "warket"
	app.Usage = "market tool for the game Warframe"
	app.Version = "0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:    "watch",
			Aliases: []string{"w"},
			Usage:   "watch a market player",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Value:   true,
					Name:    "buy",
					Aliases: []string{"b"},
					Usage:   "show buys",
				},
				&cli.BoolFlag{
					Value:   true,
					Name:    "sell",
					Aliases: []string{"s"},
					Usage:   "show sells`",
				},
				&cli.IntFlag{
					Name:    "log",
					Aliases: []string{"l"},
					Value:   10,
					Usage:   "show `N` last events`",
				},
				&cli.IntFlag{
					Name:    "time",
					Aliases: []string{"t"},
					Value:   -1,
					Usage:   "timer to loop, disabled if -1",
				},
				&cli.StringFlag{
					Name:  "sort",
					Value: "name",
					Usage: "sort by `COLUMN`",
				},
			},
			Action: func(c *cli.Context) error {
				tc := TableContext{
					Buys:        c.Bool("buy"),
					Sells:       c.Bool("sell"),
					LinesOfLogs: c.Int("log"),
					TicksSecond: c.Int("t"),
					Sorting:     c.String("sort"),
				}
				args := c.Args()
				if args.Len() == 0 {
					fmt.Println("Please provide username as parameter")
					os.Exit(1)
				} else {
					tc.Username = args.First()
				}
				tc.Run()
				return nil
			},
		},
		{
			Name:    "snipe",
			Aliases: []string{"s"},
			Usage:   "snipe market",
			Flags: []cli.Flag{
				&cli.Float64Flag{
					Value:   5,
					Name:    "gain",
					Aliases: []string{"g"},
					Usage:   "minimum gain",
				},
				&cli.BoolFlag{
					Value:   true,
					Name:    "notification",
					Aliases: []string{"n"},
					Usage:   "notification",
				},
				&cli.BoolFlag{
					Value:   true,
					Name:    "copy",
					Aliases: []string{"c"},
					Usage:   "copy message",
				},
			},
			Action: func(c *cli.Context) error {
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
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
