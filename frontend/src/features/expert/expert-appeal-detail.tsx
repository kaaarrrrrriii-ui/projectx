"use client";

import Link from "next/link";
import { FormEvent, useCallback, useEffect, useState } from "react";
import Button from "@/shared/ui/button";
import {
  createReplacementRequest,
  getExpertTicket,
  openExpertTicket,
  sendExpertAnswer,
  type TicketDetail,
} from "@/shared/api/staff-api";

const priorityStyles: Record<string, { label: string; title: string; badge: string }> = {
  urgent: { label: "Срочное", title: "text-[#d70d14]", badge: "bg-[#ffd7d9] text-[#d70d14]" },
  standard: { label: "Стандартное", title: "text-[#4562f0]", badge: "bg-[#dfe6ff] text-[#4562f0]" },
  low: { label: "Низкое", title: "text-[#087f1a]", badge: "bg-[#dff2e0] text-[#087f1a]" },
};

function RequestModal({ pending, error, onClose, onSubmit }: { pending: boolean; error: string; onClose: () => void; onSubmit: (reason: string) => void }) {
  const [reason, setReason] = useState("");
  return (
    <div className="fixed inset-0 z-[80] grid place-items-center bg-[#000828]/45 p-4" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget) onClose(); }}>
      <section role="dialog" aria-modal="true" aria-labelledby="request-dialog-title" className="w-full max-w-[680px] rounded-2xl bg-white p-6 shadow-2xl sm:p-8">
        <header className="flex items-start justify-between gap-5"><h2 id="request-dialog-title" className="text-[26px] font-extrabold text-[#4562f0]">Запросить замену исполнителя</h2><button type="button" onClick={onClose} className="grid h-9 w-9 cursor-pointer place-items-center rounded-full text-2xl text-[#4562f0] hover:bg-[#eef1ff]">×</button></header>
        <p className="mt-3 text-sm leading-6 text-[#646d86]">Опишите причину. Нового специалиста выберет оператор.</p>
        <label className="mt-5 grid gap-2 text-sm font-medium text-[#20263a]">Причина<textarea autoFocus required value={reason} onChange={(event) => setReason(event.target.value)} rows={5} placeholder="Почему обращение нужно передать другому специалисту" className="resize-y rounded-xl border border-[#4f5873] bg-[#fcfdff] p-4 text-sm font-normal leading-6 outline-none focus:border-[#4562f0]" /></label>
        {error && <p role="alert" className="mt-4 text-sm text-[#b42318]">{error}</p>}
        <Button text={pending ? "Отправляем…" : "Отправить запрос"} disabled={pending || !reason.trim()} onClick={() => onSubmit(reason.trim())} className="mt-5 w-full" />
      </section>
    </div>
  );
}

export default function ExpertAppealDetail({ trackID }: { trackID: string }) {
  const [appeal, setAppeal] = useState<TicketDetail | null>(null);
  const [answer, setAnswer] = useState("");
  const [requestOpen, setRequestOpen] = useState(false);
  const [pending, setPending] = useState(false);
  const [notice, setNotice] = useState("");
  const [error, setError] = useState("");

  const load = useCallback(async () => {
    try {
      const value = await getExpertTicket(trackID);
      setAppeal(value);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Не удалось загрузить обращение");
    }
  }, [trackID]);

  // Opening is idempotent; after it finishes the authoritative card is loaded.
  useEffect(() => {
    openExpertTicket(trackID).catch(() => undefined).finally(() => void load());
  }, [load, trackID]);

  async function sendAnswer(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const text = answer.trim();
    if (!text) return;
    setPending(true);
    setError("");
    try {
      await sendExpertAnswer(trackID, text);
      setAnswer("");
      setNotice("Ответ сохранён и отправлен заявителю");
      await load();
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Не удалось отправить ответ");
    } finally {
      setPending(false);
    }
  }

  async function requestReplacement(reason: string) {
    setPending(true);
    setError("");
    try {
      await createReplacementRequest(trackID, reason);
      setRequestOpen(false);
      setNotice("Запрос на замену отправлен оператору");
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Не удалось отправить запрос");
    } finally {
      setPending(false);
    }
  }

  if (!appeal) return <section className="min-w-0 flex-1 p-8 text-sm text-[#646d86]">{error || "Загружаем обращение…"}</section>;
  const priority = priorityStyles[appeal.priority] ?? priorityStyles.standard;
  const responsible = appeal.workers.find((worker) => worker.is_responsible);

  return (
    <div className="grid min-w-0 flex-1 grid-cols-[minmax(0,1fr)_330px] bg-white/35 max-[1099px]:grid-cols-1">
      <article className="min-w-0 px-5 py-6 sm:px-8">
        <Link href="/expert/assigned" className="rounded-sm text-sm text-[#85899b] hover:text-[#4562f0]">← Вернуться назад</Link>
        <header className="mt-7"><div className="flex flex-wrap items-center gap-4"><h1 className={`text-[clamp(22px,2.4vw,30px)] font-extrabold ${priority.title}`}>Обращение №{appeal.track_id}</h1><span className={`rounded-full px-4 py-1.5 text-[11px] ${priority.badge}`}>{priority.label}</span></div><p className="mt-1 text-sm text-[#30384f]">Дата: {new Date(appeal.created_at).toLocaleString("ru-RU")}</p></header>
        <section className="mt-7"><h2 className="text-base font-medium text-[#151515]">Тип заявителя: {appeal.applicant_type}</h2><p className="mt-1 text-sm text-[#646d86]">Категория: {appeal.category.name}</p></section>
        <section className="mt-7"><h2 className="text-base font-medium text-[#151515]">Исходный текст</h2><div className="mt-3 min-h-[140px] rounded-[14px] border border-[#4f5873] bg-[#fcfdff]/90 p-4 text-sm leading-6 text-[#30384f]">{appeal.description}</div></section>
        <section className="mt-7"><h2 className="text-base font-medium text-[#151515]">Ответы на уточняющие вопросы:</h2><dl className="mt-3 grid gap-2 text-sm">{appeal.clarifications.map((item) => <div key={item.question_id} className="grid gap-1 sm:grid-cols-[minmax(190px,auto)_1fr]"><dt>{item.question}</dt><dd className="font-medium">{item.answer}</dd></div>)}</dl></section>
        <section id="chat-history" className="mt-9"><h2 className="text-[22px] font-extrabold text-[#4562f0]">История чата обращения</h2><div className="mt-5 grid gap-4">{(appeal.messages ?? []).filter((message) => message.type !== "internal_note").map((message) => <article key={message.id} className={`max-w-[92%] rounded-2xl border p-4 ${message.type === "specialist" ? "border-[#7990ff] bg-white" : "ml-auto border-[#aeb9e9] bg-[#dfe6ff]"}`}><h3 className="text-sm font-medium">{message.type === "specialist" ? "Ответ специалиста" : "Ответ заявителя"}</h3><p className="mt-3 text-sm leading-6">{message.text}</p></article>)}{(appeal.messages ?? []).length === 0 && <p className="text-sm text-[#858da2]">Сообщений пока нет</p>}</div></section>
        <section id="answer-form" className="mt-9"><h2 className="text-[22px] font-extrabold text-[#4562f0]">Отправка ответа</h2><form onSubmit={sendAnswer} className="mt-5 overflow-hidden rounded-2xl border border-[#7990ff] bg-white"><label className="block p-4 text-sm font-medium">Ответ<textarea value={answer} onChange={(event) => setAnswer(event.target.value)} rows={7} placeholder="Напишите бережный и понятный ответ заявителю" className="mt-3 block w-full resize-y rounded-xl border border-[#4f5873] bg-[#fcfdff] p-4 text-sm font-normal leading-6 outline-none focus:border-[#4562f0]" /></label><button type="submit" disabled={pending || !answer.trim()} className="min-h-12 w-full cursor-pointer bg-[#4562f0] px-5 text-sm font-medium text-white hover:bg-[#374ecc] disabled:cursor-not-allowed disabled:bg-[#b7c0f5]">{pending ? "Отправляем…" : "Отправить ответ"}</button></form></section>
      </article>

      <aside className="border-l border-[#7990ff] bg-white/75 px-5 py-6 max-[1099px]:border-t max-[1099px]:border-l-0 sm:px-7">
        <div className="sticky top-[120px]"><h2 className="text-lg font-medium text-[#151515]">Обращение</h2><section className="mt-4 rounded-2xl border border-[#7990ff] bg-white p-4"><p className="text-xs text-[#30384f]">Статус</p><p className="mt-1 text-sm font-medium text-[#4562f0]">{appeal.status}</p><div className="mt-5"><p className="text-xs text-[#30384f]">Ответственный</p><p className="mt-1.5 text-sm font-medium text-[#4562f0]">{responsible?.full_name || "Не назначен"}</p></div></section>
          {notice && <p role="status" className="mt-4 rounded-xl border border-[#b7c2fa] bg-[#eef1ff] px-4 py-3 text-xs text-[#4562f0]">{notice}</p>}
          {error && <p role="alert" className="mt-4 rounded-xl bg-[#fff1f1] px-4 py-3 text-xs text-[#b42318]">{error}</p>}
          <div className="mt-6 grid gap-2"><a href="#answer-form" className="inline-flex min-h-10 items-center justify-center rounded-lg bg-[#4562f0] px-4 text-xs font-medium text-white">Перейти к ответу</a><Button text="Запросить замену исполнителя" variant="secondary" size="small" onClick={() => { setError(""); setRequestOpen(true); }} className="h-10 w-full rounded-lg text-xs font-normal" /></div>
        </div>
      </aside>
      {requestOpen && <RequestModal pending={pending} error={error} onClose={() => setRequestOpen(false)} onSubmit={(reason) => void requestReplacement(reason)} />}
    </div>
  );
}
