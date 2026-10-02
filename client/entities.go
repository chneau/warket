package client

import "time"

// ItemInfo is the localized item data for a single language (i18n map entry).
type ItemInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	WikiLink    string `json:"wikiLink"`
	Icon        string `json:"icon"`
	Thumb       string `json:"thumb"`
	SubIcon     string `json:"subIcon"`
}

// Item is the Warframe.market v2 item model.
type Item struct {
	ID             string              `json:"id"`
	Slug           string              `json:"slug"`
	GameRef        string              `json:"gameRef"`
	Tags           []string            `json:"tags"`
	SetRoot        bool                `json:"setRoot"`
	SetParts       []string            `json:"setParts"`
	QuantityInSet  int                 `json:"quantityInSet"`
	Rarity         string              `json:"rarity"`
	BulkTradable   bool                `json:"bulkTradable"`
	Subtypes       []string            `json:"subtypes"`
	MaxRank        int                 `json:"maxRank"`
	MaxCharges     int                 `json:"maxCharges"`
	Ducats         int                 `json:"ducats"`
	Vosfor         int                 `json:"vosfor"`
	ReqMasteryRank int                 `json:"reqMasteryRank"`
	Vaulted        bool                `json:"vaulted"`
	TradingTax     int                 `json:"tradingTax"`
	Tradable       bool                `json:"tradable"`
	I18N           map[string]ItemInfo `json:"i18n"`
}

// URLName is the URL-safe slug used by the v2 API.
func (i *Item) URLName() string {
	if i == nil {
		return ""
	}
	return i.Slug
}

// Name is the English display name.
func (i *Item) Name() string {
	if i == nil {
		return ""
	}
	return i.I18N["en"].Name
}

// Activity is the rich activity attached to a user.
type Activity struct {
	Type      string `json:"type"`
	Details   string `json:"details"`
	StartedAt string `json:"startedAt"`
}

// User is both the short (embedded in orders) and full public user model.
type User struct {
	ID          string     `json:"id"`
	IngameName  string     `json:"ingameName"`
	Slug        string     `json:"slug"`
	Reputation  int        `json:"reputation"`
	Platform    string     `json:"platform"`
	Crossplay   bool       `json:"crossplay"`
	Locale      string     `json:"locale"`
	Status      string     `json:"status"`
	Avatar      string     `json:"avatar"`
	Activity    *Activity  `json:"activity"`
	LastSeen    *time.Time `json:"lastSeen"`
	Role        string     `json:"role"`
	Tier        string     `json:"tier"`
	MasteryRank int        `json:"masteryRank"`
	Banned      bool       `json:"banned"`
	About       string     `json:"about"`
}

// Order is the Warframe.market v2 order model. Item is resolved locally from
// ItemID because v2 order responses do not embed the item.
type Order struct {
	ID         string  `json:"id"`
	OrderType  string  `json:"type"`
	Platinum   float64 `json:"platinum"`
	Quantity   int     `json:"quantity"`
	PerTrade   int     `json:"perTrade"`
	Subtype    string  `json:"subtype"`
	ModRank    int     `json:"rank"`
	Charges    int     `json:"charges"`
	AmberStars int     `json:"amberStars"`
	CyanStars  int     `json:"cyanStars"`
	Visible    bool    `json:"visible"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
	ItemID     string  `json:"itemId"`
	User       *User   `json:"user"`
	Item       *Item   `json:"-"`
}
