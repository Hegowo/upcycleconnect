<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\StorePrestationRequest;
use App\Http\Resources\Admin\PrestationResource;
use App\Models\Prestation;
use App\Services\Admin\AuditService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;

class PrestationController extends Controller
{
    public function __construct(private readonly AuditService $audit) {}

    public function index(Request $request): AnonymousResourceCollection
    {
        $query = Prestation::with(['category', 'provider'])
            ->whereNull('deleted_at');

        if ($request->filled('status')) {
            $query->where('status', $request->status);
        }

        if ($request->filled('category_id')) {
            $query->where('category_id', $request->category_id);
        }

        if ($request->filled('provider_id')) {
            $query->where('provider_id', $request->provider_id);
        }

        if ($request->filled('search')) {
            $search = '%'.$request->search.'%';
            $query->where('title', 'like', $search);
        }

        $prestations = $query->latest()->paginate($request->input('per_page', 20));

        return PrestationResource::collection($prestations);
    }

    public function store(StorePrestationRequest $request): JsonResponse
    {
        $prestation = Prestation::create($request->validated());
        $prestation->load('category', 'provider');

        $this->audit->log('prestation.created', 'Prestation', $prestation->id, null, $prestation->toArray());

        return response()->json(new PrestationResource($prestation), 201);
    }

    public function show(Prestation $prestation): JsonResponse
    {
        $prestation->load('category', 'provider');

        return response()->json(new PrestationResource($prestation));
    }

    public function update(StorePrestationRequest $request, Prestation $prestation): JsonResponse
    {
        $old = $prestation->toArray();
        $prestation->update($request->validated());
        $prestation->load('category', 'provider');

        $this->audit->log('prestation.updated', 'Prestation', $prestation->id, $old, $prestation->fresh()->toArray());

        return response()->json(new PrestationResource($prestation));
    }

    public function updateStatus(Request $request, Prestation $prestation): JsonResponse
    {
        $validated = $request->validate([
            'status' => ['required', 'in:draft,published,suspended,archived'],
        ]);

        $old = ['status' => $prestation->status];
        $prestation->update($validated);

        $this->audit->log('prestation.status_changed', 'Prestation', $prestation->id, $old, $validated);

        $prestation->load('category', 'provider');

        return response()->json(new PrestationResource($prestation));
    }

    public function destroy(Prestation $prestation): JsonResponse
    {
        $this->audit->log('prestation.deleted', 'Prestation', $prestation->id, $prestation->toArray(), null);
        $prestation->delete();

        return response()->json(null, 204);
    }
}
