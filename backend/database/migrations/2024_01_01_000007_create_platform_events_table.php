<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('platform_events', function (Blueprint $table) {
            $table->id();
            $table->foreignId('category_id')->nullable()->constrained('prestation_categories')->nullOnDelete();
            $table->string('title', 255);
            $table->text('description');
            $table->string('location', 255)->nullable();
            $table->dateTime('start_date');
            $table->dateTime('end_date');
            $table->unsignedInteger('max_participants')->nullable();
            $table->enum('status', ['draft', 'published', 'cancelled'])->default('draft');
            $table->foreignId('created_by')->constrained('users')->cascadeOnDelete();
            $table->timestamps();
            $table->softDeletes();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('platform_events');
    }
};
