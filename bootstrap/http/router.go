package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

// DefaultHttpAddr is the listen address used when NewHttpServer is given an
// empty addr.
const DefaultHttpAddr = ":8080"

// NewHttpServer starts the application HTTP router. addr is the listen address
// (e.g. ":8081" or "0.0.0.0:8081"); an empty addr falls back to DefaultHttpAddr
// so existing callers keep the historical :8080 behaviour.
func NewHttpServer(enableDebug bool, addr string) (*gin.Engine, *http.Server) {
	if enableDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if addr == "" {
		addr = DefaultHttpAddr
	}

	router := gin.New()
	router.Use(gin.Recovery(), gin.ErrorLogger())

	srv := &http.Server{
		Addr:    addr,
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
