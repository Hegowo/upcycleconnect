<?php

namespace App\Http\Resources\Admin;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class UserResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id'                => $this->id,
            'email'             => $this->email,
            'first_name'        => $this->first_name,
            'last_name'         => $this->last_name,
            'phone'             => $this->phone,
            'status'            => $this->status,
            'role'              => $this->primary_role,
            'email_verified_at' => $this->email_verified_at?->toISOString(),
            'created_at'        => $this->created_at->toISOString(),
            'updated_at'        => $this->updated_at->toISOString(),
        ];
    }
}
