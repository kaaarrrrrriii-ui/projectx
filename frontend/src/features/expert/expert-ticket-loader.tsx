"use client";

import { useEffect, useState } from "react";
import { staffRequest } from "@/shared/api/staff-api";
import ExpertAppealDetail from "./expert-appeal-detail";
import type { ExpertAppeal, ExpertPriority } from "./expert-data";

type APITicket = { track_id: string; category: { name: string }; applicant_type: string; priority: ExpertPriority; status: string; created_at: string; description: string; clarifications: Array<{ question: string; answer: string }>; attachments: Array<{ name: string }>; workers: Array<{ full_name: string; is_responsible: boolean }> };

export default function ExpertTicketLoader({ track }: { track: string }) {
  const [appeal, setAppeal] = useState<ExpertAppeal | null>(null);
  const [executor, setExecutor] = useState("");
  const [error, setError] = useState("");
  useEffect(() => {
    staffRequest<APITicket>(`/api/expert/tickets/${encodeURIComponent(track)}`).then((value) => {
      setAppeal({ track: value.track_id, category: value.category.name, applicant: value.applicant_type, priority: value.priority, status: value.status, waiting: "", submittedAt: new Date(value.created_at).toLocaleString("ru-RU"), description: value.description, clarifications: value.clarifications, attachments: value.attachments.map((item) => item.name), group: "assigned" });
      setExecutor(value.workers.find((item) => item.is_responsible)?.full_name ?? "");
    }).catch((reason) => setError(reason instanceof Error ? reason.message : "Обращение недоступно"));
  }, [track]);
  if (error) return <section className="flex-1 p-8 text-[#d70d14]">{error}. Проверьте назначение и выполните вход заново.</section>;
  if (!appeal) return <section className="flex-1 p-8 text-[#646d86]">Загрузка обращения…</section>;
  return <ExpertAppealDetail appeal={appeal} initialExecutor={executor} persist />;
}
