package models

import "time"

type Excuse struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Category  string    `json:"category"`
	Language  string    `json:"language"`
	Severity  string    `json:"severity"`
	CreatedAt time.Time `json:"created_at"`
	Rating    int       `json:"rating"`
}

type ExcuseRequest struct {
	Text     string `json:"text"`
	Category string `json:"category"`
	Language string `json:"language"`
	Severity string `json:"severity"`
}

type RatingRequest struct {
	Upvote bool `json:"upvote"`
}

type Stats struct {
	TotalExcuses               int    `json:"total_excuses"`
	MostPopularCategory        string `json:"most_popular_category"`
	ExcusesToday               int    `json:"excuses_today"`
	GlobalProcrastinationLevel string `json:"global_procrastination_level"`
}
