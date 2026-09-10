"use client";

import { useEffect, useState } from "react";
import { staffRequest } from "@/shared/api/staff-api";
import TicketDetail from "./ticket-detail";
import type { OperatorTicket, TicketPriority } from "./tickets";

type APITicket = { track_id: string; status: string; priority: TicketPriority; created_at: string; category: { name: string }; applicant_type: string; description: string; clarifications: Array<{ question: string; answer: string }>; attachments: Array<{ name: string }>; workers: Array<{ full_name: string; expert_group: string; is_responsible: boolean }> };

export default function OperatorTicketLoader({ track }: { track: string }) {
  const [ticket, setTicket] = useState<OperatorTicket | null>(null);
  const [expert, setExpert] = useState("");
  const [error, setError] = useState("");
  useEffect(() => {
    staffRequest<APITicket>(`/api/operator/tickets/${encodeURIComponent(track)}`).then((value) => {
      setTicket({ track: value.track_id, status: value.status, category: value.category.name, applicant: value.applicant_type, priority: value.priority, waiting: "", submittedAt: new Date(value.created_at).toLocaleString("ru-RU"), description: value.description, clarifications: value.clarifications, attachments: value.attachments.map((item) => item.name) });
      const responsible = value.workers.find((item) => item.is_responsible);
      setExpert(responsible ? `${responsible.full_name} — ${responsible.expert_group}` : "");
    }).catch((reason) => setError(reason instanceof Error ? reason.message : "Обращение не найдено"));
  }, [track]);
  if (error) return <section className="flex-1 p-8 text-[#d70d14]">{error}. Выполните вход сотрудника заново.</section>;
  if (!ticket) return <section className="flex-1 p-8 text-[#646d86]">Загрузка обращения…</section>;
  return <TicketDetail ticket={ticket} initialExpert={expert} persist />;
}
