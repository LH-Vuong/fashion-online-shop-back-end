package rating

type Rating struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	Rate      int    `bson:"rate,omitempty" json:"rate"`
	RateFor   string `bson:"rate_for,omitempty" json:"rate_for"`
	RateBy    string `bson:"rate_by,omitempty" json:"rate_by"`
	CreatedAt int64  `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt int64  `bson:"updated_at,omitempty" json:"updated_at"`
}
