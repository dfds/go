package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func NewHttpServer(enableDebug bool) (*gin.Engine, *http.Server) {
	if enableDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery(), gin.ErrorLogger())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// HTTP server
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return router, srv
}

func NewMetricsServer() {
	router := gin.New()
	router.Use(gin.Recovery(), gin.ErrorLogger())

	router.GET("/metrics", MetricsHandler())

	srv := &http.Server{
		Addr:    ":9090",
		Handler: router,
	}

	// HTTP server
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
