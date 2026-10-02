package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Root is the warframe.market v2 API root.
	Root = "https://api.warframe.market/v2"
	// WSRoot is the realtime WebSocket endpoint.
	WSRoot = "wss://ws.warframe.market/socket"
	// Platform is the marketplace platform used by this client.
	Platform = "pc"
)

// H ...
type H map[string]interface{}

// envelope is the common v2 response envelope:
// {"apiVersion": "...", "data": ..., "error": ...}
type envelope struct {
	APIVersion string          `json:"apiVersion"`
	Data       json.RawMessage `json:"data"`
	Error      json.RawMessage `json:"error"`
}

// getJSON performs a GET request against the v2 API, unwraps the response
// envelope and decodes data into out.
func getJSON(path string, out interface{}) error {
	req, err := http.NewRequest(http.MethodGet, Root+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Platform", Platform)
	req.Header.Set("Language", "en")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "warket (github.com/chneau/warket)")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return fmt.Errorf("GET %s: %s: %s", path, res.Status, string(b))
	}
	env := envelope{}
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		return err
	}
	if len(env.Error) > 0 && string(env.Error) != "null" {
		return fmt.Errorf("GET %s: %s", path, string(env.Error))
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(env.Data, out)
}

var (
	itemsOnce   sync.Once
	itemsErr    error
	itemsByID   map[string]*Item
	itemsBySlug map[string]*Item
)

// loadItems fetches and caches the full item list once. v2 order responses
// only carry an itemId, so the cache is used to re-attach item details.
func loadItems() error {
	itemsOnce.Do(func() {
		var items []*Item
		if err := getJSON("/items", &items); err != nil {
			itemsErr = err
			return
		}
		itemsByID = make(map[string]*Item, len(items))
		itemsBySlug = make(map[string]*Item, len(items))
		for _, it := range items {
			itemsByID[it.ID] = it
			itemsBySlug[it.Slug] = it
		}
	})
	return itemsErr
}

// resolveItem attaches the cached item to an order when possible.
func resolveItem(o *Order) {
	if o.Item == nil && o.ItemID != "" {
		o.Item = itemsByID[o.ItemID]
	}
}

// FetchItems Get all item names.
func FetchItems() ([]Item, error) {
	if err := loadItems(); err != nil {
		return nil, err
	}
	items := make([]Item, 0, len(itemsByID))
	for _, it := range itemsByID {
		items = append(items, *it)
	}
	return items, nil
}

// FetchItemInfo Get the item and the other items in its set.
func FetchItemInfo(urlName string) ([]Item, error) {
	var data struct {
		ID    string  `json:"id"`
		Items []*Item `json:"items"`
	}
	if err := getJSON("/item/"+url.PathEscape(urlName)+"/set", &data); err != nil {
		return nil, err
	}
	items := make([]Item, 0, len(data.Items))
	for _, it := range data.Items {
		items = append(items, *it)
	}
	return items, nil
}

// FetchItemOrders Get visible orders for an item.
func FetchItemOrders(urlName string) ([]Order, error) {
	if err := loadItems(); err != nil {
		return nil, err
	}
	var orders []Order
	if err := getJSON("/orders/item/"+url.PathEscape(urlName), &orders); err != nil {
		return nil, err
	}
	for i := range orders {
		resolveItem(&orders[i])
	}
	return orders, nil
}

// FetchUser Get a user's public profile.
func FetchUser(userName string) (*User, error) {
	user := &User{}
	if err := getJSON("/user/"+url.PathEscape(userName), user); err != nil {
		return nil, err
	}
	return user, nil
}

// FetchUserOrders Get a user's orders, split into buy and sell.
func FetchUserOrders(userName string) (buy []Order, sell []Order, err error) {
	if err := loadItems(); err != nil {
		return nil, nil, err
	}
	var orders []Order
	if err := getJSON("/orders/user/"+url.PathEscape(userName), &orders); err != nil {
		return nil, nil, err
	}
	for i := range orders {
		resolveItem(&orders[i])
		switch orders[i].OrderType {
		case "buy":
			buy = append(buy, orders[i])
		case "sell":
			sell = append(sell, orders[i])
		}
	}
	return buy, sell, nil
}

// SubWS subscribes to the realtime feed of newly posted orders and forwards
// each one on ch until the connection drops.
func SubWS(ch chan<- *Order) error {
	if err := loadItems(); err != nil {
		return err
	}
	dialer := websocket.Dialer{
		Subprotocols:     []string{"wfm"},
		HandshakeTimeout: 15 * time.Second,
	}
	ws, _, err := dialer.Dial(WSRoot, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	sub, err := json.Marshal(H{
		"route": "@wfm|cmd/subscribe/newOrders",
		"id":    "1",
		"payload": H{
			"platform": Platform,
		},
	})
	if err != nil {
		return err
	}
	if err := ws.WriteMessage(websocket.TextMessage, sub); err != nil {
		return err
	}

	for {
		_, b, err := ws.ReadMessage()
		if err != nil {
			return err
		}
		msg := struct {
			Route   string          `json:"route"`
			Payload json.RawMessage `json:"payload"`
		}{}
		if err := json.Unmarshal(b, &msg); err != nil {
			return err
		}
		if msg.Route != "@wfm|event/subscriptions/newOrder" {
			continue
		}
		order := &Order{}
		if err := json.Unmarshal(msg.Payload, order); err != nil {
			return err
		}
		resolveItem(order)
		ch <- order
	}
}
