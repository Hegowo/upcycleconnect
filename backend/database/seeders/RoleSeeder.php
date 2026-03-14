<?php

namespace Database\Seeders;

use App\Models\Role;
use Illuminate\Database\Seeder;

class RoleSeeder extends Seeder
{
    public function run(): void
    {
        $roles = [
            ['name' => 'admin',       'label' => 'Administrateur',       'description' => 'Accès complet au back-office'],
            ['name' => 'super_admin', 'label' => 'Super Administrateur', 'description' => 'Accès complet incluant la gestion des comptes admin'],
        ];

        foreach ($roles as $role) {
            Role::firstOrCreate(['name' => $role['name']], $role);
        }
    }
}
