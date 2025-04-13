
package model

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type List struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Username    string `bson:"username" json:"username"`
    Date      time.Time          `bson:"date" json:"date"`
    Equipment string             `bson:"equipment" json:"equipment"`
    Location  string             `bson:"location" json:"location"`
}
