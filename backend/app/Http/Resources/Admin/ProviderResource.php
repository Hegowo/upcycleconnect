<?php

namespace App\Http\Resources\Admin;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class ProviderResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        $profile = $this->providerProfile;

        return [
            'id'         => $this->id,
            'email'      => $this->email,
            'first_name' => $this->first_name,
            'last_name'  => $this->last_name,
            'phone'      => $this->phone,
            'status'     => $this->status,
            'created_at' => $this->created_at->toISOString(),
            'profile'    => $profile ? [
                'id'           => $profile->id,
                'company_name' => $profile->company_name,
                'siret'        => $profile->siret,
                'description'  => $profile->description,
                'website'      => $profile->website,
                'status'       => $profile->status,
                'approved_at'  => $profile->approved_at?->toISOString(),
            ] : null,
        ];
    }
}
