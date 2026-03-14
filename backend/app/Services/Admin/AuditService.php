<?php

namespace App\Services\Admin;

use App\Models\AuditLog;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

class AuditService
{
    public function __construct(private readonly Request $request) {}

    public function log(
        string $action,
        string $resourceType,
        ?int $resourceId = null,
        ?array $oldValues = null,
        ?array $newValues = null
    ): void {
        AuditLog::create([
            'user_id'       => Auth::id(),
            'action'        => $action,
            'resource_type' => $resourceType,
            'resource_id'   => $resourceId,
            'old_values'    => $oldValues,
            'new_values'    => $newValues,
            'ip_address'    => $this->request->ip(),
            'created_at'    => now(),
        ]);
    }
}
