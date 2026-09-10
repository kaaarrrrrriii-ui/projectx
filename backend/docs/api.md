# Backend API

The backend listens on `http://localhost:8080` by default. JSON responses use
UTF-8 and timestamps are returned in RFC 3339 format.

## Authentication

Employee endpoints use a JWT bearer token. The server requires a
`JWT_SECRET` environment variable containing at least 32 characters.

### Login

`POST /api/auth/login`

```json
{
  "username": "operator",
  "password": "password"
}
```

The response contains `token`, `expires_at`, and the authenticated `user`.
Send the token to protected endpoints:

```http
Authorization: Bearer <token>
```

`GET /api/auth/me` returns the current employee. `POST /api/auth/logout`
validates the token and returns `204 No Content`; because JWT is stateless, the
client must then discard the token.

## Operator endpoints

Operator endpoints allow employees with role `operator` or `admin`.

### Eligible experts

`GET /api/operator/tickets/{track_id}/eligible-workers`

Optional query parameters: `name`, `expert_group_id`. Experts whose active
ticket count has reached `max_tickets` are omitted. A category route is a
recommendation and does not prohibit assigning an expert from another group.

### Responsible expert

`PUT /api/operator/tickets/{track_id}/responsible-worker`

```json
{ "worker_id": 7 }
```

The previous responsible assignment is deactivated, co-workers remain active,
and the ticket status becomes `assigned`. The same operation is used to send a
returned ticket back for another round of work.

### Co-workers

- `POST /api/operator/tickets/{track_id}/workers` with `{ "worker_id": 7 }`.
- `DELETE /api/operator/tickets/{track_id}/workers/{worker_id}`.

The responsible expert cannot be removed through the co-worker endpoint.

### Reject or close by operator

`POST /api/operator/tickets/{track_id}/reject`

```json
{
  "reason_code": "spam",
  "message": "Объяснение для заявителя"
}
```

Allowed reason codes are `no_expert_help`, `spam`, and `out_of_scope`. The
ticket becomes `rejected`, `closed_at` is set, assignments are deactivated,
and the public status response includes `resolution_message`.

### Assigned and returned queues

`GET /api/operator/tickets`

Required query parameter: `queue=assigned` or `queue=returned`.

Optional parameters: `search`, `priority`, `status`, `category_id`,
`applicant_type`, `page`, and `limit`. The default limit is 20 and the maximum
is 100. Overdue calculation is intentionally not part of this endpoint.

### Analytics

`GET /api/operator/analytics?date_from=2026-09-01&date_to=2026-09-10`

Calendar dates are interpreted in `Asia/Yekaterinburg`; stored instants remain
UTC. The response contains ticket totals and distributions, average seconds to
first assignment, first specialist response and closure, aggregate expert
load, urgent share, and return share. Operator load is not calculated.

### Anonymous report

`GET /api/operator/reports?date_from=2026-09-01&date_to=2026-09-10&format=csv`

`format` may be `csv` or `xlsx`. Reports are generated immediately and are not
stored. Each row contains only dates, category, applicant type, status,
priority, return count, and timing metrics. Track codes, employee data,
message text, contacts, and attachment names are excluded.

## Persisted enum values

- Priority: `1 = standard`, `2 = urgent`.
- Applicant type: `1 = schoolchild`, `2 = parent`, `3 = teacher`.
- Message type: `1 = applicant`, `2 = specialist`, `3 = internal_note`,
  `4 = system`, `5 = return_reason`, `6 = operator`.

Ticket statuses remain numbered from 1 through 9 and are exposed through the
API as strings such as `new`, `assigned`, `answer_ready`, and `rejected`.
