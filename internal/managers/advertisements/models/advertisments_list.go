package models

import authModels "api-tests-template/internal/managers/auth/models"

type AdvertisementsList struct {
	Items      []Advertisement `json:"items"`
	Total      int             `json:"total"`
	HasNext    bool            `json:"has_next"`
	NextCursor string          `json:"next_cursor"`
}

type Advertisement struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Quantity    int             `json:"quantity"`
	Photos      []Photo         `json:"photos"`
	UserID      string          `json:"user_id"`
	User        authModels.User `json:"user"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

type Photo struct {
	ID              string `json:"id"`
	AdvertisementID string `json:"advertisement_id"`
	URL             string `json:"url"`
	SortOrder       int    `json:"sort_order"`
	CreatedAt       string `json:"created_at"`
}
