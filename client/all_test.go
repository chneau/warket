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

	slug := "ash_prime_set"

	itemInfo, err := FetchItemInfo(slug)
	if err != nil {
		panic(err)
	}
	log.Println("Size of set:", len(itemInfo))

	orders, err := FetchItemOrders(slug)
	if err != nil {
		panic(err)
	}
	log.Println("Number of orders:", len(orders))
	if len(orders) == 0 {
		t.Fatal("no orders returned")
	}
	if orders[0].Item == nil || orders[0].Item.Name() == "" {
		t.Fatal("order item was not resolved")
	}

	username := orders[0].User.Slug
	profile, err := FetchUser(username)
	if err != nil {
		panic(err)
	}
	log.Println("Profile:", profile.IngameName)

	buy, sell, err := FetchUserOrders(username)
	if err != nil {
		panic(err)
	}
	log.Println("Buy for user:", len(buy))
	log.Println("Sell for user:", len(sell))
	for _, o := range append(buy, sell...) {
		if o.Item == nil {
			t.Fatalf("order %s item %s was not resolved", o.ID, o.ItemID)
		}
	}
}
