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
	PopularityScore        float64 `json:"popularityScore"`
	SalesScore             float64 `json:"salesScore"`
	PriceRub               string  `json:"priceRub"`
	PriceScore             float64 `json:"priceScore"`
	Rating                 float64 `json:"rating"`
	RatesCount             int     `json:"ratesCount"`
	RatingScore            float64 `json:"ratingScore"`
	QueryFitScore          float64 `json:"queryFitScore"`
	PopularityTotalScore   float64 `json:"popularityTotalScore"`
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
			CurSellerItemsMinRating   float64 `json:"curSellerItemsMinRating"`
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

type ReviewsData struct {
	Products map[string]struct {
		Name       string        `json:"name"`
		CoverImage string        `json:"coverImage"`
		Uri        string        `json:"uri"`
		Variants   []interface{} `json:"variants"`
		ItemId     int           `json:"itemId"`
		Score      int           `json:"Score"`
	} `json:"products"`
	User struct {
		ClientOfficial interface{} `json:"clientOfficial"`
		AvatarUrl      string      `json:"avatarUrl"`
		Guid           string      `json:"guid"`
		FirstName      string      `json:"firstName"`
		LastName       string      `json:"lastName"`
		Fio            string      `json:"fio"`
		Id             int         `json:"id"`
	} `json:"user"`
	Filters struct {
		WithMedia bool `json:"withMedia"`
	} `json:"filters"`
	Author           interface{} `json:"author"`
	CellTrackingInfo struct {
		AddReviewSortingPublishedAtDesc struct {
			ActionType string `json:"actionType"`
			Key        string `json:"key"`
		} `json:"add-review-sorting-published_at_desc"`
		AddReviewSortingScoreAsc struct {
			ActionType string `json:"actionType"`
			Key        string `json:"key"`
		} `json:"add-review-sorting-score_asc"`
		AddReviewSortingScoreDesc struct {
			ActionType string `json:"actionType"`
			Key        string `json:"key"`
		} `json:"add-review-sorting-score_desc"`
		AddReviewSortingUsefullnessDesc struct {
			ActionType string `json:"actionType"`
			Key        string `json:"key"`
		} `json:"add-review-sorting-usefullness_desc"`
		ClickReviewSorting struct {
			ActionType string `json:"actionType"`
			Key        string `json:"key"`
		} `json:"click-review-sorting"`
		Set struct {
			ActionType string `json:"actionType"`
			Key        string `json:"key"`
		} `json:"set"`
		Unset struct {
			ActionType string `json:"actionType"`
			Key        string `json:"key"`
		} `json:"unset"`
		View struct {
			ActionType string `json:"actionType"`
			Key        string `json:"key"`
		} `json:"view"`
	} `json:"cellTrackingInfo"`
	Uri           string `json:"uri"`
	RequestedPath string `json:"requestedPath"`
	Reviews       []struct {
		Sharing *struct {
			Url string `json:"url"`
		} `json:"sharing"`
		EditUrl     string      `json:"editUrl"`
		UpdatedAt   interface{} `json:"updatedAt"`
		PublishedAt int         `json:"publishedAt"`
		Author      struct {
			ClientOfficial interface{} `json:"clientOfficial"`
			AvatarUrl      string      `json:"avatarUrl"`
			Guid           string      `json:"guid"`
			FirstName      string      `json:"firstName"`
			LastName       string      `json:"lastName"`
			Fio            string      `json:"fio"`
			Id             int         `json:"id"`
		} `json:"author"`
		Status struct {
			Name string `json:"name"`
			Id   int    `json:"id"`
		} `json:"status"`
		Usefulness struct {
			UserSelection string `json:"userSelection"`
			Useful        int    `json:"useful"`
			Unuseful      int    `json:"unuseful"`
		} `json:"usefulness,omitempty"`
		ItemId          int    `json:"itemId"`
		CreatedAt       int    `json:"createdAt"`
		Version         int    `json:"version"`
		StatusId        int    `json:"statusId"`
		Uuid            string `json:"uuid"`
		IsDeletable     bool   `json:"isDeletable"`
		IsEdited        bool   `json:"isEdited"`
		IsEditable      bool   `json:"isEditable"`
		IsVotable       bool   `json:"isVotable"`
		IsReportable    bool   `json:"isReportable"`
		IsDeleted       bool   `json:"isDeleted"`
		IsRejected      bool   `json:"isRejected"`
		IsAbuseReported bool   `json:"isAbuseReported"`
		IsAnchor        bool   `json:"isAnchor"`
		Content         struct {
			Comment  string        `json:"comment"`
			Positive string        `json:"positive"`
			Negative string        `json:"negative"`
			Videos   []interface{} `json:"videos"`
			Photos   []struct {
				GalleryParams string `json:"galleryParams"`
				Name          string `json:"name"`
				Url           string `json:"url"`
				Width         int    `json:"width"`
				Height        int    `json:"height"`
				UUID          string `json:"UUID"`
				Published     bool   `json:"published"`
			} `json:"photos"`
			ContextQuestions []interface{} `json:"contextQuestions"`
			Score            int           `json:"score"`
		} `json:"content"`
		Comments struct {
			List       []interface{} `json:"list"`
			TotalCount int           `json:"totalCount"`
		} `json:"comments"`
		IsCommentable    bool `json:"isCommentable"`
		IsAnonymous      bool `json:"isAnonymous"`
		ShowVariantImage bool `json:"showVariantImage"`
		IsUserReview     bool `json:"isUserReview"`
		IsItemPurchased  bool `json:"isItemPurchased"`
	} `json:"reviews"`
	Sortings []struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Active bool   `json:"active"`
	} `json:"sortings"`
	ItemId int `json:"itemId"`
	Paging struct {
		Total       int `json:"total"`
		CommonTotal int `json:"commonTotal"`
		Page        int `json:"page"`
		PerPage     int `json:"perPage"`
	} `json:"paging"`
	PageType     int     `json:"pageType"`
	ProductScore float64 `json:"productScore"`
	Actions      struct {
		AddComment    string `json:"addComment"`
		LoadComment   string `json:"loadComment"`
		RemoveComment string `json:"removeComment"`
		RemoveReview  string `json:"removeReview"`
		ReportComment string `json:"reportComment"`
		ReportReview  string `json:"reportReview"`
		VoteComment   string `json:"voteComment"`
		VoteReview    string `json:"voteReview"`
	} `json:"actions"`
}

type ReviewsResult struct {
	PositiveWords []Word
	NegativeWords []Word
}

type Word struct {
	Key   string
	Value int
}
