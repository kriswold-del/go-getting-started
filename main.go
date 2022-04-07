package main

import (
	"log"
	"net/http"
	"os"

    "github.com/bp3d"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type payload struct {
    bins    bins
    items   items
}

type bins struct {
    Name    string
    Width   float64
    Height  float64
    Depth   float64
    MaxWeight   float64
}

type items struct {
    Name    string
    Width   float64
    Height  float64
    Depth   float64
    Weight  float64
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
        var newpayload payload

        if err := c.BindJSON(&payload); err != nil {
            return c.IndentedJSON(http.StatusOK, err)
        }
        payload = append(payload, newpayload)
        c.IndentedJSON(http.StatusOK, payload)
    })

	router.Run(":" + port)
}
