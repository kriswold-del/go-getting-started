package main

import (
	"log"
	"net/http"
	"os"
    "encoding/json"
    "io/ioutil"
    "sort"

    "github.com/bp3d"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)
type Response struct {
	Bins []struct {
		Name   string `json:"name"`
		Width  float64    `json:"width"`
		Height float64    `json:"height"`
		Depth  float64    `json:"depth"`
		Weight float64    `json:"weight"`
	} `json:"bins"`
	Items []struct {
		Name   string `json:"name"`
		Width  float64    `json:"width"`
		Height float64    `json:"height"`
		Depth  float64    `json:"depth"`
		Weight float64    `json:"weight"`
	} `json:"items"`
}


type packingreults []bp3d.Packer

func test(w http.ResponseWriter, r *http.Request) {
    var jsonObj Response

    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Println(err.Error())
    }
    json.Unmarshal(reqBody, &jsonObj)
    sort.SliceStable(jsonObj.Bins, func(i, j int) bool {
            return (jsonObj.Bins[i].Width * jsonObj.Bins[i].Height * jsonObj.Bins[i].Depth) > (jsonObj.Bins[j].Width * jsonObj.Bins[j].Height * jsonObj.Bins[j].Depth)
    })
    for i := range jsonObj.Bins {
        p := bp3d.NewPacker()
        p.AddBin(bp3d.NewBin(
        jsonObj.Bins[i].Name,
        jsonObj.Bins[i].Width,
        jsonObj.Bins[i].Height,
        jsonObj.Bins[i].Depth,
        jsonObj.Bins[i].Weight))
        for i := range jsonObj.Items {
            p.AddItem(bp3d.NewItem(
            jsonObj.Items[i].Name,
            jsonObj.Items[i].Width,
            jsonObj.Items[i].Height,
            jsonObj.Items[i].Depth,
            jsonObj.Items[i].Weight))
        }
        if err := p.Pack(); err != nil {
            log.Fatal(err)
        }
            packingreults = append(packingreults, p)
        }


    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")



    json.NewEncoder(w).Encode(packedbins)
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
