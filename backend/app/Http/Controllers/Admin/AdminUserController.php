<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\StoreAdminRequest;
use App\Http\Resources\Admin\UserResource;
use App\Models\Role;
use App\Models\User;
use App\Services\Admin\AuditService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;
use Illuminate\Support\Facades\Hash;

class AdminUserController extends Controller
{
    public function __construct(private readonly AuditService $audit) {}

    public function index(): AnonymousResourceCollection
    {
        $admins = User::whereHas('roles', function ($q) {
            $q->whereIn('name', ['admin', 'super_admin']);
        })->with('roles')->whereNull('deleted_at')->latest()->get();

        return UserResource::collection($admins);
    }

    public function show(User $user): JsonResponse
    {
        $user->load('roles');
        return response()->json(new UserResource($user));
    }

    public function store(StoreAdminRequest $request): JsonResponse
    {
        $user = User::create([
            'email'             => $request->email,
            'password'          => Hash::make($request->password),
            'first_name'        => $request->first_name,
            'last_name'         => $request->last_name,
            'status'            => 'active',
            'email_verified_at' => now(),
        ]);

        $role = Role::where('name', $request->role)->firstOrFail();
        $user->roles()->attach($role->id);

        $this->audit->log('admin.created', 'User', $user->id, null, ['email' => $user->email, 'role' => $request->role]);

        $user->load('roles');

        return response()->json(new UserResource($user), 201);
    }

    public function update(Request $request, User $user): JsonResponse
    {
        $validated = $request->validate([
            'first_name' => ['sometimes', 'string', 'max:100'],
            'last_name'  => ['sometimes', 'string', 'max:100'],
            'role'       => ['sometimes', 'in:admin,super_admin'],
            'password'   => ['sometimes', 'string', 'min:8'],
        ]);

        if (isset($validated['password'])) {
            $validated['password'] = Hash::make($validated['password']);
        }

        if (isset($validated['role'])) {
            $role = Role::where('name', $validated['role'])->firstOrFail();
            $user->roles()->sync([$role->id]);
            unset($validated['role']);
        }

        $user->update($validated);

        $user->load('roles');

        return response()->json(new UserResource($user));
    }

    public function destroy(Request $request, User $user): JsonResponse
    {
        if ($user->id === $request->user()->id) {
            return response()->json(['message' => 'Vous ne pouvez pas supprimer votre propre compte.'], 422);
        }

        $this->audit->log('admin.deleted', 'User', $user->id, ['email' => $user->email], null);
        $user->delete();

        return response()->json(null, 204);
    }
}
