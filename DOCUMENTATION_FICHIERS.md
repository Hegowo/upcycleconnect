# Documentation des fichiers clés

Ce document décrit le rôle des fichiers applicatifs principaux du dépôt : tous les fichiers `.vue`, tous les fichiers `.go` et les fichiers d'infrastructure/configuration utiles pour lancer ou comprendre le projet.

## Points d'attention

- Le backend réellement présent dans ce dépôt est en Go (`Gin` + `GORM`), alors que plusieurs fichiers historiques (`README.md`, certaines variables d'environnement, une partie du Docker Compose) gardent des traces d'une ancienne base Laravel/PHP.
- Plusieurs pages publiques du frontend sont encore des maquettes statiques : elles servent surtout de vitrine UX et ne consomment pas encore l'API.
- `database/schema.sql` ressemble à un snapshot de base de données partiellement ancien : il ne reflète pas parfaitement tous les modèles Go actuels.
- Les répertoires générés comme `frontend/dist` et `frontend/node_modules` ne sont pas détaillés ici.

## Frontend

### Noyau frontend

- `frontend/src/main.js` : point d'entrée du SPA Vue. Monte l'application, branche Pinia, le routeur, l'i18n et les styles globaux.
- `frontend/src/App.vue` : racine minimale de l'application. Affiche le `RouterView` et le composant global de toasts.
- `frontend/src/router.js` : carte de toutes les routes publiques et admin. Gère aussi les guards d'authentification admin, la restriction `super_admin` et la mise à jour du titre de page.
- `frontend/src/assets/main.css` : couche CSS/Tailwind commune. Définit les classes utilitaires maison (`btn`, `input`, `card`, etc.) utilisées dans beaucoup d'écrans admin.
- `frontend/src/locales/fr.json` : dictionnaire principal des libellés FR de l'interface admin.
- `frontend/src/locales/en.json` : dictionnaire EN miroir pour le changement de langue dans l'admin.

### Services, stores et utilitaires frontend

- `frontend/src/services/api.js` : instance Axios de base pour l'API admin. Injecte automatiquement le token admin et redirige vers `/admin/login` en cas de `401`.
- `frontend/src/services/authService.js` : petite façade pour `login`, `logout` et `me` côté administration.
- `frontend/src/services/categoryService.js` : couche d'accès aux endpoints admin des catégories.
- `frontend/src/services/eventService.js` : couche d'accès aux endpoints admin des événements, y compris la mise à jour de statut.
- `frontend/src/services/prestationService.js` : couche d'accès aux endpoints admin des prestations.
- `frontend/src/services/userService.js` : regroupe les appels admin liés aux utilisateurs finaux, prestataires et administrateurs.
- `frontend/src/stores/auth.js` : store Pinia de session admin. Conserve le token, recharge le profil courant, expose `isAuthenticated`, `isSuperAdmin`, etc.
- `frontend/src/stores/userAuth.js` : store Pinia de session utilisateur public. Utilise `fetch` directement sur `/api/v1` pour inscription, connexion, profil et restauration de session.
- `frontend/src/utils/i18n.js` : initialise `vue-i18n`, choisit la langue depuis `localStorage` et expose `t`.
- `frontend/src/utils/useToast.js` : bus de toasts global basé sur un tableau réactif partagé.

### Composants Vue réutilisables

- `frontend/src/components/AppButton.vue` : bouton réutilisable avec variantes (`primary`, `secondary`, `danger`, `ghost`), tailles et état de chargement avec spinner.
- `frontend/src/components/AppBadge.vue` : badge d'état générique. Transforme des statuts métier (`active`, `pending`, `banned`, etc.) en couleurs cohérentes.
- `frontend/src/components/AppConfirmDialog.vue` : modale de confirmation téléportée dans le `body`. Sert aux suppressions et actions sensibles.
- `frontend/src/components/AppPagination.vue` : pagination générique avec ellipses et émission de l'événement `page-change`.
- `frontend/src/components/AppToast.vue` : rendu visuel des toasts du composable `useToast`, avec transition d'entrée/sortie.
- `frontend/src/components/Header.vue` : barre supérieure de l'admin. Affiche le titre courant, un champ de recherche surtout visuel, les infos du compte et le bouton de déconnexion.
- `frontend/src/components/Sidebar.vue` : navigation latérale desktop de l'admin. Construit les liens selon le rôle et affiche un bloc "Super Admin" quand nécessaire.
- `frontend/src/components/AdminBottomNav.vue` : navigation admin mobile. Reprend la logique du menu latéral en version compacte.
- `frontend/src/components/NavItem.vue` : wrapper simple autour de `RouterLink` pour styliser un lien actif. Semble surtout prévu comme primitive générique.
- `frontend/src/components/StatCard.vue` : carte KPI cliquable avec palette de couleurs et badge d'alerte éventuel.

### Layouts Vue

- `frontend/src/layouts/AdminLayout.vue` : squelette principal des pages admin. Assemble sidebar, header, contenu routé et navigation mobile ; respecte le mode `fullscreen` défini sur certaines routes.
- `frontend/src/layouts/AuthLayout.vue` : layout minimal de connexion admin centré avec logo et fallback d'image.
- `frontend/src/layouts/PublicLayout.vue` : layout du site public. Gère header, menu mobile, menu profil connecté, footer, restauration de session utilisateur et directive locale `click-outside`.

### Pages admin

- `frontend/src/pages/AdminLoginPage.vue` : écran de connexion des administrateurs. Utilise le store `auth` et route vers `/admin/dashboard` après succès.
- `frontend/src/pages/DashboardPage.vue` : tableau de bord admin. Affiche les KPIs, un score environnemental, une tendance mensuelle, un flux d'activité récent, une modale de détail de log et un export CSV.
- `frontend/src/pages/UserListPage.vue` : liste paginée des utilisateurs finaux. Gère recherche, filtre de statut, bannissement/réactivation, navigation vers la fiche et quelques cartes d'analyse ; certains boutons (`Ajouter`, `Filtres avancés`, crayon) restent surtout UI.
- `frontend/src/pages/UserDetailPage.vue` : fiche d'un utilisateur final. Affiche les informations de compte et permet de bannir ou réactiver le compte.
- `frontend/src/pages/ProviderListPage.vue` : listing des prestataires avec recherche, filtre de statut, pagination et actions rapides d'approbation/suspension. Les boutons d'export et d'ajout sont surtout visuels à ce stade.
- `frontend/src/pages/ProviderDetailPage.vue` : fiche détaillée d'un prestataire. Affiche SIRET, site, description et permet le changement de statut.
- `frontend/src/pages/AdminListPage.vue` : liste des comptes admin/super-admin. Permet création, édition et suppression, avec protection contre la suppression de soi-même côté backend.
- `frontend/src/pages/AdminFormPage.vue` : formulaire de création/édition d'administrateur. Gère prénom, nom, email, mot de passe et rôle.
- `frontend/src/pages/CategoryListPage.vue` : galerie de catégories avec compteur de prestations, activation/désactivation, édition et suppression.
- `frontend/src/pages/CategoryFormPage.vue` : formulaire CRUD de catégorie. Auto-génère le slug à partir du nom et permet de régler ordre/activation.
- `frontend/src/pages/PrestationListPage.vue` : écran de gestion du catalogue de prestations. Ajoute filtres par recherche, statut, catégorie et type de prix, plus publication/suspension/suppression.
- `frontend/src/pages/PrestationFormPage.vue` : formulaire CRUD d'une prestation. Choisit catégorie, prestataire, prix, type de prix et statut.
- `frontend/src/pages/EventListPage.vue` : liste des événements admin avec KPI de remplissage, prochain événement, filtre de statut, pagination et actions publier/annuler/supprimer.
- `frontend/src/pages/EventFormPage.vue` : formulaire CRUD d'événement. Gère titre, description, lieu, dates, capacité, catégorie et validation cohérente des dates.
- `frontend/src/pages/LogsPage.vue` : écran d'audit admin. Mélange logs d'administration, activité de plateforme, export CSV, détail modal JSON et section de gouvernance des administrateurs.
- `frontend/src/pages/BoxRequestsPage.vue` : page admin de traitement des demandes de dépôt. Affiche file d'attente, panneau de détail, note de validation et action d'approbation/rejet avec génération de QR code côté backend.
- `frontend/src/pages/SettingsPage.vue` : page de paramètres admin. Gère surtout des préférences locales d'interface comme la langue et quelques toggles non persistés côté serveur.
- `frontend/src/pages/SwaggerPage.vue` : interface Swagger UI embarquée. Charge les assets depuis `unpkg`, récupère le spec OpenAPI sur `/docs/spec` et injecte le token admin pour tester les endpoints.
- `frontend/src/pages/DatabasePage.vue` : page plein écran qui embarque un `iframe` vers `https://dbgate.upcycleconnect.xyz`. Elle ne pointe pas sur le `dbgate` local exposé par le Docker Compose.

### Pages publiques connectées à l'API

- `frontend/src/pages/PublicLoginPage.vue` : connexion utilisateur public. Gère aussi le cas de vérification secondaire si la connexion vient d'une IP inconnue.
- `frontend/src/pages/PublicRegisterPage.vue` : inscription multi-étapes. Permet de choisir particulier ou prestataire, de préremplir via SIRET, de mesurer la force du mot de passe puis d'envoyer le payload d'inscription.
- `frontend/src/pages/PublicVerifyEmailPage.vue` : page de validation du lien d'activation de compte. Gère les états succès, lien déjà utilisé et erreur.
- `frontend/src/pages/PublicVerifyLoginPage.vue` : page de confirmation de connexion depuis un nouvel emplacement. Récupère le JWT final et connecte automatiquement l'utilisateur.
- `frontend/src/pages/PublicProfilPage.vue` : tableau de bord personnel connecté. Consomme `/api/v1/profile` pour afficher score, dépôts, réservations et badges.
- `frontend/src/pages/PublicProfilEditPage.vue` : édition complète du profil utilisateur connecté. Gère l'upload d'avatar, la mise à jour des infos, le changement de mot de passe et le changement d'email en double vérification.

### Pages publiques surtout statiques / vitrine

- `frontend/src/pages/PublicHomePage.vue` : landing page marketing d'UpcycleConnect avec catégories mises en avant et CTA vers prestations/inscription.
- `frontend/src/pages/PublicPrestationsPage.vue` : vitrine publique des prestations, orientée contenu marketing plus qu'appel API.
- `frontend/src/pages/PublicEvenementsPage.vue` : catalogue public d'événements/ateliers sous forme de vitrine. Le contenu est aujourd'hui majoritairement statique.
- `frontend/src/pages/PublicEvenementDetailPage.vue` : détail d'un atelier public type. Sert surtout de page d'exemple visuelle.
- `frontend/src/pages/PublicCommunautePage.vue` : page de communauté mettant en avant projets, forum et actualités ; contenu principalement statique.
- `frontend/src/pages/PublicCategoriesForumPage.vue` : index des catégories de discussion du forum communautaire. Statique, orienté navigation UX.
- `frontend/src/pages/PublicDepotPage.vue` : formulaire public de dépôt d'objet côté UX. Les photos et la soumission sont simulées localement, sans branchement réel au backend.
- `frontend/src/pages/PublicConfirmationDepotPage.vue` : confirmation visuelle d'une demande de dépôt avec QR code et bon de dépôt d'exemple. C'est une maquette statique.
- `frontend/src/pages/PublicProfilProfessionnelPage.vue` : profil public d'un artisan/prestataire avec services, projets et boutique. Contenu de démonstration non alimenté par API.

## Backend Go

### Bootstrap, config et routing

- `backend/main.go` : point d'entrée du serveur Go. Charge la config, ouvre la base, construit le routeur et lance Gin sur le port demandé.
- `backend/config/config.go` : centralise la lecture des variables d'environnement. Contient surtout les paramètres DB, JWT, SMTP, URL d'application et clé API Société.com.
- `backend/database/database.go` : ouvre la connexion MySQL via GORM, configure le niveau de logs SQL et le pool de connexions.
- `backend/router/router.go` : compose tout le backend HTTP. Instancie les handlers/services, déclare les groupes de routes public/utilisateur/admin, branche CORS, auth et RBAC.

### Handlers HTTP

- `backend/handlers/auth.go` : authentification de l'interface admin. Vérifie email/mot de passe, impose un rôle admin/super-admin, émet un JWT puis expose `logout` et `me`.
- `backend/handlers/user_auth.go` : authentification et inscription côté utilisateur public. Gère création de compte, création d'un profil prestataire éventuel, email d'activation, connexion avec vérification d'IP inconnue, validation des liens et émission du JWT final.
- `backend/handlers/profile.go` : endpoints du profil utilisateur public. Retourne les stats/badges, permet la mise à jour des infos, l'upload d'avatar, le changement de mot de passe, le changement d'email en deux étapes et l'inscription/désinscription aux événements.
- `backend/handlers/siret.go` : proxy vers l'API Société.com. Valide le format SIRET, appelle le service externe et remonte proprement les erreurs fonctionnelles.
- `backend/handlers/dashboard.go` : agrégateur de statistiques pour le dashboard admin. Calcule utilisateurs, prestataires, prestations, événements, dépôts et total CO2 économisé.
- `backend/handlers/users.go` : gestion des utilisateurs finaux côté admin. Liste paginée, détail, mise à jour simple, bannissement, réactivation et suppression logique.
- `backend/handlers/providers.go` : gestion admin des prestataires. Liste filtrée, détail et changement de statut du `ProviderProfile`.
- `backend/handlers/admins.go` : CRUD des comptes d'administration réservé aux super-admins. Gère le hachage du mot de passe, l'affectation de rôle et l'interdiction de s'auto-supprimer.
- `backend/handlers/categories.go` : CRUD des catégories de prestations. Gère slug automatique, activation/désactivation, ordre d'affichage et audit.
- `backend/handlers/prestations.go` : CRUD admin des prestations. Supporte pagination, filtres, changement de statut et audit des modifications.
- `backend/handlers/events.go` : CRUD admin des événements. Gère les dates, la capacité, les compteurs d'inscriptions et les changements de statut.
- `backend/handlers/deposits.go` : traitement admin des demandes de dépôt. Fournit liste paginée, détail et validation/rejet ; l'approbation génère un code QR textuel.
- `backend/handlers/logs.go` : expose deux flux pour l'admin : les logs d'audit détaillés et l'activité récente de la plateforme (utilisateurs/prestataires), plus les admins utilisés par l'écran de logs.
- `backend/handlers/swagger.go` : construit à la main un document OpenAPI 3 en Go et le sert en JSON pour Swagger UI.

### Middleware

- `backend/middleware/auth.go` : middleware JWT. Extrait le bearer token, vérifie la signature avec `APP_KEY`, recharge l'utilisateur et le place dans le contexte Gin.
- `backend/middleware/rbac.go` : garde-fous de rôle. Refuse l'accès si l'utilisateur courant n'est pas admin ou super-admin selon le groupe.
- `backend/middleware/cors.go` : CORS permissif en `*`, avec méthodes et headers de base pour le frontend.

### Modèles

- `backend/models/user.go` : entité utilisateur principale. Porte les champs d'identité, le statut, les rôles, le profil prestataire lié et les helpers `IsAdmin` / `IsSuperAdmin`.
- `backend/models/role.go` : modèle simple de rôle (`admin`, `super_admin`, etc.).
- `backend/models/provider_profile.go` : extension métier d'un utilisateur prestataire. Contient la société, le SIRET, la description, le site et le statut de validation.
- `backend/models/prestation_category.go` : catégorie de prestations avec slug, ordre, activation et compteur de prestations associé dans la réponse.
- `backend/models/prestation.go` : prestation de catalogue liée à une catégorie et à un prestataire.
- `backend/models/event.go` : événement de plateforme avec dates, lieu, capacité, catégorie et créateur.
- `backend/models/event_registration.go` : table de liaison entre utilisateurs et événements.
- `backend/models/deposit_request.go` : demande de dépôt d'objet avec historique, poids, économie carbone, validation, QR code et relations vers utilisateur/catégorie.
- `backend/models/audit_log.go` : log d'audit admin avec action, ressource, ancien/nouveau JSON et IP.
- `backend/models/email_verification.go` : stocke les tokens de vérification email/connexion et la table `user_known_ips` pour reconnaître les IP déjà vues.
- `backend/models/email_change_request.go` : stocke les codes temporaires utilisés lors du changement d'email.

### Services backend

- `backend/services/audit.go` : service transversal d'écriture dans `audit_logs`. Sérialise les anciennes/nouvelles valeurs en JSON et ajoute l'IP/la personne courante.
- `backend/services/mailer.go` : client SMTP simple en texte brut. Gère l'envoi standard, le fallback TLS explicite sur le port `465` et trois templates : activation de compte, changement d'email et validation de connexion.

### Dépendances backend

- `backend/go.mod` : manifeste Go du backend. Déclare Gin, JWT, GORM, le driver MySQL, `bcrypt` et `slug`.

## Infrastructure et configuration

- `docker-compose.yml` : pile principale `mysql + backend + frontend + dbgate + apache`. Le fichier marche comme orchestrateur, mais il garde des variables et volumes hérités de Laravel (`SANCTUM_*`, `/var/www/html/storage`, etc.) qui ne sont pas tous utiles au backend Go actuel.
- `docker-compose.dev.yml` : surcharge de développement. Monte le code frontend pour Vite, mais le montage backend vers `/var/www/html` et la référence à `vendor` semblent datés d'une ancienne structure PHP et ne correspondent pas au `WORKDIR /app` du Dockerfile Go.
- `backend/Dockerfile` : build multi-stage du backend Go. Compile un binaire statique `server` puis le copie dans une image Alpine légère.
- `frontend/Dockerfile` : build de l'application Vite puis service du `dist` via Apache avec configuration SPA.
- `frontend/Dockerfile.dev` : image de développement frontend qui lance `vite` en écoutant sur `0.0.0.0:5173`.
- `apache/Dockerfile` : image Apache reverse proxy. Active les modules `proxy`, `proxy_http`, `proxy_wstunnel`, `headers` et `rewrite`, puis inclut `vhost.conf`.
- `apache/vhost.conf` : reverse proxy principal. Route `/api/public`, `/api/v1`, `/api/admin` et `/up` vers le backend, `/dbgate` vers DbGate et le reste vers le frontend.
- `frontend/apache.conf` : configuration Apache côté frontend pur. Réécrit toutes les routes inconnues vers `index.html` pour que le SPA fonctionne en direct URL.
- `.env.example` : base de configuration à copier. Le fichier documente les variables générales et MySQL, mais il garde aussi des clés Laravel/Sanctum historiques et ne mentionne pas toutes les variables réellement lues par le backend Go (`SMTP_*`, `MAIL_FROM`, `SOCIETE_API_KEY`).
- `.env` : valeurs locales d'exécution. Sert au lancement réel sur la machine/en Docker ; il contient des secrets et ne doit pas être utilisé comme documentation métier.
- `database/schema.sql` : dump SQL du schéma MySQL. Il documente les tables principales (`users`, `provider_profiles`, `prestations`, `platform_events`, `audit_logs`, etc.), mais contient aussi des traces historiques (`migrations`, `personal_access_tokens`) et semble en décalage avec certains modèles actuels.
- `database/init/.gitkeep` : simple fichier technique pour conserver le dossier `database/init` dans Git, dossier monté par `docker-compose.yml`.
- `frontend/package.json` : manifeste Node du frontend. Déclare Vue 3, Pinia, Vue Router, Axios, Vue I18n et les scripts `dev/build/preview`.
- `frontend/vite.config.js` : configuration Vite. Définit l'alias `@`, expose le serveur sur `0.0.0.0:5173` et proxyfie `/api` vers `localhost:8080`.
- `frontend/tailwind.config.js` : personnalisation Tailwind. Définit la palette de couleurs de marque et la fonte `Plus Jakarta Sans`.
- `frontend/postcss.config.js` : branche `tailwindcss` et `autoprefixer`.
- `frontend/index.html` : shell HTML de Vite. Charge la police Google Fonts et le point d'entrée `/src/main.js`.
- `README.md` : README du projet. Il donne une vue d'ensemble utile, mais la section stack/backend est obsolète puisqu'elle parle encore de Laravel/PHP alors que le dépôt contient aujourd'hui un backend Go.

