package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/gocolly/colly"
	"time"
	"fmt"
	"math/rand"
)

type ScrappedData struct {
	Data [][]string
	Address string
	StartTime string
	Time string
	RegTime string
}

func main() {
	fmt.Println("Starting...")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	d := scrapeGoogle()
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", d)
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/the/real/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.GET("/google-charts", func(c *gin.Context) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		offsetNumber := r.Intn(7)
		offsetAmount := 0.2

		c.HTML(http.StatusOK, "charts.tmpl.html", map[string]interface{}{
			"pie_cool": r.Intn(50),
			"pie_battles": r.Intn(50),
			"pie_sleep": r.Intn(50),
			"pie_council": r.Intn(50),
			"pie_eat": r.Intn(50),
			"pie_commute": r.Intn(50),
			"pie_tv": r.Intn(50),
			"offset_number": offsetNumber,
			"offset_amount": offsetAmount,
		})
	})

	router.Run(":" + port)
}

func scrapeGoogle() ScrappedData {
	var url string
	var startTime string
	var unixTime string
	var regTime string

	startTime = gettime(time.Now().Unix())
	google := "www.google.com"

	co := colly.NewCollector(
		colly.AllowedDomains(google),
		//colly.CacheDir("./macewindu_cache"),
	)
	co.OnRequest(func(r *colly.Request) {
		url = fmt.Sprintf("%v", r.URL)
	})
	var ret [][]string
	// Find all links
	co.OnHTML("a[href^='/url?q']", func(e *colly.HTMLElement) {
		if len(e.Text) > 0 && e.Text != "Cached"  { // fix ??
			var record = make([]string, 2)
			record[0] = e.Text
			link := e.Attr("href")
			record[1] = e.Request.AbsoluteURL(link)
			ret = append(ret, record)
		}
	})

	co.OnScraped(func(r *colly.Response) {
		t := time.Now()
		unixTime = fmt.Sprintf("%s",gettime(t.Unix()))
		regTime = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d\n",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	})

	co.Visit("https://www.google.com/search?q=mace+windu")
	d := ScrappedData{Data: ret, Address: url, StartTime: startTime, Time: unixTime, RegTime: regTime}
	return d
}

func gettime(i int64) string {
	return fmt.Sprintf("%v", i)
}
