<?php

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;

class StoreEventRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    public function rules(): array
    {
        return [
            'category_id'      => ['nullable', 'exists:prestation_categories,id'],
            'title'            => ['required', 'string', 'max:255'],
            'description'      => ['required', 'string'],
            'location'         => ['nullable', 'string', 'max:255'],
            'start_date'       => ['required', 'date'],
            'end_date'         => ['required', 'date', 'after:start_date'],
            'max_participants' => ['nullable', 'integer', 'min:1'],
            'status'           => ['required', 'in:draft,published,cancelled'],
        ];
    }

    public function messages(): array
    {
        return [
            'title.required'       => 'Le titre est requis.',
            'description.required' => 'La description est requise.',
            'start_date.required'  => 'La date de début est requise.',
            'end_date.required'    => 'La date de fin est requise.',
            'end_date.after'       => 'La date de fin doit être postérieure à la date de début.',
        ];
    }
}
