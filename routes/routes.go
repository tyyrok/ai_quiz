package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Answer struct {
	Id int `json:"id"`
	Title string `json:"text"`
	Likes int `json:"likes"`
	Dislikes int `json:"dislikes"`
	Users_answered int `json:"users_answered"`
}

type Question struct {
	Id int `json:"id"`
	Title string `json:"text"`
	Likes int `json:"likes"`
	Dislikes int `json:"dislikes"`
	Answers []Answer `json:"answers"`
}

func NewRouter() *gin.Engine {
	// Set the router as the default one shipped with Gin
	router := gin.Default()
	expectedHost := os.Getenv("ORIGIN")
	port := os.Getenv("PORT")
	var host string
	if port == "80" {
		host = expectedHost
	} else {
		host = fmt.Sprintf("%s:%s", expectedHost, port)
	}

	// Setup Security Headers
	router.Use(func(c *gin.Context) {
		if c.Request.Host != host {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src 'self' 'unsafe-inline' 'unsafe-eval' *; script-src-elem * 'unsafe-inline' *; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	})

	return router
}


func Run(dbpool *pgxpool.Pool) {
	httpPort := os.Getenv("PORT")
	router := NewRouter()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("db", dbpool)
		ctx.Next()
	})
	router.LoadHTMLGlob("templates/*")

	v1 := router.Group("/api", CheckOrigin())

	router.GET("/", mainPageHandler)

	v1.POST("/:question_id/:answer_id", answerHandler)
	
	v1.PATCH("/:question_id/:answer_id", answerLikeHandler)

	v1.PATCH("/:question_id", questionLikeHandler)

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}


func CheckOrigin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		schema := os.Getenv("SCHEMA")
		expectedHost := os.Getenv("ORIGIN")
		port := os.Getenv("PORT")
		var host string
		if port == "80" {
			host = fmt.Sprintf("%s%s", schema, expectedHost)
		} else {
			host = fmt.Sprintf("%s%s:%s", schema, expectedHost, port)
		}
		origin := ctx.GetHeader("Origin")
		referer := ctx.GetHeader("Referer")

		if origin != "" && origin != host {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid origin"})
			return
		}
		if referer != "" && !strings.HasPrefix(referer, host) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid referer"})
			return
		}

		ctx.Next()
	}
}