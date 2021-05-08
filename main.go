package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/chneau/warket/cmd"
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
	app.Usage = "market tool for Warframe"
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
					Usage:   "show sells",
				},
				&cli.IntFlag{
					Name:    "log",
					Aliases: []string{"l"},
					Value:   10,
					Usage:   "shows last logs",
				},
				&cli.IntFlag{
					Name:    "ticker",
					Aliases: []string{"t"},
					Value:   -1,
					Usage:   "ticker in seconds",
				},
				&cli.StringFlag{
					Name:  "sort",
					Value: "name",
					Usage: "sort by column",
				},
			},
			Action: cmd.Watch,
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
			Action: cmd.Snipe,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
