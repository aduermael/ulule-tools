package clientapi

import ()

// ListProjectResponse represents a response from
// Ulule's API to a GET */projects request.
type ListProjectResponse struct {
	Meta     *Metadata
	Projects []*Project `json:"projects"`
}

// Project represents a Ulule project
type Project struct {
	Id              int    `json:"id"`
	Url             string `json:"absolute_url"`
	Goal            int    `json:"goal"`
	GoalRaised      bool   `json:"goal_raised"`
	AmountRaised    int    `json:"amount_raised"`
	CommentCount    int    `json:"comments_count"`
	Committed       int    `json:"committed"`
	Currency        string `json:"currency"`
	CurrencyDisplay string `json:"currency_display"`
	DateEnd         string `json:"date_end"`
	DateStart       string `json:"date_start"`
	Finished        bool   `json:"finished"`
	Slug            string `json:"slug"`
	SupportersCount int    `json:"supporters_count"`
	TimeZone        string `json:"timezone"`
}

// Metadata
type Metadata struct {
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
	TotalCount int `json:"total_count"`
	// Next ?
	// Previous ?
}
