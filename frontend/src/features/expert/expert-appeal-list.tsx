"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { getExpertRequests, getExpertTickets, type ExpertTicket, type WorkerRequest } from "@/shared/api/staff-api";

type ListMode = "queue" | "assigned" | "returns" | "requests";

const titles: Record<ListMode, { title: string; subtitle: string }> = {
  queue: { title: "Очередь новых", subtitle: "Назначенные и ещё не открытые обращения" },
  assigned: { title: "Распределённые", subtitle: "Обращения в работе" },
  returns: { title: "Возвраты", subtitle: "Повторно назначенные обращения" },
  requests: { title: "Запросы", subtitle: "История запросов оператору" },
};

const priorityLabels: Record<string, string> = { urgent: "Срочное", standard: "Стандартное", low: "Низкое" };
const applicantLabels: Record<string, string> = { schoolchild: "Школьник", parent: "Родитель", teacher: "Педагог" };

export default function ExpertAppealList({ mode }: { mode: ListMode; initialSearch?: string }) {
  const [tickets, setTickets] = useState<ExpertTicket[]>([]);
  const [requests, setRequests] = useState<WorkerRequest[]>([]);
  const [error, setError] = useState("");
  const heading = titles[mode];

  useEffect(() => {
    if (mode === "requests") {
      getExpertRequests().then((page) => setRequests(page.items)).catch((reason) => setError(reason instanceof Error ? reason.message : "Не удалось загрузить запросы"));
      return;
    }
    const queue = mode === "returns" ? "returned" : mode;
    getExpertTickets(queue).then((page) => setTickets(page.items)).catch((reason) => setError(reason instanceof Error ? reason.message : "Не удалось загрузить обращения"));
  }, [mode]);

  return (
    <section className="min-w-0 flex-1 px-5 py-6 sm:px-8" aria-labelledby="expert-list-heading">
      <header><h1 id="expert-list-heading" className="text-[26px] font-extrabold text-[#4562f0]">{heading.title}</h1><p className="mt-1 text-sm text-[#646d86]">{heading.subtitle}</p></header>
      {error && <p role="alert" className="mt-5 rounded-xl bg-[#fff1f1] px-4 py-3 text-sm text-[#b42318]">{error}</p>}

      {mode === "requests" ? (
        <div className="mt-6 grid gap-3">
          {requests.map((request) => <article key={request.id} className="rounded-[13px] border border-[#7990ff] bg-white p-4"><div className="flex flex-wrap justify-between gap-3"><strong className="text-[#4562f0]">{request.track_id}</strong><span className={`rounded-full px-3 py-1 text-xs ${request.status === "completed" ? "bg-[#dff2e0] text-[#087f1a]" : "bg-[#fff0b3] text-[#8a6500]"}`}>{request.status === "completed" ? "Выполнен" : "Отправлен"}</span></div><p className="mt-3 text-sm text-[#30384f]">{request.reason}</p></article>)}
          {requests.length === 0 && !error && <p className="rounded-[13px] border border-[#b7c2fa] bg-white/80 p-10 text-center text-sm text-[#646d86]">Запросов пока нет</p>}
        </div>
      ) : (
        <div className="mt-6 overflow-x-auto rounded-[13px] border border-[#4562f0] bg-white/85">
          <table className="w-full min-w-[720px] table-fixed border-collapse">
            <thead className="bg-[#dfe6ff] text-[#4562f0]"><tr className="h-[52px]"><th className="w-[30%] border-r border-[#4562f0] px-4 font-normal">Трек-номер</th><th className="w-[28%] border-r border-[#4562f0] px-4 font-normal">Категория</th><th className="w-[22%] border-r border-[#4562f0] px-4 font-normal">Заявитель</th><th className="w-[20%] px-4 font-normal">Приоритет</th></tr></thead>
            <tbody>
              {tickets.map((ticket) => <tr key={ticket.track_id} className="h-[52px] border-t border-[#4562f0] hover:bg-[#f7f8ff]"><td className="border-r border-[#4562f0] px-4"><Link href={`/expert/appeals/${encodeURIComponent(ticket.track_id)}`} className="font-medium text-[#4562f0] hover:underline">{ticket.track_id}</Link></td><td className="border-r border-[#4562f0] px-4 text-center text-sm">{ticket.category.name}</td><td className="border-r border-[#4562f0] px-4 text-center text-sm">{applicantLabels[ticket.applicant_type] ?? ticket.applicant_type}</td><td className="px-4 text-center text-sm">{priorityLabels[ticket.priority] ?? ticket.priority}</td></tr>)}
              {tickets.length === 0 && !error && <tr className="h-24 border-t border-[#4562f0]"><td colSpan={4} className="text-center text-sm text-[#646d86]">Обращений пока нет</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
