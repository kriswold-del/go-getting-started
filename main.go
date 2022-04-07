package main

import (
	"log"
	"net/http"
	"os"
    "encoding/json"

    "github.com/bp3d"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)
type Response struct {
    bins    []struct {
        Name    string  `json:"name" binding:"required"`
        Width   float64  `json:"width" binding:"required"`
        Height  float64  `json:"height" binding:"required"`
        Depth   float64  `json:"depth" binding:"required"`
        Weight   float64  `json:"weight" binding:"required"`
    } `json:"data"`
    items    []struct {
        Name    string  `json:"name" binding:"required"`
        Width   float64 `json:"width" binding:"required"`
        Height  float64 `json:"height" binding:"required"`
        Depth   float64 `json:"depth" binding:"required"`
        Weight  float64 `json:"weight" binding:"required"`
    } `json:"data"`
}

func test(rw http.ResponseWriter, req *http.Request) {
    req.ParseForm()
    log.Println(req.Form)
    //LOG: map[{"test": "that"}:[]]
    var r Response
    for key, _ := range req.Form {
        log.Println(key)
        //LOG: {"test": "that"}
        err := json.Unmarshal([]byte(key), &r)
        if err != nil {
            log.Println(err.Error())
        }
    }
    log.Println(r)
    //LOG: that
}

func main() {
    port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

    http.HandleFunc("/", test)

    log.Fatal(http.ListenAndServe(":" + port, nil))
}

func mainold() {
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
    var result Response
        if err := c.BindJSON(&result); err != nil {
            c.IndentedJSON(http.StatusOK, err.Error())
            return
        }
        c.IndentedJSON(http.StatusOK, result)
    })
	router.Run(":" + port)
}
