package bot

import "go.mongodb.org/mongo-driver/bson"

type Chat struct {
	Name string `json:"name" bson:"name,omitempty"`
	ID int64 `json:"id" bson:"id,omitempty"`
	IsActive bool `json:"isactive" bson:"isactive,omitempty"`
}

func (c *Chat) ToDoc() bson.D {
	return bson.D{
		{ Key: "name", Value: c.Name},
		{ Key: "id", Value: c.ID },
		{ Key: "isactive", Value: c.IsActive },
	}
}
