"use client";

import Button from "@/shared/ui/button";
import Link from "next/link";
import { useEffect, useState } from "react";
import {
  assignResponsible,
  assignAdminResponsible,
  closeOperatorTicket,
  getAdminExpertWorkers,
  getAdminTicket,
  getEligibleWorkers,
  getOperatorTicket,
  updateOperatorTicket,
  updateAdminTicket,
  type TicketDetail as TicketDetailData,
  type Worker,
} from "@/shared/api/staff-api";
import { CloseTicketModal } from "./ticket-detail-modals";

const priorities = [
  { value: "urgent" as const, label: "Срочное" },
  { value: "standard" as const, label: "Стандарт" },
  { value: "low" as const, label: "Вопрос" },
];

const statuses = [
  { value: "new", label: "Новое" },
  { value: "assigned", label: "Распределено" },
  { value: "in_progress", label: "В работе" },
  { value: "needs_clarification", label: "Нужно уточнение" },
  { value: "answer_ready", label: "Ответ готов" },
  { value: "returned", label: "Возвращено" },
  { value: "rejected", label: "Отклонено" },
  { value: "completed", label: "Завершено" },
  { value: "closed_without_answer", label: "Закрыто" },
];

const applicantLabels: Record<string, string> = {
  schoolchild: "Школьник",
  parent: "Родитель",
  teacher: "Педагог",
};

const colors = {
  urgent: { text: "text-[#e5141b]", soft: "bg-[#f7b6b8] text-[#d70d14]" },
  standard: { text: "text-[#4562f0]", soft: "bg-[#dfe6ff] text-[#4562f0]" },
  low: { text: "text-[#087f1a]", soft: "bg-[#dff2e0] text-[#087f1a]" },
};

function AssignmentDialog({ trackID, mode, onClose, onAssigned }: { trackID: string; mode: "operator" | "admin"; onClose: () => void; onAssigned: () => void }) {
  const [workers, setWorkers] = useState<Worker[]>([]);
  const [selected, setSelected] = useState<number | null>(null);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    (mode === "admin" ? getAdminExpertWorkers() : getEligibleWorkers(trackID).then((result) => result.workers))
      .then(setWorkers)
      .catch(() => setError("Не удалось загрузить список специалистов"));
  }, [mode, trackID]);

  async function assign() {
    if (!selected) return;
    setPending(true);
    setError("");
    try {
      if (mode === "admin") await assignAdminResponsible(trackID, selected);
      else await assignResponsible(trackID, selected);
      onAssigned();
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Не удалось назначить специалиста");
    } finally {
      setPending(false);
    }
  }

  return (
    <div className="fixed inset-0 z-[80] grid place-items-center bg-[#000828]/45 p-4" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget) onClose(); }}>
      <section role="dialog" aria-modal="true" aria-labelledby="assignment-heading" className="w-full max-w-2xl rounded-2xl bg-white p-6 shadow-2xl">
        <div className="flex items-center justify-between gap-4">
          <h2 id="assignment-heading" className="text-xl font-extrabold text-[#4562f0]">Назначение исполнителя</h2>
          <button type="button" onClick={onClose} className="h-9 w-9 cursor-pointer rounded-full text-2xl text-[#4562f0] hover:bg-[#eef1ff]">×</button>
        </div>
        <div className="mt-5 grid max-h-[360px] gap-2 overflow-y-auto">
          {workers.map((worker) => (
            <button key={worker.id} type="button" disabled={!worker.available} onClick={() => setSelected(worker.id)} className={`flex cursor-pointer items-center justify-between rounded-xl border px-4 py-3 text-left text-sm ${selected === worker.id ? "border-[#4562f0] bg-[#dfe6ff]" : "border-[#c7cee8] bg-white hover:bg-[#f7f8ff]"} disabled:cursor-not-allowed disabled:opacity-45`}>
              <span><strong className="block text-[#000828]">{worker.full_name}</strong><span className="text-xs text-[#646d86]">{worker.expert_group.title}{worker.recommended ? " · рекомендован" : ""}</span></span>
              <span className="text-xs text-[#646d86]">{worker.active_tickets}/{worker.max_tickets}</span>
            </button>
          ))}
          {!error && workers.length === 0 && <p className="py-6 text-center text-sm text-[#646d86]">Загружаем специалистов…</p>}
        </div>
        {error && <p role="alert" className="mt-4 text-sm text-[#b42318]">{error}</p>}
        <Button text={pending ? "Назначаем…" : "Назначить выбранного специалиста"} disabled={!selected || pending} onClick={assign} className="mt-5 w-full" />
      </section>
    </div>
  );
}

export default function TicketDetail({ trackID, returnBasePath = "/operator/queueNew.tsx", mode = "operator" }: { trackID: string; returnBasePath?: string; mode?: "operator" | "admin" }) {
  const [ticket, setTicket] = useState<TicketDetailData | null>(null);
  const [priority, setPriority] = useState<TicketDetailData["priority"]>("standard");
  const [status, setStatus] = useState("new");
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);
  const [assignmentOpen, setAssignmentOpen] = useState(false);
  const [closeOpen, setCloseOpen] = useState(false);

  async function load() {
    setError("");
    try {
      const value = mode === "admin" ? await getAdminTicket(trackID) : await getOperatorTicket(trackID);
      setTicket(value);
      setPriority(value.priority);
      setStatus(value.status);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Не удалось загрузить обращение");
    }
  }

  // eslint-disable-next-line react-hooks/set-state-in-effect, react-hooks/exhaustive-deps
  useEffect(() => { void load(); }, [trackID]);

  async function save() {
    if (!ticket) return;
    const patch: { priority?: string; status?: string } = {};
    if (priority !== ticket.priority) patch.priority = priority;
    if (status !== ticket.status) patch.status = status;
    if (!patch.priority && !patch.status) {
      setNotice("Изменений нет");
      return;
    }
    setPending(true);
    setError("");
    try {
      const updated = mode === "admin" ? await updateAdminTicket(trackID, patch) : await updateOperatorTicket(trackID, patch);
      setTicket(updated);
      setPriority(updated.priority);
      setStatus(updated.status);
      setNotice("Изменения сохранены");
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Не удалось сохранить изменения");
    } finally {
      setPending(false);
    }
  }

  async function closeTicket(message: string) {
    setPending(true);
    setError("");
    try {
      await closeOperatorTicket(trackID, message);
      setCloseOpen(false);
      setNotice("Обращение закрыто");
      await load();
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Не удалось закрыть обращение");
    } finally {
      setPending(false);
    }
  }

  if (!ticket) return <section className="min-w-0 flex-1 p-8 text-sm text-[#646d86]">{error || "Загружаем обращение…"}</section>;
  const accent = colors[priority];
  const responsible = ticket.workers.find((worker) => worker.is_responsible);

  return (
    <div className="grid min-w-0 flex-1 grid-cols-[minmax(0,1fr)_268px] max-[899px]:grid-cols-1">
      <article className="flex min-w-0 flex-col px-3 py-4 sm:px-4">
        <Link href={returnBasePath} className="w-fit rounded-sm text-sm text-[#85899b] hover:text-[#4562f0]">Вернуться назад</Link>
        <header className="mt-6 flex flex-wrap items-center gap-3">
          <h1 className={`text-[clamp(20px,2.3vw,25px)] leading-tight font-extrabold ${accent.text}`}>Обращение №{ticket.track_id}</h1>
          <span className={`rounded-full px-4 py-1.5 text-[10px] ${accent.soft}`}>{priorities.find((item) => item.value === priority)?.label}</span>
          {ticket.crisis_detected && <span className="rounded-full bg-[#fff1f1] px-3 py-1 text-[10px] font-semibold text-[#b42318]">Кризисное</span>}
        </header>
        <p className="mt-1 text-sm font-medium text-[#151515]">Дата: {new Date(ticket.created_at).toLocaleString("ru-RU")}</p>
        <section className="mt-7"><h2 className="text-base font-medium text-[#151515]">Тип заявителя: {applicantLabels[ticket.applicant_type] ?? ticket.applicant_type}</h2><p className="mt-1 text-sm text-[#646d86]">Категория: {ticket.category.name}</p></section>
        <section className="mt-6"><h2 className="text-base font-medium text-[#151515]">Исходный текст</h2><div className="mt-2 min-h-[130px] rounded-[11px] border border-[#333] bg-white/75 px-3 py-2 text-xs leading-5 text-[#30384f]">{ticket.description}</div></section>
        <section className="mt-6"><h2 className="text-base font-medium text-[#151515]">Ответы на уточняющие вопросы:</h2><dl className="mt-2 grid gap-1 text-xs text-[#3f475d]">{ticket.clarifications.map((item) => <div key={item.question_id} className="flex flex-wrap gap-1"><dt>{item.question}</dt><dd className="font-medium text-[#151515]">{item.answer}</dd></div>)}</dl></section>
        <section className="mt-6"><h2 className="text-base font-medium text-[#151515]">Прикреплённые файлы:</h2>{ticket.attachments.length ? <ul className="mt-2 flex flex-wrap gap-2">{ticket.attachments.map((file) => <li key={file.id} className="text-xs text-[#30384f]">{file.name}</li>)}</ul> : <p className="mt-2 text-xs text-[#7b849b]">Файлы не приложены</p>}</section>
      </article>

      <aside aria-labelledby="edit-ticket-heading" className="flex min-h-[calc(100dvh-100px)] flex-col border border-[#4562f0] bg-white/70 px-3.5 py-5 max-[899px]:min-h-0 max-[899px]:border-r-0 max-[899px]:border-b-0 max-[899px]:border-l-0">
        <h2 id="edit-ticket-heading" className="text-sm font-medium text-[#151515]">Редактировать обращение</h2>
        <section className="mt-3 rounded-[12px] border border-[#4562f0] bg-white/80 p-2.5">
          <fieldset><legend className="text-[11px] text-[#30384f]">Приоритет</legend><div className="mt-2 grid grid-cols-3 gap-1">{priorities.map((item) => <button key={item.value} type="button" aria-pressed={priority === item.value} onClick={() => setPriority(item.value)} className={`min-h-7 cursor-pointer rounded-full border px-1 text-[8px] whitespace-nowrap ${priority === item.value ? colors[item.value].soft : "border-[#e5e8f0] bg-[#fafbff] text-[#30384f]"}`}>{item.label}</button>)}</div></fieldset>
          <label className="mt-4 block text-[11px] text-[#30384f]">Статус<select value={status} onChange={(event) => setStatus(event.target.value)} className="mt-2 h-8 w-full cursor-pointer rounded-[9px] border border-[#333] bg-white px-2 text-center text-[10px] outline-none focus:border-[#4562f0]">{statuses.map((item) => <option key={item.value} value={item.value}>{item.label}</option>)}</select></label>
          <div className="mt-4"><h3 className="text-[11px] text-[#30384f]">Исполнитель</h3><p className="mt-1.5 min-h-4 text-[10px] text-[#4562f0]">{responsible?.full_name || "Не назначен"}</p><Button text={responsible ? "Изменить исполнителя" : "Добавить исполнителя"} variant={responsible ? "secondary" : "primary"} size="small" onClick={() => setAssignmentOpen(true)} className="mt-2 h-8 w-full rounded-[7px] px-2 text-[10px] font-normal" /></div>
        </section>
        {notice && <p role="status" className="mt-3 rounded-lg bg-[#eef1ff] px-3 py-2 text-[11px] text-[#4562f0]">{notice}</p>}
        {error && <p role="alert" className="mt-3 rounded-lg bg-[#fff1f1] px-3 py-2 text-[11px] text-[#b42318]">{error}</p>}
        <div className="mt-auto grid gap-2 pt-6">
          {mode === "operator" && <button type="button" disabled={pending} onClick={() => setCloseOpen(true)} className="h-9 w-full cursor-pointer rounded-[7px] border border-[#e5141b] bg-[#e5141b] text-[11px] font-medium text-white hover:bg-[#c81017] disabled:opacity-50">Закрыть обращение</button>}
          <button type="button" disabled={pending} onClick={save} className="h-9 w-full cursor-pointer rounded-[7px] border border-[#e5141b] bg-white text-[11px] font-medium text-[#e5141b] hover:bg-[#fff1f1] disabled:opacity-50">Сохранить изменения</button>
        </div>
      </aside>
      {closeOpen && <CloseTicketModal onClose={() => setCloseOpen(false)} onSubmit={(message) => void closeTicket(message)} />}
      {assignmentOpen && <AssignmentDialog trackID={trackID} mode={mode} onClose={() => setAssignmentOpen(false)} onAssigned={() => { setAssignmentOpen(false); setNotice("Исполнитель назначен"); void load(); }} />}
    </div>
  );
}
