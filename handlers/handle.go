package handlers

import (
	"context"
	"net/http"
	"recipes-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRecipeHandler interface {
	ListRecipeHandler(c *gin.Context)
	NewRecipeHandler(c *gin.Context)
	UpdateRecipeHandler(c *gin.Context)
}

type RecipeHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

// Constructor for RecipeHandler
func NewRecipeHandler(ctx context.Context, collection *mongo.Collection) *RecipeHandler {
	return &RecipeHandler{
		ctx:        ctx,
		collection: collection,
	}
}

func (handler *RecipeHandler) ListRecipeHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cur.Close(handler.ctx)
	recipes := make([]models.Recipe, 0)

	for cur.Next(handler.ctx) {
		var recipe models.Recipe

		cur.Decode(&recipe)

		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, recipes)
}

func (mp *RecipeHandler) NewRecipeHandler(c *gin.Context, arg models.Recipe) {
	panic("not implemented") // TODO: Implement
}

func (mp *RecipeHandler) UpdateRecipeHandler(c *gin.Context, arg models.Recipe) {
	panic("not implemented") // TODO: Implement
}
