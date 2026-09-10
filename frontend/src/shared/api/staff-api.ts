import { publicAPIURL } from "./public-api";

export type StaffRole = "operator" | "expert" | "admin";
export type StaffUser = { id: number; username: string; full_name: string; role: StaffRole };

const TOKEN_KEY = "otklik_staff_token";
const USER_KEY = "otklik_staff_user";

export function staffToken() {
  return typeof window === "undefined" ? "" : localStorage.getItem(TOKEN_KEY) ?? "";
}

export function currentStaffUser(): StaffUser | null {
  if (typeof window === "undefined") return null;
  try { return JSON.parse(localStorage.getItem(USER_KEY) ?? "null") as StaffUser | null; } catch { return null; }
}

export async function staffRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  const token = staffToken();
  const headers = new Headers(init.headers);
  if (token) headers.set("Authorization", `Bearer ${token}`);
  if (init.body && !(init.body instanceof FormData)) headers.set("Content-Type", "application/json");
  const response = await fetch(publicAPIURL(path), { ...init, headers, cache: "no-store" });
  if (response.status === 204) return undefined as T;
  const body = await response.json().catch(() => ({}));
  if (!response.ok) throw new Error(body?.error?.message ?? "Не удалось выполнить запрос");
  return body as T;
}

export async function loginStaff(username: string, password: string) {
  const result = await staffRequest<{ token: string; user: StaffUser }>("/api/auth/login", {
    method: "POST", body: JSON.stringify({ username, password }),
  });
  localStorage.setItem(TOKEN_KEY, result.token);
  localStorage.setItem(USER_KEY, JSON.stringify(result.user));
  return result.user;
}

export function clearStaffSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
}
