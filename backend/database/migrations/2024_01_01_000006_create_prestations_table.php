<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('prestations', function (Blueprint $table) {
            $table->id();
            $table->foreignId('category_id')->nullable()->constrained('prestation_categories')->nullOnDelete();
            $table->foreignId('provider_id')->constrained('users')->cascadeOnDelete();
            $table->string('title', 255);
            $table->text('description');
            $table->decimal('price', 10, 2)->nullable();
            $table->enum('price_type', ['fixed', 'hourly', 'quote'])->default('fixed');
            $table->enum('status', ['draft', 'published', 'suspended', 'archived'])->default('draft');
            $table->timestamps();
            $table->softDeletes();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('prestations');
    }
};
