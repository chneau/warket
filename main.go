package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
)

func init() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		println()
		os.Exit(0)
	}()
}

func main() {
	d := data{}
	flag.BoolVar(&d.buys, "buy", true, "show buys")
	flag.BoolVar(&d.sells, "sell", true, "show sells")
	flag.IntVar(&d.logging, "log", 10, "show 10 last events")
	flag.IntVar(&d.t, "t", -1, "timer to loop")
	flag.StringVar(&d.sorting, "sort", "name", "sort by name or place")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Please provide username as parameter")
		os.Exit(1)
	} else {
		d.username = args[0]
	}
	d.run()
}
