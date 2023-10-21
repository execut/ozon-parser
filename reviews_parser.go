package main

import (
	"encoding/json"
	"fmt"
	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/chromedp"
	"github.com/redis/go-redis/v9"
	"log"
	"math"
	"time"
)

func ParseReviews(productId string) []ReviewsData {
	// New creates a new context for use with chromedp. With this context
	// you can use chromedp as you normally would.
	config := cu.NewConfig(
		// Remove this if you want to see a browser window.
		cu.WithHeadless(),

		// If the webelement is not found within 10 seconds, timeout.
		cu.WithTimeout(10*time.Second),
	)
	ctx, cancel, err := cu.New(config)
	if err != nil {
		panic(err)
	}
	defer cancel()

	var ok bool
	lastPage := 999
	var result []ReviewsData
	for page := 1; page <= lastPage; page++ {
		url := fmt.Sprintf("https://www.ozon.ru/product/%s/reviews/?page=%v", productId, page)
		dataState, err := GetCachedValue(url)
		if err == redis.Nil {
			err = chromedp.Run(ctx,
				chromedp.Navigate(url),
				chromedp.WaitVisible(`//*[@id="state-webListReviews-3231710-default-1"]`),
				chromedp.AttributeValue(`#state-webListReviews-3231710-default-1`, `data-state`, &dataState, &ok),
			)
			if err != nil {
				log.Fatal(err)
			}

			SetCachedValue(url, dataState)
		}

		var reviewsData ReviewsData
		err = json.Unmarshal([]byte(dataState), &reviewsData)
		if err != nil {
			panic(err)
		}

		if lastPage == 999 {
			lastPage = int(math.Ceil(float64(reviewsData.Paging.Total) / float64(reviewsData.Paging.PerPage)))
		}

		if lastPage > 10 {
			lastPage = 10
		}

		result = append(result, reviewsData)
	}

	return result
}

func ExtractTextFromReviews(reviewsPages []ReviewsData) string {
	var text string
	for _, reviewsPage := range reviewsPages {
		for _, reviewsList := range reviewsPage.Reviews {
			content := reviewsList.Content
			if content.Positive == "" && content.Negative == "" && content.Comment == "" {
				continue
			}

			if content.Comment != "" {
				text += " " + content.Comment
			}

			if content.Positive != "" {
				text += " " + content.Positive
			}

			if content.Negative != "" {
				text += " " + content.Negative
			}
		}
	}

	return text
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
