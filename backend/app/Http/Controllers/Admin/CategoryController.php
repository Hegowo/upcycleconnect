<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\StoreCategoryRequest;
use App\Http\Resources\Admin\CategoryResource;
use App\Models\PrestationCategory;
use App\Services\Admin\AuditService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;

class CategoryController extends Controller
{
    public function __construct(private readonly AuditService $audit) {}

    public function index(): AnonymousResourceCollection
    {
        $categories = PrestationCategory::withCount('prestations')
            ->orderBy('sort_order')
            ->orderBy('name')
            ->get();

        return CategoryResource::collection($categories);
    }

    public function store(StoreCategoryRequest $request): JsonResponse
    {
        $category = PrestationCategory::create($request->validated());

        $this->audit->log('category.created', 'PrestationCategory', $category->id, null, $category->toArray());

        return response()->json(new CategoryResource($category), 201);
    }

    public function show(PrestationCategory $category): JsonResponse
    {
        $category->loadCount('prestations');

        return response()->json(new CategoryResource($category));
    }

    public function update(StoreCategoryRequest $request, PrestationCategory $category): JsonResponse
    {
        $old = $category->toArray();
        $category->update($request->validated());

        $this->audit->log('category.updated', 'PrestationCategory', $category->id, $old, $category->fresh()->toArray());

        return response()->json(new CategoryResource($category));
    }

    public function destroy(PrestationCategory $category): JsonResponse
    {
        $this->audit->log('category.deleted', 'PrestationCategory', $category->id, $category->toArray(), null);
        $category->delete();

        return response()->json(null, 204);
    }

    public function toggle(PrestationCategory $category): JsonResponse
    {
        $old = ['is_active' => $category->is_active];
        $category->update(['is_active' => !$category->is_active]);

        $this->audit->log('category.toggled', 'PrestationCategory', $category->id, $old, ['is_active' => $category->is_active]);

        return response()->json(new CategoryResource($category));
    }
}
