package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SwaggerHandler struct{}

func (h *SwaggerHandler) Spec(c *gin.Context) {
	spec := buildOpenAPISpec()
	c.JSON(http.StatusOK, spec)
}

func buildOpenAPISpec() map[string]any {
	bearerAuth := map[string]any{
		"type":         "http",
		"scheme":       "bearer",
		"bearerFormat": "JWT",
		"description":  "JWT obtenu via POST /api/admin/v1/auth/login. Format : `Bearer <token>`",
	}

	userSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id":                  map[string]any{"type": "integer", "example": 1},
			"email":               map[string]any{"type": "string", "format": "email", "example": "user@example.com"},
			"first_name":          map[string]any{"type": "string", "example": "Jean"},
			"last_name":           map[string]any{"type": "string", "example": "Dupont"},
			"phone":               map[string]any{"type": "string", "example": "+33612345678", "nullable": true},
			"avatar_url":          map[string]any{"type": "string", "nullable": true},
			"status":              map[string]any{"type": "string", "enum": []string{"active", "pending", "banned"}, "example": "active"},
			"role":                map[string]any{"type": "string", "example": "admin"},
			"email_verified_at":   map[string]any{"type": "string", "format": "date-time", "nullable": true},
			"created_at":          map[string]any{"type": "string", "format": "date-time"},
			"updated_at":          map[string]any{"type": "string", "format": "date-time"},
		},
	}

	providerProfileSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id":           map[string]any{"type": "integer"},
			"company_name": map[string]any{"type": "string", "example": "EcoRépar SARL"},
			"siret":        map[string]any{"type": "string", "example": "12345678901234", "nullable": true},
			"description":  map[string]any{"type": "string", "nullable": true},
			"website":      map[string]any{"type": "string", "nullable": true},
			"status":       map[string]any{"type": "string", "enum": []string{"pending", "approved", "suspended"}, "example": "pending"},
			"approved_at":  map[string]any{"type": "string", "format": "date-time", "nullable": true},
		},
	}

	categorySchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id":                map[string]any{"type": "integer", "example": 1},
			"name":              map[string]any{"type": "string", "example": "Électronique"},
			"slug":              map[string]any{"type": "string", "example": "electronique"},
			"description":       map[string]any{"type": "string", "nullable": true},
			"is_active":         map[string]any{"type": "boolean", "example": true},
			"sort_order":        map[string]any{"type": "integer", "example": 0},
			"prestations_count": map[string]any{"type": "integer", "example": 5},
			"created_at":        map[string]any{"type": "string", "format": "date-time"},
			"updated_at":        map[string]any{"type": "string", "format": "date-time"},
		},
	}

	prestationSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id":          map[string]any{"type": "integer", "example": 1},
			"title":       map[string]any{"type": "string", "example": "Réparation smartphone"},
			"description": map[string]any{"type": "string", "nullable": true},
			"price":       map[string]any{"type": "number", "format": "float", "example": 49.90, "nullable": true},
			"price_type":  map[string]any{"type": "string", "enum": []string{"fixed", "hourly", "quote"}, "example": "fixed"},
			"status":      map[string]any{"type": "string", "enum": []string{"draft", "published", "suspended", "archived"}, "example": "published"},
			"category":    map[string]any{"$ref": "#/components/schemas/Category"},
			"provider":    map[string]any{"$ref": "#/components/schemas/User"},
			"created_at":  map[string]any{"type": "string", "format": "date-time"},
			"updated_at":  map[string]any{"type": "string", "format": "date-time"},
		},
	}

	eventSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id":                   map[string]any{"type": "integer", "example": 1},
			"title":                map[string]any{"type": "string", "example": "Repair Café Paris"},
			"description":          map[string]any{"type": "string", "nullable": true},
			"location":             map[string]any{"type": "string", "example": "15 rue de la Paix, Paris", "nullable": true},
			"start_date":           map[string]any{"type": "string", "format": "date-time"},
			"end_date":             map[string]any{"type": "string", "format": "date-time"},
			"max_participants":     map[string]any{"type": "integer", "nullable": true, "example": 50},
			"registrations_count":  map[string]any{"type": "integer", "example": 12},
			"status":               map[string]any{"type": "string", "enum": []string{"draft", "published", "cancelled"}, "example": "published"},
			"category":             map[string]any{"$ref": "#/components/schemas/Category"},
			"creator":              map[string]any{"$ref": "#/components/schemas/User"},
			"created_at":           map[string]any{"type": "string", "format": "date-time"},
			"updated_at":           map[string]any{"type": "string", "format": "date-time"},
		},
	}

	depositSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id":               map[string]any{"type": "integer", "example": 1},
			"title":            map[string]any{"type": "string", "example": "Aspirateur Dyson"},
			"description":      map[string]any{"type": "string", "nullable": true},
			"condition":        map[string]any{"type": "string", "enum": []string{"good", "fair", "poor"}, "example": "fair"},
			"history":          map[string]any{"type": "string", "nullable": true},
			"estimated_weight": map[string]any{"type": "number", "format": "float", "nullable": true, "example": 2.5},
			"carbon_savings":   map[string]any{"type": "number", "format": "float", "nullable": true},
			"status":           map[string]any{"type": "string", "enum": []string{"pending", "approved", "rejected", "accepted", "validated"}, "example": "pending"},
			"validation_note":  map[string]any{"type": "string", "nullable": true},
			"qr_code":          map[string]any{"type": "string", "nullable": true, "example": "UP-1-a3f9b2"},
			"validated_at":     map[string]any{"type": "string", "format": "date-time", "nullable": true},
			"user":             map[string]any{"$ref": "#/components/schemas/User"},
			"category":         map[string]any{"$ref": "#/components/schemas/Category"},
			"created_at":       map[string]any{"type": "string", "format": "date-time"},
		},
	}

	auditLogSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"id":            map[string]any{"type": "integer"},
			"action":        map[string]any{"type": "string", "example": "user.banned"},
			"resource_type": map[string]any{"type": "string", "example": "User"},
			"resource_id":   map[string]any{"type": "integer", "nullable": true},
			"admin_name":    map[string]any{"type": "string", "example": "Admin User"},
			"ip_address":    map[string]any{"type": "string", "nullable": true, "example": "192.168.1.1"},
			"old_values":    map[string]any{"type": "object", "nullable": true},
			"new_values":    map[string]any{"type": "object", "nullable": true},
			"created_at":    map[string]any{"type": "string", "format": "date-time"},
		},
	}

	paginationMeta := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"current_page": map[string]any{"type": "integer", "example": 1},
			"last_page":    map[string]any{"type": "integer", "example": 5},
			"per_page":     map[string]any{"type": "integer", "example": 20},
			"total":        map[string]any{"type": "integer", "example": 93},
		},
	}

	errorSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"error": map[string]any{"type": "string", "example": "Unauthorized"},
		},
	}

	bearerSecurityReq := []map[string]any{{"bearerAuth": []string{}}}

	paginationParams := []map[string]any{
		{
			"name":        "page",
			"in":          "query",
			"schema":      map[string]any{"type": "integer", "default": 1},
			"description": "Numéro de page",
		},
		{
			"name":        "per_page",
			"in":          "query",
			"schema":      map[string]any{"type": "integer", "default": 20},
			"description": "Nombre d'éléments par page",
		},
	}

	searchParam := map[string]any{
		"name":        "search",
		"in":          "query",
		"schema":      map[string]any{"type": "string"},
		"description": "Recherche textuelle",
	}

	paths := map[string]any{

		"/up": map[string]any{
			"get": map[string]any{
				"tags":    []string{"Système"},
				"summary": "Health check",
				"responses": map[string]any{
					"200": map[string]any{"description": "Service opérationnel", "content": map[string]any{"text/plain": map[string]any{"schema": map[string]any{"type": "string", "example": "OK"}}}},
				},
			},
		},

		"/api/public/v1/siret/{siret}": map[string]any{
			"get": map[string]any{
				"tags":    []string{"Public"},
				"summary": "Recherche d'entreprise par SIRET",
				"description": "Consulte l'API Société.com pour obtenir les informations légales d'une entreprise à partir de son numéro SIRET (14 chiffres).",
				"parameters": []map[string]any{
					{
						"name": "siret", "in": "path", "required": true,
						"schema": map[string]any{"type": "string", "pattern": "^[0-9]{14}$"},
						"description": "Numéro SIRET (14 chiffres)",
						"example": "12345678901234",
					},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Informations de l'entreprise"},
					"400": map[string]any{"description": "Format SIRET invalide"},
					"404": map[string]any{"description": "Entreprise introuvable"},
					"503": map[string]any{"description": "Service externe indisponible"},
				},
			},
		},

		"/api/v1/auth/register": map[string]any{
			"post": map[string]any{
				"tags":    []string{"Auth Utilisateur"},
				"summary": "Inscription",
				"description": "Crée un nouveau compte utilisateur. Un email de vérification est envoyé.",
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{
								"type": "object",
								"required": []string{"first_name", "last_name", "email", "password"},
								"properties": map[string]any{
									"first_name":   map[string]any{"type": "string", "example": "Marie"},
									"last_name":    map[string]any{"type": "string", "example": "Martin"},
									"email":        map[string]any{"type": "string", "format": "email", "example": "marie@example.com"},
									"password":     map[string]any{"type": "string", "minLength": 8, "example": "motdepasse123"},
									"phone":        map[string]any{"type": "string", "example": "+33612345678"},
									"account_type": map[string]any{"type": "string", "enum": []string{"default", "provider"}, "example": "default"},
									"company_name": map[string]any{"type": "string", "example": "EcoRépar SARL"},
									"siret":        map[string]any{"type": "string", "example": "12345678901234"},
									"activity":     map[string]any{"type": "string", "example": "Réparation électronique"},
									"address":      map[string]any{"type": "string", "example": "15 rue de la Paix, Paris"},
								},
							},
						},
					},
				},
				"responses": map[string]any{
					"201": map[string]any{
						"description": "Compte créé — email de vérification envoyé",
						"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"message": map[string]any{"type": "string", "example": "Account created. Verify your email."}}}}},
					},
					"422": map[string]any{"description": "Validation échouée ou email déjà utilisé"},
				},
			},
		},

		"/api/v1/auth/login": map[string]any{
			"post": map[string]any{
				"tags":    []string{"Auth Utilisateur"},
				"summary": "Connexion utilisateur",
				"description": "Authentifie un utilisateur. Si l'IP est inconnue, un email de vérification est envoyé et la réponse est 202.",
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{
								"type": "object",
								"required": []string{"email", "password"},
								"properties": map[string]any{
									"email":    map[string]any{"type": "string", "format": "email", "example": "user@example.com"},
									"password": map[string]any{"type": "string", "example": "motdepasse123"},
								},
							},
						},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Connexion réussie — token JWT retourné",
						"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"token": map[string]any{"type": "string"}, "user": map[string]any{"$ref": "#/components/schemas/User"}}}}},
					},
					"202": map[string]any{"description": "IP inconnue — email de vérification envoyé", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"pending_verification": map[string]any{"type": "boolean", "example": true}, "message": map[string]any{"type": "string"}}}}}},
					"401": map[string]any{"description": "Identifiants invalides"},
					"403": map[string]any{"description": "Compte banni ou inactif"},
				},
			},
		},

		"/api/v1/auth/verify-email": map[string]any{
			"get": map[string]any{
				"tags":    []string{"Auth Utilisateur"},
				"summary": "Vérification de l'email d'inscription",
				"parameters": []map[string]any{{"name": "token", "in": "query", "required": true, "schema": map[string]any{"type": "string"}, "description": "Token reçu par email"}},
				"responses": map[string]any{
					"200": map[string]any{"description": "Email vérifié"},
					"400": map[string]any{"description": "Token invalide ou expiré"},
				},
			},
		},

		"/api/v1/auth/verify-login": map[string]any{
			"get": map[string]any{
				"tags":    []string{"Auth Utilisateur"},
				"summary": "Confirmation de connexion depuis une nouvelle IP",
				"parameters": []map[string]any{{"name": "token", "in": "query", "required": true, "schema": map[string]any{"type": "string"}, "description": "Token reçu par email"}},
				"responses": map[string]any{
					"200": map[string]any{"description": "IP validée — token JWT retourné", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"token": map[string]any{"type": "string"}, "user": map[string]any{"$ref": "#/components/schemas/User"}}}}}},
					"400": map[string]any{"description": "Token invalide ou expiré"},
				},
			},
		},

		"/api/v1/auth/logout": map[string]any{
			"post": map[string]any{
				"tags":     []string{"Auth Utilisateur"},
				"summary":  "Déconnexion",
				"security": bearerSecurityReq,
				"responses": map[string]any{
					"200": map[string]any{"description": "Déconnecté"},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/v1/auth/me": map[string]any{
			"get": map[string]any{
				"tags":     []string{"Auth Utilisateur"},
				"summary":  "Profil de l'utilisateur connecté",
				"security": bearerSecurityReq,
				"responses": map[string]any{
					"200": map[string]any{"description": "Données de l'utilisateur", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/User"}}}},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/v1/profile": map[string]any{
			"get": map[string]any{
				"tags":        []string{"Profil Utilisateur"},
				"summary":     "Statistiques et données du profil",
				"description": "Retourne le score, CO2 économisé, dépôts, réservations d'événements et badges.",
				"security":    bearerSecurityReq,
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Données de profil complètes",
						"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"score":          map[string]any{"type": "integer", "example": 72},
								"co2_saved":      map[string]any{"type": "number", "format": "float", "example": 12.5},
								"objects_saved":  map[string]any{"type": "integer", "example": 4},
								"deposits":       map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/Deposit"}},
								"reservations":   map[string]any{"type": "array", "items": map[string]any{"type": "object", "properties": map[string]any{"id": map[string]any{"type": "integer"}, "event_id": map[string]any{"type": "integer"}, "title": map[string]any{"type": "string"}, "start_date": map[string]any{"type": "string", "format": "date-time"}, "location": map[string]any{"type": "string", "nullable": true}, "past": map[string]any{"type": "boolean"}}}},
								"badges":         map[string]any{"type": "array", "items": map[string]any{"type": "object", "properties": map[string]any{"key": map[string]any{"type": "string"}, "label": map[string]any{"type": "string"}, "desc": map[string]any{"type": "string"}, "earned": map[string]any{"type": "boolean"}}}},
							},
						}}},
					},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/v1/profile/info": map[string]any{
			"put": map[string]any{
				"tags":     []string{"Profil Utilisateur"},
				"summary":  "Mettre à jour les informations personnelles",
				"security": bearerSecurityReq,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
						"type": "object",
						"required": []string{"first_name", "last_name"},
						"properties": map[string]any{
							"first_name": map[string]any{"type": "string", "example": "Marie"},
							"last_name":  map[string]any{"type": "string", "example": "Martin"},
							"phone":      map[string]any{"type": "string", "example": "+33612345678"},
						},
					}}},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Informations mises à jour", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/User"}}}},
					"401": map[string]any{"description": "Non authentifié"},
					"422": map[string]any{"description": "Validation échouée"},
				},
			},
		},

		"/api/v1/profile/avatar": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Profil Utilisateur"},
				"summary":     "Upload d'avatar",
				"description": "Formats acceptés : JPG, PNG, WebP. Taille max : 2 MB.",
				"security":    bearerSecurityReq,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{"multipart/form-data": map[string]any{"schema": map[string]any{"type": "object", "required": []string{"avatar"}, "properties": map[string]any{"avatar": map[string]any{"type": "string", "format": "binary"}}}}},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Avatar mis à jour", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"avatar_url": map[string]any{"type": "string"}}}}}},
					"400": map[string]any{"description": "Fichier invalide ou trop lourd"},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/v1/profile/change-password": map[string]any{
			"post": map[string]any{
				"tags":     []string{"Profil Utilisateur"},
				"summary":  "Changer le mot de passe",
				"security": bearerSecurityReq,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
						"type": "object",
						"required": []string{"current", "new"},
						"properties": map[string]any{
							"current": map[string]any{"type": "string", "example": "ancien_mdp"},
							"new":     map[string]any{"type": "string", "minLength": 8, "example": "nouveau_mdp_123"},
						},
					}}},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Mot de passe modifié"},
					"400": map[string]any{"description": "Mot de passe actuel incorrect"},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/v1/profile/email/start": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Profil Utilisateur"},
				"summary":     "Démarrer le changement d'email (étape 1/3)",
				"description": "Envoie un code de vérification à l'email actuel.",
				"security":    bearerSecurityReq,
				"responses": map[string]any{
					"200": map[string]any{"description": "Code envoyé à l'email actuel"},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/v1/profile/email/verify-current": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Profil Utilisateur"},
				"summary":     "Vérifier l'email actuel (étape 2/3)",
				"description": "Valide le code reçu sur l'email actuel et envoie un code au nouvel email.",
				"security":    bearerSecurityReq,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
						"type": "object",
						"required": []string{"code", "new_email"},
						"properties": map[string]any{
							"code":      map[string]any{"type": "string", "pattern": "^[0-9]{6}$", "example": "482910"},
							"new_email": map[string]any{"type": "string", "format": "email", "example": "newemail@example.com"},
						},
					}}},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Code envoyé au nouvel email"},
					"400": map[string]any{"description": "Code invalide ou expiré"},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/v1/profile/email/verify-new": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Profil Utilisateur"},
				"summary":     "Vérifier le nouvel email (étape 3/3)",
				"description": "Valide le code reçu sur le nouvel email et finalise le changement.",
				"security":    bearerSecurityReq,
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
						"type": "object",
						"required": []string{"code"},
						"properties": map[string]any{
							"code": map[string]any{"type": "string", "pattern": "^[0-9]{6}$", "example": "193847"},
						},
					}}},
				},
				"responses": map[string]any{
					"200": map[string]any{"description": "Email modifié avec succès"},
					"400": map[string]any{"description": "Code invalide ou expiré"},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/v1/events/{id}/register": map[string]any{
			"post": map[string]any{
				"tags":     []string{"Profil Utilisateur"},
				"summary":  "S'inscrire à un événement",
				"security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}, "description": "ID de l'événement"}},
				"responses": map[string]any{
					"200": map[string]any{"description": "Inscription enregistrée"},
					"400": map[string]any{"description": "Déjà inscrit ou événement complet"},
					"404": map[string]any{"description": "Événement introuvable"},
				},
			},
			"delete": map[string]any{
				"tags":     []string{"Profil Utilisateur"},
				"summary":  "Se désinscrire d'un événement",
				"security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}, "description": "ID de l'événement"}},
				"responses": map[string]any{
					"200": map[string]any{"description": "Désinscription effectuée"},
					"404": map[string]any{"description": "Inscription introuvable"},
				},
			},
		},

		"/api/admin/v1/auth/login": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Auth Admin"},
				"summary":     "Connexion administrateur",
				"description": "Authentifie un administrateur (rôle `admin` ou `super_admin` requis). Retourne un JWT valide 480 minutes.",
				"requestBody": map[string]any{
					"required": true,
					"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
						"type": "object",
						"required": []string{"email", "password"},
						"properties": map[string]any{
							"email":    map[string]any{"type": "string", "format": "email", "example": "admin@upcycleconnect.fr"},
							"password": map[string]any{"type": "string", "example": "admin_password"},
						},
					}}},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Token JWT admin",
						"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"token": map[string]any{"type": "string", "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}, "user": map[string]any{"$ref": "#/components/schemas/User"}}}}},
					},
					"401": map[string]any{"description": "Identifiants invalides"},
					"403": map[string]any{"description": "Compte non administrateur"},
				},
			},
		},

		"/api/admin/v1/auth/logout": map[string]any{
			"post": map[string]any{
				"tags": []string{"Auth Admin"}, "summary": "Déconnexion admin", "security": bearerSecurityReq,
				"responses": map[string]any{"200": map[string]any{"description": "Déconnecté"}, "401": map[string]any{"description": "Non authentifié"}},
			},
		},

		"/api/admin/v1/auth/me": map[string]any{
			"get": map[string]any{
				"tags": []string{"Auth Admin"}, "summary": "Profil admin connecté", "security": bearerSecurityReq,
				"responses": map[string]any{
					"200": map[string]any{"description": "Données de l'admin", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/User"}}}},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/admin/v1/dashboard/stats": map[string]any{
			"get": map[string]any{
				"tags":        []string{"Dashboard"},
				"summary":     "Statistiques générales",
				"description": "Retourne les métriques clés de la plateforme.",
				"security":    bearerSecurityReq,
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Statistiques",
						"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{
							"users_count":           map[string]any{"type": "integer", "example": 342},
							"providers_count":        map[string]any{"type": "integer", "example": 28},
							"providers_pending":      map[string]any{"type": "integer", "example": 5},
							"prestations_count":      map[string]any{"type": "integer", "example": 87},
							"events_count":           map[string]any{"type": "integer", "example": 14},
							"prestations_by_status":  map[string]any{"type": "object", "example": map[string]any{"draft": 12, "published": 70, "suspended": 5}},
							"deposits_pending":       map[string]any{"type": "integer", "example": 9},
							"deposits_accepted":      map[string]any{"type": "integer", "example": 63},
							"co2_total":              map[string]any{"type": "number", "format": "float", "example": 1240.5},
						}}}},
					},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/admin/v1/logs": map[string]any{
			"get": map[string]any{
				"tags": []string{"Journaux"}, "summary": "Journaux d'audit (50 derniers)", "security": bearerSecurityReq,
				"responses": map[string]any{
					"200": map[string]any{"description": "Liste des logs", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/AuditLog"}}}}}}},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/admin/v1/logs/activity": map[string]any{
			"get": map[string]any{
				"tags": []string{"Journaux"}, "summary": "Activité récente utilisateurs & prestataires (50 derniers)", "security": bearerSecurityReq,
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Activité récente",
						"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"type": "object", "properties": map[string]any{
							"id":         map[string]any{"type": "integer"},
							"type":       map[string]any{"type": "string", "enum": []string{"user", "provider"}},
							"name":       map[string]any{"type": "string"},
							"email":      map[string]any{"type": "string"},
							"company":    map[string]any{"type": "string", "nullable": true},
							"status":     map[string]any{"type": "string", "nullable": true},
							"created_at": map[string]any{"type": "string", "format": "date-time"},
						}}}}}}},
					},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/admin/v1/users": map[string]any{
			"get": map[string]any{
				"tags":        []string{"Utilisateurs"},
				"summary":     "Liste des utilisateurs",
				"description": "Retourne les utilisateurs réguliers (non admins) avec pagination, recherche et filtre de statut.",
				"security":    bearerSecurityReq,
				"parameters": append(paginationParams, searchParam, map[string]any{
					"name": "status", "in": "query",
					"schema": map[string]any{"type": "string", "enum": []string{"active", "pending", "banned"}},
					"description": "Filtrer par statut",
				}),
				"responses": map[string]any{
					"200": map[string]any{"description": "Liste paginée", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/User"}}, "meta": map[string]any{"$ref": "#/components/schemas/PaginationMeta"}}}}}},
					"401": map[string]any{"description": "Non authentifié"},
				},
			},
		},

		"/api/admin/v1/users/{id}": map[string]any{
			"get": map[string]any{
				"tags": []string{"Utilisateurs"}, "summary": "Détail d'un utilisateur", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{
					"200": map[string]any{"description": "Détail utilisateur", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/User"}}}},
					"404": map[string]any{"description": "Introuvable"},
				},
			},
			"put": map[string]any{
				"tags": []string{"Utilisateurs"}, "summary": "Modifier un utilisateur", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"first_name": map[string]any{"type": "string"}, "last_name": map[string]any{"type": "string"}, "phone": map[string]any{"type": "string"}}}}}},
				"responses": map[string]any{"200": map[string]any{"description": "Utilisateur mis à jour"}, "404": map[string]any{"description": "Introuvable"}},
			},
			"delete": map[string]any{
				"tags": []string{"Utilisateurs"}, "summary": "Supprimer un utilisateur (soft delete)", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Supprimé"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/users/{id}/ban": map[string]any{
			"post": map[string]any{
				"tags": []string{"Utilisateurs"}, "summary": "Bannir un utilisateur", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Banni — action auditée"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/users/{id}/activate": map[string]any{
			"post": map[string]any{
				"tags": []string{"Utilisateurs"}, "summary": "Activer un utilisateur", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Activé — action auditée"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/providers": map[string]any{
			"get": map[string]any{
				"tags": []string{"Prestataires"}, "summary": "Liste des prestataires", "security": bearerSecurityReq,
				"parameters": append(paginationParams, searchParam, map[string]any{
					"name": "status", "in": "query",
					"schema": map[string]any{"type": "string", "enum": []string{"pending", "approved", "suspended"}},
					"description": "Filtrer par statut de validation",
				}),
				"responses": map[string]any{
					"200": map[string]any{"description": "Liste paginée des prestataires", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/ProviderProfile"}}, "meta": map[string]any{"$ref": "#/components/schemas/PaginationMeta"}}}}}},
				},
			},
		},

		"/api/admin/v1/providers/{id}": map[string]any{
			"get": map[string]any{
				"tags": []string{"Prestataires"}, "summary": "Détail d'un prestataire", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Prestataire + profil", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/ProviderProfile"}}}}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/providers/{id}/status": map[string]any{
			"put": map[string]any{
				"tags": []string{"Prestataires"}, "summary": "Modifier le statut d'un prestataire", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "required": []string{"status"}, "properties": map[string]any{"status": map[string]any{"type": "string", "enum": []string{"pending", "approved", "suspended"}, "example": "approved"}}}}}},
				"responses": map[string]any{"200": map[string]any{"description": "Statut mis à jour — action auditée"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/categories": map[string]any{
			"get": map[string]any{
				"tags": []string{"Catégories"}, "summary": "Liste de toutes les catégories", "security": bearerSecurityReq,
				"responses": map[string]any{"200": map[string]any{"description": "Catégories avec compteur de prestations", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/Category"}}}}}}}},
			},
			"post": map[string]any{
				"tags": []string{"Catégories"}, "summary": "Créer une catégorie", "security": bearerSecurityReq,
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
					"type": "object", "required": []string{"name"},
					"properties": map[string]any{
						"name":        map[string]any{"type": "string", "example": "Électronique"},
						"description": map[string]any{"type": "string", "example": "Appareils électroniques et accessoires"},
						"sort_order":  map[string]any{"type": "integer", "example": 1},
						"is_active":   map[string]any{"type": "boolean", "example": true},
					},
				}}}},
				"responses": map[string]any{"201": map[string]any{"description": "Catégorie créée", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/Category"}}}}, "422": map[string]any{"description": "Validation échouée"}},
			},
		},

		"/api/admin/v1/categories/{id}": map[string]any{
			"get": map[string]any{
				"tags": []string{"Catégories"}, "summary": "Détail d'une catégorie", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Catégorie", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/Category"}}}}, "404": map[string]any{"description": "Introuvable"}},
			},
			"put": map[string]any{
				"tags": []string{"Catégories"}, "summary": "Modifier une catégorie", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string"}, "description": map[string]any{"type": "string"}, "sort_order": map[string]any{"type": "integer"}, "is_active": map[string]any{"type": "boolean"}}}}}},
				"responses": map[string]any{"200": map[string]any{"description": "Catégorie mise à jour"}, "404": map[string]any{"description": "Introuvable"}},
			},
			"delete": map[string]any{
				"tags": []string{"Catégories"}, "summary": "Supprimer une catégorie (soft delete)", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Supprimée"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/categories/{id}/toggle": map[string]any{
			"post": map[string]any{
				"tags": []string{"Catégories"}, "summary": "Activer / désactiver une catégorie", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Statut basculé"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/prestations": map[string]any{
			"get": map[string]any{
				"tags": []string{"Prestations"}, "summary": "Liste des prestations", "security": bearerSecurityReq,
				"parameters": append(paginationParams, searchParam,
					map[string]any{"name": "status", "in": "query", "schema": map[string]any{"type": "string", "enum": []string{"draft", "published", "suspended", "archived"}}},
					map[string]any{"name": "category_id", "in": "query", "schema": map[string]any{"type": "integer"}},
					map[string]any{"name": "provider_id", "in": "query", "schema": map[string]any{"type": "integer"}},
				),
				"responses": map[string]any{"200": map[string]any{"description": "Liste paginée", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/Prestation"}}, "meta": map[string]any{"$ref": "#/components/schemas/PaginationMeta"}}}}}}},
			},
			"post": map[string]any{
				"tags": []string{"Prestations"}, "summary": "Créer une prestation", "security": bearerSecurityReq,
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
					"type": "object", "required": []string{"title"},
					"properties": map[string]any{
						"category_id":  map[string]any{"type": "integer"},
						"provider_id":  map[string]any{"type": "integer"},
						"title":        map[string]any{"type": "string", "example": "Réparation smartphone"},
						"description":  map[string]any{"type": "string"},
						"price":        map[string]any{"type": "number", "format": "float", "example": 49.90},
						"price_type":   map[string]any{"type": "string", "enum": []string{"fixed", "hourly", "quote"}, "default": "fixed"},
						"status":       map[string]any{"type": "string", "enum": []string{"draft", "published", "suspended", "archived"}, "default": "draft"},
					},
				}}}},
				"responses": map[string]any{"201": map[string]any{"description": "Prestation créée", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/Prestation"}}}}, "422": map[string]any{"description": "Validation échouée"}},
			},
		},

		"/api/admin/v1/prestations/{id}": map[string]any{
			"get": map[string]any{
				"tags": []string{"Prestations"}, "summary": "Détail d'une prestation", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Prestation", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/Prestation"}}}}, "404": map[string]any{"description": "Introuvable"}},
			},
			"put": map[string]any{
				"tags": []string{"Prestations"}, "summary": "Modifier une prestation", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/Prestation"}}}},
				"responses": map[string]any{"200": map[string]any{"description": "Prestation mise à jour"}, "404": map[string]any{"description": "Introuvable"}},
			},
			"delete": map[string]any{
				"tags": []string{"Prestations"}, "summary": "Supprimer une prestation (soft delete)", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Supprimée"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/prestations/{id}/status": map[string]any{
			"put": map[string]any{
				"tags": []string{"Prestations"}, "summary": "Changer le statut d'une prestation", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "required": []string{"status"}, "properties": map[string]any{"status": map[string]any{"type": "string", "enum": []string{"draft", "published", "suspended", "archived"}}}}}}},
				"responses": map[string]any{"200": map[string]any{"description": "Statut mis à jour — action auditée"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/events": map[string]any{
			"get": map[string]any{
				"tags": []string{"Événements"}, "summary": "Liste des événements", "security": bearerSecurityReq,
				"parameters": append(paginationParams,
					map[string]any{"name": "status", "in": "query", "schema": map[string]any{"type": "string", "enum": []string{"draft", "published", "cancelled"}}},
					map[string]any{"name": "start_date_from", "in": "query", "schema": map[string]any{"type": "string", "format": "date"}, "example": "2025-01-01"},
					map[string]any{"name": "start_date_to", "in": "query", "schema": map[string]any{"type": "string", "format": "date"}, "example": "2025-12-31"},
				),
				"responses": map[string]any{"200": map[string]any{"description": "Liste paginée avec compteur d'inscriptions", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/Event"}}, "meta": map[string]any{"$ref": "#/components/schemas/PaginationMeta"}}}}}}},
			},
			"post": map[string]any{
				"tags": []string{"Événements"}, "summary": "Créer un événement", "security": bearerSecurityReq,
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
					"type": "object", "required": []string{"title", "start_date", "end_date"},
					"properties": map[string]any{
						"category_id":      map[string]any{"type": "integer"},
						"title":            map[string]any{"type": "string", "example": "Repair Café Paris"},
						"description":      map[string]any{"type": "string"},
						"location":         map[string]any{"type": "string", "example": "15 rue de la Paix, Paris"},
						"start_date":       map[string]any{"type": "string", "example": "2025-06-15T14:00:00", "description": "ISO 8601 ou YYYY-MM-DD HH:mm"},
						"end_date":         map[string]any{"type": "string", "example": "2025-06-15T18:00:00"},
						"max_participants": map[string]any{"type": "integer", "example": 50},
						"status":           map[string]any{"type": "string", "enum": []string{"draft", "published", "cancelled"}, "default": "draft"},
					},
				}}}},
				"responses": map[string]any{"201": map[string]any{"description": "Événement créé"}, "422": map[string]any{"description": "Validation échouée"}},
			},
		},

		"/api/admin/v1/events/{id}": map[string]any{
			"get": map[string]any{
				"tags": []string{"Événements"}, "summary": "Détail d'un événement", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Événement", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/Event"}}}}, "404": map[string]any{"description": "Introuvable"}},
			},
			"put": map[string]any{
				"tags": []string{"Événements"}, "summary": "Modifier un événement", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/Event"}}}},
				"responses": map[string]any{"200": map[string]any{"description": "Événement mis à jour"}, "404": map[string]any{"description": "Introuvable"}},
			},
			"delete": map[string]any{
				"tags": []string{"Événements"}, "summary": "Supprimer un événement (soft delete)", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Supprimé"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/events/{id}/status": map[string]any{
			"put": map[string]any{
				"tags": []string{"Événements"}, "summary": "Changer le statut d'un événement", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "required": []string{"status"}, "properties": map[string]any{"status": map[string]any{"type": "string", "enum": []string{"draft", "published", "cancelled"}}}}}}},
				"responses": map[string]any{"200": map[string]any{"description": "Statut mis à jour — action auditée"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/deposits": map[string]any{
			"get": map[string]any{
				"tags": []string{"Dépôts"}, "summary": "Liste des demandes de dépôt", "security": bearerSecurityReq,
				"parameters": append(paginationParams, searchParam, map[string]any{
					"name": "status", "in": "query",
					"schema": map[string]any{"type": "string", "enum": []string{"pending", "approved", "rejected", "accepted", "validated"}},
				}),
				"responses": map[string]any{"200": map[string]any{"description": "Liste paginée", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/Deposit"}}, "meta": map[string]any{"$ref": "#/components/schemas/PaginationMeta"}}}}}}},
			},
		},

		"/api/admin/v1/deposits/{id}": map[string]any{
			"get": map[string]any{
				"tags": []string{"Dépôts"}, "summary": "Détail d'un dépôt", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Dépôt", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/Deposit"}}}}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/deposits/{id}/status": map[string]any{
			"put": map[string]any{
				"tags":        []string{"Dépôts"},
				"summary":     "Approuver ou rejeter un dépôt",
				"description": "Si approuvé, un QR code est généré au format `UP-{id}-{hex}`. Action auditée.",
				"security":    bearerSecurityReq,
				"parameters":  []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
					"type": "object", "required": []string{"status"},
					"properties": map[string]any{
						"status":          map[string]any{"type": "string", "enum": []string{"approved", "rejected"}, "example": "approved"},
						"validation_note": map[string]any{"type": "string", "example": "Objet en bon état, accepté."},
					},
				}}}},
				"responses": map[string]any{"200": map[string]any{"description": "Statut mis à jour — QR code généré si approuvé"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},

		"/api/admin/v1/admins": map[string]any{
			"get": map[string]any{
				"tags":        []string{"Administrateurs (super_admin)"},
				"summary":     "Liste des administrateurs",
				"description": "Route réservée aux `super_admin`.",
				"security":    bearerSecurityReq,
				"responses":   map[string]any{"200": map[string]any{"description": "Liste des admins", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"data": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/User"}}}}}}}, "403": map[string]any{"description": "Accès refusé — super_admin requis"}},
			},
			"post": map[string]any{
				"tags": []string{"Administrateurs (super_admin)"}, "summary": "Créer un administrateur", "security": bearerSecurityReq,
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
					"type": "object", "required": []string{"email", "password", "first_name", "last_name", "role"},
					"properties": map[string]any{
						"email":      map[string]any{"type": "string", "format": "email", "example": "newadmin@upcycleconnect.fr"},
						"password":   map[string]any{"type": "string", "minLength": 8, "example": "secure_password"},
						"first_name": map[string]any{"type": "string", "example": "Sophie"},
						"last_name":  map[string]any{"type": "string", "example": "Laurent"},
						"role":       map[string]any{"type": "string", "enum": []string{"admin", "super_admin"}, "example": "admin"},
					},
				}}}},
				"responses": map[string]any{"201": map[string]any{"description": "Admin créé"}, "422": map[string]any{"description": "Validation échouée"}},
			},
		},

		"/api/admin/v1/admins/{id}": map[string]any{
			"get": map[string]any{
				"tags": []string{"Administrateurs (super_admin)"}, "summary": "Détail d'un admin", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses": map[string]any{"200": map[string]any{"description": "Admin", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/User"}}}}, "404": map[string]any{"description": "Introuvable"}},
			},
			"put": map[string]any{
				"tags": []string{"Administrateurs (super_admin)"}, "summary": "Modifier un administrateur", "security": bearerSecurityReq,
				"parameters": []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"requestBody": map[string]any{"required": true, "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"type": "object", "properties": map[string]any{"first_name": map[string]any{"type": "string"}, "last_name": map[string]any{"type": "string"}, "password": map[string]any{"type": "string", "minLength": 8}, "role": map[string]any{"type": "string", "enum": []string{"admin", "super_admin"}}}}}}},
				"responses": map[string]any{"200": map[string]any{"description": "Admin mis à jour"}, "404": map[string]any{"description": "Introuvable"}},
			},
			"delete": map[string]any{
				"tags":        []string{"Administrateurs (super_admin)"},
				"summary":     "Supprimer un administrateur",
				"description": "Soft delete. Impossible de supprimer son propre compte.",
				"security":    bearerSecurityReq,
				"parameters":  []map[string]any{{"name": "id", "in": "path", "required": true, "schema": map[string]any{"type": "integer"}}},
				"responses":   map[string]any{"200": map[string]any{"description": "Supprimé"}, "400": map[string]any{"description": "Impossible de supprimer son propre compte"}, "404": map[string]any{"description": "Introuvable"}},
			},
		},
	}

	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "UpcycleConnect API",
			"description": "Documentation complète de l'API UpcycleConnect. Authentifiez-vous via **POST /api/admin/v1/auth/login** puis cliquez sur **Authorize** et saisissez le token JWT obtenu.",
			"version":     "1.0.0",
			"contact": map[string]any{
				"name": "UpcycleConnect",
			},
		},
		"servers": []map[string]any{
			{"url": "/", "description": "Serveur courant"},
		},
		"tags": []map[string]any{
			{"name": "Système", "description": "Health check & statut serveur"},
			{"name": "Public", "description": "Routes publiques sans authentification"},
			{"name": "Auth Utilisateur", "description": "Authentification et gestion de session utilisateur"},
			{"name": "Profil Utilisateur", "description": "Profil, avatar, mot de passe, email, événements"},
			{"name": "Auth Admin", "description": "Authentification administrateur"},
			{"name": "Dashboard", "description": "Statistiques globales de la plateforme"},
			{"name": "Journaux", "description": "Logs d'audit et activité récente"},
			{"name": "Utilisateurs", "description": "Gestion des utilisateurs (admin)"},
			{"name": "Prestataires", "description": "Validation et gestion des prestataires (admin)"},
			{"name": "Catégories", "description": "Gestion du catalogue de catégories (admin)"},
			{"name": "Prestations", "description": "Gestion des prestations & services (admin)"},
			{"name": "Événements", "description": "Gestion des événements (admin)"},
			{"name": "Dépôts", "description": "Gestion des demandes de dépôt d'objets (admin)"},
			{"name": "Administrateurs (super_admin)", "description": "Gestion des comptes admin — réservé aux super_admin"},
		},
		"paths": paths,
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": bearerAuth,
			},
			"schemas": map[string]any{
				"User":           userSchema,
				"ProviderProfile": providerProfileSchema,
				"Category":       categorySchema,
				"Prestation":     prestationSchema,
				"Event":          eventSchema,
				"Deposit":        depositSchema,
				"AuditLog":       auditLogSchema,
				"PaginationMeta": paginationMeta,
				"Error":          errorSchema,
			},
		},
	}
}
