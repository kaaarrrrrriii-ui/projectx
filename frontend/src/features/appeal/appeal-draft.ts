import type { CreatedTicket } from "@/shared/api/public-api";

export type DraftAnswer = {
  answerId: number;
  text: string;
};

export type AppealDraft = {
  role: string;
  categoryId?: number;
  categoryName?: string;
  customTopic?: string;
  description?: string;
  answers?: Record<string, DraftAnswer>;
};

const draftKey = "otklik:appeal-draft:v1";
const submissionKey = "otklik:last-submission:v1";

export function readAppealDraft(): AppealDraft {
  if (typeof window === "undefined") return { role: "student" };
  try {
    return JSON.parse(window.sessionStorage.getItem(draftKey) ?? "{}") as AppealDraft;
  } catch {
    return { role: "student" };
  }
}

export function updateAppealDraft(patch: Partial<AppealDraft>) {
  if (typeof window === "undefined") return;
  const current = readAppealDraft();
  window.sessionStorage.setItem(draftKey, JSON.stringify({ ...current, ...patch }));
}

export function startAppealDraft(role: string) {
  if (typeof window === "undefined") return;
  const current = readAppealDraft();
  if (current.role !== role) {
    window.sessionStorage.setItem(draftKey, JSON.stringify({ role }));
  } else {
    updateAppealDraft({ role });
  }
}

export function clearAppealDraft() {
  if (typeof window !== "undefined") window.sessionStorage.removeItem(draftKey);
}

export function saveSubmissionResult(result: CreatedTicket) {
  if (typeof window !== "undefined") window.sessionStorage.setItem(submissionKey, JSON.stringify(result));
}

export function readSubmissionResult(): CreatedTicket | null {
  if (typeof window === "undefined") return null;
  try {
    return JSON.parse(window.sessionStorage.getItem(submissionKey) ?? "null") as CreatedTicket | null;
  } catch {
    return null;
  }
}

const frontendCrisisMarkers = [
  "хочу умереть", "не хочу жить", "незачем жить", "нет смысла жить", "покончить с собой",
  "убить себя", "наложить на себя руки", "совершить суицид", "суицид", "самоубийств",
  "вскрыть вены", "порезать вены", "прыгнуть с крыши", "прыгнуть из окна", "повеситься",
  "отравиться таблетками", "таблетки чтобы умереть", "меня бьют", "меня избивают",
  "избивает меня", "избили меня", "угрожают убить",
  "угрожает убить", "угроза жизни", "меня убьют", "хочу убить", "физическое насилие",
  "убьет меня", "убьёт меня", "убью его", "убью ее", "убью её", "домашнее насилие",
  "сексуальное насилие", "изнасилова", "принуждают к сексу", "принуждает к сексу",
  "домогается", "держат силой", "удерживают силой", "похитили", "угрожает ножом",
  "угрожают ножом", "угрожает оружием", "угрожают оружием",
];

export function draftLooksCritical(draft: AppealDraft) {
  const answerText = Object.values(draft.answers ?? {}).map((answer) => answer.text).join("\n");
  const text = `${draft.customTopic ?? ""}\n${draft.description ?? ""}\n${answerText}`.toLocaleLowerCase("ru").replaceAll("ё", "е");
  return frontendCrisisMarkers.some((marker) => text.includes(marker.replaceAll("ё", "е")));
}
