package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type Recipes struct {
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

func main() {
	router := gin.Default()

	err := router.Run()
	if err != nil {
		log.Fatal("cannot running server:", err)

	}
}
