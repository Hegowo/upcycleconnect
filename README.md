# UpcycleConnect — Administration

Interface d'administration de la plateforme UpcycleConnect.

## Stack

| Couche | Technologie |
|---|---|
| Frontend | Vue 3 + Vite + Tailwind CSS + Pinia |
| Backend | PHP 8.3 + Laravel 11 + Sanctum |
| Base de données | MySQL 8.0 |
| Reverse proxy | Apache 2.4 |
| Conteneurisation | Docker + Docker Compose |

---

## Démarrage rapide (Docker)

### Prérequis
- Docker >= 24
- Docker Compose >= 2.20

### 1. Variables d'environnement

```bash
cp .env.example .env
```

Modifier `.env` avec des mots de passe sécurisés. Générer une `APP_KEY` Laravel :

```bash
docker run --rm php:8.3-cli php -r "echo 'base64:'.base64_encode(random_bytes(32));"
```

Copier la valeur dans `APP_KEY=` dans `.env`.

### 2. Lancer les conteneurs

```bash
docker compose up -d --build
```

### 3. Migrations et données de démo

```bash
docker compose exec backend php artisan migrate --seed
```

### 4. Accéder à l'application

- **Interface admin** : http://localhost:8080
- **API** : http://localhost:8080/api/admin/v1

---

## Comptes de démonstration

| Email | Mot de passe | Rôle |
|---|---|---|
| superadmin@upcycleconnect.fr | Admin1234! | Super Administrateur |
| admin@upcycleconnect.fr | Admin1234! | Administrateur |

---

## Développement local (sans Docker)

### Backend

```bash
cd backend
composer install
cp .env.example .env
php artisan key:generate
php artisan migrate --seed
php artisan serve
```

Backend disponible sur : http://localhost:8000

### Frontend

```bash
cd frontend
npm install
cp .env.example .env
# Mettre VITE_API_BASE_URL=http://localhost:8000/api/admin/v1
npm run dev
```

Frontend disponible sur : http://localhost:5173

---

## Structure du projet

```
upcycleconnect-admin/
├── frontend/              # Vue 3 SPA Admin
│   ├── src/
│   │   ├── pages/         # Pages par module
│   │   ├── components/    # Composants réutilisables
│   │   ├── stores/        # État global (Pinia)
│   │   ├── services/      # Couche API (Axios)
│   │   ├── layouts/       # Layouts admin / auth
│   │   └── router/        # Routes + guards
│   └── Dockerfile
├── backend/               # Laravel 11 API REST
│   ├── app/
│   │   ├── Http/
│   │   │   ├── Controllers/Admin/
│   │   │   ├── Middleware/
│   │   │   ├── Requests/Admin/
│   │   │   └── Resources/Admin/
│   │   ├── Models/
│   │   └── Services/Admin/
│   ├── database/
│   │   ├── migrations/
│   │   └── seeders/
│   └── Dockerfile
├── apache/
│   └── httpd.conf         # Reverse proxy config
├── database/init/         # Scripts SQL optionnels
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
| Administrateurs | /admin/admins | super_admin |

---

## API

Base URL : `GET /api/admin/v1`

Authentification : `Authorization: Bearer {token}`

Voir la documentation complète dans `/docs/API.md`.
