package client

import "time"

// DropHow explains how to drop something
type DropHow struct {
	Link string `json:"link"`
	Name string `json:"name"`
}

// ItemInfo human info
type ItemInfo struct {
	Drop        []DropHow `json:"drop"`
	Codex       string    `json:"codex"`
	ItemName    string    `json:"item_name"`
	Description string    `json:"description"`
	WikiLink    string    `json:"wiki_link"`
}

// Item is any warframe item
type Item struct {
	URLName      string    `json:"url_name"`
	ID           string    `json:"id"`
	ItemName     string    `json:"item_name"`
	Thumb        string    `json:"thumb"`
	SubIcon      string    `json:"sub_icon"`
	MasteryLevel string    `json:"mastery_level"`
	Tags         []string  `json:"tags"`
	SetRoot      bool      `json:"set_root"`
	Icon         string    `json:"icon"`
	Ducats       int       `json:"ducats"`
	IconFormat   string    `json:"icon_format"`
	TradingTax   int       `json:"trading_tax"`
	Info         *ItemInfo `json:"en"`
}

// User ...
type User struct {
	IngameName      string     `json:"ingame_name"`
	LastSeen        *time.Time `json:"last_seen"`
	ReputationBonus int        `json:"reputation_bonus"`
	Reputation      int        `json:"reputation"`
	Region          string     `json:"region"`
	Status          string     `json:"status"`
	ID              string     `json:"id"`
	Avatar          string     `json:"avatar"`
}

// Order ...
type Order struct {
	Visible      bool       `json:"visible"`
	CreationDate *time.Time `json:"creation_date"`
	Quantity     int        `json:"quantity"`
	User         *User      `json:"user"`
	LastUpdate   *time.Time `json:"last_update"`
	ClosedDate   *time.Time `json:"closed_date"`
	Platinum     float64    `json:"platinum"`
	OrderType    string     `json:"order_type"`
	Region       string     `json:"region"`
	Platform     string     `json:"platform"`
	ID           string     `json:"id"`
	ModRank      int        `json:"mod_rank"`
	Item         *Item      `json:"item"`
}

// Stat ...
type Stat struct {
	Datetime    *time.Time `json:"datetime"`
	Volume      int        `json:"volume"`
	MinPrice    int        `json:"min_price"`
	MaxPrice    int        `json:"max_price"`
	OpenPrice   int        `json:"open_price"`
	ClosedPrice int        `json:"closed_price"`
	AvgPrice    float64    `json:"avg_price"`
	WaPrice     float64    `json:"wa_price"`
	Median      float64    `json:"median"`
	MovingAvg   float64    `json:"moving_avg"`
	DonchTop    int        `json:"donch_top"`
	DonchBot    int        `json:"donch_bot"`
	ID          string     `json:"id"`
}

// Profile ...
type Profile struct {
	OwnProfile   bool       `json:"own_profile"`
	Status       string     `json:"status"`
	Background   string     `json:"background"`
	Avatar       string     `json:"avatar"`
	IngameName   string     `json:"ingame_name"`
	About        string     `json:"about"`
	Reputation   int        `json:"reputation"`
	LastSeen     *time.Time `json:"last_seen"`
	ID           string     `json:"id"`
	Region       string     `json:"region"`
	Achievements []string   `json:"achievements"` // Not yet sure if it is string
	Platform     string     `json:"platform"`
}

// Review ...
type Review struct {
	Date       time.Time `json:"date"`
	UserFrom   *User     `json:"user_from"`
	Text       string    `json:"text"`
	ID         string    `json:"id"`
	ReviewType int       `json:"review_type"`
	Hidden     bool      `json:"hidden"`
}
