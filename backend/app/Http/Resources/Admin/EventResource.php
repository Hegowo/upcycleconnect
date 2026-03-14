<?php

namespace App\Http\Resources\Admin;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class EventResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id'               => $this->id,
            'title'            => $this->title,
            'description'      => $this->description,
            'location'         => $this->location,
            'start_date'       => $this->start_date->toISOString(),
            'end_date'         => $this->end_date->toISOString(),
            'max_participants' => $this->max_participants,
            'status'           => $this->status,
            'category'         => $this->whenLoaded('category', fn() => new CategoryResource($this->category)),
            'creator'          => $this->whenLoaded('creator', fn() => new UserResource($this->creator)),
            'created_at'       => $this->created_at->toISOString(),
            'updated_at'       => $this->updated_at->toISOString(),
        ];
    }
}
