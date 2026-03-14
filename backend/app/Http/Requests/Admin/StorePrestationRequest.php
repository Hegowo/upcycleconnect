<?php

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;

class StorePrestationRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    public function rules(): array
    {
        return [
            'category_id' => ['nullable', 'exists:prestation_categories,id'],
            'provider_id' => ['required', 'exists:users,id'],
            'title'       => ['required', 'string', 'max:255'],
            'description' => ['required', 'string'],
            'price'       => ['nullable', 'numeric', 'min:0'],
            'price_type'  => ['required', 'in:fixed,hourly,quote'],
            'status'      => ['required', 'in:draft,published,suspended,archived'],
        ];
    }

    public function messages(): array
    {
        return [
            'provider_id.required' => 'Le prestataire est requis.',
            'provider_id.exists'   => 'Le prestataire sélectionné est invalide.',
            'title.required'       => 'Le titre est requis.',
            'description.required' => 'La description est requise.',
            'price_type.required'  => 'Le type de tarif est requis.',
            'status.required'      => 'Le statut est requis.',
        ];
    }
}
