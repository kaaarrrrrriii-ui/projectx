"use client";

import Button from "@/shared/ui/button";
import Link from "next/link";
import { useState } from "react";
import { applicantLabels, type OperatorTicket, type TicketPriority } from "./tickets";

const priorities: Array<{ value: TicketPriority; label: string }> = [
  { value: "urgent", label: "Срочное" },
  { value: "standard", label: "Стандартное" },
  { value: "low", label: "Низкое" },
];

const statuses = [
  { value: "new", label: "Новое" },
  { value: "assigned", label: "Распределено" },
  { value: "in-progress", label: "В работе" },
  { value: "clarification", label: "Нужно уточнение" },
  { value: "answer-ready", label: "Ответ готов" },
];

const experts = [
  "Анна Петрова — психолог",
  "Мария Соколова — социальный педагог",
  "Илья Воронов — юрист",
];

const colors: Record<TicketPriority, { text: string; soft: string; solid: string; outline: string; focus: string }> = {
  urgent: {
    text: "text-[#d70d14]",
    soft: "border-[#ef8b8f] bg-[#ffdfe0] text-[#d70d14]",
    solid: "border-[#e5141b] bg-[#e5141b] text-white hover:bg-[#c81017]",
    outline: "border-[#e5141b] bg-white text-[#e5141b] hover:bg-[#fff1f1]",
    focus: "focus-visible:outline-[#e5141b]",
  },
  standard: {
    text: "text-[#946100]",
    soft: "border-[#e4bd52] bg-[#fff0b3] text-[#946100]",
    solid: "border-[#e4ad2b] bg-[#e4ad2b] text-[#302100] hover:bg-[#d09a18]",
    outline: "border-[#e4ad2b] bg-white text-[#946100] hover:bg-[#fff8df]",
    focus: "focus-visible:outline-[#e4ad2b]",
  },
  low: {
    text: "text-[#087f1a]",
    soft: "border-[#72bf78] bg-[#dff2e0] text-[#087f1a]",
    solid: "border-[#15952a] bg-[#15952a] text-white hover:bg-[#087f1a]",
    outline: "border-[#15952a] bg-white text-[#087f1a] hover:bg-[#effaf0]",
    focus: "focus-visible:outline-[#15952a]",
  },
};

const actionBase = "h-10 w-full cursor-pointer rounded-[8px] border px-4 text-sm font-medium transition-colors focus-visible:outline-3 focus-visible:outline-offset-3";

export default function TicketDetail({ ticket }: { ticket: OperatorTicket }) {
  const [priority, setPriority] = useState<TicketPriority>(ticket.priority);
  const [status, setStatus] = useState(ticket.status);
  const [expertDraft, setExpertDraft] = useState("");
  const [assignedExpert, setAssignedExpert] = useState("");
  const [notice, setNotice] = useState("");
  const accent = colors[priority];

  function saveExpert() {
    if (!expertDraft) return;
    setAssignedExpert(expertDraft);
    setNotice(`Исполнитель назначен: ${expertDraft}`);
  }

  return (
    <div className="grid min-w-0 flex-1 grid-cols-[minmax(0,1fr)_480px] max-[1249px]:grid-cols-1">
      <article className="flex min-w-0 flex-col px-5 py-6 sm:px-8">
        <Link href="/operator/queueNew.tsx" className="w-fit rounded-sm text-sm text-[#85899b] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[#4562f0]">
          Вернуться назад
        </Link>

        <header className="mt-8 flex flex-wrap items-center gap-4">
          <h1 className={`text-[clamp(22px,2.3vw,30px)] leading-tight font-extrabold ${accent.text}`}>
            Обращение №{ticket.track}
          </h1>
          <span className={`rounded-full border px-5 py-1.5 text-xs font-medium ${accent.soft}`}>
            {priorities.find((item) => item.value === priority)?.label}
          </span>
        </header>
        <p className="mt-1.5 text-base text-[#151515]">Дата: {ticket.submittedAt}</p>

        <section className="mt-8">
          <h2 className="text-lg font-medium text-[#151515]">Тип заявителя: {applicantLabels[ticket.applicant] ?? ticket.applicant}</h2>
          <p className="mt-2 text-sm text-[#646d86]">Тема: {ticket.category.charAt(0).toUpperCase() + ticket.category.slice(1)}</p>
        </section>

        <section aria-labelledby="original-text-heading" className="mt-8">
          <h2 id="original-text-heading" className="text-lg font-medium text-[#151515]">Исходный текст</h2>
          <div className="mt-3 min-h-[180px] rounded-[13px] border border-[#333] bg-white/75 px-4 py-3 text-sm leading-6 text-[#30384f]">
            {ticket.description}
          </div>
        </section>

        <section aria-labelledby="answers-heading" className="mt-8">
          <h2 id="answers-heading" className="text-lg font-medium text-[#151515]">Ответы на уточняющие вопросы:</h2>
          <dl className="mt-3 grid gap-1 text-sm text-[#3f475d]">
            {ticket.clarifications.map((item) => (
              <div key={item.question} className="flex flex-wrap gap-1">
                <dt>{item.question}</dt>
                <dd className="font-medium text-[#151515]">{item.answer}</dd>
              </div>
            ))}
          </dl>
        </section>

        <section aria-labelledby="files-heading" className="mt-8">
          <h2 id="files-heading" className="text-lg font-medium text-[#151515]">Прикреплённые файлы:</h2>
          {ticket.attachments.length ? (
            <ul className="mt-3 flex flex-wrap gap-3">
              {ticket.attachments.map((file) => (
                <li key={file} className="flex items-center gap-2 rounded-lg border border-[#bdc7f8] bg-white/75 px-3 py-2 text-sm text-[#30384f]">
                  <svg viewBox="0 0 20 20" aria-hidden="true" className={`h-4 w-4 ${accent.text}`} fill="none"><path d="m7 10 5-5a3 3 0 1 1 4 4l-7 7a4 4 0 0 1-6-6l7-7" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" /></svg>
                  {file}
                </li>
              ))}
            </ul>
          ) : <p className="mt-3 text-sm text-[#7b849b]">Файлы не приложены</p>}
        </section>

        <Link href="/operator/queueNew.tsx" className="mt-auto w-fit rounded-sm pt-12 text-sm text-[#a0a6b7] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[#4562f0]">
          Вернуться назад
        </Link>
      </article>

      <aside aria-labelledby="edit-ticket-heading" className="flex min-h-[calc(100dvh-100px)] flex-col border-l border-[#4562f0] bg-white/55 px-6 py-6 max-[1249px]:min-h-0 max-[1249px]:border-t max-[1249px]:border-l-0">
        <h2 id="edit-ticket-heading" className="text-xl font-medium text-[#151515]">Редактировать обращение</h2>

        <section aria-label="Параметры обращения" className="mt-4 rounded-[14px] border border-[#4562f0] bg-white/80 p-3">
          <fieldset>
            <legend className="text-sm text-[#30384f]">Приоритет</legend>
            <div className="mt-2 grid grid-cols-3 gap-2">
              {priorities.map((item) => (
                <button key={item.value} type="button" aria-pressed={priority === item.value} onClick={() => setPriority(item.value)} className={`min-h-8 cursor-pointer rounded-full border px-3 text-[11px] whitespace-nowrap transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] ${priority === item.value ? colors[item.value].soft : "border-transparent bg-[#fafbff] text-[#30384f] hover:border-[#bbc5f5]"}`}>
                  {item.label}
                </button>
              ))}
            </div>
          </fieldset>

          <label className="mt-6 block text-sm text-[#30384f]">
            Статус
            <select value={status} onChange={(event) => setStatus(event.target.value)} className="mt-2 h-10 w-full rounded-[10px] border border-[#333] bg-white px-3 text-center text-sm outline-none focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/15">
              {statuses.map((item) => <option key={item.value} value={item.value}>{item.label}</option>)}
            </select>
          </label>
          <Button text="Изменить" variant="primary" size="small" onClick={() => setNotice("Статус обращения изменён")} className="mt-2.5 h-9 w-full rounded-[8px] text-sm font-normal" />

          <div className="mt-7">
            <h3 className="text-sm text-[#30384f]">Исполнитель</h3>
            <p className="mt-1.5 text-sm text-[#4562f0]">{assignedExpert || "Не назначен"}</p>
            <select value={expertDraft} onChange={(event) => setExpertDraft(event.target.value)} aria-label="Выбрать исполнителя" className="mt-2 h-10 w-full rounded-[8px] border border-[#4562f0] bg-white px-3 text-sm outline-none focus:ring-2 focus:ring-[#4562f0]/15">
              <option value="">Выберите исполнителя</option>
              {experts.map((item) => <option key={item} value={item}>{item}</option>)}
            </select>
            <Button text="Изменить исполнителя" variant="secondary" size="small" disabled={!expertDraft} onClick={saveExpert} className="mt-2.5 h-9 w-full rounded-[8px] text-sm font-normal" />
            <Button text="Добавить исполнителя" variant="primary" size="small" disabled={!expertDraft || Boolean(assignedExpert)} onClick={saveExpert} className="mt-2.5 h-9 w-full rounded-[8px] text-sm font-normal" />
          </div>
        </section>

        <label className="mt-5 block text-sm text-[#30384f]">
          Справка от исполнителя
          <textarea className="mt-2 block min-h-[138px] w-full resize-y rounded-[13px] border border-[#333] bg-white/80 p-3 text-sm leading-5 outline-none focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/15" />
        </label>

        {notice && <p role="status" className={`mt-4 rounded-lg border px-3 py-2 text-sm ${accent.soft}`}>{notice}</p>}

        <div className="mt-auto grid gap-2.5 pt-6">
          <button type="button" onClick={() => setNotice("Обращение закрыто")} className={`${actionBase} ${accent.solid} ${accent.focus}`}>Закрыть обращение</button>
          <button type="button" onClick={() => setNotice("Обращение вернуто на доработку")} className={`${actionBase} ${accent.outline} ${accent.focus}`}>Вернуть на доработку</button>
        </div>
      </aside>
    </div>
  );
}
