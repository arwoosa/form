package ezapi

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/arwoosa/vulpes/log"

	"github.com/gin-gonic/gin"
)

var (
	// engine is the singleton gin.Engine instance.
	engine *gin.Engine
	// routers holds all the registered routes before the engine is initialized.
	routers = newRouterGroup()
	// defaultMiddelware is the set of default middleware used by the gin engine.
	defaultMiddelware = []gin.HandlerFunc{
		gin.Recovery(),
		gin.Logger(),
	}
)

// RegisterGinApi allows for the registration of API routes using a function.
// This function can be called from anywhere to add routes to the central routerGroup.
func RegisterGinApi(f func(router Router)) {
	f(routers)
}

// initEngin initializes the gin engine as a singleton.
// It sets up the default middleware and registers all the routes that have been collected.
func initEngin() {
	once.Do(func() {
		engine = gin.New()
		engine.Use(defaultMiddelware...)
		routers.register(engine)
	})
}

// server creates and configures an *http.Server with the gin engine as its handler.
func server(port int) *http.Server {
	portStr := fmt.Sprintf(":%d", port)
	log.Info("api service listen on port " + portStr)
	initEngin()
	return &http.Server{
		Addr:              portStr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           engine,
	}
}

// once ensures that the engine is initialized only once.
var once sync.Once

// GetHttpHandler returns the singleton gin.Engine instance as an http.Handler.
// This allows the gin engine to be used with an existing http.Server.
func GetHttpHandler() http.Handler {
	initEngin()
	return engine
}

// RunGin starts the gin server on the specified port and handles graceful shutdown.
// It blocks until the provided context is canceled.
func RunGin(ctx context.Context, port int) error {
	ser := server(port)
	var apiWait sync.WaitGroup
	apiWait.Add(1)
	go func(srv *http.Server) {
		defer apiWait.Done()
		for {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Info("api service listen failed: " + err.Error())
				time.Sleep(5 * time.Second)
			} else if err == http.ErrServerClosed {
				return
			}
		}
	}(ser)
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ser.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown failed: " + err.Error())
	}
	apiWait.Wait()
	return nil
}
