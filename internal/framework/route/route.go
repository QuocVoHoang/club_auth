package route

import (
	"github.com/your-org/go-base/internal/config"
	"github.com/your-org/go-base/internal/infrastructure/handler"
	middlewarepkg "github.com/your-org/go-base/pkg/middleware"

	_ "github.com/your-org/go-base/internal/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @title           Go Base Monolith API
// @version         1.0
// @description     Base monolith service for Go projects
// @BasePath        /
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer token. Format: "Bearer {token}"
func NewRouter(
	cfg *config.Config,
	db *gorm.DB,
	jwtService middlewarepkg.JWT,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) *gin.Engine {
	if cfg.ENV == config.ENVProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(
		middlewarepkg.RequestID(),
		gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: middlewarepkg.LogFormatterJSON,
			Output:    gin.DefaultWriter,
			SkipPaths: []string{"/", "/api/healthz", "/api/readyz"},
		}),

		middlewarepkg.Recovery(),
		middlewarepkg.Secure(),
		middlewarepkg.Headers,
	)

	if cfg.ENV == config.ENVDevelopment || cfg.ENV == config.ENVProduction {
		router.Use(middlewarepkg.CorsMiddleware(cfg.CORS.AllowHosts))
	}

	router.GET("/api/healthz", middlewarepkg.Health)
	router.GET("/api/readyz", middlewarepkg.Readiness(db))
	router.GET("/api/v1/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		publicGroup := v1.Group("")
		{
			authGroup := publicGroup.Group("/auth")
			{
				authGroup.POST("/register", authHandler.Register)
				authGroup.POST("/login", authHandler.Login)
			}
		}

		protectedGroup := v1.Group("")
		protectedGroup.Use(middlewarepkg.AuthRequired(jwtService))
		{
			userGroup := protectedGroup.Group("/users")
			{
				userGroup.GET("/me", userHandler.GetMe)
				userGroup.PATCH("/me", userHandler.UpdateMe)
			}
		}

	}

	return router
}
