package handlers

import (
	"rokhelper/db"
	"rokhelper/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TitleHandler struct {
	DB *db.Mongo
}

// Helper function for retrieving Kingdom by ID or Discord Channel ID
func (t *TitleHandler) findKingdomByID(id string, c *fiber.Ctx) (model.Kingdom, error) {
	var kd model.Kingdom
	db := t.DB.Client.Database("rokhelper")
	ObjectId, err := primitive.ObjectIDFromHex(id)
	var filter bson.M
	if err != nil {
		filter = bson.M{
			"$or": []bson.M{
				{"_id": id},
				{"discord_channel_id": id},
			},
		}
	} else {
		filter = bson.M{
			"$or": []bson.M{
				{"_id": ObjectId},
				{"discord_channel_id": id},
			},
		}
	}
	err = db.Collection("kingdoms").FindOne(c.Context(), filter).Decode(&kd)
	return kd, err
}

// General error response function

func handleError(c *fiber.Ctx, err error, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{"message": message, "error": err.Error()})
}

// GetTitles godoc
// @Summary Get titles
// @Description Retrieve all titles for a Kingdom or Discord Channel
// @Tags title
// @Accept json
// @Produce json
// @Param id path string true "Kingdom ID or Discord Channel ID"
// @Success 200 {object} model.TitleAssignment
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/title/{id} [get]
func (t *TitleHandler) GetTitle(c *fiber.Ctx) error {
	id := c.Params("id")
	kd, err := t.findKingdomByID(id, c)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"message": "Not Found"})
		}
		return handleError(c, err, 500, "Internal Server Error")
	}
	result, cheker := kd.Title.GetTitleAssignment()
	if cheker == false {
		return c.Status(404).JSON(fiber.Map{"message": "Not Found"})
	}
	if result.Local.Map == "Home Kingdom" {
		result.Local.Map = kd.Title.HomeKingdomMap
	} else {
		result.Local.Map = kd.Title.LostKingdomMap
	}
	return c.JSON(result)
}

// AddTitle godoc
// @Summary Add title
// @Description Add a new title assignment
// @Tags title
// @Accept json
// @Produce json
// @Param id path string true "Kingdom ID or Discord Channel ID"
// @Param titleAssignment body model.TitleAssignment true "Title Assignment"
// @Success 200 {object} model.TitleAssignment
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/title/{id} [post]
func (t *TitleHandler) AddTitle(c *fiber.Ctx) error {
	id := c.Params("id")
	kd, err := t.findKingdomByID(id, c)
	if err != nil {
		return handleError(c, err, 404, "Not Found")
	}
	var titleAssig model.TitleAssignment
	if err := c.BodyParser(&titleAssig); err != nil {
		return handleError(c, err, 400, "Bad Request")
	}
	if !kd.Title.AddTitle(titleAssig) {
		return c.Status(400).JSON(fiber.Map{"message": "User already in queue or title is not valid"})
	}
	_, err = t.DB.Client.Database("rokhelper").Collection("kingdoms").UpdateOne(
		c.Context(), bson.M{"_id": kd.ID}, bson.M{"$set": bson.M{"title": kd.Title}},
	)
	if err != nil {
		return handleError(c, err, 500, "Internal Server Error")
	}
	return c.JSON(fiber.Map{"message": "Success"})
}

// FinishTitle godoc
// @Summary Finish title
// @Description Mark a title assignment as finished
// @Tags title
// @Accept json
// @Produce json
// @Param id path string true "Kingdom ID or Discord Channel ID"
// @Param titleAssignment body model.TitleAssignment true "Title Assignment"
// @Success 200 {object} model.TitleAssignment
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/title/finish/{id} [post]
func (t *TitleHandler) FinishTitle(c *fiber.Ctx) error {
	id := c.Params("id")
	kd, err := t.findKingdomByID(id, c)
	if err != nil {
		return handleError(c, err, 404, "Not Found")
	}
	var titleAssig model.TitleAssignment
	if err := c.BodyParser(&titleAssig); err != nil {
		return handleError(c, err, 400, "Bad Request")
	}
	kd.Title.Finish(titleAssig)
	_, err = t.DB.Client.Database("rokhelper").Collection("kingdoms").UpdateOne(
		c.Context(), bson.M{"_id": kd.ID}, bson.M{"$set": bson.M{"title": kd.Title}},
	)
	if err != nil {
		return handleError(c, err, 500, "Internal Server Error")
	}
	return c.JSON(fiber.Map{"message": "Success"})
}

// DoneTitle godoc
// @Summary Done title
// @Description Mark a title assignment as done
// @Tags title
// @Accept json
// @Produce json
// @Param id path string true "Kingdom ID or Discord Channel ID"
// @Param titleAssignment body model.TitleAssignment true "Title Assignment"
// @Success 200 {object} model.TitleAssignment
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/title/done/{id} [post]
func (t *TitleHandler) DoneTitle(c *fiber.Ctx) error {
	id := c.Params("id")
	kd, err := t.findKingdomByID(id, c)
	if err != nil {
		return handleError(c, err, 404, "Not Found")
	}
	var titleAssig model.TitleAssignment
	if err := c.BodyParser(&titleAssig); err != nil {
		return handleError(c, err, 400, "Bad Request")
	}
	kd.Title.Done(titleAssig)
	_, err = t.DB.Client.Database("rokhelper").Collection("kingdoms").UpdateOne(
		c.Context(), bson.M{"_id": kd.ID}, bson.M{"$set": bson.M{"title": kd.Title}},
	)
	if err != nil {
		return handleError(c, err, 500, "Internal Server Error")
	}
	return c.JSON(fiber.Map{"message": "Success"})
}

// EditConfig godoc
// @Summary Edit config
// @Description Edit title configuration
// @Tags config
// @Accept json
// @Produce json
// @Param id path string true "Kingdom ID or Discord Channel ID"
// @Param config body model.Config true "Config"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/config/{id} [put]
func (t *TitleHandler) EditConfig(c *fiber.Ctx) error {
	var config model.Config
	// chuyêrn body về dạng map[string]interface{}
	jsonData := make(map[string]interface{})
	if err := c.BodyParser(&jsonData); err != nil {
		return handleError(c, err, 400, "Bad Request")
	}

	id := c.Params("id")
	kd, err := t.findKingdomByID(id, c)
	if err != nil {
		return handleError(c, err, 404, "Not Found")
	}

	if jsonData["duke"] != nil {
		dukeStr := jsonData["duke"].(string)
		dukeValue, err := strconv.ParseInt(dukeStr, 10, 64)
		if err != nil {
			return handleError(c, err, 400, "Invalid Duke value")
		}
		config.Duke = dukeValue
	} else {
		config.Duke = kd.Title.Config.Duke
	}

	if jsonData["architect"] != nil {
		architectStr := jsonData["architect"].(string)
		architectValue, err := strconv.ParseInt(architectStr, 10, 64)
		if err != nil {
			return handleError(c, err, 400, "Invalid Architect value")
		}
		config.Architect = architectValue
	} else {
		config.Architect = kd.Title.Config.Architect
	}

	if jsonData["scientist"] != nil {
		scientistStr := jsonData["scientist"].(string)
		scientistValue, err := strconv.ParseInt(scientistStr, 10, 64)
		if err != nil {
			return handleError(c, err, 400, "Invalid Scientist value")
		}
		config.Scientist = scientistValue
	} else {
		config.Scientist = kd.Title.Config.Scientist
	}

	if jsonData["justice"] != nil {
		justiceStr := jsonData["justice"].(string)
		justiceValue, err := strconv.ParseInt(justiceStr, 10, 64)
		if err != nil {
			return handleError(c, err, 400, "Invalid Justice value")
		}
		config.Justice = justiceValue
	} else {
		config.Justice = kd.Title.Config.Justice
	}

	_, err = t.DB.Client.Database("rokhelper").Collection("kingdoms").UpdateOne(
		c.Context(), bson.M{"_id": kd.ID}, bson.M{"$set": bson.M{"title.config": config}},
	)
	if err != nil {
		return handleError(c, err, 500, "Internal Server Error")
	}
	return c.JSON(fiber.Map{"message": "Success"})
}

// GetConfig godoc
// @Summary Get config
// @Description Get the title configuration
// @Tags config
// @Accept json
// @Produce json
// @Param id path string true "Kingdom ID or Discord Channel ID"
// @Success 200 {object} model.Config
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/config/{id} [get]
func (t *TitleHandler) GetConfig(c *fiber.Ctx) error {
	id := c.Params("id")
	kd, err := t.findKingdomByID(id, c)
	if err != nil {
		return handleError(c, err, 404, "Not Found")
	}
	configWithMaps := fiber.Map{
		"config":           kd.Title.Config,
		"home_kingdom_map": kd.Title.HomeKingdomMap,
		"lost_kingdom_map": kd.Title.LostKingdomMap,
	}

	return c.JSON(configWithMaps)
}

// EditMap godoc
// @Summary Edit maps
// @Description Edit home and lost kingdom maps
// @Tags map
// @Accept json
// @Produce json
// @Param id path string true "Kingdom ID or Discord Channel ID"
// @Param body body map[string]interface{} true "Maps"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /api/maps/{id} [put]
func (t *TitleHandler) EditMap(c *fiber.Ctx) error {
	id := c.Params("id")
	kd, err := t.findKingdomByID(id, c)
	if err != nil {
		return handleError(c, err, 404, "Not Found")
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return handleError(c, err, 400, "Bad Request")
	}

	updateFields := bson.M{}
	if homeMap, ok := body["home_kingdom_map"].(string); ok && homeMap != "" {
		updateFields["title.home_kingdom_map"] = homeMap
	}
	if lostMap, ok := body["lost_kingdom_map"].(string); ok && lostMap != "" {
		updateFields["title.lost_kingdom_map"] = lostMap
	}

	_, err = t.DB.Client.Database("rokhelper").Collection("kingdoms").UpdateOne(
		c.Context(), bson.M{"_id": kd.ID}, bson.M{"$set": updateFields},
	)
	if err != nil {
		return handleError(c, err, 500, "Internal Server Error")
	}
	return c.JSON(fiber.Map{"message": "Success"})
}
