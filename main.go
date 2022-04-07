package main

import (
	"log"
	"net/http"
	"os"
    "encoding/json"
    "io/ioutil"

    "github.com/bp3d"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)
type Response struct {
	Bins []struct {
		Name   string `json:"name"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Depth  int    `json:"depth"`
		Weight int    `json:"weight"`
	} `json:"bins"`
	Items []struct {
		Name   string `json:"name"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Depth  int    `json:"depth"`
		Weight int    `json:"weight"`
	} `json:"items"`
}

func test(w http.ResponseWriter, r *http.Request) {
    var jsonObj Response
    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Println(err.Error())
    }
    json.Unmarshal(reqBody, &jsonObj)

    log.Println(jsonObj.Bins)
//     for i := range jsonObj.Bins {
//         log.Println(jsonObj[i].Name)
//     }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(jsonObj)
    //log.Println(t.bins)
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
