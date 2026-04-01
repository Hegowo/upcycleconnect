# UpcycleConnect

Plateforme complète UpcycleConnect — site public, espace utilisateur et interface d'administration.

## Architecture

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

## API

- **Base URL admin** : `/api/admin/v1`
- **Base URL utilisateur** : `/api/v1`
- **Authentification** : `Authorization: Bearer {token}`
- **Documentation interactive** : `/admin/docs` (Swagger UI, admin requis)

Backend go port 80 en interne. Apache (port `APP_PORT`) proxifie :
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
