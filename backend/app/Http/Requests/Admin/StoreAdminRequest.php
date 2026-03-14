<?php

namespace App\Http\Requests\Admin;

use Illuminate\Foundation\Http\FormRequest;

class StoreAdminRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    public function rules(): array
    {
        return [
            'email'      => ['required', 'email', 'unique:users,email'],
            'password'   => ['required', 'string', 'min:8'],
            'first_name' => ['required', 'string', 'max:100'],
            'last_name'  => ['required', 'string', 'max:100'],
            'role'       => ['required', 'in:admin,super_admin'],
        ];
    }

    public function messages(): array
    {
        return [
            'email.unique'       => 'Cette adresse email est déjà utilisée.',
            'password.min'       => 'Le mot de passe doit contenir au moins 8 caractères.',
            'role.required'      => 'Le rôle est requis.',
            'role.in'            => 'Le rôle doit être "admin" ou "super_admin".',
        ];
    }
}
