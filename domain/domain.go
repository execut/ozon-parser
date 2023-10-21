package domain

import "time"

type Keyword struct {
	Name string
}

type KeywordAnalyticsResult struct {
	Time    time.Time
	Keyword Keyword
	Data    AnalyticsData
}

type AnalyticsItem struct {
	Position               int     `json:"position"`
	IsTraforetto           bool    `json:"isTraforetto"`
	Sku                    string  `json:"sku"`
	Name                   string  `json:"name"`
	ImageUrl               string  `json:"imageUrl"`
	SellerName             string  `json:"sellerName"`
	IsCurSellerItem        bool    `json:"isCurSellerItem"`
	IsInPromo              bool    `json:"isInPromo"`
	Delivery               string  `json:"delivery"`
	DeliveryBoost          int     `json:"deliveryBoost"`
	PopularityScore        int     `json:"popularityScore"`
	SalesScore             int     `json:"salesScore"`
	PriceRub               string  `json:"priceRub"`
	PriceScore             int     `json:"priceScore"`
	Rating                 float64 `json:"rating"`
	RatesCount             int     `json:"ratesCount"`
	RatingScore            int     `json:"ratingScore"`
	QueryFitScore          int     `json:"queryFitScore"`
	PopularityTotalScore   int     `json:"popularityTotalScore"`
	DeliverySpeed          string  `json:"deliverySpeed"`
	FinalResult            float64 `json:"finalResult"`
	SearchPromotionBoost   float64 `json:"searchPromotionBoost"`
	SearchPromotionEnabled bool    `json:"searchPromotionEnabled"`
	PriceDefectRateBoost   float64 `json:"priceDefectRateBoost"`
	IsInTraforettoCampaign bool    `json:"isInTraforettoCampaign"`
	DeliverySpeedBoostSlot *struct {
		FromDays int `json:"fromDays"`
		ToDays   int `json:"toDays"`
	} `json:"deliverySpeedBoostSlot"`
}

type AnalyticsData struct {
	Items     []AnalyticsItem `json:"items"`
	Analytics struct {
		CurSellerItems struct {
			ItemsInTopQnty int           `json:"itemsInTopQnty"`
			ItemsTotalQnty int           `json:"itemsTotalQnty"`
			PagesQueried   int           `json:"pagesQueried"`
			PagesResult    int           `json:"pagesResult"`
			ItemsOutOfTop  []interface{} `json:"itemsOutOfTop"`
		} `json:"curSellerItems"`
		ExpressDelivery struct {
			CurSellerItemsQnty      int `json:"curSellerItemsQnty"`
			CompetitorsItemsQnty    int `json:"competitorsItemsQnty"`
			CompetitorsItemsTopSize int `json:"competitorsItemsTopSize"`
			CurSellerItemsTotalQnty int `json:"curSellerItemsTotalQnty"`
		} `json:"expressDelivery"`
		LocalStore struct {
			CurSellerItemsQnty      int `json:"curSellerItemsQnty"`
			CompetitorsItemsQnty    int `json:"competitorsItemsQnty"`
			CompetitorsItemsTopSize int `json:"competitorsItemsTopSize"`
			CurSellerItemsTotalQnty int `json:"curSellerItemsTotalQnty"`
		} `json:"localStore"`
		Delivery struct {
			DeliveryType            string `json:"deliveryType"`
			CurSellerItemsQnty      int    `json:"curSellerItemsQnty"`
			CompetitorsItemsQnty    int    `json:"competitorsItemsQnty"`
			CompetitorsItemsTopSize int    `json:"competitorsItemsTopSize"`
			CurSellerItemsTotalQnty int    `json:"curSellerItemsTotalQnty"`
		} `json:"delivery"`
		Promo struct {
			CurSellerItemsQnty      int `json:"curSellerItemsQnty"`
			CompetitorsItemsQnty    int `json:"competitorsItemsQnty"`
			CompetitorsItemsTopSize int `json:"competitorsItemsTopSize"`
			CurSellerItemsTotalQnty int `json:"curSellerItemsTotalQnty"`
		} `json:"promo"`
		ItemsRating struct {
			CurSellerItemsAvgRating   float64 `json:"curSellerItemsAvgRating"`
			CompetitorsItemsAvgRating float64 `json:"competitorsItemsAvgRating"`
			CompetitorsItemsTopSize   int     `json:"competitorsItemsTopSize"`
			CurSellerItemsMinRating   int     `json:"curSellerItemsMinRating"`
			CompetitorsItemsMaxRating float64 `json:"competitorsItemsMaxRating"`
			CurSellerMinRatingItem    string  `json:"curSellerMinRatingItem"`
		} `json:"itemsRating"`
		SearchPromotion struct {
			CurSellerItemsQnty            int      `json:"curSellerItemsQnty"`
			CompetitorsItemsQnty          int      `json:"competitorsItemsQnty"`
			MaxCompetitorsSearchPromotion float64  `json:"maxCompetitorsSearchPromotion"`
			CurSellerNotPromotedItems     []string `json:"curSellerNotPromotedItems"`
			CurSellerItemsTotalQnty       int      `json:"curSellerItemsTotalQnty"`
		} `json:"searchPromotion"`
		QueryFit struct {
			CurSellerItemsMinScore   int    `json:"curSellerItemsMinScore"`
			CompetitorsItemsMaxScore int    `json:"competitorsItemsMaxScore"`
			CurSellerMinScoreItem    string `json:"curSellerMinScoreItem"`
		} `json:"queryFit"`
	} `json:"analytics"`
}
