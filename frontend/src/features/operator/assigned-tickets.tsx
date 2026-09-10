"use client";

import { useEffect, useState } from "react";
import Button from "@/shared/ui/button";
import {
  completeOperatorRequest,
  getEligibleWorkers,
  getOperatorRequests,
  type Worker,
  type WorkerRequest,
} from "@/shared/api/staff-api";

export default function AssignedTickets() {
  const [requests, setRequests] = useState<WorkerRequest[]>([]);
  const [workers, setWorkers] = useState<Record<number, Worker[]>>({});
  const [selected, setSelected] = useState<Record<number, number>>({});
  const [pending, setPending] = useState<number | null>(null);
  const [message, setMessage] = useState("");

  async function load() {
    try {
      const page = await getOperatorRequests();
      setRequests(page.items);
      const entries = await Promise.all(page.items.map(async (request) => {
        try {
          const result = await getEligibleWorkers(request.track_id);
          return [request.id, result.workers.filter((worker) => worker.available)] as const;
        } catch {
          return [request.id, []] as const;
        }
      }));
      setWorkers(Object.fromEntries(entries));
    } catch (reason) {
      setMessage(reason instanceof Error ? reason.message : "Не удалось загрузить запросы");
    }
  }

  // The initial request synchronizes this client view with the backend.
  // eslint-disable-next-line react-hooks/set-state-in-effect
  useEffect(() => { void load(); }, []);

  async function complete(request: WorkerRequest) {
    const workerID = selected[request.id];
    if (!workerID) return;
    setPending(request.id);
    setMessage("");
    try {
      await completeOperatorRequest(request.id, workerID);
      setMessage("Исполнитель заменён, запрос завершён");
      await load();
    } catch (reason) {
      setMessage(reason instanceof Error ? reason.message : "Не удалось завершить запрос");
    } finally {
      setPending(null);
    }
  }

  return (
    <section className="min-w-0 flex-1 px-5 pt-5 pb-8 sm:px-[22px]" aria-labelledby="requests-heading">
      <header><h1 id="requests-heading" className="text-[26px] font-extrabold text-[#4562f0]">Запросы от экспертов</h1><p className="mt-1 text-base text-[#151515]">Запросы на замену ответственного исполнителя</p></header>
      {message && <p role="status" className="mt-5 rounded-xl border border-[#b7c2fa] bg-[#eef1ff] px-4 py-3 text-sm text-[#30384f]">{message}</p>}
      <div className="mt-6 grid gap-4">
        {requests.map((request) => (
          <article key={request.id} className="rounded-[13px] border border-[#4562f0] bg-white p-5">
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div><h2 className="font-semibold text-[#4562f0]">Обращение {request.track_id}</h2><p className="mt-1 text-sm text-[#646d86]">{request.category?.name} · {request.request_type === "replace_responsible" ? "Замена исполнителя" : "Добавление соисполнителя"}</p></div>
              <span className="rounded-full bg-[#fff0b3] px-3 py-1 text-xs text-[#8a6500]">Ожидает решения</span>
            </div>
            <p className="mt-4 rounded-xl bg-[#f7f9fe] px-4 py-3 text-sm text-[#30384f]">{request.reason}</p>
            <div className="mt-4 flex flex-wrap items-end gap-3">
              <label className="min-w-[260px] flex-1 text-xs text-[#30384f]">Новый исполнитель
                <select value={selected[request.id] ?? ""} onChange={(event) => setSelected((current) => ({ ...current, [request.id]: Number(event.target.value) }))} className="mt-2 h-11 w-full rounded-xl border border-[#4562f0] bg-white px-3 text-sm outline-none">
                  <option value="">Выберите эксперта</option>
                  {(workers[request.id] ?? []).map((worker) => <option key={worker.id} value={worker.id}>{worker.full_name} — {worker.expert_group.title} ({worker.active_tickets}/{worker.max_tickets})</option>)}
                </select>
              </label>
              <Button text={pending === request.id ? "Сохраняем…" : "Заменить исполнителя"} disabled={!selected[request.id] || pending !== null} onClick={() => void complete(request)} />
            </div>
          </article>
        ))}
        {requests.length === 0 && <div className="rounded-[13px] border border-[#b7c2fa] bg-white/80 px-6 py-12 text-center text-sm text-[#646d86]">Новых запросов нет</div>}
      </div>
    </section>
  );
}
