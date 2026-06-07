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
	mailer := services.NewMailer(cfg)
	notifications := services.NewNotificationService(db, cfg)

	authHandler := &handlers.AuthHandler{DB: db, Cfg: cfg}
	dashHandler := &handlers.DashboardHandler{DB: db}
	swaggerHandler := &handlers.SwaggerHandler{}
	notificationHandler := &handlers.NotificationHandler{DB: db}
	userHandler := &handlers.UserHandler{DB: db, Audit: audit, Mailer: mailer, Cfg: cfg, Notifications: notifications}
	providerHandler := &handlers.ProviderHandler{DB: db, Audit: audit, Notifications: notifications}
	categoryHandler := &handlers.CategoryHandler{DB: db, Audit: audit}
	prestationHandler := &handlers.PrestationHandler{DB: db, Audit: audit}
	eventHandler := &handlers.EventHandler{DB: db, Audit: audit}
	adminHandler := &handlers.AdminHandler{DB: db, Audit: audit}
	logsHandler := &handlers.LogsHandler{DB: db}
	siretHandler    := &handlers.SiretHandler{Cfg: cfg}
	depositHandler         := &handlers.DepositHandler{DB: db, Audit: audit}
	collectionPointHandler := &handlers.CollectionPointHandler{DB: db, Audit: audit}
	userAuthHandler := &handlers.UserAuthHandler{DB: db, Cfg: cfg, Mailer: mailer}
	profileHandler  := &handlers.ProfileHandler{DB: db, Mailer: mailer}
	publicHandler   := &handlers.PublicHandler{DB: db}
	userDepositHandler := &handlers.UserDepositHandler{DB: db, Audit: audit}
	userProviderHandler := &handlers.UserProviderHandler{DB: db, Audit: audit, Notifications: notifications}
	calendarHandler := &handlers.CalendarHandler{DB: db}
	providerEventHandler := &handlers.ProviderEventHandler{DB: db, Audit: audit}
	eventMessageHandler := &handlers.EventMessageHandler{DB: db, Audit: audit}
	forumHandler := &handlers.ForumHandler{DB: db, Audit: audit}
	searchHandler := &handlers.SearchHandler{DB: db}
	oauthHandler := &handlers.OAuthHandler{DB: db, Cfg: cfg}

	passkeyHandler, err := handlers.NewPasskeyHandler(db, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize WebAuthn: %v", err)
	}

	stripeService := services.NewStripeService(cfg)
	pdfService, err := services.NewPDFService("/var/lib/upcycleconnect/invoices")
	if err != nil {
		log.Fatalf("Failed to initialize PDF service: %v", err)
	}
	paymentHandler := &handlers.PaymentHandler{
		DB:            db,
		Audit:         audit,
		Stripe:        stripeService,
		PDF:           pdfService,
		Mailer:        mailer,
		Notifications: notifications,
	}
	invoiceHandler := &handlers.InvoiceHandler{DB: db}
	contractHandler := &handlers.ContractHandler{DB: db, PDF: pdfService, Audit: audit}
	tipHandler := &handlers.TipHandler{DB: db, Audit: audit}

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

		public.GET("/calendar.ics", calendarHandler.Feed)
		public.GET("/search", searchHandler.Search)

		public.GET("/forum/categories", forumHandler.ListCategories)
		public.GET("/forum/categories/:slug", forumHandler.ShowCategory)
		public.GET("/forum/threads/:id", forumHandler.ShowThread)

		public.GET("/tips", tipHandler.PublicIndex)
		public.GET("/tips/:slug", tipHandler.PublicShow)
	}

	userAPI := r.Group("/api/v1")
	{
		userAPI.POST("/auth/register", userAuthHandler.Register)
		userAPI.POST("/auth/login", userAuthHandler.Login)
		userAPI.POST("/auth/google", oauthHandler.GoogleAuth)
		userAPI.POST("/auth/apple", oauthHandler.AppleAuth)
		userAPI.POST("/auth/oauth/link", oauthHandler.OAuthLink)
		userAPI.GET("/auth/verify-email", userAuthHandler.VerifyEmail)
		userAPI.POST("/auth/reset-password", userAuthHandler.ResetPassword)

		userAPI.POST("/passkeys/authenticate/begin", passkeyHandler.AuthBegin)
		userAPI.POST("/passkeys/authenticate/complete", passkeyHandler.AuthComplete)
		userAPI.GET("/auth/verify-login", userAuthHandler.VerifyLogin)

		userProtected := userAPI.Group("")
		userProtected.Use(middleware.AuthMiddleware(db, cfg))
		{
			userProtected.POST("/auth/logout", userAuthHandler.Logout)
			userProtected.GET("/auth/me", userAuthHandler.Me)
			userProtected.POST("/onboarding/complete", userAuthHandler.CompleteOnboarding)

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
			userProtected.GET("/prestations/:id/contract-preview", contractHandler.Preview)
			userProtected.GET("/payments/session", paymentHandler.SessionStatus)
			userProtected.GET("/reservations/:id", paymentHandler.ShowReservation)
			userProtected.GET("/reservations/:id/contract", contractHandler.ByReservation)
			userProtected.GET("/reservations/:id/quote-contract-preview", contractHandler.QuotePreview)
			userProtected.POST("/reservations/:id/accept-quote", paymentHandler.AcceptQuote)

			userProtected.GET("/invoices", invoiceHandler.Index)
			userProtected.GET("/invoices/:id", invoiceHandler.Show)
			userProtected.GET("/invoices/:id/download", invoiceHandler.Download)
			userProtected.GET("/contracts/:id/download", contractHandler.Download)
			userProtected.GET("/provider/contracts", contractHandler.ProviderContracts)

			userProtected.GET("/notifications", notificationHandler.Index)
			userProtected.GET("/notifications/unread-count", notificationHandler.UnreadCount)
			userProtected.POST("/notifications/:id/read", notificationHandler.MarkRead)
			userProtected.POST("/notifications/read-all", notificationHandler.MarkAllRead)
			userProtected.POST("/notifications/push-token", notificationHandler.RegisterPushToken)
			userProtected.DELETE("/notifications/push-token", notificationHandler.UnregisterPushToken)

			userProtected.GET("/calendar/token", calendarHandler.GetToken)
			userProtected.POST("/calendar/token/regenerate", calendarHandler.RegenerateToken)

			userProtected.GET("/events/:id/registration", profileHandler.CheckEventRegistration)
			userProtected.GET("/events/:id/messages", eventMessageHandler.Index)
			userProtected.POST("/events/:id/messages", eventMessageHandler.Store)
			userProtected.POST("/events/:id/messages/image", eventMessageHandler.UploadImage)

			userProtected.POST("/forum/threads", forumHandler.CreateThread)
			userProtected.PUT("/forum/threads/:id", forumHandler.UpdateThread)
			userProtected.DELETE("/forum/threads/:id", forumHandler.DeleteThread)
			userProtected.POST("/forum/threads/:id/close", forumHandler.CloseThread)
			userProtected.POST("/forum/threads/:id/ban/:userId", forumHandler.BanUser)
			userProtected.DELETE("/forum/threads/:id/ban/:userId", forumHandler.UnbanUser)
			userProtected.POST("/forum/threads/:id/replies", forumHandler.CreateReply)
			userProtected.PUT("/forum/replies/:id", forumHandler.UpdateReply)
			userProtected.DELETE("/forum/replies/:id", forumHandler.DeleteReply)
			userProtected.POST("/forum/reports", forumHandler.CreateReport)
			userProtected.POST("/forum/media", forumHandler.UploadMedia)

			userProtected.GET("/passkeys", passkeyHandler.List)
			userProtected.POST("/passkeys/register/begin", passkeyHandler.RegisterBegin)
			userProtected.POST("/passkeys/register/complete", passkeyHandler.RegisterComplete)
			userProtected.PUT("/passkeys/:id", passkeyHandler.Rename)
			userProtected.DELETE("/passkeys/:id", passkeyHandler.Delete)

			userProtected.GET("/provider/events", providerEventHandler.List)
			userProtected.POST("/provider/events", providerEventHandler.Store)
			userProtected.PUT("/provider/events/:id", providerEventHandler.Update)
			userProtected.PUT("/provider/events/:id/submit", providerEventHandler.Submit)
			userProtected.DELETE("/provider/events/:id", providerEventHandler.Destroy)

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

			protected.GET("/notifications", notificationHandler.Index)
			protected.GET("/notifications/unread-count", notificationHandler.UnreadCount)
			protected.POST("/notifications/:id/read", notificationHandler.MarkRead)
			protected.POST("/notifications/read-all", notificationHandler.MarkAllRead)

			protected.GET("/docs/spec", swaggerHandler.Spec)

			protected.GET("/dashboard/stats", dashHandler.Stats)

			protected.GET("/logs", logsHandler.Index)
			protected.GET("/logs/activity", logsHandler.Activity)

			protected.GET("/users", userHandler.Index)
			protected.GET("/users/:id", userHandler.Show)
			protected.PUT("/users/:id", userHandler.Update)
			protected.PUT("/users/:id/email", userHandler.UpdateEmail)
			protected.POST("/users/:id/ban", userHandler.Ban)
			protected.POST("/users/:id/activate", userHandler.Activate)
			protected.POST("/users/:id/send-reset", userHandler.SendPasswordReset)
			protected.GET("/users/:id/export", userHandler.Export)
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

			protected.GET("/tips", tipHandler.AdminIndex)
			protected.POST("/tips", tipHandler.AdminStore)
			protected.GET("/tips/:id", tipHandler.AdminShow)
			protected.PUT("/tips/:id", tipHandler.AdminUpdate)
			protected.DELETE("/tips/:id", tipHandler.AdminDestroy)

			protected.GET("/forum/categories", forumHandler.ListCategories)
			protected.POST("/forum/categories", forumHandler.AdminCreateCategory)
			protected.PUT("/forum/categories/:id", forumHandler.AdminUpdateCategory)
			protected.DELETE("/forum/categories/:id", forumHandler.AdminDeleteCategory)
			protected.GET("/forum/threads", forumHandler.AdminListThreads)
			protected.DELETE("/forum/threads/:id", forumHandler.AdminDeleteThread)
			protected.PUT("/forum/threads/:id/pin", forumHandler.AdminPinThread)
			protected.PUT("/forum/threads/:id/lock", forumHandler.AdminLockThread)
			protected.DELETE("/forum/replies/:id", forumHandler.AdminDeleteReply)
			protected.GET("/forum/reports", forumHandler.AdminListReports)
			protected.PUT("/forum/reports/:id/resolve", forumHandler.AdminResolveReport)

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
