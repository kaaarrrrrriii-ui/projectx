"use client";

import Link from "next/link";
import { FormEvent, useEffect, useState } from "react";
import Button from "@/shared/ui/button";
import type { ExpertAppeal, ExpertPriority } from "./expert-data";

const priorityStyles: Record<ExpertPriority, { label: string; title: string; badge: string }> = {
  urgent: { label: "Срочное", title: "text-[#d70d14]", badge: "bg-[#ffd7d9] text-[#d70d14]" },
  standard: { label: "Стандартное", title: "text-[#4562f0]", badge: "bg-[#dfe6ff] text-[#4562f0]" },
  low: { label: "Низкое", title: "text-[#087f1a]", badge: "bg-[#dff2e0] text-[#087f1a]" },
};

type ChatMessage = { id: number; author: "expert" | "applicant"; text: string };

function RequestModal({ onClose, onSent }: { onClose: () => void; onSent: () => void }) {
  useEffect(() => {
    function closeOnEscape(event: KeyboardEvent) {
      if (event.key === "Escape") onClose();
    }
    document.addEventListener("keydown", closeOnEscape);
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", closeOnEscape);
      document.body.style.overflow = "";
    };
  }, [onClose]);

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    onSent();
  }

  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-[#000828]/45 p-4 backdrop-blur-[2px]" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget) onClose(); }}>
      <section role="dialog" aria-modal="true" aria-labelledby="request-dialog-title" className="w-full max-w-[680px] rounded-2xl bg-white p-6 shadow-[0_28px_80px_rgba(0,8,40,0.28)] sm:p-8">
        <header className="flex items-start justify-between gap-5">
          <h2 id="request-dialog-title" className="text-[clamp(22px,3vw,28px)] font-extrabold text-[#4562f0]">Отправка запроса</h2>
          <button type="button" onClick={onClose} aria-label="Закрыть окно" className="grid h-9 w-9 cursor-pointer place-items-center rounded-full text-2xl leading-none text-[#4562f0] transition-colors hover:bg-[#eef1ff]">×</button>
        </header>

        <form onSubmit={handleSubmit} className="mt-5 grid gap-5">
          <label className="grid gap-2 text-sm font-medium text-[#20263a]">
            Фамилия и инициалы эксперта
            <input autoFocus required defaultValue="Зетник О. А." className="h-11 rounded-xl border border-[#4f5873] bg-[#fcfdff] px-4 text-sm font-normal outline-none transition-shadow focus:border-[#4562f0] focus:ring-3 focus:ring-[#4562f0]/15" />
          </label>
          <label className="grid gap-2 text-sm font-medium text-[#20263a]">
            Тип запроса
            <select required defaultValue="consultation" className="h-11 rounded-xl border border-[#4f5873] bg-[#fcfdff] px-4 text-sm font-normal outline-none transition-shadow focus:border-[#4562f0] focus:ring-3 focus:ring-[#4562f0]/15">
              <option value="consultation">Консультация по обращению</option>
              <option value="legal">Юридическая консультация</option>
              <option value="operator">Уточнение у оператора</option>
            </select>
          </label>
          <label className="grid gap-2 text-sm font-medium text-[#20263a]">
            Причина
            <textarea required rows={5} placeholder="Опишите, какая информация или помощь необходима" className="resize-y rounded-xl border border-[#4f5873] bg-[#fcfdff] p-4 text-sm font-normal leading-6 outline-none transition-shadow placeholder:text-[#9aa0b3] focus:border-[#4562f0] focus:ring-3 focus:ring-[#4562f0]/15" />
          </label>
          <Button text="Отправить запрос" type="submit" variant="primary" className="w-full" />
        </form>
      </section>
    </div>
  );
}

export default function ExpertAppealDetail({ appeal, initialExecutor = "" }: { appeal: ExpertAppeal; initialExecutor?: string }) {
  const [status, setStatus] = useState(appeal.status);
  const executor = initialExecutor || "Не назначен";
  const [notes, setNotes] = useState("");
  const [answer, setAnswer] = useState("");
  const [requestOpen, setRequestOpen] = useState(false);
  const [notice, setNotice] = useState("");
  const [messages, setMessages] = useState<ChatMessage[]>([
    { id: 1, author: "expert", text: "Спасибо, что рассказали об этом. Вы не виноваты в происходящем, и с этой ситуацией не нужно оставаться один на один." },
    { id: 2, author: "applicant", text: "Спасибо. Я хочу понять, к кому можно обратиться в школе и как сделать это безопасно." },
  ]);
  const priority = priorityStyles[appeal.priority];

  function sendAnswer(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const text = answer.trim();
    if (!text) return;
    setMessages((current) => [...current, { id: Date.now(), author: "expert", text }]);
    setAnswer("");
    setStatus("Ответ отправлен");
    setNotice("Ответ отправлен заявителю");
  }

  return (
    <div className="grid min-w-0 flex-1 grid-cols-[minmax(0,1fr)_390px] bg-white/35 max-[1199px]:grid-cols-1">
      <article className="min-w-0 px-5 py-6 sm:px-8 lg:px-10">
        <div className="w-full">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <Link href={appeal.group === "return" ? "/expert/returns" : appeal.group === "assigned" ? "/expert/assigned" : "/expert/queue"} className="rounded-sm text-sm text-[#85899b] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[#4562f0]">← Вернуться назад</Link>
            <a href="#chat-history" className="rounded-sm text-xs text-[#30384f] underline-offset-4 hover:text-[#4562f0] hover:underline">Просмотреть историю чата</a>
          </div>

          <header className="mt-7">
            <div className="flex flex-wrap items-center gap-4">
              <h1 className={`text-[clamp(22px,2.4vw,30px)] font-extrabold leading-tight ${priority.title}`}>Обращение №{appeal.track}</h1>
              <span className={`rounded-full px-4 py-1.5 text-[11px] font-medium ${priority.badge}`}>{priority.label}</span>
            </div>
            <p className="mt-1 text-sm text-[#30384f]">Дата: {appeal.submittedAt}</p>
          </header>

          <section className="mt-7">
            <h2 className="text-base font-medium text-[#151515]">Тип заявителя: {appeal.applicant}</h2>
            <p className="mt-1 text-sm text-[#646d86]">Категория: {appeal.category}</p>
          </section>

          <section className="mt-7" aria-labelledby="source-text-title">
            <h2 id="source-text-title" className="text-base font-medium text-[#151515]">Исходный текст</h2>
            <div className="mt-3 min-h-[160px] rounded-[14px] border border-[#4f5873] bg-[#fcfdff]/90 p-4 text-sm leading-6 text-[#30384f] shadow-[inset_0_1px_3px_rgba(0,8,40,0.04)]">{appeal.description}</div>
          </section>

          <section className="mt-7" aria-labelledby="clarifications-title">
            <h2 id="clarifications-title" className="text-base font-medium text-[#151515]">Ответы на уточняющие вопросы:</h2>
            <dl className="mt-3 grid gap-2 text-sm text-[#30384f]">
              {appeal.clarifications.map((item) => (
                <div key={item.question} className="grid gap-0.5 sm:grid-cols-[minmax(190px,auto)_1fr] sm:gap-3">
                  <dt>{item.question}</dt>
                  <dd className="font-medium text-[#151515]">{item.answer}</dd>
                </div>
              ))}
            </dl>
          </section>

          <section className="mt-7" aria-labelledby="attachments-title">
            <h2 id="attachments-title" className="text-base font-medium text-[#151515]">Прикреплённые файлы:</h2>
            {appeal.attachments.length > 0 ? (
              <ul className="mt-3 flex flex-wrap gap-2">
                {appeal.attachments.map((file) => <li key={file} className="rounded-lg border border-[#bdc7f8] bg-white px-3 py-2 text-sm text-[#4562f0]">📎 {file}</li>)}
              </ul>
            ) : <p className="mt-2 text-sm text-[#858da2]">Файлы не приложены</p>}
          </section>

          <section id="chat-history" className="mt-9 scroll-mt-6" aria-labelledby="chat-title">
            <h2 id="chat-title" className="text-[22px] font-extrabold text-[#4562f0]">История чата обращения</h2>
            <div className="mt-5 grid gap-4">
              {messages.map((message) => (
                <article key={message.id} className={`max-w-[92%] rounded-2xl border p-4 sm:p-5 ${message.author === "expert" ? "border-[#7990ff] bg-white" : "ml-auto border-[#aeb9e9] bg-[#dfe6ff]"}`}>
                  <h3 className="text-sm font-medium text-[#30384f]">{message.author === "expert" ? "Ответ эксперта" : "Ответ заявителя"}</h3>
                  <p className="mt-3 rounded-xl border border-[#8c9af0] bg-white/70 p-4 text-sm leading-6 text-[#30384f]">{message.text}</p>
                </article>
              ))}
            </div>
          </section>

          <section id="answer-form" className="mt-9 scroll-mt-6" aria-labelledby="answer-title">
            <h2 id="answer-title" className="text-[22px] font-extrabold text-[#4562f0]">Отправка ответа</h2>
            <form onSubmit={sendAnswer} className="mt-5 overflow-hidden rounded-2xl border border-[#7990ff] bg-white shadow-[0_12px_28px_rgba(69,98,240,0.07)]">
              <label className="block p-4 text-sm font-medium text-[#30384f]">
                Ответ
                <textarea value={answer} onChange={(event) => setAnswer(event.target.value)} rows={8} placeholder="Напишите бережный и понятный ответ заявителю" className="mt-3 block w-full resize-y rounded-xl border border-[#4f5873] bg-[#fcfdff] p-4 text-sm font-normal leading-6 outline-none placeholder:text-[#9aa0b3] focus:border-[#4562f0] focus:ring-3 focus:ring-[#4562f0]/15" />
              </label>
              <button type="submit" disabled={!answer.trim()} className="min-h-12 w-full cursor-pointer bg-[#4562f0] px-5 text-sm font-medium text-white transition-colors hover:bg-[#374ecc] disabled:cursor-not-allowed disabled:bg-[#b7c0f5]">Отправить ответ</button>
            </form>
          </section>
        </div>
      </article>

      <aside className="border-l border-[#7990ff] bg-white/75 px-5 py-6 max-[1199px]:border-t max-[1199px]:border-l-0 sm:px-7" aria-labelledby="appeal-settings-title">
        <div className="sticky top-5">
          <h2 id="appeal-settings-title" className="text-lg font-medium text-[#151515]">Редактировать обращение</h2>
          <section className="mt-4 rounded-2xl border border-[#7990ff] bg-white p-4 shadow-[0_10px_24px_rgba(69,98,240,0.06)]">
            <label className="grid gap-2 text-xs text-[#30384f]">
              Статус
              <select value={status} onChange={(event) => setStatus(event.target.value)} className="h-10 rounded-xl border border-[#4f5873] bg-white px-3 text-center text-sm outline-none focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/15">
                <option>Новое</option>
                <option>В работе</option>
                <option>Нужна консультация</option>
                <option>Нужен повторный ответ</option>
                <option>Ответ отправлен</option>
                <option>Закрыто</option>
              </select>
            </label>
            <div className="mt-5">
              <p className="text-xs text-[#30384f]">Исполнитель</p>
              <p className={`mt-1.5 text-sm font-medium ${executor === "Не назначен" ? "text-[#8c93a8]" : "text-[#4562f0]"}`}>{executor}</p>
              <Button
                text={executor === "Не назначен" ? "Добавить исполнителя" : "Изменить исполнителя"}
                variant={executor === "Не назначен" ? "primary" : "secondary"}
                size="small"
                link={`/expert/appeals/${encodeURIComponent(appeal.track)}/assign`}
                className="mt-3 h-9 w-full rounded-lg text-xs font-normal"
              />
            </div>
          </section>

          <label className="mt-5 grid gap-2 text-xs text-[#30384f]">
            Заметки эксперта
            <textarea value={notes} onChange={(event) => setNotes(event.target.value)} rows={5} placeholder="Эти заметки видны только сотрудникам" className="resize-y rounded-xl border border-[#4f5873] bg-white p-3 text-sm leading-5 outline-none placeholder:text-[#a0a6b7] focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/15" />
          </label>

          {notice && <p role="status" className="mt-4 rounded-xl border border-[#b7c2fa] bg-[#eef1ff] px-4 py-3 text-xs text-[#4562f0]">{notice}</p>}

          <div className="mt-6 grid gap-2">
            <a href="#answer-form" className="inline-flex min-h-10 items-center justify-center rounded-lg bg-[#4562f0] px-4 text-xs font-medium text-white transition-colors hover:bg-[#374ecc]">Перейти к ответу</a>
            <Button text="Отправить запрос" variant="secondary" size="small" onClick={() => setRequestOpen(true)} className="h-10 w-full rounded-lg text-xs font-normal" />
          </div>
        </div>
      </aside>

      {requestOpen && <RequestModal onClose={() => setRequestOpen(false)} onSent={() => { setRequestOpen(false); setNotice("Запрос отправлен оператору"); }} />}
    </div>
  );
}
