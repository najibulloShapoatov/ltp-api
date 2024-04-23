package server

import (
	"context"
	"fmt"
	"github.com/samber/do"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"ltp-api/docs"
	"ltp-api/internal/config"
	"ltp-api/internal/transport/http/handlers"
	"ltp-api/internal/transport/http/middlewares"
	"ltp-api/pkg/redis"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	ctx    context.Context
	wg     *sync.WaitGroup
	server *http.Server
	router *gin.Engine
	cache  *redis.Cache
}

// New
// @title           LTP API
// @version         1.0
// @description     LTP API
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
// @SecurityDefinitions.apikey X-API-Key
// @in header
// @name X-API-Key
// @SecurityDefinitions.apikey Authorization
// @in header
// @name Authorization
func New(i *do.Injector) (*Server, error) {
	cfg := do.MustInvoke[*config.Config](i).Server
	ctx := do.MustInvoke[context.Context](i)
	wg := do.MustInvoke[*sync.WaitGroup](i)
	cache := do.MustInvoke[*redis.Cache](i)
	publicHandlers := do.MustInvokeNamed[[]handlers.Handler](i, "publicHandlers")

	docs.SwaggerInfo.Host = cfg.SwaggerAddress
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	s := &Server{
		ctx: ctx,
		wg:  wg,
		server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:           nil,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       30 * time.Second,
		},
		router: gin.New(),
		cache:  cache,
	}

	s.router.Use(middlewares.Recover())
	s.router.Use(middlewares.CORSMiddleware())
	s.router.Use(middlewares.AccessLogMiddleware())
	s.router.NoRoute(handlers.NotFoundHandler)
	s.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.registerHandlers(publicHandlers...)

	return s, nil
}

func (s *Server) Run() {
	s.wg.Add(1)
	zap.S().Infof("server listining: %s", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.S().Error(err.Error())
	}
}

func (s *Server) registerHandlers(handlers ...handlers.Handler) {
	api := s.router.Group("api/v1")
	for _, h := range handlers {
		h.Register(api)
	}

	s.server.Handler = s.router
}

func (s *Server) Shutdown(ctx context.Context) error {
	zap.S().Info("Shutdown server...")
	zap.S().Info("Stopping http server...")

	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer func() {
		cancel()
		s.wg.Done()
	}()

	if err := s.server.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server forced to shutdown:", zap.Error(err))
		return err
	}

	zap.S().Info("Server successfully stopped.")

	return nil
}
