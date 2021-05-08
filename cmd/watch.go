package cmd

import (
	"fmt"
	"os"

	"github.com/chneau/warket/display"
	"github.com/urfave/cli/v2"
)

func Watch(c *cli.Context) error {
	tc := display.TableContext{
		Buys:        c.Bool("buy"),
		Sells:       c.Bool("sell"),
		LinesOfLogs: c.Int("log"),
		TicksSecond: c.Int("ticker"),
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
}
