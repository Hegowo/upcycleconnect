<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Http\Requests\Admin\StoreEventRequest;
use App\Http\Resources\Admin\EventResource;
use App\Models\Event;
use App\Services\Admin\AuditService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;

class EventController extends Controller
{
    public function __construct(private readonly AuditService $audit) {}

    public function index(Request $request): AnonymousResourceCollection
    {
        $query = Event::with(['category', 'creator'])
            ->whereNull('deleted_at');

        if ($request->filled('status')) {
            $query->where('status', $request->status);
        }

        if ($request->filled('start_date_from')) {
            $query->where('start_date', '>=', $request->start_date_from);
        }

        if ($request->filled('start_date_to')) {
            $query->where('start_date', '<=', $request->start_date_to);
        }

        $events = $query->orderBy('start_date')->paginate($request->input('per_page', 20));

        return EventResource::collection($events);
    }

    public function store(StoreEventRequest $request): JsonResponse
    {
        $data              = $request->validated();
        $data['created_by'] = $request->user()->id;

        $event = Event::create($data);
        $event->load('category', 'creator');

        $this->audit->log('event.created', 'Event', $event->id, null, $event->toArray());

        return response()->json(new EventResource($event), 201);
    }

    public function show(Event $event): JsonResponse
    {
        $event->load('category', 'creator');

        return response()->json(new EventResource($event));
    }

    public function update(StoreEventRequest $request, Event $event): JsonResponse
    {
        $old = $event->toArray();
        $event->update($request->validated());
        $event->load('category', 'creator');

        $this->audit->log('event.updated', 'Event', $event->id, $old, $event->fresh()->toArray());

        return response()->json(new EventResource($event));
    }

    public function updateStatus(Request $request, Event $event): JsonResponse
    {
        $validated = $request->validate([
            'status' => ['required', 'in:draft,published,cancelled'],
        ]);

        $old = ['status' => $event->status];
        $event->update($validated);

        $this->audit->log('event.status_changed', 'Event', $event->id, $old, $validated);

        $event->load('category', 'creator');

        return response()->json(new EventResource($event));
    }

    public function destroy(Event $event): JsonResponse
    {
        $this->audit->log('event.deleted', 'Event', $event->id, $event->toArray(), null);
        $event->delete();

        return response()->json(null, 204);
    }
}
