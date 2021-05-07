package client

import (
	"log"
	"testing"
)

func Test_All(t *testing.T) {
	items, err := FetchItems()
	if err != nil {
		panic(err)
	}
	log.Println("items:", len(items))
	log.Println("First item", items[0].URLName)
	itemInfo, err := FetchItemInfo(items[0].URLName)
	if err != nil {
		panic(err)
	}
	log.Println("Size of set:", len(itemInfo))
	orders, err := FetchItemOrders(items[0].URLName)
	if err != nil {
		panic(err)
	}
	log.Println("Number of orders:", len(orders))
	_, _, lh48, _, err := FetchItemStats(items[0].URLName)
	if err != nil {
		panic(err)
	}
	log.Println("Stats?", len(lh48))
	profile, err := FetchUser(orders[0].User.IngameName)
	if err != nil {
		panic(err)
	}
	log.Println("Profile:", profile)
	buy, sell, err := FetchUserOrders(orders[0].User.IngameName)
	if err != nil {
		panic(err)
	}
	log.Println("Buy for user:", len(buy))
	log.Println("Sell for user:", len(sell))
	stats, err := FetchUserStats(orders[0].User.IngameName)
	if err != nil {
		panic(err)
	}
	log.Println("Stats for user:", len(stats))
	log.Println(stats)
	reviews, err := FetchUserReview(orders[0].User.IngameName)
	if err != nil {
		panic(err)
	}
	log.Println("User reviews:", len(reviews))
}
