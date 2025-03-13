package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Cấu trúc của Kingdom
type Kingdom struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty" `              // ID của kingdom (sử dụng ObjectID của MongoDB)
	Name             string               `bson:"name" json:"name"`                                 // Tên của kingdom
	KingID           primitive.ObjectID   `bson:"king_id,omitempty" json:"king_id,omitempty"`       // ID của King
	R4IDs            []primitive.ObjectID `bson:"r4_ids,omitempty" json:"r4_ids,omitempty"`         // Danh sách ID của các R4
	DiscordChannelID string               `bson:"discord_channel_id" json:"discord_channel_id"`     // ID của discordChannel
	PlayerIDs        []primitive.ObjectID `bson:"player_ids,omitempty" json:"player_ids,omitempty"` // Danh sách ID của các player trong kingdom
	Title            Title                `bson:"title" json:"title"`                               // Title của kingdom

}

func NewKingdom() Kingdom {
	var kingdom Kingdom
	kingdom.Title = *kingdom.Title.NewTitle()
	return kingdom
}
