package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"

	"github.com/chneau/warket/pkg/client"

	"github.com/urfave/cli"
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
	app.Commands = []cli.Command{
		{
			Name:    "watch",
			Aliases: []string{"w"},
			Usage:   "watch a market player",
			Flags: []cli.Flag{
				cli.BoolTFlag{
					Name:  "buy, b",
					Usage: "show buys",
				},
				cli.BoolTFlag{
					Name:  "sell, s",
					Usage: "show sells`",
				},
				cli.IntFlag{
					Name:  "log, l",
					Value: 10,
					Usage: "show `N` last events`",
				},
				cli.IntFlag{
					Name:  "time, t",
					Value: -1,
					Usage: "timer to loop, disabled if -1",
				},
				cli.StringFlag{
					Name:  "sort",
					Value: "name",
					Usage: "sort by `COLUMN`",
				},
			},
			Action: func(c *cli.Context) error {
				d := data{
					buys:    c.BoolT("buy"),
					sells:   c.BoolT("sell"),
					logging: c.Int("log"),
					t:       c.Int("t"),
					sorting: c.String("sort"),
				}
				args := c.Args()
				if len(args) == 0 {
					fmt.Println("Please provide username as parameter")
					os.Exit(1)
				} else {
					d.username = args[0]
				}
				d.run()
				return nil
			},
		},
		{
			Name:    "snipe",
			Aliases: []string{"s"},
			Usage:   "snipe market",
			Action: func(c *cli.Context) error {
				log.Println("hello world !")
				ch := make(chan *client.Order)
				go func() {
					err := client.SubWS(ch)
					if err != nil {
						log.Fatalln(err)
					}
				}()
				for order := range ch {
					if order.OrderType != "sell" {
						continue
					}
					orders, err := client.FetchItemOrders(order.Item.URLName)
					if err != nil {
						log.Fatalln(err)
					}
					// TODO remove copy pasta and create a funciton
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
					// TODO nice formating
					// TODO potential gain for buy/sell
					log.Println(order.OrderType, order.Platinum, order.Item.URLName, position, all)
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
