const apiBase = (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080").replace(/\/$/, "");

const sessionKey = "otklik.staff.session";

export type StaffRole = "operator" | "expert" | "admin";

export type StaffUser = {
  id: number;
  username: string;
  full_name: string;
  role: StaffRole;
};

export type StaffSession = {
  token: string;
  expires_at: string;
  user: StaffUser;
};

type APIErrorBody = { error?: { code?: string; message?: string } };

export class APIError extends Error {
  constructor(public status: number, public code: string, message: string) {
    super(message);
  }
}

export function getStaffSession(): StaffSession | null {
  if (typeof window === "undefined") return null;
  try {
    const value = window.localStorage.getItem(sessionKey);
    if (!value) return null;
    const session = JSON.parse(value) as StaffSession;
    if (!session.token || !session.user?.role || Date.parse(session.expires_at) <= Date.now()) {
      window.localStorage.removeItem(sessionKey);
      return null;
    }
    return session;
  } catch {
    window.localStorage.removeItem(sessionKey);
    return null;
  }
}

export function saveStaffSession(session: StaffSession) {
  window.localStorage.setItem(sessionKey, JSON.stringify(session));
}

export function clearStaffSession() {
  if (typeof window !== "undefined") window.localStorage.removeItem(sessionKey);
}

async function decodeResponse<T>(response: Response): Promise<T> {
  if (response.ok) {
    if (response.status === 204) return undefined as T;
    return response.json() as Promise<T>;
  }
  let body: APIErrorBody = {};
  try {
    body = (await response.json()) as APIErrorBody;
  } catch {
    // Keep a useful generic message for non-JSON infrastructure errors.
  }
  throw new APIError(
    response.status,
    body.error?.code ?? "request_failed",
    body.error?.message ?? `Ошибка запроса (${response.status})`,
  );
}

export async function login(username: string, password: string): Promise<StaffSession> {
  const response = await fetch(`${apiBase}/api/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  const session = await decodeResponse<StaffSession>(response);
  saveStaffSession(session);
  return session;
}

export async function authenticatedRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  const session = getStaffSession();
  if (!session) {
    // eslint-disable-next-line @next/next/no-location-assign-relative-destination
    if (typeof window !== "undefined") window.location.assign("/staff");
    throw new APIError(401, "unauthorized", "Войдите в кабинет сотрудника");
  }
  const headers = new Headers(init.headers);
  headers.set("Authorization", `Bearer ${session.token}`);
  if (init.body && !(init.body instanceof FormData)) headers.set("Content-Type", "application/json");
  const response = await fetch(`${apiBase}${path}`, { ...init, headers, cache: "no-store" });
  if (response.status === 401 || response.status === 403) {
    clearStaffSession();
    // eslint-disable-next-line @next/next/no-location-assign-relative-destination
    if (typeof window !== "undefined") window.location.assign("/staff");
  }
  return decodeResponse<T>(response);
}

export type Category = { id: number; name: string };
export type Worker = {
  id: number;
  full_name: string;
  expert_group: { id: number; title: string };
  active_tickets: number;
  max_tickets: number;
  recommended: boolean;
  available: boolean;
};

export type OperatorTicket = {
  track_id: string;
  category: Category;
  status: string;
  applicant_type: string;
  priority: "urgent" | "standard" | "low";
  created_at: string;
  waiting_seconds: number;
  responsible?: { id: number; full_name: string };
  return_count: number;
  return_reason?: string;
  crisis_detected: boolean;
  is_overdue: boolean;
};

export type TicketPage<T> = { items: T[]; total: number; page: number; limit: number };

export type TicketDetail = {
  track_id: string;
  status: string;
  priority: "urgent" | "standard" | "low";
  created_at: string;
  category: Category;
  applicant_type: string;
  description: string;
  return_count: number;
  return_reason?: string;
  crisis_detected?: boolean;
  crisis_contact?: string;
  clarifications: Array<{ question_id: number; question: string; answer_id: number; answer: string }>;
  attachments: Array<{ id: number; name: string; mime_type: string; size_bytes: number }>;
  workers: Array<{ id: number; full_name: string; expert_group: string; is_responsible: boolean }>;
  messages?: Array<{ id: number; text: string; type: string; created_at: string }>;
  notes?: Array<{ id: number; text: string; created_at: string }>;
  allowed_actions?: string[];
};

export type WorkerRequest = {
  id: number;
  track_id: string;
  category: Category;
  applicant_type: string;
  priority: string;
  request_type: string;
  reason: string;
  status: string;
  created_at: string;
  created_by?: number;
};

export function getOperatorTickets(queue: "new" | "assigned" | "returned") {
  return authenticatedRequest<TicketPage<OperatorTicket>>(`/api/operator/tickets?queue=${queue}&limit=100`);
}

export function getOperatorDashboard() {
  return authenticatedRequest<{ new_count: number; assigned_count: number; returned_count: number; crisis_count: number; overdue_count: number }>("/api/operator/dashboard");
}

export function getOperatorTicket(trackID: string) {
  return authenticatedRequest<TicketDetail>(`/api/operator/tickets/${encodeURIComponent(trackID)}`);
}

export function getEligibleWorkers(trackID: string) {
  return authenticatedRequest<{ workers: Worker[] }>(`/api/operator/tickets/${encodeURIComponent(trackID)}/eligible-workers?only_available=false`);
}

export function assignResponsible(trackID: string, workerID: number) {
  return authenticatedRequest<{ worker_id: number; status: string }>(`/api/operator/tickets/${encodeURIComponent(trackID)}/responsible-worker`, {
    method: "PUT",
    body: JSON.stringify({ worker_id: workerID }),
  });
}

export function updateOperatorTicket(trackID: string, patch: { category_id?: number; priority?: string; status?: string; reason?: string }) {
  return authenticatedRequest<TicketDetail>(`/api/operator/tickets/${encodeURIComponent(trackID)}`, {
    method: "PATCH",
    body: JSON.stringify(patch),
  });
}

export function closeOperatorTicket(trackID: string, message: string) {
  return authenticatedRequest<{ status: string }>(`/api/operator/tickets/${encodeURIComponent(trackID)}/close`, {
    method: "POST",
    body: JSON.stringify({ message }),
  });
}

export function getOperatorRequests() {
  return authenticatedRequest<TicketPage<WorkerRequest>>("/api/operator/worker-requests?status=sent&page=1&limit=100");
}

export function completeOperatorRequest(requestID: number, workerID: number) {
  return authenticatedRequest<{ request_id: number; status: string; worker_id: number }>(`/api/operator/worker-requests/${requestID}/complete`, {
    method: "POST",
    body: JSON.stringify({ worker_id: workerID }),
  });
}

export type ExpertTicket = OperatorTicket & { is_responsible: boolean };

export function getExpertTickets(queue: "queue" | "assigned" | "returned") {
  return authenticatedRequest<TicketPage<ExpertTicket>>(`/api/expert/tickets?queue=${queue}&limit=100`);
}

export function openExpertTicket(trackID: string) {
  return authenticatedRequest<{ status: string }>(`/api/expert/tickets/${encodeURIComponent(trackID)}/open`, { method: "POST" });
}

export function getExpertTicket(trackID: string) {
  return authenticatedRequest<TicketDetail>(`/api/expert/tickets/${encodeURIComponent(trackID)}`);
}

export function createReplacementRequest(trackID: string, reason: string) {
  return authenticatedRequest<{ request: WorkerRequest }>(`/api/expert/tickets/${encodeURIComponent(trackID)}/requests`, {
    method: "POST",
    body: JSON.stringify({ request_type: "replace_responsible", reason }),
  });
}

export function getExpertRequests() {
  return authenticatedRequest<TicketPage<WorkerRequest>>("/api/expert/requests?page=1&limit=100");
}

export function sendExpertAnswer(trackID: string, text: string) {
  return authenticatedRequest<{ message: { id: number; text: string; type: string; created_at: string }; ticket_status: string }>(`/api/expert/tickets/${encodeURIComponent(trackID)}/messages`, {
    method: "POST",
    body: JSON.stringify({ text }),
  });
}

export type AdminCategory = {
  id: number;
  name: string;
  expert_groups: Array<{ id: number; title: string; employee_count: number; category_count: number }>;
  questions: Array<{ id: number; text: string; answers: Array<{ id: number; text: string }> }>;
};

export function getAdminCategories() {
  return authenticatedRequest<TicketPage<AdminCategory>>("/api/admin/categories?page=1&limit=100");
}

export function getAdminDashboard() {
  return authenticatedRequest<{ new_count: number; active_count: number; urgent_count: number; returned_count: number; routing_issue_count: number; experts_at_capacity_count: number }>("/api/admin/dashboard");
}

type AdminTicketRaw = {
  track_id: string;
  category: Category;
  status: string;
  applicant_type: string;
  priority: "urgent" | "standard" | "low";
  created_at: string;
  responsible?: { id: number; full_name: string };
  return_count: number;
  workers?: Array<{ id: number; full_name: string; expert_group: { id: number; title: string }; is_responsible: boolean }>;
};

function adminTicketToOperator(ticket: AdminTicketRaw): OperatorTicket {
  return {
    ...ticket,
    waiting_seconds: Math.max(0, Math.floor((Date.now() - Date.parse(ticket.created_at)) / 1000)),
    crisis_detected: false,
    is_overdue: false,
  };
}

export async function getAdminTickets() {
  const page = await authenticatedRequest<TicketPage<AdminTicketRaw>>("/api/admin/tickets?page=1&limit=100");
  return { ...page, items: page.items.map(adminTicketToOperator) };
}

export async function getAdminTicket(trackID: string): Promise<TicketDetail> {
  const ticket = await authenticatedRequest<AdminTicketRaw>(`/api/admin/tickets/${encodeURIComponent(trackID)}`);
  return {
    track_id: ticket.track_id,
    status: ticket.status,
    priority: ticket.priority,
    created_at: ticket.created_at,
    category: ticket.category,
    applicant_type: ticket.applicant_type,
    description: "Содержимое обращения скрыто от администратора",
    return_count: ticket.return_count,
    clarifications: [],
    attachments: [],
    workers: (ticket.workers ?? []).map((worker) => ({ id: worker.id, full_name: worker.full_name, expert_group: worker.expert_group.title, is_responsible: worker.is_responsible })),
  };
}

export async function getAdminExpertWorkers(): Promise<Worker[]> {
  const page = await authenticatedRequest<TicketPage<{ id: number; full_name: string; expert_group?: { id: number; title: string }; active_tickets: number; max_tickets: number; active: boolean }>>("/api/admin/employees?role=expert&page=1&limit=100");
  return page.items.map((employee) => ({ id: employee.id, full_name: employee.full_name, expert_group: employee.expert_group ?? { id: 0, title: "Эксперт" }, active_tickets: employee.active_tickets, max_tickets: employee.max_tickets, recommended: false, available: employee.active && employee.active_tickets < employee.max_tickets }));
}

export function assignAdminResponsible(trackID: string, workerID: number) {
  return authenticatedRequest<AdminTicketRaw>(`/api/admin/tickets/${encodeURIComponent(trackID)}/responsible-worker`, { method: "PUT", body: JSON.stringify({ worker_id: workerID }) });
}

export async function updateAdminTicket(trackID: string, patch: { priority?: string; status?: string }) {
  if (patch.priority) await authenticatedRequest<AdminTicketRaw>(`/api/admin/tickets/${encodeURIComponent(trackID)}/priority`, { method: "PATCH", body: JSON.stringify({ priority: patch.priority }) });
  if (patch.status) await authenticatedRequest<AdminTicketRaw>(`/api/admin/tickets/${encodeURIComponent(trackID)}/status`, { method: "PATCH", body: JSON.stringify({ status: patch.status }) });
  return getAdminTicket(trackID);
}

export function deleteAdminCategory(id: number) {
  return authenticatedRequest<void>(`/api/admin/categories/${id}`, { method: "DELETE" });
}
