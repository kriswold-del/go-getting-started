package main

import (
	"log"
	"net/http"
	"os"

    "github.com/bp3d"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type bin struct {
    Name    string  `json:"name" binding:"required"`
    Width   float64  `json:"width" binding:"required"`
    Height  float64  `json:"height" binding:"required"`
    Depth   float64  `json:"depth" binding:"required"`
    Weight   float64  `json:"weight" binding:"required"`
}
var bins []bin

type item struct {
    Name    string  `json:"name" binding:"required"`
    Width   float64 `json:"width" binding:"required"`
    Height  float64 `json:"height" binding:"required"`
    Depth   float64 `json:"depth" binding:"required"`
    Weight  float64 `json:"weight" binding:"required"`
}

var items []item

type jsoninput struct {
    bins    bin `json:"bins"`
    items   item    `json:"items"`
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

    p := bp3d.NewPacker()
    p.AddBin(bp3d.NewBin("Small Bin", 10, 15, 20, 100))
	p.AddBin(bp3d.NewBin("Medium Bin", 100, 150, 200, 1000))

	// Add items.
	p.AddItem(bp3d.NewItem("Item 1", 2, 2, 1, 2))
	p.AddItem(bp3d.NewItem("Item 2", 3, 3, 2, 3))

	// Pack items to bins.
	if err := p.Pack(); err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		 c.JSON(http.StatusOK, p)
	})

    router.POST("/", func(c *gin.Context) {
    var requestBody jsoninput
        if err := c.BindJSON(&requestBody); err != nil {
            c.IndentedJSON(http.StatusOK, err.Error())
            return
        }
        c.IndentedJSON(http.StatusOK, requestBody.bins)
    })
	router.Run(":" + port)
}
