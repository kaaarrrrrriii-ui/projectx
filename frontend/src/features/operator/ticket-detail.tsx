"use client";

import Button from "@/shared/ui/button";
import Link from "next/link";
import { useState } from "react";
import { applicantLabels, type OperatorTicket, type TicketPriority } from "./tickets";
import { CloseTicketModal, SpecialistAssignmentModal } from "./ticket-detail-modals";
import { staffRequest } from "@/shared/api/staff-api";

const priorities: Array<{ value: TicketPriority; label: string }> = [
  { value: "urgent", label: "Срочное" },
  { value: "standard", label: "Стандартное" },
  { value: "low", label: "Вопрос" },
];

const statuses = [
  { value: "new", label: "Новое" },
  { value: "assigned", label: "Распределено" },
  { value: "in_progress", label: "В работе" },
  { value: "needs_clarification", label: "Нужно уточнение" },
  { value: "answer_ready", label: "Ответ готов" },
  { value: "returned", label: "Возвращено" },
  { value: "rejected", label: "Отклонено" },
  { value: "closed_without_answer", label: "Закрыто" },
];

const colors: Record<TicketPriority, { text: string; soft: string }> = {
  urgent: { text: "text-[#e5141b]", soft: "bg-[#f7b6b8] text-[#d70d14]" },
  standard: { text: "text-[#4562f0]", soft: "bg-[#dfe6ff] text-[#4562f0]" },
  low: { text: "text-[#087f1a]", soft: "bg-[#dff2e0] text-[#087f1a]" },
};

export default function TicketDetail({
  ticket,
  initialExpert = "",
  returnBasePath = "/operator/queueNew.tsx",
  persist = false,
}: {
  ticket: OperatorTicket;
  initialExpert?: string;
  returnBasePath?: string;
  persist?: boolean;
}) {
  const [priority, setPriority] = useState<TicketPriority>(ticket.priority);
  const [status, setStatus] = useState(ticket.status);
  const [assignedExpert, setAssignedExpert] = useState(initialExpert);
  const [notice, setNotice] = useState("");
  const [assignmentOpen, setAssignmentOpen] = useState(false);
  const [closeOpen, setCloseOpen] = useState(false);
  const accent = colors[priority];

  async function assignSpecialist(specialist: string) {
    if (persist) {
      try {
        const result = await staffRequest<{ workers: Array<{ id: number; full_name: string; expert_group: { title: string }; available: boolean }> }>(`/api/operator/tickets/${encodeURIComponent(ticket.track)}/eligible-workers?only_available=true`);
        const worker = result.workers.find((item) => item.available && specialist.includes(item.full_name)) ?? result.workers.find((item) => item.available);
        if (!worker) throw new Error("Нет доступных экспертов");
        await staffRequest(`/api/operator/tickets/${encodeURIComponent(ticket.track)}/responsible-worker`, { method: "PUT", body: JSON.stringify({ worker_id: worker.id }) });
        specialist = `${worker.full_name} — ${worker.expert_group.title}`;
      } catch (reason) { setNotice(reason instanceof Error ? reason.message : "Не удалось назначить эксперта"); return; }
    }
    setAssignedExpert(specialist);
    setStatus("assigned");
    setNotice("Исполнитель назначен");
    setAssignmentOpen(false);
  }

  async function closeTicket(reason: string) {
    if (persist) {
      try { await staffRequest(`/api/operator/tickets/${encodeURIComponent(ticket.track)}/close`, { method: "POST", body: JSON.stringify({ message: reason }) }); }
      catch (cause) { setNotice(cause instanceof Error ? cause.message : "Не удалось закрыть обращение"); return; }
    }
    setStatus("closed");
    setNotice(`Обращение закрыто. Причина: ${reason}`);
    setCloseOpen(false);
  }

  async function saveChanges() {
    if (persist) {
      try { await staffRequest(`/api/operator/tickets/${encodeURIComponent(ticket.track)}`, { method: "PATCH", body: JSON.stringify({ priority, status }) }); }
      catch (cause) { setNotice(cause instanceof Error ? cause.message : "Не удалось сохранить изменения"); return; }
    }
    setNotice("Изменения сохранены");
  }

  return (
    <div className="grid min-w-0 flex-1 grid-cols-[minmax(0,1fr)_268px] max-[899px]:grid-cols-1">
      <article className="flex min-w-0 flex-col px-3 py-4 sm:px-4">
        <Link href={returnBasePath} className="w-fit rounded-sm text-sm text-[#85899b] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[#4562f0]">Вернуться назад</Link>

        <header className="mt-6 flex flex-wrap items-center gap-3">
          <h1 className={`text-[clamp(20px,2.3vw,25px)] leading-tight font-extrabold ${accent.text}`}>Обращение №{ticket.track}</h1>
          <span className={`rounded-full px-4 py-1.5 text-[10px] ${accent.soft}`}>{priorities.find((item) => item.value === priority)?.label}</span>
        </header>
        <p className="mt-1 text-sm font-medium text-[#151515]">Дата: {ticket.submittedAt}</p>

        <section className="mt-7">
          <h2 className="text-base font-medium text-[#151515]">Тип заявителя: {applicantLabels[ticket.applicant] ?? ticket.applicant}</h2>
        </section>

        <section aria-labelledby="original-text-heading" className="mt-6">
          <h2 id="original-text-heading" className="text-base font-medium text-[#151515]">Исходный текст</h2>
          <div className="mt-2 min-h-[130px] rounded-[11px] border border-[#333] bg-white/75 px-3 py-2 text-xs leading-5 text-[#30384f]">{ticket.description}</div>
        </section>

        <section aria-labelledby="answers-heading" className="mt-6">
          <h2 id="answers-heading" className="text-base font-medium text-[#151515]">Ответы на уточняющие вопросы:</h2>
          <dl className="mt-2 grid gap-1 text-xs text-[#3f475d]">
            {ticket.clarifications.map((item) => (
              <div key={item.question} className="flex flex-wrap gap-1"><dt>{item.question}</dt><dd className="font-medium text-[#151515]">{item.answer}</dd></div>
            ))}
          </dl>
        </section>

        <section aria-labelledby="files-heading" className="mt-6">
          <h2 id="files-heading" className="text-base font-medium text-[#151515]">Прикреплённые файлы:</h2>
          {ticket.attachments.length ? (
            <ul className="mt-2 flex flex-wrap gap-2">
              {ticket.attachments.map((file) => <li key={file} className="text-xs text-[#30384f]">{file}</li>)}
            </ul>
          ) : <p className="mt-2 text-xs text-[#7b849b]">Файлы не приложены</p>}
        </section>

        <section aria-labelledby="chat-history-heading" className="mt-7 max-w-[520px]">
          <h2 id="chat-history-heading" className="text-[22px] font-extrabold text-[#4562f0]">История чата обращения</h2>
          <div className="mt-4 rounded-[12px] border border-[#4562f0] bg-white p-4">
            <h3 className="text-sm text-[#4562f0]">Ответ эксперта</h3>
            <div className="mt-3 min-h-[82px] rounded-[10px] border border-[#7990ff] bg-white" />
          </div>
          <div className="mt-5 rounded-[12px] border border-[#000828] bg-[#dfe6ff] p-4">
            <h3 className="text-sm text-[#000828]">Ответ заявителя</h3>
            <div className="mt-3 min-h-[94px] rounded-[10px] border border-[#000828] bg-[#dfe6ff]" />
          </div>
        </section>
      </article>

      <aside aria-labelledby="edit-ticket-heading" className="flex flex-col border border-[#4562f0] bg-white/70 px-3.5 py-5 min-[900px]:sticky min-[900px]:top-[100px] min-[900px]:h-[calc(100dvh-100px)] min-[900px]:self-start min-[900px]:overflow-y-auto max-[899px]:min-h-0 max-[899px]:border-r-0 max-[899px]:border-b-0 max-[899px]:border-l-0">
        <h2 id="edit-ticket-heading" className="text-base font-medium text-[#151515]">Редактировать обращение</h2>

        <section aria-label="Параметры обращения" className="mt-3 rounded-[12px] border border-[#4562f0] bg-white/80 p-2.5">
          <fieldset>
            <legend className="text-xs text-[#30384f]">Приоритет</legend>
            <div className="mt-2 grid grid-cols-3 gap-1.5">
              {priorities.map((item) => (
                <button key={item.value} type="button" aria-pressed={priority === item.value} onClick={() => setPriority(item.value)} className={`min-h-7 cursor-pointer rounded-full border border-transparent px-1.5 text-[8px] leading-none whitespace-nowrap transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] ${priority === item.value ? colors[item.value].soft : "bg-[#fafbff] text-[#30384f] hover:border-[#bbc5f5]"}`}>{item.label}</button>
              ))}
            </div>
          </fieldset>

          <label className="mt-5 block text-xs text-[#30384f]">
            Статус
            <select value={status} onChange={(event) => setStatus(event.target.value)} className="mt-2 h-8 w-full cursor-pointer rounded-[9px] border border-[#333] bg-white px-2 text-center text-[11px] outline-none focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/15">
              {statuses.map((item) => <option key={item.value} value={item.value}>{item.label}</option>)}
            </select>
          </label>

          <div className="mt-5">
            <h3 className="text-xs text-[#30384f]">Исполнитель</h3>
            <p className="mt-1.5 min-h-4 text-[11px] text-[#4562f0]">{assignedExpert || "Не назначен"}</p>
            <Button text="Изменить исполнителя" variant="secondary" size="small" onClick={() => setAssignmentOpen(true)} className="mt-2 h-8 w-full rounded-[7px] px-2 text-[10px] font-normal" />
            <Button text="Добавить исполнителя" variant="primary" size="small" onClick={() => setAssignmentOpen(true)} disabled={Boolean(assignedExpert)} className="mt-2 h-8 w-full rounded-[7px] px-2 text-[10px] font-normal" />
          </div>
        </section>

        {notice && <p role="status" className="mt-3 rounded-lg bg-[#eef1ff] px-3 py-2 text-[11px] leading-4 text-[#4562f0]">{notice}</p>}

        <div className="mt-auto grid gap-2 pt-6">
          <button type="button" onClick={() => setCloseOpen(true)} className="h-9 w-full cursor-pointer rounded-[7px] border border-[#e5141b] bg-[#e5141b] text-[11px] font-medium text-white transition-colors hover:bg-[#c81017] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#e5141b]">Закрыть обращение</button>
          <button type="button" onClick={saveChanges} className="h-9 w-full cursor-pointer rounded-[7px] border border-[#e5141b] bg-white text-[11px] font-medium text-[#e5141b] transition-colors hover:bg-[#fff1f1] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#e5141b]">Сохранить изменения</button>
        </div>
      </aside>

      {closeOpen && <CloseTicketModal onClose={() => setCloseOpen(false)} onSubmit={closeTicket} />}
      {assignmentOpen && <SpecialistAssignmentModal category={ticket.category} onClose={() => setAssignmentOpen(false)} onAssign={assignSpecialist} />}
    </div>
  );
}
