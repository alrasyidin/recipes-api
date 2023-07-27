package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
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

func init() {
	recipes = make([]Recipe, 0)

	file, err := ioutil.ReadFile("recipes.json")
	if err != nil {
		log.Fatal("cannot open file")
	}
	err = json.Unmarshal([]byte(file), &recipes)
	if err != nil {
		log.Fatal("failed unmarhshal json recipes")
	}
}

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

func ListRecipeHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, recipes)
}

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

func main() {
	router := gin.Default()

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)

	err := router.Run()
	if err != nil {
		log.Fatal("cannot running server:", err)
	}
}
