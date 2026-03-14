<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Models\Event;
use App\Models\Prestation;
use App\Models\ProviderProfile;
use App\Models\User;
use Illuminate\Http\JsonResponse;

class DashboardController extends Controller
{
    public function stats(): JsonResponse
    {
        $usersCount = User::whereDoesntHave('roles')
            ->whereNull('deleted_at')
            ->count();

        $providersCount = ProviderProfile::where('status', 'approved')->count();

        $providersPending = ProviderProfile::where('status', 'pending')->count();

        $prestationsCount = Prestation::whereNull('deleted_at')->count();

        $prestationsByStatus = Prestation::whereNull('deleted_at')
            ->selectRaw('status, COUNT(*) as count')
            ->groupBy('status')
            ->pluck('count', 'status')
            ->toArray();

        $eventsCount = Event::whereNull('deleted_at')
            ->whereIn('status', ['draft', 'published'])
            ->count();

        return response()->json([
            'users_count'          => $usersCount,
            'providers_count'      => $providersCount,
            'providers_pending'    => $providersPending,
            'prestations_count'    => $prestationsCount,
            'events_count'         => $eventsCount,
            'prestations_by_status' => $prestationsByStatus,
        ]);
    }
}
