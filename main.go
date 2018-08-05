package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"html/template"
	"github.com/gocolly/colly"
	"time"
	"fmt"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", map[string]interface{}{
			"results": template.HTML(scrapeGoogle()),
		})
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
		c.HTML(http.StatusOK, "charts.tmpl.html", nil)
	})

	router.Run(":" + port)
}

func scrapeGoogle() string {
	startTime := gettime(time.Now().Unix())
	result := ""
	google := "www.google.com"
	co := colly.NewCollector(
		colly.AllowedDomains(google),
		//colly.CacheDir("./macewindu_cache"),
	)
	co.OnRequest(func(r *colly.Request) {
		result += startTime + fmt.Sprintf("Visiting: <a href=\"%s\" target=\"_blank\" rel=\"noreferrer\">%s</a><br><br>\n", r.URL, r.URL)
	})

	// Find all links
	co.OnHTML("a[href^='/url?q']", func(e *colly.HTMLElement) {
		if len(e.Text) > 0 && e.Text != "Cached"  { // fix ??
			result += "<a href=\"" + "https://" + google + e.Attr("href") +
				"\" target=\"_blank\" rel=\"noreferrer\">" + e.Text + "</a><br>\n"
		}
	})

	co.OnScraped(func(r *colly.Response) {
		t := time.Now()
		result += fmt.Sprintf("Finished: %s<br>\n%d-%02d-%02dT%02d:%02d:%02d\n",
			gettime(t.Unix()), t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	})

	co.Visit("https://www.google.com/search?q=mace+windu")
	return result
}

func gettime(i int64) string {
	return fmt.Sprintf("Current Unix Time: %v<br>", i)
}
