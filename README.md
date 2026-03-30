# UpcycleConnect

Plateforme complète UpcycleConnect — site public, espace utilisateur et interface d'administration.

Le frontend couvre l'intégralité du site : pages publiques (accueil, prestations, événements, communauté), espace utilisateur connecté (profil, dépôt d'objets, inscriptions aux événements) et back-office admin. Le backend Go expose une API REST unique consommée par toutes ces parties.

## Stack

| Couche | Technologie |
|---|---|
| Frontend | Vue 3 + Vite + Tailwind CSS + Pinia |
| Backend | Go 1.22 + Gin + GORM + JWT |
| Base de données | MySQL 8.4 |
| Reverse proxy | Apache 2.4 |
| Conteneurisation | Docker + Docker Compose |

---

## Démarrage rapide (Docker)

Le script `upcycleinit` (macOS/Linux) ou `upcycleinit.ps1` (Windows) automatise l'intégralité du premier démarrage : build des images, attente MySQL, import SQL, démarrage de tous les services.

### macOS / Linux

```bash
chmod +x upcycleinit
./upcycleinit
```

### Windows (PowerShell)

```powershell
.\upcycleinit.ps1
```

---

### Démarrage manuel

#### 1. Variables d'environnement

```bash
cp .env.example .env
```

Modifier `.env` avec des valeurs sécurisées (mots de passe, clé JWT, SMTP).

#### 2. Lancer les conteneurs

```bash
docker compose up -d --build
```

#### 3. Importer la base de données

```bash
docker exec -i upcycleconnect_mysql mysql -u root -p"${DB_ROOT_PASSWORD}" upcycleconnect < upcycleconnect.sql
```

#### 4. Accéder à l'application

- **Site public** : http://localhost:8080/
- **Interface admin** : http://localhost:8080/admin/login
- **API** : http://localhost:8080/api/admin/v1
- **Documentation API** : http://localhost:8080/admin/docs *(admin requis)*
- **DBGate (UI base de données)** : http://localhost:8080/dbgate/

---

## Développement local (sans Docker)

### Backend (Go)

```bash
cd backend
go mod download
cp ../.env.example ../.env
go run main.go
```

Backend disponible sur : http://localhost:80

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend disponible sur : http://localhost:5173

> Configurer `VITE_API_BASE_URL=http://localhost:80/api/admin/v1` dans le `.env` pour le dev local.

---

## Structure du projet

```
upcycleconnect/
├── frontend/                  # Vue 3 SPA (site complet)
│   ├── src/
│   │   ├── pages/             # Pages publiques (Public*) + admin
│   │   ├── components/        # Composants réutilisables
│   │   ├── stores/            # État global (Pinia)
│   │   ├── services/          # Couche API (Axios)
│   │   ├── layouts/           # PublicLayout / AdminLayout / AuthLayout
│   │   ├── locales/           # Traductions FR / EN
│   │   └── router.js          # Routes + guards
│   └── Dockerfile
├── backend/                   # API REST Go
│   ├── config/                # Chargement des variables d'environnement
│   ├── database/              # Connexion MySQL + AutoMigrate
│   ├── handlers/              # Handlers Gin (un fichier par ressource)
│   ├── middleware/            # Auth JWT, RBAC, CORS
│   ├── models/                # Modèles GORM
│   ├── router/                # Déclaration de toutes les routes
│   ├── services/              # Mailer, AuditService
│   ├── main.go
│   └── Dockerfile
├── apache/
│   └── vhost.conf             # Reverse proxy Apache
├── database/
│   └── init/                  # Initialisé automatiquement par MySQL au démarrage
├── docs/
│   └── arborescence.html      # Documentation complète de l'arborescence
├── upcycleconnect.sql         # Dump SQL de la base de données
├── upcycleinit                # Script de démarrage macOS/Linux
├── upcycleinit.ps1            # Script de démarrage Windows
├── docker-compose.yml
└── .env.example
```

---

## Modules admin disponibles

| Module | URL | Rôle requis |
|---|---|---|
| Tableau de bord | /admin/dashboard | admin |
| Utilisateurs | /admin/users | admin |
| Prestataires | /admin/providers | admin |
| Catégories | /admin/categories | admin |
| Prestations | /admin/prestations | admin |
| Événements | /admin/events | admin |
| Journaux | /admin/logs | admin |
| Demandes de dépôt | /admin/box-requests | admin |
| Paramètres | /admin/settings | admin |
| Documentation API | /admin/docs | admin |
| Administrateurs | /admin/admins | super_admin |
| Base de données | /admin/database | super_admin |

---

## API

- **Base URL admin** : `/api/admin/v1`
- **Base URL utilisateur** : `/api/v1`
- **Authentification** : `Authorization: Bearer {token}`
- **Documentation interactive** : `/admin/docs` (Swagger UI, admin requis)

Le backend Go démarre sur le port 80 en interne. Apache (port `APP_PORT`) proxifie :
- `/api/*` → backend Go
- `/dbgate/*` → DBGate
- `/*` → frontend Vue

---

## Services Docker

| Conteneur | Image | Rôle |
|---|---|---|
| `upcycleconnect_mysql` | mysql:8.4 | Base de données |
| `upcycleconnect_backend` | Go 1.22 / Alpine | API REST |
| `upcycleconnect_frontend` | Node 20 / Nginx | SPA Vue 3 |
| `upcycleconnect_apache` | Apache 2.4 | Reverse proxy |
| `upcycleconnect_dbgate` | dbgate/dbgate | Interface DB |
