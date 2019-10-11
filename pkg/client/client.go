package client

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	// Root url.
	Root = "https://api.warframe.market/v1"
)

// H ...
type H map[string]interface{}

// FetchItems Get all item names
func FetchItems() ([]Item, error) {
	res, err := http.Get(Root + "/items")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data := struct {
		Payload struct {
			Items struct {
				En []Item
			}
		}
		Error string
	}{}
	json.NewDecoder(res.Body).Decode(&data)
	if data.Error != "" {
		return nil, errors.New(data.Error)
	}
	return data.Payload.Items.En, nil
}

// FetchItemInfo Get item information.
func FetchItemInfo(urlName string) ([]Item, error) {
	res, err := http.Get(Root + "/items/" + urlName)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data := struct {
		Payload struct {
			Item struct {
				ItemInSet []Item `json:"items_in_set"`
			}
		}
		Error string
	}{}
	json.NewDecoder(res.Body).Decode(&data)
	if data.Error != "" {
		return nil, errors.New(data.Error)
	}
	return data.Payload.Item.ItemInSet, nil
}

// FetchItemOrders Get orders for item.
func FetchItemOrders(urlName string) ([]Order, error) {
	res, err := http.Get(Root + "/items/" + urlName + "/orders")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data := struct {
		Payload struct {
			Orders []Order
		}
		Error string
	}{}
	json.NewDecoder(res.Body).Decode(&data)
	if data.Error != "" {
		return nil, errors.New(data.Error)
	}
	return data.Payload.Orders, nil
}

// FetchItemStats Get orders for item.
// In order:
// Closed H48
// Closed D90
// Live H48
// Live D90
func FetchItemStats(urlName string) (closedHours48 []Stat, closedDays90 []Stat, LiveHours48 []Stat, LiveDays90 []Stat, err error) {
	res, err := http.Get(Root + "/items/" + urlName + "/statistics")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer res.Body.Close()
	data := struct {
		Payload struct {
			StatisticsClosed struct {
				Hours48 []Stat `json:"48hours"`
				Days90  []Stat `json:"90days"`
			} `json:"statistics_closed"`
			StatisticsLive struct {
				Hours48 []Stat `json:"48hours"`
				Days90  []Stat `json:"90days"`
			} `json:"statistics_live"`
		}
		Error string
	}{}
	json.NewDecoder(res.Body).Decode(&data)
	if data.Error != "" {
		return nil, nil, nil, nil, errors.New(data.Error)
	}
	return data.Payload.StatisticsClosed.Hours48, data.Payload.StatisticsClosed.Days90,
		data.Payload.StatisticsLive.Hours48, data.Payload.StatisticsLive.Days90, nil
}

// FetchUser Get user profile.
func FetchUser(userName string) (*Profile, error) {
	res, err := http.Get(Root + "/profile/" + userName)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data := struct {
		Payload struct {
			Profile *Profile
		}
		Error string
	}{}
	json.NewDecoder(res.Body).Decode(&data)
	if data.Error != "" {
		return nil, errors.New(data.Error)
	}
	return data.Payload.Profile, nil
}

// FetchUserOrders Get user orders.
// Buy, then, Sell
func FetchUserOrders(userName string) (buy []Order, sell []Order, err error) {
	res, err := http.Get(Root + "/profile/" + userName + "/orders")
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	data := struct {
		Payload struct {
			BuyOrders  []Order `json:"buy_orders"`
			SellOrders []Order `json:"sell_orders"`
		}
		Error string
	}{}
	json.NewDecoder(res.Body).Decode(&data)
	if data.Error != "" {
		return nil, nil, errors.New(data.Error)
	}
	return data.Payload.BuyOrders, data.Payload.SellOrders, nil
}

// FetchUserStats Get user statistics.
func FetchUserStats(userName string) ([]Order, error) {
	res, err := http.Get(Root + "/profile/" + userName + "/statistics")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data := struct {
		Payload struct {
			ClosedOrders []Order `json:"closed_orders"`
		}
		Error string
	}{}
	json.NewDecoder(res.Body).Decode(&data)
	if data.Error != "" {
		return nil, errors.New(data.Error)
	}
	return data.Payload.ClosedOrders, nil
}

// FetchUserReview Get user reviews.
func FetchUserReview(userName string) ([]Review, error) {
	res, err := http.Get(Root + "/profile/" + userName + "/reviews")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data := struct {
		Payload struct {
			Reviews []Review `json:"reviews"`
		}
		Error string
	}{}
	json.NewDecoder(res.Body).Decode(&data)
	if data.Error != "" {
		return nil, errors.New(data.Error)
	}
	return data.Payload.Reviews, nil
}

func SubWS(ch chan<- *Order) error {
	h := http.Header{}
	h.Set("Pragma", "no-cache")
	h.Set("Origin", "https://warframe.market")
	h.Set("Accept-Encoding", "gzip, deflate, br")
	h.Set("Accept-Language", "en-GB,en;q=0.9,fr-FR;q=0.8,fr;q=0.7,en-US;q=0.6")
	h.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36")
	h.Set("Cache-Control", "no-cache")

	ws, _, err := websocket.DefaultDialer.Dial("wss://warframe.market/socket?platform=pc", h)
	if err != nil {
		return err
	}
	err = ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"@WS/SUBSCRIBE/MOST_RECENT"}`))
	for {
		result := &struct {
			Type string
		}{}
		_, b, err := ws.ReadMessage()
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, result)
		if err != nil {
			return err
		}
		if result.Type != `@WS/SUBSCRIPTIONS/MOST_RECENT/NEW_ORDER` {
			continue
		}
		obj := &struct {
			Payload struct {
				Order Order
			}
		}{}
		err = json.Unmarshal(b, obj)
		if err != nil {
			return err
		}
		ch <- &obj.Payload.Order
	}
	return nil
}
