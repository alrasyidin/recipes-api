// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/alrasyidin/recipes-api
//
// Schemes: http
// Host: localhost:8080
// Basepath: /
// Version: 1.0.0
// Contact: Hafidh Pradipta<hamstergeek38@gmail.com> https://github.com/alrasyidin
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"recipes-api/handlers"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

var ctx context.Context
var err error
var client *mongo.Client

var recipesHandler handlers.IRecipeHandler

func init() {
	ctx = context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("failed connect to mongo db:", err)
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to mongo db")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	recipesHandler = handlers.NewRecipeHandler(ctx, collection)
}

// swagger:operation post /recipes recipes newRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	  description: Succesfull operation
//	'400':
//	  description: invalid input
func NewRecipeHandler(ctx *gin.Context) {
	var recipe Recipe

	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})

		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)

	ctx.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// Return list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	  description: Succesfull operation
func ListRecipeHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, recipes)
}

// swagger:operation PUT /recipes/:id recipes updateRecipes
// Update existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: id of recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//	  description: Succesfull operation
//	'400':
//	  description: invalid input
//	'404':
//	  description: invalid recipe ID
func UpdateRecipeHandler(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	var recipe Recipe
	if err := ctx.ShouldBindJSON(&recipe); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})

		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = 1
		}
	}

	if index == -1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"errors": "recipe not found",
		})
		return
	}

	recipes[index] = recipe

	ctx.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/:id recipes deleteRecipes
// Delete existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: id of recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//	  description: Succesfull operation
//	'404':
//	  description: invalid recipe ID
func DeleteRecipeHandler(ctx *gin.Context) {
	id := ctx.Params.ByName("id")

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = 1
		}
	}

	if index == -1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"errors": "recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

// swagger:operation GET /recipes/search recipes searchRecipes
// Return searched of recipes
// ---
// parameters:
//   - name: tag
//     in: query
//     description: recipe tag
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//	  description: Succesfull operation
func SearchRecipeHandler(ctx *gin.Context) {
	tag := ctx.Query("tag")

	var listRecipes []Recipe
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listRecipes = append(listRecipes, recipes[i])
		}
	}
	ctx.JSON(http.StatusOK, listRecipes)
}

// swagger:operation GET /recipes/{id} recipes oneRecipe
// Get one recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'404':
//	    description: Invalid recipe ID
func GetRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			c.JSON(http.StatusOK, recipes[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

func main() {
	router := gin.Default()

	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipeHandler)
	// router.GET("/recipes/:id", recipesHandler.GetRecipeHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)

	err := router.Run()
	if err != nil {
		log.Fatal("cannot running server:", err)
	}
}
