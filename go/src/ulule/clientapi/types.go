package clientapi

import ()

// ListProjectResponse represents a response from
// Ulule's API to a GET */projects request.
type ListProjectResponse struct {
	Meta     *Metadata  `json:"meta"`
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

// ListSupporterResponse represents a response from
// Ulule's API to a GET /projects/:id/supporters request.
type ListSupporterResponse struct {
	Meta       *Metadata    `json:"meta"`
	Supporters []*Supporter `json:"supporters"`
}

type Supporter struct {
	Id         int    `json:"id"`
	Url        string `json:"absolute_url"`
	DateJoined string `json:"date_joined"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Name       string `json:"name"`
	UserName   string `json:"username"`
	TimeZone   string `json:"timezone"`
	IsStaff    bool   `json:"is_staff"`
}

// Metadata
type Metadata struct {
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	TotalCount int    `json:"total_count"`
	Next       string `json:"next"`
	Previous   string `json:"previous"`
}
