<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\SoftDeletes;

class Event extends Model
{
    use SoftDeletes;

    protected $table = 'platform_events';

    protected $fillable = [
        'category_id',
        'title',
        'description',
        'location',
        'start_date',
        'end_date',
        'max_participants',
        'status',
        'created_by',
    ];

    protected function casts(): array
    {
        return [
            'start_date' => 'datetime',
            'end_date'   => 'datetime',
        ];
    }

    public function category(): BelongsTo
    {
        return $this->belongsTo(PrestationCategory::class, 'category_id');
    }

    public function creator(): BelongsTo
    {
        return $this->belongsTo(User::class, 'created_by');
    }
}
