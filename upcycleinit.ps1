#Requires -Version 5.1
$ErrorActionPreference = "Stop"

function Write-Ok   { param($msg) Write-Host "  " -NoNewline; Write-Host "v" -ForegroundColor Green -NoNewline; Write-Host "  $msg" }
function Write-Info { param($msg) Write-Host "  " -NoNewline; Write-Host "i" -ForegroundColor Cyan  -NoNewline; Write-Host "  $msg" }
function Write-Warn { param($msg) Write-Host "  " -NoNewline; Write-Host "!" -ForegroundColor Yellow -NoNewline; Write-Host "  $msg" }
function Write-Step { param($msg) Write-Host "`n  >> $msg" -ForegroundColor Blue }
function Write-Fail {
  param($msg)
  Write-Host "  " -NoNewline
  Write-Host "X  $msg" -ForegroundColor Red
  exit 1
}

function Get-EnvValue {
  param([string]$Key, [string]$Default = "")
  $envPath = Join-Path $PSScriptRoot ".env"
  if (-not (Test-Path $envPath)) { return $Default }
  $line = Get-Content $envPath | Where-Object { $_ -match "^${Key}=" } | Select-Object -First 1
  if (-not $line) { return $Default }
  return ($line -split "=", 2)[1].Trim().Trim('"').Trim("'")
}

Write-Host ""
Write-Host "  UpcycleConnect" -ForegroundColor Green
Write-Host "  Initialisation premier demarrage (Windows)" -ForegroundColor DarkGray
Write-Host "  -----------------------------------------------" -ForegroundColor DarkGray
Write-Host ""

$ROOT          = $PSScriptRoot
$SQL_FILE      = Join-Path $ROOT "upcycleconnect.sql"
$ENV_FILE      = Join-Path $ROOT ".env"
$ENV_EXAMPLE   = Join-Path $ROOT ".env.example"
$COMPOSE_FILE  = Join-Path $ROOT "docker-compose.yml"

Write-Step "Verification des prerequis"

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
  Write-Fail "Docker n'est pas installe. Telechargez Docker Desktop : https://docs.docker.com/get-docker/"
}
Write-Ok "Docker trouve : $(docker --version)"

$COMPOSE_CMD = $null
try {
  docker compose version | Out-Null
  $COMPOSE_CMD = "docker compose"
  Write-Ok "Docker Compose (plugin v2) disponible"
} catch {
  if (Get-Command docker-compose -ErrorAction SilentlyContinue) {
    $COMPOSE_CMD = "docker-compose"
    Write-Ok "docker-compose (v1) disponible"
  } else {
    Write-Fail "Docker Compose introuvable. Installez Docker Desktop avec Compose."
  }
}

try {
  docker info | Out-Null
  Write-Ok "Demon Docker actif"
} catch {
  Write-Fail "Le demon Docker n'est pas demarre. Lancez Docker Desktop et reessayez."
}

if (-not (Test-Path $COMPOSE_FILE)) { Write-Fail "docker-compose.yml introuvable dans $ROOT" }
Write-Ok "docker-compose.yml present"

if (-not (Test-Path $SQL_FILE)) { Write-Fail "Fichier SQL introuvable : $SQL_FILE" }
Write-Ok "Fichier SQL present : upcycleconnect.sql"

Write-Step "Configuration de l'environnement"

if (-not (Test-Path $ENV_FILE)) {
  if (-not (Test-Path $ENV_EXAMPLE)) { Write-Fail ".env.example introuvable, impossible de creer .env" }
  Copy-Item $ENV_EXAMPLE $ENV_FILE
  Write-Ok ".env cree a partir de .env.example"
  Write-Warn "Verifiez et personnalisez les valeurs dans .env si necessaire"
} else {
  Write-Ok ".env deja present"
}

$DB_ROOT_PASSWORD = Get-EnvValue "DB_ROOT_PASSWORD" "rootpassword"
$DB_DATABASE      = Get-EnvValue "DB_DATABASE"      "upcycleconnect"
$DB_USERNAME      = Get-EnvValue "DB_USERNAME"      "ucadmin"
$DB_PASSWORD      = Get-EnvValue "DB_PASSWORD"      "strongpassword"
$APP_PORT         = Get-EnvValue "APP_PORT"          "8080"

Write-Info "Base : $DB_DATABASE  |  Utilisateur : $DB_USERNAME  |  Port app : $APP_PORT"

Write-Step "Verification des donnees existantes"

$VOLUME_NAME  = "upcycleconnect_mysql_data"
$SKIP_IMPORT  = $false
$volumeExists = $false

try {
  docker volume inspect $VOLUME_NAME | Out-Null
  $volumeExists = $true
} catch {}

if ($volumeExists) {
  Write-Warn "Le volume MySQL '$VOLUME_NAME' existe deja."
  Write-Host ""
  Write-Host "  Que souhaitez-vous faire ?" -ForegroundColor Yellow
  Write-Host "  [1] Reinitialiser la base (supprime les donnees) — recommande pour premier demarrage propre"
  Write-Host "  [2] Conserver les donnees existantes et demarrer seulement les services"
  Write-Host "  [q] Quitter"
  Write-Host ""
  $choice = Read-Host "  Votre choix [1/2/q]"

  switch ($choice) {
    "1" {
      Write-Info "Suppression du volume MySQL existant..."
      if ($COMPOSE_CMD -eq "docker compose") {
        docker compose -f $COMPOSE_FILE down --volumes 2>$null
      } else {
        docker-compose -f $COMPOSE_FILE down --volumes 2>$null
      }
      docker volume rm $VOLUME_NAME 2>$null | Out-Null
      Write-Ok "Volume supprime — reinitialisation propre"
      $SKIP_IMPORT = $false
    }
    "2" {
      Write-Ok "Conservation des donnees existantes"
      $SKIP_IMPORT = $true
    }
    default {
      Write-Host ""; Write-Info "Annule."; exit 0
    }
  }
} else {
  Write-Ok "Aucune donnee existante — demarrage vierge"
  $SKIP_IMPORT = $false
}

Write-Step "Construction des images Docker"

Write-Info "Build en cours... (peut prendre quelques minutes la premiere fois)"

if ($COMPOSE_CMD -eq "docker compose") {
  docker compose -f $COMPOSE_FILE build --parallel
} else {
  docker-compose -f $COMPOSE_FILE build --parallel
}

Write-Ok "Images construites"

Write-Step "Demarrage de MySQL"

if ($COMPOSE_CMD -eq "docker compose") {
  docker compose -f $COMPOSE_FILE up -d mysql
} else {
  docker-compose -f $COMPOSE_FILE up -d mysql
}

Write-Info "Attente que MySQL soit pret (jusqu'a 90s)..."

$elapsed = 0
$timeout = 90
$mysqlReady = $false

while ($elapsed -lt $timeout) {
  try {
    $result = docker exec upcycleconnect_mysql mysqladmin ping -h localhost --silent 2>&1
    if ($LASTEXITCODE -eq 0) {
      $mysqlReady = $true
      break
    }
  } catch {}
  Write-Host "  . " -NoNewline -ForegroundColor DarkGray
  Start-Sleep -Seconds 3
  $elapsed += 3
}

Write-Host ""

if (-not $mysqlReady) {
  Write-Fail "MySQL n'a pas demarre apres ${timeout}s. Verifiez avec : docker logs upcycleconnect_mysql"
}
Write-Ok "MySQL est pret (${elapsed}s)"

if (-not $SKIP_IMPORT) {
  Write-Step "Import de la base de donnees"
  Write-Info "Import de upcycleconnect.sql en cours..."

  $importStart = Get-Date
  $importSuccess = $false

  try {
    $sqlContent = Get-Content $SQL_FILE -Raw -Encoding UTF8
    $sqlContent | docker exec -i upcycleconnect_mysql `
      mysql -u root -p"$DB_ROOT_PASSWORD" `
      --default-character-set=utf8mb4 `
      $DB_DATABASE
    if ($LASTEXITCODE -eq 0) { $importSuccess = $true }
  } catch {
    $importSuccess = $false
  }

  if (-not $importSuccess) {
    Write-Warn "Tentative avec l'utilisateur $DB_USERNAME..."
    try {
      $sqlContent = Get-Content $SQL_FILE -Raw -Encoding UTF8
      $sqlContent | docker exec -i upcycleconnect_mysql `
        mysql -u "$DB_USERNAME" -p"$DB_PASSWORD" `
        --default-character-set=utf8mb4 `
        $DB_DATABASE
      if ($LASTEXITCODE -eq 0) { $importSuccess = $true }
    } catch {
      $importSuccess = $false
    }
  }

  if (-not $importSuccess) {
    Write-Fail "Echec de l'import SQL. Verifiez les credentials dans .env"
  }

  $duration = [int]((Get-Date) - $importStart).TotalSeconds
  Write-Ok "Import termine en ${duration}s"
} else {
  Write-Info "Import ignore (donnees conservees)"
}

Write-Step "Demarrage de tous les services"

if ($COMPOSE_CMD -eq "docker compose") {
  docker compose -f $COMPOSE_FILE up -d
} else {
  docker-compose -f $COMPOSE_FILE up -d
}

Write-Ok "Tous les services demarres"

Write-Step "Verification de l'etat des services"

Write-Info "Attente de la stabilisation (15s)..."
Start-Sleep -Seconds 15

Write-Host ""
Write-Host "  Conteneurs :" -ForegroundColor White
docker ps --filter "name=upcycleconnect_" --format "  {{.Names}} | {{.Status}} | {{.Ports}}"

Write-Host ""
Write-Host "  Endpoints :" -ForegroundColor White

try {
  $response = Invoke-WebRequest -Uri "http://localhost:$APP_PORT/up" -TimeoutSec 5 -UseBasicParsing -ErrorAction Stop
  if ($response.StatusCode -eq 200) {
    Write-Ok "Backend  -> http://localhost:$APP_PORT/up (OK)"
  }
} catch {
  Write-Warn "Backend  -> http://localhost:$APP_PORT/up (pas encore joignable)"
}

try {
  $response = Invoke-WebRequest -Uri "http://localhost:$APP_PORT/" -TimeoutSec 5 -UseBasicParsing -ErrorAction Stop
  if ($response.StatusCode -eq 200) {
    Write-Ok "Frontend -> http://localhost:$APP_PORT/ (OK)"
  }
} catch {
  Write-Warn "Frontend -> http://localhost:$APP_PORT/ (pas encore joignable)"
}

Write-Host ""
Write-Host "  -----------------------------------------------" -ForegroundColor DarkGray
Write-Host ""
Write-Host "  [OK] UpcycleConnect est pret !" -ForegroundColor Green
Write-Host ""
Write-Host "  Acces :" -ForegroundColor White
Write-Host "    Application    ->  http://localhost:$APP_PORT" -ForegroundColor Cyan
Write-Host "    Admin panel    ->  http://localhost:$APP_PORT/admin/login" -ForegroundColor Cyan
Write-Host "    DBGate (DB UI) ->  http://localhost:3000" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Commandes utiles :" -ForegroundColor White
Write-Host "    Logs en direct  :  docker compose logs -f" -ForegroundColor DarkGray
Write-Host "    Arreter         :  docker compose down" -ForegroundColor DarkGray
Write-Host "    Redemarrer      :  docker compose restart" -ForegroundColor DarkGray
Write-Host ""
