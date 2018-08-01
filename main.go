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
		result += fmt.Sprintf("Finished: %s", gettime(time.Now().Unix()))
	})

	co.Visit("https://www.google.com/search?q=mace+windu")
	return result
}

func gettime(i int64) string {
	return fmt.Sprintf("Current Unix Time: %v<br>", i)
}
