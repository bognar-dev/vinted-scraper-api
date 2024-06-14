package vinted_scraper

// User represents the user object inside each Item
type User struct {
	ID         int         `json:"id"`
	Login      string      `json:"login"`
	Business   bool        `json:"business"`
	ProfileURL string      `json:"profile_url"`
	Photo      interface{} `json:"photo"` // Assuming photo can be null
}
type Thumbnail struct {
	Type         string      `json:"type"`
	URL          string      `json:"url"`
	Width        int         `json:"width"`
	Height       int         `json:"height"`
	OriginalSize interface{} `json:"original_size"` // Assuming original_size can be null
}

// Photo represents the photo object inside each Item
type Photo struct {
	ID                  int         `json:"id"`
	ImageNo             int         `json:"image_no"`
	Width               int         `json:"width"`
	Height              int         `json:"height"`
	DominantColor       string      `json:"dominant_color"`
	DominantColorOpaque string      `json:"dominant_color_opaque"`
	URL                 string      `json:"url"`
	IsMain              bool        `json:"is_main"`
	Thumbnails          []Thumbnail `json:"thumbnails"`
	HighResolution      struct {
		ID          string      `json:"id"`
		Timestamp   int         `json:"timestamp"`
		Orientation interface{} `json:"orientation"` // Assuming orientation can be null
	} `json:"high_resolution"`
	IsSuspicious bool     `json:"is_suspicious"`
	FullSizeURL  string   `json:"full_size_url"`
	IsHidden     bool     `json:"is_hidden"`
	Extra        struct{} `json:"extra"`
}

// Item represents each item in the JSON response
type Item struct {
	ID                    int           `json:"id"`
	Title                 string        `json:"title"`
	Price                 string        `json:"price"`
	IsVisible             int           `json:"is_visible"`
	Discount              interface{}   `json:"discount"` // Assuming discount can be null
	Currency              string        `json:"currency"`
	BrandTitle            string        `json:"brand_title"`
	User                  User          `json:"user"`
	URL                   string        `json:"url"`
	Promoted              bool          `json:"promoted"`
	Photo                 Photo         `json:"photo"`
	FavouriteCount        int           `json:"favourite_count"`
	IsFavourite           bool          `json:"is_favourite"`
	Badge                 interface{}   `json:"badge"`      // Assuming badge can be null
	Conversion            interface{}   `json:"conversion"` // Assuming conversion can be null
	ServiceFee            string        `json:"service_fee"`
	TotalItemPrice        string        `json:"total_item_price"`
	TotalItemPriceRounded interface{}   `json:"total_item_price_rounded"` // Assuming total_item_price_rounded can be null
	ViewCount             int           `json:"view_count"`
	SizeTitle             string        `json:"size_title"`
	ContentSource         string        `json:"content_source"`
	Status                string        `json:"status"`
	IconBadges            []interface{} `json:"icon_badges"` // Assuming icon_badges can be empty or null
	SearchTrackingParams  struct {
		Score          float64  `json:"score"`
		MatchedQueries []string `json:"matched_queries"`
	} `json:"search_tracking_params"`
}

// Response represents the entire JSON response structure
type VintedApi_Response struct {
	Items                []Item      `json:"items"`
	DominantBrand        interface{} `json:"dominant_brand"` // Assuming dominant_brand can be null
	SearchTrackingParams struct {
		SearchCorrelationID   string `json:"search_correlation_id"`
		SearchSessionID       string `json:"search_session_id"`
		GlobalSearchSessionID string `json:"global_search_session_id"`
	} `json:"search_tracking_params"`
	Pagination struct {
		CurrentPage  int `json:"current_page"`
		TotalPages   int `json:"total_pages"`
		TotalEntries int `json:"total_entries"`
		PerPage      int `json:"per_page"`
		Time         int `json:"time"`
	} `json:"pagination"`
	Code int `json:"code"`
}
