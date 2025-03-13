package handlers

import (
	"rokhelper/db"
	"rokhelper/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type KingdomHandler struct {
	DB *db.Mongo
}

// CreateKingdom godoc
// @Summary Create a kingdom
// @Description Create a kingdom
// @Tags kingdom
// @Accept json
// @Produce json
// @Param kingdom body model.Kingdom true "Kingdom object"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/kingdom [post]

func (k *KingdomHandler) CreateKingdom(c *fiber.Ctx) error {
	var kingdom = model.NewKingdom()
	if err := c.BodyParser(&kingdom); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Bad Request"})
	}

	db := k.DB.Client.Database("rokhelper")
	_, err := db.Collection("kingdoms").InsertOne(c.Context(), kingdom)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Kingdom created successfully"})
}

// GetKingdoms godoc
// @Summary Get kingdoms
// @Description Get kingdoms
// @Tags kingdom
// @Accept json
// @Produce json
// @Success 200 {object} []model.Kingdom
// @Failure 500 {object} string
// @Router /api/kingdom [get]

func (k *KingdomHandler) GetKingdoms(c *fiber.Ctx) error {
	db := k.DB.Client.Database("rokhelper")
	var kingdoms []model.Kingdom

	cursor, err := db.Collection("kingdoms").Find(c.Context(), nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	defer cursor.Close(c.Context())
	if err := cursor.All(c.Context(), &kingdoms); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	return c.JSON(kingdoms)
}

// GetKingdom godoc
// @Summary Get a kingdom
// @Description Get a kingdom
// @Tags kingdom
// @Accept json
// @Produce json
// @Param id path string true "Kingdom ID or Discord Channel ID"
// @Success 200 {object} model.Kingdom
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/kingdom/{id} [get]

func (k *KingdomHandler) GetKingdom(c *fiber.Ctx) error {
	id := c.Params("id")
	db := k.DB.Client.Database("rokhelper")
	var kd model.Kingdom

	// Define filter for MongoDB query
	var ObjectId, _ = primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"$or": []bson.M{
			{"_id": ObjectId},
			{"discord_channel_id": id},
		},
	}

	// Find one document based on the filter
	err := db.Collection("kingdoms").FindOne(c.Context(), filter).Decode(&kd)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"message": "Not Found"})
		}
		return c.Status(500).JSON(fiber.Map{"message": "Internal Server Error"})
	}

	return c.JSON(kd)
}
