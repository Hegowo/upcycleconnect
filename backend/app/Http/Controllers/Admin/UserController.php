<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Resources\Admin\UserResource;
use App\Models\User;
use App\Services\Admin\AuditService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;

class UserController extends Controller
{
    public function __construct(private readonly AuditService $audit) {}

    public function index(Request $request): AnonymousResourceCollection
    {
        $query = User::whereDoesntHave('roles')
            ->whereNull('deleted_at');

        if ($request->filled('search')) {
            $search = '%'.$request->search.'%';
            $query->where(function ($q) use ($search) {
                $q->where('email', 'like', $search)
                  ->orWhere('first_name', 'like', $search)
                  ->orWhere('last_name', 'like', $search);
            });
        }

        if ($request->filled('status')) {
            $query->where('status', $request->status);
        }

        $users = $query->latest()->paginate($request->input('per_page', 20));

        return UserResource::collection($users);
    }

    public function show(User $user): JsonResponse
    {
        $user->load('roles', 'providerProfile');

        return response()->json(new UserResource($user));
    }

    public function update(Request $request, User $user): JsonResponse
    {
        $validated = $request->validate([
            'first_name' => ['sometimes', 'string', 'max:100'],
            'last_name'  => ['sometimes', 'string', 'max:100'],
            'phone'      => ['nullable', 'string', 'max:20'],
        ]);

        $user->update($validated);

        return response()->json(new UserResource($user));
    }

    public function ban(User $user): JsonResponse
    {
        $old = ['status' => $user->status];
        $user->update(['status' => 'banned']);
        $this->audit->log('user.banned', 'User', $user->id, $old, ['status' => 'banned']);

        return response()->json(new UserResource($user));
    }

    public function activate(User $user): JsonResponse
    {
        $old = ['status' => $user->status];
        $user->update(['status' => 'active']);
        $this->audit->log('user.activated', 'User', $user->id, $old, ['status' => 'active']);

        return response()->json(new UserResource($user));
    }

    public function destroy(User $user): JsonResponse
    {
        $this->audit->log('user.deleted', 'User', $user->id, $user->toArray(), null);
        $user->delete();

        return response()->json(null, 204);
    }
}
