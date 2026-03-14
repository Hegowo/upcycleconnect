<?php

namespace App\Http\Resources\Admin;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class PrestationResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id'          => $this->id,
            'title'       => $this->title,
            'description' => $this->description,
            'price'       => $this->price,
            'price_type'  => $this->price_type,
            'status'      => $this->status,
            'category'    => $this->whenLoaded('category', fn() => new CategoryResource($this->category)),
            'provider'    => $this->whenLoaded('provider', fn() => new UserResource($this->provider)),
            'created_at'  => $this->created_at->toISOString(),
            'updated_at'  => $this->updated_at->toISOString(),
        ];
    }
}
