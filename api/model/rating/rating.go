package rating

import "time"

type Rating struct {
	Date    time.Time `bson:"date"`
	Rate    int       `bson:"rate"`
	RateFor string    `bson:"rate_for"`
	RateBy  string    `bson:"rate_by"`
}
