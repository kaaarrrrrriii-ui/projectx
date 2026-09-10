export type PublicCategory = {
  id: number;
  name: string;
};

export type PublicAnswer = {
  id: number;
  text: string;
};

export type PublicQuestion = {
  id: number;
  text: string;
  answers: PublicAnswer[];
};

export type CrisisContact = {
  title: string;
  phone: string;
  description: string;
  url?: string;
};

export type CreatedTicket = {
  track_id: string;
  status: string;
  created_at: string;
  crisis_detected: boolean;
  crisis_contacts: CrisisContact[];
};

export type TicketStatus = {
  track_id: string;
  status: string;
  created_at: string;
  category: PublicCategory;
  can_open_chat: boolean;
  can_complete: boolean;
  can_return: boolean;
  resolution_message?: string;
};

export type TicketAttachment = {
  id: number;
  name: string;
  mime_type: string;
  size_bytes: number;
};

export type TicketMessage = {
  id: number;
  text: string;
  type: "applicant" | "specialist" | string;
  created_at: string;
  attachments: TicketAttachment[];
};

export type TicketChat = {
  track_id: string;
  status: string;
  specialist?: {
    label: string;
    expert_group: string;
  };
  messages: TicketMessage[];
  attachments: TicketAttachment[];
};

export type CreatedTicketMessage = {
  message: TicketMessage;
  ticket_status: string;
};

type APIErrorBody = {
  error?: {
    code?: string;
    message?: string;
  };
};

export class PublicAPIError extends Error {
  constructor(
    message: string,
    public readonly status: number,
    public readonly code?: string,
  ) {
    super(message);
  }
}

export function publicAPIURL(path: string) {
  const base = (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080").replace(/\/$/, "");
  return `${base}${path}`;
}

async function parseResponse<T>(response: Response): Promise<T> {
  if (response.ok) return response.json() as Promise<T>;

  let body: APIErrorBody = {};
  try {
    body = (await response.json()) as APIErrorBody;
  } catch {
    // Keep the human-friendly fallback below for non-JSON proxy errors.
  }
  throw new PublicAPIError(
    body.error?.message ?? "Не удалось связаться с сервисом. Попробуйте ещё раз.",
    response.status,
    body.error?.code,
  );
}

export async function getPublicCategories(signal?: AbortSignal) {
  const response = await fetch(publicAPIURL("/api/categories"), { signal, cache: "no-store" });
  return parseResponse<{ categories: PublicCategory[] }>(response);
}

export async function getPublicQuestions(categoryID: number, signal?: AbortSignal) {
  const response = await fetch(publicAPIURL(`/api/categories/${categoryID}/questions`), { signal, cache: "no-store" });
  return parseResponse<{ questions: PublicQuestion[] }>(response);
}

export async function createPublicTicket(form: FormData) {
  const response = await fetch(publicAPIURL("/api/tickets"), { method: "POST", body: form });
  return parseResponse<CreatedTicket>(response);
}

export async function getTicketStatus(trackID: string) {
  const response = await fetch(publicAPIURL(`/api/tickets/${encodeURIComponent(trackID)}/status`), { cache: "no-store" });
  return parseResponse<TicketStatus>(response);
}

export async function getTicketChat(trackID: string, signal?: AbortSignal) {
  const response = await fetch(publicAPIURL(`/api/tickets/${encodeURIComponent(trackID)}/chat`), {
    signal,
    cache: "no-store",
  });
  return parseResponse<TicketChat>(response);
}

export async function postTicketMessage(trackID: string, text: string, attachments: File[]) {
  const form = new FormData();
  form.set("text", text);
  for (const attachment of attachments) form.append("attachments[]", attachment);
  const response = await fetch(publicAPIURL(`/api/tickets/${encodeURIComponent(trackID)}/messages`), {
    method: "POST",
    body: form,
  });
  return parseResponse<CreatedTicketMessage>(response);
}

export function ticketAttachmentURL(trackID: string, attachmentID: number) {
  return publicAPIURL(`/api/tickets/${encodeURIComponent(trackID)}/attachments/${attachmentID}`);
}

export async function completeTicket(trackID: string) {
  const response = await fetch(publicAPIURL(`/api/tickets/${encodeURIComponent(trackID)}/complete`), { method: "POST" });
  return parseResponse<{ status: string; closed_at: string }>(response);
}

export async function returnTicket(trackID: string, reason: string) {
  const response = await fetch(publicAPIURL(`/api/tickets/${encodeURIComponent(trackID)}/return`), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ reason }),
  });
  return parseResponse<{ status: string; return_count: number }>(response);
}

export async function createTicketReview(trackID: string, rating: number, text = "") {
  const response = await fetch(publicAPIURL(`/api/tickets/${encodeURIComponent(trackID)}/review`), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ rating, text }),
  });
  return parseResponse<{ id: number; rating: number; text: string }>(response);
}
