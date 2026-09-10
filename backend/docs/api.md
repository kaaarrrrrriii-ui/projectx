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

Operator endpoints allow employees with role `operator` only. Administrators
use the separate `/api/admin/*` endpoints documented below.

### Eligible experts

`GET /api/operator/tickets/{track_id}/eligible-workers`

Optional query parameters: `name`, `expert_group_id`, and `only_available`
(default `true`). When `only_available=false`, overloaded experts are returned
with `available=false`; assigning them is still prohibited. A category route is a
recommendation and does not prohibit assigning an expert from another group.

### Dashboard and ticket card

`GET /api/operator/dashboard` returns current new, distributed, returned,
crisis and overdue counts. A new ticket is overdue after
`OPERATOR_NEW_SLA_HOURS` (default 2); a distributed ticket without a specialist
answer is overdue after `OPERATOR_RESPONSE_SLA_HOURS` (default 24).

`GET /api/operator/tickets/{track_id}` returns the initial description,
clarifying answers, intake attachments, active workers, routing groups, expert
notes, crisis contact (only for a detected crisis), and audit events. It never
returns the applicant/expert chat. Operator attachment download is available at
`GET /api/operator/tickets/{track_id}/attachments/{attachment_id}`.

`PATCH /api/operator/tickets/{track_id}` atomically changes any supplied
`category_id`, `priority`, and `status`; `reason` is stored with audit events.
Category and priority changes are allowed while a ticket is `new` or `returned`.
Status changes follow the operator transition rules and `assigned` requires an
active responsible expert.

`POST /api/operator/tickets/{track_id}/close` accepts `{ "message": "..." }`.
For a `new` or `returned` ticket it stores the public operator message,
deactivates assignments and changes the status to `completed`.

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

Required query parameter: `queue=new`, `queue=assigned`, or `queue=returned`.

Optional parameters: `search`, `priority`, `status`, `category_id`,
`applicant_type`, `page`, and `limit`. The default limit is 20 and the maximum
is 100. `sort` may be `created_at_asc`, `created_at_desc`, `priority_desc`, or
`waiting_desc`. Items include `crisis_detected` and `is_overdue`. The default
new-queue order is crisis first, then urgent, then oldest.

### Analytics

`GET /api/operator/analytics?date_from=2026-09-01&date_to=2026-09-10`

Calendar dates are interpreted in `Asia/Yekaterinburg`; stored instants remain
UTC. The response contains ticket totals and distributions, average seconds to
first assignment, first specialist response and closure, aggregate expert
load, urgent share, return share, and per-operator action, assignment and close
counts.

### Anonymous report

`GET /api/operator/reports?date_from=2026-09-01&date_to=2026-09-10&format=csv`

`format` may be `csv` or `xlsx`. Reports are generated immediately and are not
stored. Each row contains only dates, category, applicant type, status,
priority, return count, and timing metrics. Track codes, employee data,
message text, contacts, and attachment names are excluded.

### Expert worker requests

`GET /api/operator/worker-requests?status=sent&page=1&limit=20` lists requests
created by experts. `status` may be `sent` or `completed`.

`POST /api/operator/worker-requests/{request_id}/complete` performs the action
and marks the request as completed:

```json
{ "worker_id": 7 }
```

For `add_coworker`, the worker is added as a co-worker. For
`replace_responsible`, the worker becomes the responsible expert. A request is
marked completed only after the assignment operation succeeds.

## Expert endpoints

Expert endpoints require a JWT for an employee with role `expert`. An expert
can access only tickets that have an active `tickets_workers` assignment for
that employee. The worker id is always taken from the JWT.

### Profile and dashboard

- `GET /api/expert/me` returns the expert group, rating, capacity, and current
  active ticket count.
- `GET /api/expert/dashboard` returns personal queue, in-progress, returned,
  urgent, and processed counts. A ticket is processed after the expert has
  sent at least one answer.

### Ticket lists

`GET /api/expert/tickets` requires one of these `queue` values:

- `queue` — assigned to this expert, not previously returned, and not yet
  opened;
- `assigned` — the expert's tickets in progress, awaiting clarification, or
  with an answer ready;
- `returned` — previously returned tickets reassigned by the operator and not
  yet opened by the expert.

Optional parameters are `search`, `priority`, `category_id`,
`applicant_type`, `page`, and `limit`. The default limit is 20 and the maximum
is 100.

### Ticket card and automatic start

`POST /api/expert/tickets/{track_id}/open` is called by the frontend when the
card opens. It idempotently changes `assigned` to `in_progress`. The UI does
not require a separate “start work” button.

`GET /api/expert/tickets/{track_id}` returns the ticket, initial description,
clarifying answers, public chat, applicant attachments, active workers,
internal expert notes, and allowed actions.

`GET /api/expert/tickets/{track_id}/attachments/{attachment_id}` downloads an
applicant attachment after checking the expert's active assignment. Experts
cannot attach files to their own answers.

### Notes and answers

`POST /api/expert/tickets/{track_id}/notes`:

```json
{ "text": "Внутренняя заметка" }
```

`POST /api/expert/tickets/{track_id}/messages`:

```json
{ "text": "Ответ заявителю" }
```

The answer is saved as a `specialist` message and atomically moves the ticket
to `answer_ready`. Internal notes are never included in the applicant chat.

### Requests to the operator

`POST /api/expert/tickets/{track_id}/requests`:

```json
{
  "request_type": "add_coworker",
  "reason": "Нужна помощь коллеги"
}
```

`request_type` is `add_coworker` or `replace_responsible`. The operator chooses
the new worker. `GET /api/expert/requests?status=sent&page=1&limit=20` returns
the current expert's request history. Experts see only the statuses `sent` and
`completed`.

The database schema is intentionally unchanged. Requests, status events, and
expert notes use `messages.type = 3` (`internal_note`) with versioned JSON in
`messages.text`. The initial message id is the request id; completion is an
additional internal message that references that id.

### Personal analytics and report

- `GET /api/expert/analytics?date_from=2026-09-01&date_to=2026-09-10`
- `GET /api/expert/reports?date_from=2026-09-01&date_to=2026-09-10&format=csv`

The response schema and calculations match the operator analytics and report,
but only tickets assigned to the current expert are included. Reports are
generated immediately and are not stored. `format` may be `csv` or `xlsx`.

## Persisted enum values

- Priority: `1 = standard`, `2 = urgent`, `3 = low`.
- Applicant type: `1 = schoolchild`, `2 = parent`, `3 = teacher`.
- Message type: `1 = applicant`, `2 = specialist`, `3 = internal_note`,
  `4 = system`, `5 = return_reason`, `6 = operator`.

Ticket statuses remain numbered from 1 through 9 and are exposed through the
API as strings such as `new`, `assigned`, `answer_ready`, and `rejected`.

## Administrator endpoints

Administrator endpoints require a JWT whose current database role is exactly
`admin`. Administrator and operator permissions are intentionally separate:
an administrator cannot use `/api/operator/*` endpoints. Admin ticket
responses contain metadata and audit events only; they never contain message
text, clarifying answers, notes, contacts, or attachments.

### Profile and dashboard

- `GET /api/admin/me`
- `GET /api/admin/dashboard`

The dashboard returns counts for new, active, urgent and returned tickets,
routing issues, and experts whose active workload reached their limit.

### Employees

- `GET /api/admin/employees?role=expert&search=&page=1&limit=20`
- `POST /api/admin/employees`
- `PATCH /api/admin/employees/{employee_id}`
- `POST /api/admin/employees/{employee_id}/deactivate`
- `POST /api/admin/employees/{employee_id}/restore`

Creation example:

```json
{
  "username": "e.ivanova",
  "password": "temporary-password",
  "full_name": "Иванова Е.С.",
  "role": "expert",
  "expert_group_id": 2,
  "max_tickets": 10
}
```

Only `operator` and `expert` accounts can be created or edited. Because the
current schema requires `users.expert_group_id` for every user, an operator
also receives an administrator-selected operator group, although that group
is hidden in operator employee responses.

Deactivation stores the reserved role value `200`. Such a user cannot log in,
and an already issued JWT stops working because its role no longer matches the
database. An expert with active tickets cannot be deactivated. Restoration
requires the target role and group because the schema does not store the old
role separately.

### Expert groups, categories, and questions

- `GET|POST /api/admin/expert-groups`
- `PATCH|DELETE /api/admin/expert-groups/{group_id}`
- `GET|POST /api/admin/categories`
- `GET|PATCH|DELETE /api/admin/categories/{category_id}`
- `PUT /api/admin/categories/{category_id}/expert-groups`
- `POST /api/admin/categories/{category_id}/questions`
- `PATCH|DELETE /api/admin/questions/{question_id}`

Routing replacement body:

```json
{ "expert_group_ids": [1, 3] }
```

Each question belongs to exactly one category and is created with a non-empty
array of unique answer strings. A question referenced by `QA` cannot be edited
or deleted because doing so would change historical answers. Referenced
groups and categories cannot be deleted. The required category
`Не знаю, как это назвать` cannot be renamed or deleted.

### Ticket supervision

- `GET /api/admin/tickets`
- `GET /api/admin/tickets/{track_id}`
- `PATCH /api/admin/tickets/{track_id}/priority` with `{ "priority": "urgent" }`
- `PATCH /api/admin/tickets/{track_id}/status` with `{ "status": "in_progress" }`
- `PUT /api/admin/tickets/{track_id}/responsible-worker` with `{ "worker_id": 7 }`
- `GET /api/admin/tickets/{track_id}/events?page=1&limit=20`

Ticket filters are `search`, `status`, `priority`, `category_id`,
`applicant_type`, `responsible_worker_id`, `routing_issue`, `page`, and
`limit`. `routing_issue` is `no_group` or `no_available_expert`.

An administrator may set any persisted ticket status. Terminal statuses set
`closed_at` and deactivate assignments. `new` and `returned` also deactivate
assignments. Active work statuses require an active responsible expert.
Reassigning a responsible expert changes the status to `assigned`. Every
actual admin change is recorded in `ticket_events`; no reason is required.

### Analytics and report

- `GET /api/admin/analytics?date_from=2026-09-01&date_to=2026-09-10`
- `GET /api/admin/reports?date_from=2026-09-01&date_to=2026-09-10&format=csv`

Admin analytics reuse the system-wide calculations and add an employee load
list. For experts, load is active tickets divided by `max_tickets`; operator
`action_count` is the number of their audit events in the selected period.
Reports are the same anonymous in-memory CSV/XLSX files as operator reports.
They are not stored.
