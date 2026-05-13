package router

import (
	"log"
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
	swaggerHandler := &handlers.SwaggerHandler{}
	userHandler := &handlers.UserHandler{DB: db, Audit: audit}
	providerHandler := &handlers.ProviderHandler{DB: db, Audit: audit}
	categoryHandler := &handlers.CategoryHandler{DB: db, Audit: audit}
	prestationHandler := &handlers.PrestationHandler{DB: db, Audit: audit}
	eventHandler := &handlers.EventHandler{DB: db, Audit: audit}
	adminHandler := &handlers.AdminHandler{DB: db, Audit: audit}
	logsHandler := &handlers.LogsHandler{DB: db}
	siretHandler    := &handlers.SiretHandler{Cfg: cfg}
	depositHandler         := &handlers.DepositHandler{DB: db, Audit: audit}
	collectionPointHandler := &handlers.CollectionPointHandler{DB: db, Audit: audit}
	mailer                 := services.NewMailer(cfg)
	userAuthHandler := &handlers.UserAuthHandler{DB: db, Cfg: cfg, Mailer: mailer}
	profileHandler  := &handlers.ProfileHandler{DB: db, Mailer: mailer}
	publicHandler   := &handlers.PublicHandler{DB: db}
	userDepositHandler := &handlers.UserDepositHandler{DB: db, Audit: audit}
	userProviderHandler := &handlers.UserProviderHandler{DB: db, Audit: audit}

	stripeService := services.NewStripeService(cfg)
	pdfService, err := services.NewPDFService("/var/lib/upcycleconnect/invoices")
	if err != nil {
		log.Fatalf("Failed to initialize PDF service: %v", err)
	}
	paymentHandler := &handlers.PaymentHandler{
		DB:     db,
		Audit:  audit,
		Stripe: stripeService,
		PDF:    pdfService,
		Mailer: mailer,
	}
	invoiceHandler := &handlers.InvoiceHandler{DB: db}

	r.Static("/uploads", "/uploads")

	r.GET("/up", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	public := r.Group("/api/public/v1")
	{
		public.GET("/siret/:siret", siretHandler.Lookup)

		public.GET("/categories", publicHandler.Categories)
		public.GET("/prestations", publicHandler.Prestations)
		public.GET("/prestations/:id", publicHandler.ShowPrestation)
		public.GET("/events", publicHandler.Events)
		public.GET("/events/:id", publicHandler.ShowEvent)
		public.GET("/providers", publicHandler.Providers)
		public.GET("/providers/:id", publicHandler.ShowProvider)
		public.GET("/collection-points", collectionPointHandler.PublicIndex)

		public.POST("/payments/webhook", paymentHandler.Webhook)
	}

	userAPI := r.Group("/api/v1")
	{
		userAPI.POST("/auth/register", userAuthHandler.Register)
		userAPI.POST("/auth/login", userAuthHandler.Login)
		userAPI.GET("/auth/verify-email", userAuthHandler.VerifyEmail)
		userAPI.GET("/auth/verify-login", userAuthHandler.VerifyLogin)

		userProtected := userAPI.Group("")
		userProtected.Use(middleware.AuthMiddleware(db, cfg))
		{
			userProtected.POST("/auth/logout", userAuthHandler.Logout)
			userProtected.GET("/auth/me", userAuthHandler.Me)

			userProtected.GET("/profile", profileHandler.Stats)
			userProtected.PUT("/profile/info", profileHandler.UpdateInfo)
			userProtected.POST("/profile/avatar", profileHandler.UploadAvatar)
			userProtected.POST("/profile/change-password", profileHandler.ChangePassword)
			userProtected.POST("/profile/email/start", profileHandler.EmailChangeStart)
			userProtected.POST("/profile/email/verify-current", profileHandler.EmailVerifyCurrent)
			userProtected.POST("/profile/email/verify-new", profileHandler.EmailVerifyNew)

			userProtected.POST("/events/:id/register", profileHandler.RegisterForEvent)
			userProtected.DELETE("/events/:id/register", profileHandler.UnregisterFromEvent)

			userProtected.GET("/deposits", userDepositHandler.Index)
			userProtected.POST("/deposits", userDepositHandler.Store)
			userProtected.GET("/deposits/:id", userDepositHandler.Show)
			userProtected.DELETE("/deposits/:id", userDepositHandler.Destroy)

			userProtected.GET("/upcycling-score", userDepositHandler.Score)

			userProtected.POST("/prestations/:id/reserve", paymentHandler.Reserve)
			userProtected.GET("/payments/session", paymentHandler.SessionStatus)
			userProtected.GET("/reservations/:id", paymentHandler.ShowReservation)

			userProtected.GET("/invoices", invoiceHandler.Index)
			userProtected.GET("/invoices/:id", invoiceHandler.Show)
			userProtected.GET("/invoices/:id/download", invoiceHandler.Download)

			userProtected.POST("/provider/apply", userProviderHandler.Apply)
			userProtected.GET("/provider/profile", userProviderHandler.GetProfile)
			userProtected.PUT("/provider/profile", userProviderHandler.UpdateProfile)
			userProtected.GET("/provider/prestations", userProviderHandler.ListPrestations)
			userProtected.POST("/provider/prestations", userProviderHandler.CreatePrestation)
			userProtected.PUT("/provider/prestations/:id", userProviderHandler.UpdatePrestation)
			userProtected.PUT("/provider/prestations/:id/submit", userProviderHandler.SubmitPrestation)
			userProtected.DELETE("/provider/prestations/:id", userProviderHandler.DestroyPrestation)
		}
	}

	api := r.Group("/api/admin/v1")
	{
		api.POST("/auth/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(db, cfg))
		protected.Use(middleware.IsAdmin())
		{
			protected.POST("/auth/logout", authHandler.Logout)
			protected.GET("/auth/me", authHandler.Me)

			protected.GET("/docs/spec", swaggerHandler.Spec)

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

			protected.GET("/deposits", depositHandler.Index)
			protected.GET("/deposits/:id", depositHandler.Show)
			protected.PUT("/deposits/:id/status", depositHandler.UpdateStatus)

			protected.GET("/collection-points", collectionPointHandler.Index)
			protected.POST("/collection-points", collectionPointHandler.Store)
			protected.PUT("/collection-points/:id", collectionPointHandler.Update)
			protected.DELETE("/collection-points/:id", collectionPointHandler.Destroy)

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
