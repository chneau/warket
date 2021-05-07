package client

import "time"

type DropHow struct {
	Link string `json:"link"`
	Name string `json:"name"`
}

type ItemInfo struct {
	Drop        []DropHow `json:"drop"`
	Codex       string    `json:"codex"`
	ItemName    string    `json:"item_name"`
	Description string    `json:"description"`
	WikiLink    string    `json:"wiki_link"`
}

type Item struct {
	URLName      string    `json:"url_name"`
	ID           string    `json:"id"`
	ItemName     string    `json:"item_name"`
	Thumb        string    `json:"thumb"`
	SubIcon      string    `json:"sub_icon"`
	MasteryLevel int       `json:"mastery_level"`
	Tags         []string  `json:"tags"`
	SetRoot      bool      `json:"set_root"`
	Icon         string    `json:"icon"`
	Ducats       int       `json:"ducats"`
	IconFormat   string    `json:"icon_format"`
	TradingTax   int       `json:"trading_tax"`
	Info         *ItemInfo `json:"en"`
}

type User struct {
	IngameName      string         `json:"ingame_name"`
	LastSeen        *time.Time     `json:"last_seen"`
	ReputationBonus int            `json:"reputation_bonus"`
	Reputation      int            `json:"reputation"`
	Region          string         `json:"region"`
	Avatar          *string        `json:"avatar"`
	Status          string         `json:"status"`
	ID              string         `json:"id"`
	Anonymous       bool           `json:"anonymous"`
	Role            string         `json:"role"`
	Banned          bool           `json:"banned"`
	HasMail         bool           `json:"has_mail"`
	Verification    bool           `json:"verification"`
	LinkedAccounts  LinkedAccounts `json:"linked_accounts"`
	WrittenReviews  int            `json:"written_reviews"`
	UnreadMessages  int            `json:"unread_messages"`
	CheckCode       string         `json:"check_code"`
	Platform        string         `json:"platform"`
}

type LinkedAccounts struct {
	SteamProfile   bool `json:"steam_profile"`
	PatreonProfile bool `json:"patreon_profile"`
	XboxProfile    bool `json:"xbox_profile"`
}

type Order struct {
	Visible      bool       `json:"visible"`
	CreationDate string     `json:"creation_date"`
	Quantity     int        `json:"quantity"`
	User         User       `json:"user"`
	LastUpdate   string     `json:"last_update"`
	Platinum     float64    `json:"platinum"`
	OrderType    string     `json:"order_type"`
	Region       string     `json:"region"`
	Platform     string     `json:"platform"`
	ID           string     `json:"id"`
	ClosedDate   *time.Time `json:"closed_date"`
	ModRank      int        `json:"mod_rank"`
	Item         *Item      `json:"item"`
}

type Stat struct {
	Datetime    *time.Time `json:"datetime"`
	Volume      float64    `json:"volume"`
	MinPrice    float64    `json:"min_price"`
	MaxPrice    float64    `json:"max_price"`
	OpenPrice   float64    `json:"open_price"`
	ClosedPrice float64    `json:"closed_price"`
	AvgPrice    float64    `json:"avg_price"`
	WaPrice     float64    `json:"wa_price"`
	Median      float64    `json:"median"`
	MovingAvg   float64    `json:"moving_avg"`
	DonchTop    float64    `json:"donch_top"`
	DonchBot    float64    `json:"donch_bot"`
	ID          string     `json:"id"`
}

type Profile struct {
	OwnProfile bool       `json:"own_profile"`
	Status     string     `json:"status"`
	Background string     `json:"background"`
	Avatar     string     `json:"avatar"`
	IngameName string     `json:"ingame_name"`
	About      string     `json:"about"`
	Reputation int        `json:"reputation"`
	LastSeen   *time.Time `json:"last_seen"`
	ID         string     `json:"id"`
	Region     string     `json:"region"`
	Platform   string     `json:"platform"`
}

type Review struct {
	Text       string   `json:"text"`
	ReviewType int      `json:"review_type"`
	ID         string   `json:"id"`
	UserFrom   UserFrom `json:"user_from"`
	Hidden     bool     `json:"hidden"`
	Date       string   `json:"date"`
}

type UserFrom struct {
	Reputation int     `json:"reputation"`
	Region     string  `json:"region"`
	ID         string  `json:"id"`
	Avatar     *string `json:"avatar"`
	IngameName string  `json:"ingame_name"`
	Status     string  `json:"status"`
}
