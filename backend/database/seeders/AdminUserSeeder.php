<?php

namespace Database\Seeders;

use App\Models\Role;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;

class AdminUserSeeder extends Seeder
{
    public function run(): void
    {
        $superAdminRole = Role::where('name', 'super_admin')->firstOrFail();
        $adminRole      = Role::where('name', 'admin')->firstOrFail();

        $superAdmin = User::firstOrCreate(
            ['email' => 'superadmin@upcycleconnect.fr'],
            [
                'password'           => Hash::make('Admin1234!'),
                'first_name'         => 'Super',
                'last_name'          => 'Admin',
                'status'             => 'active',
                'email_verified_at'  => now(),
            ]
        );
        $superAdmin->roles()->syncWithoutDetaching([$superAdminRole->id]);

        $admin = User::firstOrCreate(
            ['email' => 'admin@upcycleconnect.fr'],
            [
                'password'           => Hash::make('Admin1234!'),
                'first_name'         => 'Admin',
                'last_name'          => 'Demo',
                'status'             => 'active',
                'email_verified_at'  => now(),
            ]
        );
        $admin->roles()->syncWithoutDetaching([$adminRole->id]);
    }
}
