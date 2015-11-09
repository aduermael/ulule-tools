package clientapi

import ()

// ListProjectResponse represents a response from
// Ulule's API to a GET */projects request.
type ListProjectResponse struct {
	Meta     *Metadata  `json:"meta"`
	Projects []*Project `json:"projects"`
}

// Project represents an Ulule project
type Project struct {
	Id              float64 `json:"id"`
	Url             string  `json:"absolute_url"`
	Goal            int     `json:"goal"`
	GoalRaised      bool    `json:"goal_raised"`
	AmountRaised    int     `json:"amount_raised"`
	CommentCount    int     `json:"comments_count"`
	Committed       int     `json:"committed"`
	Currency        string  `json:"currency"`
	CurrencyDisplay string  `json:"currency_display"`
	DateEnd         string  `json:"date_end"`
	DateStart       string  `json:"date_start"`
	Finished        bool    `json:"finished"`
	Slug            string  `json:"slug"`
	SupportersCount int     `json:"supporters_count"`
	TimeZone        string  `json:"timezone"`
}

// ListSupporterResponse represents a response from
// Ulule's API to a GET /projects/:id/supporters request.
type ListSupporterResponse struct {
	Meta       *Metadata    `json:"meta"`
	Supporters []*Supporter `json:"supporters"`
}

// Supporter represents an Ulule project supporter
type Supporter struct {
	Id         float64 `json:"id"`
	Url        string  `json:"absolute_url"`
	DateJoined string  `json:"date_joined"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Name       string  `json:"name"`
	UserName   string  `json:"username"`
	TimeZone   string  `json:"timezone"`
	IsStaff    bool    `json:"is_staff"`
}

// ListOrderResponse represents a response from
// Ulule's API to a GET /projects/:id/orders request.
type ListOrderResponse struct {
	Meta   *Metadata `json:"meta"`
	Orders []*Order  `json:"orders"`
}

// Order represents an Ulule project order
type Order struct {
	Id              float64      `json:"id"`
	Url             string       `json:"absolute_url"`
	Subtotal        string       `json:"order_subtotal"`
	Total           string       `json:"order_total"`
	PaymentMethod   string       `json:"payment_method"`
	Status          OrderStatus  `json:"status"`
	StatusDisplay   string       `json:"status_display"`
	Items           []*OrderItem `json:"items"`
	User            *Supporter   `json:"user"`
	ShippingAddress *Address     `json:"shipping_address,omitempty"`
	BillingAddress  *Address     `json:"billing_address,omitempty"`
}

type OrderStatus int8

const (
	OrderStatusAwaiting  OrderStatus = 3
	OrderStatusCompleted OrderStatus = 4
	OrderStatusCancelled OrderStatus = 6
)

// OrderItem represents an Ulule project order item
type OrderItem struct {
	UnitPrice string `json:"unit_price"`
	Quantity  int    `json:"quantity"`
	Product   int    `json:"product"`
	LineTotal string `json:"line_total"`
}

// Address represents a postal address
type Address struct {
	Address1   string `json:"address1,omitempty"`
	Address2   string `json:"address2,omitempty"`
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode int    `json:"postal_code,omitempty"`
}

// Metadata
type Metadata struct {
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	TotalCount int    `json:"total_count"`
	Next       string `json:"next"`
	Previous   string `json:"previous"`
}
