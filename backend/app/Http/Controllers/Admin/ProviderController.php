<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Resources\Admin\ProviderResource;
use App\Models\User;
use App\Services\Admin\AuditService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;

class ProviderController extends Controller
{
    public function __construct(private readonly AuditService $audit) {}

    public function index(Request $request): AnonymousResourceCollection
    {
        $query = User::has('providerProfile')
            ->with('providerProfile')
            ->whereNull('deleted_at');

        if ($request->filled('status')) {
            $query->whereHas('providerProfile', function ($q) use ($request) {
                $q->where('status', $request->status);
            });
        }

        if ($request->filled('search')) {
            $search = '%'.$request->search.'%';
            $query->where(function ($q) use ($search) {
                $q->where('email', 'like', $search)
                  ->orWhere('first_name', 'like', $search)
                  ->orWhere('last_name', 'like', $search)
                  ->orWhereHas('providerProfile', fn($q2) => $q2->where('company_name', 'like', $search));
            });
        }

        $providers = $query->latest()->paginate($request->input('per_page', 20));

        return ProviderResource::collection($providers);
    }

    public function show(User $user): JsonResponse
    {
        $user->load('providerProfile', 'prestations.category');

        return response()->json(new ProviderResource($user));
    }

    public function updateStatus(Request $request, User $user): JsonResponse
    {
        $validated = $request->validate([
            'status' => ['required', 'in:pending,approved,suspended'],
        ]);

        $profile = $user->providerProfile;

        if (!$profile) {
            return response()->json(['message' => 'Profil prestataire introuvable.'], 404);
        }

        $old = ['status' => $profile->status];

        $updateData = ['status' => $validated['status']];
        if ($validated['status'] === 'approved' && !$profile->approved_at) {
            $updateData['approved_at'] = now();
            $updateData['approved_by'] = $request->user()->id;
        }

        $profile->update($updateData);

        $this->audit->log(
            'provider.status_changed',
            'ProviderProfile',
            $profile->id,
            $old,
            ['status' => $validated['status']]
        );

        $user->load('providerProfile');

        return response()->json(new ProviderResource($user));
    }
}
