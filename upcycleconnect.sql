-- MySQL dump 10.13  Distrib 8.4.8, for Linux (x86_64)
--
-- Host: localhost    Database: upcycleconnect
-- ------------------------------------------------------
-- Server version	8.4.8

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `audit_logs`
--

DROP TABLE IF EXISTS `audit_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `audit_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `action` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `resource_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `resource_id` bigint unsigned DEFAULT NULL,
  `old_values` json DEFAULT NULL,
  `new_values` json DEFAULT NULL,
  `ip_address` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `audit_logs_user_id_foreign` (`user_id`),
  CONSTRAINT `audit_logs_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `audit_logs`
--

LOCK TABLES `audit_logs` WRITE;
/*!40000 ALTER TABLE `audit_logs` DISABLE KEYS */;
INSERT INTO `audit_logs` VALUES (1,1,'admin.created','User',18,NULL,'{\"role\": \"super_admin\", \"email\": \"alexiachang@upcycleconnect.xyz\"}','88.190.118.230','2026-03-15 19:34:02'),(2,1,'admin.created','User',19,NULL,'{\"role\": \"super_admin\", \"email\": \"camillefarjon@upcycleconnect.xyz\"}','88.190.118.230','2026-03-15 19:35:27'),(3,1,'event.updated','Event',1,'{\"title\": \"Atelier upcycling textile\", \"status\": \"published\"}','{\"title\": \"Atelier upcycling textile\"}','46.193.58.47','2026-03-26 14:04:16');
/*!40000 ALTER TABLE `audit_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `deposit_requests`
--

DROP TABLE IF EXISTS `deposit_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `deposit_requests` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `category_id` bigint unsigned DEFAULT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `condition` enum('new','good','fair','poor') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'good',
  `history` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `estimated_weight` decimal(6,2) DEFAULT NULL,
  `carbon_savings` decimal(8,2) DEFAULT NULL,
  `status` enum('pending','approved','rejected') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
  `validation_note` text COLLATE utf8mb4_unicode_ci,
  `validated_by` bigint unsigned DEFAULT NULL,
  `validated_at` timestamp NULL DEFAULT NULL,
  `qr_code` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `deposit_requests_user_id_foreign` (`user_id`),
  KEY `deposit_requests_category_id_foreign` (`category_id`),
  KEY `deposit_requests_validated_by_foreign` (`validated_by`),
  CONSTRAINT `deposit_requests_category_id_foreign` FOREIGN KEY (`category_id`) REFERENCES `prestation_categories` (`id`) ON DELETE SET NULL,
  CONSTRAINT `deposit_requests_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `deposit_requests_validated_by_foreign` FOREIGN KEY (`validated_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `deposit_requests`
--

LOCK TABLES `deposit_requests` WRITE;
/*!40000 ALTER TABLE `deposit_requests` DISABLE KEYS */;
/*!40000 ALTER TABLE `deposit_requests` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `email_verifications`
--

DROP TABLE IF EXISTS `email_verifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `email_verifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `token` varchar(64) NOT NULL,
  `type` enum('register','login') NOT NULL DEFAULT 'register',
  `ip_address` varchar(45) DEFAULT NULL,
  `used_at` timestamp NULL DEFAULT NULL,
  `expires_at` timestamp NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `email_verifications`
--

LOCK TABLES `email_verifications` WRITE;
/*!40000 ALTER TABLE `email_verifications` DISABLE KEYS */;
INSERT INTO `email_verifications` VALUES (1,20,'c3d491bb3c1167ce64696b2ff9a41333aab4bada93a961b362662fc62bd628fd','register',NULL,'2026-03-26 03:57:14','2026-03-27 03:56:55','2026-03-26 03:56:55');
/*!40000 ALTER TABLE `email_verifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `migrations`
--

DROP TABLE IF EXISTS `migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `migrations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `migration` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `batch` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `migrations`
--

LOCK TABLES `migrations` WRITE;
/*!40000 ALTER TABLE `migrations` DISABLE KEYS */;
INSERT INTO `migrations` VALUES (1,'2024_01_01_000001_create_users_table',1),(2,'2024_01_01_000002_create_roles_table',1),(3,'2024_01_01_000003_create_user_roles_table',1),(4,'2024_01_01_000004_create_provider_profiles_table',1),(5,'2024_01_01_000005_create_prestation_categories_table',1),(6,'2024_01_01_000006_create_prestations_table',1),(7,'2024_01_01_000007_create_platform_events_table',1),(8,'2024_01_01_000008_create_audit_logs_table',1),(9,'2024_01_01_000009_create_personal_access_tokens_table',1);
/*!40000 ALTER TABLE `migrations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `personal_access_tokens`
--

DROP TABLE IF EXISTS `personal_access_tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `personal_access_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `tokenable_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `tokenable_id` bigint unsigned NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `abilities` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `last_used_at` timestamp NULL DEFAULT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `personal_access_tokens_token_unique` (`token`),
  KEY `personal_access_tokens_tokenable_type_tokenable_id_index` (`tokenable_type`,`tokenable_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `personal_access_tokens`
--

LOCK TABLES `personal_access_tokens` WRITE;
/*!40000 ALTER TABLE `personal_access_tokens` DISABLE KEYS */;
INSERT INTO `personal_access_tokens` VALUES (6,'App\\Models\\User',1,'admin-token','21592631245fb2a960b1e1f77fa3294b3f2aa49476d3384c78642342a5d4cd27','[\"*\"]','2026-03-14 22:12:29','2026-03-15 06:00:28','2026-03-14 22:00:28','2026-03-14 22:12:29'),(8,'App\\Models\\User',1,'admin-token','efda0ec71599cb3a877d4964a5e8f2aef9a8f4dc26c032aea7034ca0391267bb','[\"*\"]','2026-03-14 22:17:36','2026-03-15 06:17:06','2026-03-14 22:17:06','2026-03-14 22:17:36');
/*!40000 ALTER TABLE `personal_access_tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `platform_events`
--

DROP TABLE IF EXISTS `platform_events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `platform_events` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `category_id` bigint unsigned DEFAULT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `start_date` datetime NOT NULL,
  `end_date` datetime NOT NULL,
  `max_participants` int unsigned DEFAULT NULL,
  `status` enum('draft','published','cancelled') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'draft',
  `created_by` bigint unsigned NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `platform_events_category_id_foreign` (`category_id`),
  KEY `platform_events_created_by_foreign` (`created_by`),
  CONSTRAINT `platform_events_category_id_foreign` FOREIGN KEY (`category_id`) REFERENCES `prestation_categories` (`id`) ON DELETE SET NULL,
  CONSTRAINT `platform_events_created_by_foreign` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `platform_events`
--

LOCK TABLES `platform_events` WRITE;
/*!40000 ALTER TABLE `platform_events` DISABLE KEYS */;
INSERT INTO `platform_events` VALUES (1,2,'Atelier upcycling textile','Atelier pratique pour apprendre à transformer vos vieux vêtements en nouvelles pièces tendance.','Paris 11e — Maison des associations','2026-03-30 10:00:00','2026-03-30 13:00:00',15,'published',1,'2026-03-14 20:52:46','2026-03-26 14:04:16',NULL),(2,1,'Repair Café Mobilier','Venez avec votre mobilier abîmé, nos artisans vous aident à le réparer sur place.','Lyon 7e — Centre communautaire','2026-04-13 09:00:00','2026-04-13 17:00:00',30,'published',1,'2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(3,NULL,'Conférence : Économie circulaire et upcycling','Conférence sur les enjeux économiques et environnementaux de l\'upcycling.','Bordeaux — CCI Bordeaux Gironde','2026-05-13 14:00:00','2026-05-13 17:00:00',NULL,'draft',1,'2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(4,NULL,'Salon du réemploi','Événement annulé — report à une date ultérieure.','Nantes','2026-03-04 09:00:00','2026-03-04 18:00:00',200,'cancelled',1,'2026-03-14 20:52:46','2026-03-14 20:52:46',NULL);
/*!40000 ALTER TABLE `platform_events` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `prestation_categories`
--

DROP TABLE IF EXISTS `prestation_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `prestation_categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `sort_order` int NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `prestation_categories_slug_unique` (`slug`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `prestation_categories`
--

LOCK TABLES `prestation_categories` WRITE;
/*!40000 ALTER TABLE `prestation_categories` DISABLE KEYS */;
INSERT INTO `prestation_categories` VALUES (1,'Mobilier','mobilier','Réparation et customisation de mobilier',1,1,'2026-03-14 20:52:44','2026-03-14 20:52:44',NULL),(2,'Textile','textile','Upcycling de vêtements et tissu',1,2,'2026-03-14 20:52:44','2026-03-14 20:52:44',NULL),(3,'Électronique','electronique','Réparation et reconditionnement électronique',1,3,'2026-03-14 20:52:44','2026-03-14 20:52:44',NULL),(4,'Décoration','decoration','Objets de décoration recyclés',0,4,'2026-03-14 20:52:44','2026-03-14 20:52:44',NULL);
/*!40000 ALTER TABLE `prestation_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `prestations`
--

DROP TABLE IF EXISTS `prestations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `prestations` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `category_id` bigint unsigned DEFAULT NULL,
  `provider_id` bigint unsigned NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `price` decimal(10,2) DEFAULT NULL,
  `price_type` enum('fixed','hourly','quote') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'fixed',
  `status` enum('draft','published','suspended','archived') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'draft',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `prestations_category_id_foreign` (`category_id`),
  KEY `prestations_provider_id_foreign` (`provider_id`),
  CONSTRAINT `prestations_category_id_foreign` FOREIGN KEY (`category_id`) REFERENCES `prestation_categories` (`id`) ON DELETE SET NULL,
  CONSTRAINT `prestations_provider_id_foreign` FOREIGN KEY (`provider_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `prestations`
--

LOCK TABLES `prestations` WRITE;
/*!40000 ALTER TABLE `prestations` DISABLE KEYS */;
INSERT INTO `prestations` VALUES (1,1,13,'Restauration de chaise ancienne','Description complète de la prestation d\'upcycling proposée par le prestataire.',120.00,'fixed','published','2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(2,1,13,'Réparation table basse','Description complète de la prestation d\'upcycling proposée par le prestataire.',80.00,'fixed','published','2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(3,2,14,'Transformation jean en sac','Description complète de la prestation d\'upcycling proposée par le prestataire.',45.00,'fixed','published','2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(4,2,14,'Customisation veste vintage','Description complète de la prestation d\'upcycling proposée par le prestataire.',60.00,'hourly','published','2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(5,3,15,'Réparation smartphone','Description complète de la prestation d\'upcycling proposée par le prestataire.',50.00,'fixed','published','2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(6,3,15,'Reconditionnement PC portable','Description complète de la prestation d\'upcycling proposée par le prestataire.',NULL,'quote','draft','2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(7,4,16,'Création lampe à partir de bouteilles','Description complète de la prestation d\'upcycling proposée par le prestataire.',90.00,'fixed','draft','2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(8,1,13,'Vernissage de meuble','Description complète de la prestation d\'upcycling proposée par le prestataire.',40.00,'hourly','suspended','2026-03-14 20:52:46','2026-03-14 20:52:46',NULL);
/*!40000 ALTER TABLE `prestations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `provider_profiles`
--

DROP TABLE IF EXISTS `provider_profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `provider_profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `company_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `siret` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `website` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` enum('pending','approved','suspended') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
  `approved_at` timestamp NULL DEFAULT NULL,
  `approved_by` bigint unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `provider_profiles_user_id_unique` (`user_id`),
  KEY `provider_profiles_approved_by_foreign` (`approved_by`),
  CONSTRAINT `provider_profiles_approved_by_foreign` FOREIGN KEY (`approved_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  CONSTRAINT `provider_profiles_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `provider_profiles`
--

LOCK TABLES `provider_profiles` WRITE;
/*!40000 ALTER TABLE `provider_profiles` DISABLE KEYS */;
INSERT INTO `provider_profiles` VALUES (1,13,'Atelier du Bois','12345678901234','Prestataire spécialisé en upcycling.',NULL,'approved','2026-03-14 20:52:46',1,'2026-03-14 20:52:46','2026-03-14 20:52:46'),(2,14,'Verde Couture','12345678901234','Prestataire spécialisé en upcycling.',NULL,'approved','2026-03-14 20:52:46',1,'2026-03-14 20:52:46','2026-03-14 20:52:46'),(3,15,'TechRep Services','12345678901234','Prestataire spécialisé en upcycling.',NULL,'approved','2026-03-14 20:52:46',1,'2026-03-14 20:52:46','2026-03-14 20:52:46'),(4,16,'Artisan Créatif','12345678901234','Prestataire spécialisé en upcycling.',NULL,'pending',NULL,NULL,'2026-03-14 20:52:46','2026-03-14 20:52:46'),(5,17,'Eco Design Studio','12345678901234','Prestataire spécialisé en upcycling.',NULL,'pending',NULL,NULL,'2026-03-14 20:52:46','2026-03-14 20:52:46');
/*!40000 ALTER TABLE `provider_profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `label` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `roles_name_unique` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (1,'admin','Administrateur','Accès complet au back-office','2026-03-14 19:52:43'),(2,'super_admin','Super Administrateur','Accès complet incluant la gestion des comptes admin','2026-03-14 19:52:43');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_known_ips`
--

DROP TABLE IF EXISTS `user_known_ips`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_known_ips` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `ip_address` varchar(45) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_ip` (`user_id`,`ip_address`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_known_ips`
--

LOCK TABLES `user_known_ips` WRITE;
/*!40000 ALTER TABLE `user_known_ips` DISABLE KEYS */;
INSERT INTO `user_known_ips` VALUES (1,20,'176.186.78.195','2026-03-26 03:57:14');
/*!40000 ALTER TABLE `user_known_ips` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_roles`
--

DROP TABLE IF EXISTS `user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_roles` (
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `user_roles_role_id_foreign` (`role_id`),
  CONSTRAINT `user_roles_role_id_foreign` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `user_roles_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_roles`
--

LOCK TABLES `user_roles` WRITE;
/*!40000 ALTER TABLE `user_roles` DISABLE KEYS */;
INSERT INTO `user_roles` VALUES (1,1),(2,1),(1,2),(18,2),(19,2);
/*!40000 ALTER TABLE `user_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `first_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `last_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `avatar_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` enum('active','inactive','banned','pending') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'active',
  `email_verified_at` timestamp NULL DEFAULT NULL,
  `remember_token` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_email_unique` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'superadmin@upcycleconnect.fr','$2y$12$sKnGm42fO3t8Clqu4mdAD.O5YN31v4jvLlzQ.7wVlnhH2k/oZcdSe','Arthur','KETCHEIAN',NULL,NULL,'active','2026-03-14 20:52:43',NULL,'2026-03-14 20:52:43','2026-03-14 21:55:52',NULL),(2,'admin@upcycleconnect.fr','$2y$12$xVIA341KjMQho9prxyLRlubNvRGLDvh7PoZtghoVOBTxbsDR9vOGG','Admin','Demo',NULL,NULL,'active','2026-03-14 20:52:44',NULL,'2026-03-14 20:52:44','2026-03-14 20:52:44',NULL),(3,'marie.dupont@example.fr','$2y$12$EKYcuMGo68xE3mhw4VOrtu8TAzxhFYVRrmbIPz.zzgj3bV2gmiYse','Marie','Dupont',NULL,NULL,'active','2026-03-14 20:52:44',NULL,'2026-03-14 20:52:44','2026-03-14 20:52:44',NULL),(4,'pierre.martin@example.fr','$2y$12$P/Tls9gCj0NWMpJBz7SwXeKzjtXbZcTalDMLYs2fHsuDd0pmFdT6K','Pierre','Martin',NULL,NULL,'active','2026-03-14 20:52:44',NULL,'2026-03-14 20:52:44','2026-03-14 20:52:44',NULL),(5,'sophie.bernard@example.fr','$2y$12$IDp/Yl3wtdnjr08pSIIo6.HcMKW4M8HPh7n0DIS1W771pnrEVyybO','Sophie','Bernard',NULL,NULL,'active','2026-03-14 20:52:44',NULL,'2026-03-14 20:52:44','2026-03-14 20:52:44',NULL),(6,'lucas.petit@example.fr','$2y$12$GJt3C82uzcIkqdX9NPx.6e62jyIXd.79fszpw1IqfUq/V3cdJl58i','Lucas','Petit',NULL,NULL,'active','2026-03-14 20:52:44',NULL,'2026-03-14 20:52:44','2026-03-14 20:52:44',NULL),(7,'emma.robert@example.fr','$2y$12$ntRY5kUSMpcSRrtYp/fJLOErymaOGtQTgDXjj0uIZpQkJ1vU49Buu','Emma','Robert',NULL,NULL,'active','2026-03-14 20:52:45',NULL,'2026-03-14 20:52:45','2026-03-14 20:52:45',NULL),(8,'hugo.thomas@example.fr','$2y$12$fw82he0HcSS7YMhCWSOXou1beUCCx68B3dZiZr8Z46mfVKWs5xzXG','Hugo','Thomas',NULL,NULL,'active','2026-03-14 20:52:45',NULL,'2026-03-14 20:52:45','2026-03-14 20:52:45',NULL),(9,'lea.richard@example.fr','$2y$12$5rCTk95TCYHErbNGc6Lc/ut0qXY3xHYTKZgDftsvxi9TFWjYtQpli','Léa','Richard',NULL,NULL,'active','2026-03-14 20:52:45',NULL,'2026-03-14 20:52:45','2026-03-14 20:52:45',NULL),(10,'alex.durand@example.fr','$2y$12$bv4l9U3EaEg9GUZjmjpBl.0wn37aDtfjgV3TdhVWfvFh6ZASf16NK','Alex','Durand',NULL,NULL,'active','2026-03-14 20:52:45',NULL,'2026-03-14 20:52:45','2026-03-14 20:52:45',NULL),(11,'chloe.moreau@example.fr','$2y$12$RkDGebb60pWT4C4DPiRLKeIgg.eYUEU/l//t35.Bze41mFOFGVL6C','Chloé','Moreau',NULL,NULL,'inactive','2026-03-14 20:52:45',NULL,'2026-03-14 20:52:45','2026-03-14 20:52:45',NULL),(12,'tom.simon@example.fr','$2y$12$MRbtOZwRMWY6MrUfIY7uk.5KthuB2Z/M10BKBLw85ro4KG42HhgCK','Tom','Simon',NULL,NULL,'banned','2026-03-14 20:52:45',NULL,'2026-03-14 20:52:45','2026-03-14 20:52:45',NULL),(13,'atelier.bois@example.fr','$2y$12$SfeFCA4Z/BOMPdWqFGjNR.gn8PdaSYgiuc2Cg3T1qCT4rasdcDvlC','Atelier','Bois',NULL,NULL,'active','2026-03-14 20:52:46',NULL,'2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(14,'couture.verde@example.fr','$2y$12$PQ0D2CbaMui0h7Xp/RswU.NxgNn39TeUzASIfBAYWvI9wk5Oc4i9W','Verde','Couture',NULL,NULL,'active','2026-03-14 20:52:46',NULL,'2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(15,'repare.tech@example.fr','$2y$12$PY4u6D58NWjBx7EloVZsDOaOEMeJgDhUaPSs5Hfs7fTSIYiVy7t5C','TechRep','Services',NULL,NULL,'active','2026-03-14 20:52:46',NULL,'2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(16,'new.artisan@example.fr','$2y$12$AT.mUO0lEnzOfVIB67Grs.7VuIsI/QHFRdajIRNgXvt/NQlEjYxEK','Nouveau','Artisan',NULL,NULL,'active','2026-03-14 20:52:46',NULL,'2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(17,'eco.design@example.fr','$2y$12$JEQiaz0elyVvTA7l/E1XBOGKlq8OU.6eGZ5CJmJYezc0bDrZcVGbC','Eco','Design',NULL,NULL,'active','2026-03-14 20:52:46',NULL,'2026-03-14 20:52:46','2026-03-14 20:52:46',NULL),(18,'alexiachang@upcycleconnect.xyz','$2a$10$eFclyMJ/B.T84sjvY45StO1U.NXtdYEwtUrbh79rp3QvmIl3LIwom','Alexia','CHANG',NULL,NULL,'active','2026-03-15 19:34:02',NULL,'2026-03-15 19:34:02','2026-03-15 19:34:02',NULL),(19,'camillefarjon@upcycleconnect.xyz','$2a$10$9EsHu.TFaH.PLOEBEC5wm.hP89GTplC5xBYsJ0WgsM3flipxaxh0q','Camille','FARJON',NULL,NULL,'active','2026-03-15 19:35:27',NULL,'2026-03-15 19:35:27','2026-03-15 19:35:27',NULL),(20,'arthur.ketcheian@outlook.com','$2a$10$er3CWTGIdqLQfSgFjb/3zOyda60lAVGG4UVzj98AezP2SJ91.0Sw6','Arthur','KETCHEIAN',NULL,'/uploads/avatars/20.png','active','2026-03-26 03:57:14',NULL,'2026-03-26 03:56:55','2026-03-26 04:22:23',NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-03-26 15:40:38
