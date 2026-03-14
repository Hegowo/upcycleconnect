<?php

namespace Database\Seeders;

use App\Models\Event;
use App\Models\Prestation;
use App\Models\PrestationCategory;
use App\Models\ProviderProfile;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;

class DemoDataSeeder extends Seeder
{
    public function run(): void
    {
        $superAdmin = User::where('email', 'superadmin@upcycleconnect.fr')->first();

        // Catégories
        $categories = [
            ['name' => 'Mobilier',      'slug' => 'mobilier',      'description' => 'Réparation et customisation de mobilier',      'is_active' => true,  'sort_order' => 1],
            ['name' => 'Textile',       'slug' => 'textile',       'description' => 'Upcycling de vêtements et tissu',              'is_active' => true,  'sort_order' => 2],
            ['name' => 'Électronique',  'slug' => 'electronique',  'description' => 'Réparation et reconditionnement électronique', 'is_active' => true,  'sort_order' => 3],
            ['name' => 'Décoration',    'slug' => 'decoration',    'description' => 'Objets de décoration recyclés',                'is_active' => false, 'sort_order' => 4],
        ];

        $createdCategories = [];
        foreach ($categories as $cat) {
            $createdCategories[] = PrestationCategory::firstOrCreate(['slug' => $cat['slug']], $cat);
        }

        // Utilisateurs particuliers
        $particuliers = [];
        $particulierData = [
            ['email' => 'marie.dupont@example.fr',   'first_name' => 'Marie',    'last_name' => 'Dupont'],
            ['email' => 'pierre.martin@example.fr',  'first_name' => 'Pierre',   'last_name' => 'Martin'],
            ['email' => 'sophie.bernard@example.fr', 'first_name' => 'Sophie',   'last_name' => 'Bernard'],
            ['email' => 'lucas.petit@example.fr',    'first_name' => 'Lucas',    'last_name' => 'Petit'],
            ['email' => 'emma.robert@example.fr',    'first_name' => 'Emma',     'last_name' => 'Robert'],
            ['email' => 'hugo.thomas@example.fr',    'first_name' => 'Hugo',     'last_name' => 'Thomas'],
            ['email' => 'lea.richard@example.fr',    'first_name' => 'Léa',      'last_name' => 'Richard'],
            ['email' => 'alex.durand@example.fr',    'first_name' => 'Alex',     'last_name' => 'Durand'],
            ['email' => 'chloe.moreau@example.fr',   'first_name' => 'Chloé',    'last_name' => 'Moreau', 'status' => 'inactive'],
            ['email' => 'tom.simon@example.fr',      'first_name' => 'Tom',      'last_name' => 'Simon',  'status' => 'banned'],
        ];

        foreach ($particulierData as $data) {
            $particuliers[] = User::firstOrCreate(
                ['email' => $data['email']],
                array_merge($data, [
                    'password'          => Hash::make('User1234!'),
                    'status'            => $data['status'] ?? 'active',
                    'email_verified_at' => now(),
                ])
            );
        }

        // Prestataires
        $prestataireData = [
            ['email' => 'atelier.bois@example.fr',   'first_name' => 'Atelier',  'last_name' => 'Bois',      'company' => 'Atelier du Bois',     'status' => 'approved'],
            ['email' => 'couture.verde@example.fr',  'first_name' => 'Verde',    'last_name' => 'Couture',   'company' => 'Verde Couture',        'status' => 'approved'],
            ['email' => 'repare.tech@example.fr',    'first_name' => 'TechRep',  'last_name' => 'Services',  'company' => 'TechRep Services',     'status' => 'approved'],
            ['email' => 'new.artisan@example.fr',    'first_name' => 'Nouveau',  'last_name' => 'Artisan',   'company' => 'Artisan Créatif',      'status' => 'pending'],
            ['email' => 'eco.design@example.fr',     'first_name' => 'Eco',      'last_name' => 'Design',    'company' => 'Eco Design Studio',   'status' => 'pending'],
        ];

        $providers = [];
        foreach ($prestataireData as $data) {
            $user = User::firstOrCreate(
                ['email' => $data['email']],
                [
                    'password'          => Hash::make('Provider1234!'),
                    'first_name'        => $data['first_name'],
                    'last_name'         => $data['last_name'],
                    'status'            => 'active',
                    'email_verified_at' => now(),
                ]
            );

            ProviderProfile::firstOrCreate(
                ['user_id' => $user->id],
                [
                    'company_name' => $data['company'],
                    'siret'        => '12345678901234',
                    'description'  => 'Prestataire spécialisé en upcycling.',
                    'status'       => $data['status'],
                    'approved_at'  => $data['status'] === 'approved' ? now() : null,
                    'approved_by'  => $data['status'] === 'approved' ? $superAdmin?->id : null,
                ]
            );

            $providers[] = $user;
        }

        // Prestations
        $prestationsData = [
            ['provider' => 0, 'category' => 0, 'title' => 'Restauration de chaise ancienne',   'price' => 120, 'price_type' => 'fixed',  'status' => 'published'],
            ['provider' => 0, 'category' => 0, 'title' => 'Réparation table basse',             'price' => 80,  'price_type' => 'fixed',  'status' => 'published'],
            ['provider' => 1, 'category' => 1, 'title' => 'Transformation jean en sac',         'price' => 45,  'price_type' => 'fixed',  'status' => 'published'],
            ['provider' => 1, 'category' => 1, 'title' => 'Customisation veste vintage',        'price' => 60,  'price_type' => 'hourly', 'status' => 'published'],
            ['provider' => 2, 'category' => 2, 'title' => 'Réparation smartphone',              'price' => 50,  'price_type' => 'fixed',  'status' => 'published'],
            ['provider' => 2, 'category' => 2, 'title' => 'Reconditionnement PC portable',      'price' => null,'price_type' => 'quote',  'status' => 'draft'],
            ['provider' => 3, 'category' => 3, 'title' => 'Création lampe à partir de bouteilles', 'price' => 90, 'price_type' => 'fixed', 'status' => 'draft'],
            ['provider' => 0, 'category' => 0, 'title' => 'Vernissage de meuble',               'price' => 40,  'price_type' => 'hourly', 'status' => 'suspended'],
        ];

        foreach ($prestationsData as $data) {
            Prestation::firstOrCreate(
                ['title' => $data['title'], 'provider_id' => $providers[$data['provider']]->id],
                [
                    'category_id' => $createdCategories[$data['category']]->id,
                    'provider_id' => $providers[$data['provider']]->id,
                    'description' => 'Description complète de la prestation d\'upcycling proposée par le prestataire.',
                    'price'       => $data['price'],
                    'price_type'  => $data['price_type'],
                    'status'      => $data['status'],
                ]
            );
        }

        // Événements
        $eventsData = [
            [
                'title'            => 'Atelier upcycling textile',
                'description'      => 'Atelier pratique pour apprendre à transformer vos vieux vêtements en nouvelles pièces tendance.',
                'location'         => 'Paris 11e — Maison des associations',
                'start_date'       => now()->addDays(15)->setTime(10, 0),
                'end_date'         => now()->addDays(15)->setTime(13, 0),
                'max_participants' => 15,
                'status'           => 'published',
                'category'         => 1,
            ],
            [
                'title'            => 'Repair Café Mobilier',
                'description'      => 'Venez avec votre mobilier abîmé, nos artisans vous aident à le réparer sur place.',
                'location'         => 'Lyon 7e — Centre communautaire',
                'start_date'       => now()->addDays(30)->setTime(9, 0),
                'end_date'         => now()->addDays(30)->setTime(17, 0),
                'max_participants' => 30,
                'status'           => 'published',
                'category'         => 0,
            ],
            [
                'title'            => 'Conférence : Économie circulaire et upcycling',
                'description'      => 'Conférence sur les enjeux économiques et environnementaux de l\'upcycling.',
                'location'         => 'Bordeaux — CCI Bordeaux Gironde',
                'start_date'       => now()->addDays(60)->setTime(14, 0),
                'end_date'         => now()->addDays(60)->setTime(17, 0),
                'max_participants' => null,
                'status'           => 'draft',
                'category'         => null,
            ],
            [
                'title'            => 'Salon du réemploi',
                'description'      => 'Événement annulé — report à une date ultérieure.',
                'location'         => 'Nantes',
                'start_date'       => now()->subDays(10)->setTime(9, 0),
                'end_date'         => now()->subDays(10)->setTime(18, 0),
                'max_participants' => 200,
                'status'           => 'cancelled',
                'category'         => null,
            ],
        ];

        foreach ($eventsData as $data) {
            Event::firstOrCreate(
                ['title' => $data['title']],
                [
                    'category_id'      => $data['category'] !== null ? $createdCategories[$data['category']]->id : null,
                    'title'            => $data['title'],
                    'description'      => $data['description'],
                    'location'         => $data['location'],
                    'start_date'       => $data['start_date'],
                    'end_date'         => $data['end_date'],
                    'max_participants' => $data['max_participants'],
                    'status'           => $data['status'],
                    'created_by'       => $superAdmin?->id ?? 1,
                ]
            );
        }
    }
}
