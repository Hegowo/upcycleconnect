package router

import (
	"net/http"

	"upcycleconnect/backend/config"
	"upcycleconnect/backend/handlers"
	"upcycleconnect/backend/middleware"
	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	audit := services.NewAuditService(db)

	authHandler := &handlers.AuthHandler{DB: db, Cfg: cfg}
	dashHandler := &handlers.DashboardHandler{DB: db}
	userHandler := &handlers.UserHandler{DB: db, Audit: audit}
	providerHandler := &handlers.ProviderHandler{DB: db, Audit: audit}
	categoryHandler := &handlers.CategoryHandler{DB: db, Audit: audit}
	prestationHandler := &handlers.PrestationHandler{DB: db, Audit: audit}
	eventHandler := &handlers.EventHandler{DB: db, Audit: audit}
	adminHandler := &handlers.AdminHandler{DB: db, Audit: audit}
	logsHandler := &handlers.LogsHandler{DB: db}

	r.GET("/up", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	api := r.Group("/api/admin/v1")
	{
		api.POST("/auth/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(db, cfg))
		protected.Use(middleware.IsAdmin())
		{
			protected.POST("/auth/logout", authHandler.Logout)
			protected.GET("/auth/me", authHandler.Me)

			protected.GET("/dashboard/stats", dashHandler.Stats)

			protected.GET("/logs", logsHandler.Index)
			protected.GET("/logs/activity", logsHandler.Activity)

			protected.GET("/users", userHandler.Index)
			protected.GET("/users/:id", userHandler.Show)
			protected.PUT("/users/:id", userHandler.Update)
			protected.POST("/users/:id/ban", userHandler.Ban)
			protected.POST("/users/:id/activate", userHandler.Activate)
			protected.DELETE("/users/:id", userHandler.Destroy)

			protected.GET("/providers", providerHandler.Index)
			protected.GET("/providers/:id", providerHandler.Show)
			protected.PUT("/providers/:id/status", providerHandler.UpdateStatus)

			protected.GET("/categories", categoryHandler.Index)
			protected.POST("/categories", categoryHandler.Store)
			protected.GET("/categories/:id", categoryHandler.Show)
			protected.PUT("/categories/:id", categoryHandler.Update)
			protected.DELETE("/categories/:id", categoryHandler.Destroy)
			protected.POST("/categories/:id/toggle", categoryHandler.Toggle)

			protected.GET("/prestations", prestationHandler.Index)
			protected.POST("/prestations", prestationHandler.Store)
			protected.GET("/prestations/:id", prestationHandler.Show)
			protected.PUT("/prestations/:id", prestationHandler.Update)
			protected.DELETE("/prestations/:id", prestationHandler.Destroy)
			protected.PUT("/prestations/:id/status", prestationHandler.UpdateStatus)

			protected.GET("/events", eventHandler.Index)
			protected.POST("/events", eventHandler.Store)
			protected.GET("/events/:id", eventHandler.Show)
			protected.PUT("/events/:id", eventHandler.Update)
			protected.DELETE("/events/:id", eventHandler.Destroy)
			protected.PUT("/events/:id/status", eventHandler.UpdateStatus)

			superAdmin := protected.Group("")
			superAdmin.Use(middleware.IsSuperAdmin())
			{
				superAdmin.GET("/admins", adminHandler.Index)
				superAdmin.POST("/admins", adminHandler.Store)
				superAdmin.GET("/admins/:id", adminHandler.Show)
				superAdmin.PUT("/admins/:id", adminHandler.Update)
				superAdmin.DELETE("/admins/:id", adminHandler.Destroy)
			}
		}
	}

	return r
}
