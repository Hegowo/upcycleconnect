<?php

use App\Http\Controllers\Admin\AdminUserController;
use App\Http\Controllers\Admin\AuthController;
use App\Http\Controllers\Admin\CategoryController;
use App\Http\Controllers\Admin\DashboardController;
use App\Http\Controllers\Admin\EventController;
use App\Http\Controllers\Admin\PrestationController;
use App\Http\Controllers\Admin\ProviderController;
use App\Http\Controllers\Admin\UserController;
use Illuminate\Support\Facades\Route;

Route::prefix('admin/v1')->group(function () {

    // ── Auth (publique) ────────────────────────────────────────────────
    Route::post('auth/login', [AuthController::class, 'login']);

    // ── Routes protégées ──────────────────────────────────────────────
    Route::middleware(['auth:sanctum', 'admin'])->group(function () {

        Route::post('auth/logout', [AuthController::class, 'logout']);
        Route::get('auth/me', [AuthController::class, 'me']);

        // Dashboard
        Route::get('dashboard/stats', [DashboardController::class, 'stats']);

        // Utilisateurs particuliers
        Route::get('users', [UserController::class, 'index']);
        Route::get('users/{user}', [UserController::class, 'show']);
        Route::put('users/{user}', [UserController::class, 'update']);
        Route::post('users/{user}/ban', [UserController::class, 'ban']);
        Route::post('users/{user}/activate', [UserController::class, 'activate']);
        Route::delete('users/{user}', [UserController::class, 'destroy']);

        // Prestataires
        Route::get('providers', [ProviderController::class, 'index']);
        Route::get('providers/{user}', [ProviderController::class, 'show']);
        Route::put('providers/{user}/status', [ProviderController::class, 'updateStatus']);

        // Catégories
        Route::get('categories', [CategoryController::class, 'index']);
        Route::post('categories', [CategoryController::class, 'store']);
        Route::get('categories/{category}', [CategoryController::class, 'show']);
        Route::put('categories/{category}', [CategoryController::class, 'update']);
        Route::delete('categories/{category}', [CategoryController::class, 'destroy']);
        Route::post('categories/{category}/toggle', [CategoryController::class, 'toggle']);

        // Prestations
        Route::get('prestations', [PrestationController::class, 'index']);
        Route::post('prestations', [PrestationController::class, 'store']);
        Route::get('prestations/{prestation}', [PrestationController::class, 'show']);
        Route::put('prestations/{prestation}', [PrestationController::class, 'update']);
        Route::delete('prestations/{prestation}', [PrestationController::class, 'destroy']);
        Route::put('prestations/{prestation}/status', [PrestationController::class, 'updateStatus']);

        // Événements
        Route::get('events', [EventController::class, 'index']);
        Route::post('events', [EventController::class, 'store']);
        Route::get('events/{event}', [EventController::class, 'show']);
        Route::put('events/{event}', [EventController::class, 'update']);
        Route::delete('events/{event}', [EventController::class, 'destroy']);
        Route::put('events/{event}/status', [EventController::class, 'updateStatus']);

        // Comptes administrateurs (super_admin seulement)
        Route::middleware('super_admin')->group(function () {
            Route::get('admins', [AdminUserController::class, 'index']);
            Route::post('admins', [AdminUserController::class, 'store']);
            Route::get('admins/{user}', [AdminUserController::class, 'show']);
            Route::put('admins/{user}', [AdminUserController::class, 'update']);
            Route::delete('admins/{user}', [AdminUserController::class, 'destroy']);
        });
    });
});
