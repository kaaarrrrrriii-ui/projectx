"use client";

import MetricCard from "./metric-card";
import WeeklyChart from "./weekly-chart";
import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getAdminDashboard, getOperatorDashboard } from "@/shared/api/staff-api";

export default function OperatorDashboard({ mode = "operator" }: { mode?: "operator" | "admin" }) {
  const router = useRouter();
  const [search, setSearch] = useState("");
  const [dashboard, setDashboard] = useState({ new_count: 0, assigned_count: 0, returned_count: 0, crisis_count: 0, overdue_count: 0 });
  const [adminDashboard, setAdminDashboard] = useState({ new_count: 0, active_count: 0, urgent_count: 0, returned_count: 0, routing_issue_count: 0, experts_at_capacity_count: 0 });

  useEffect(() => {
    if (mode === "admin") getAdminDashboard().then(setAdminDashboard).catch(() => undefined);
    else getOperatorDashboard().then(setDashboard).catch(() => undefined);
  }, [mode]);

  function submitSearch(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const track = search.trim().toLocaleUpperCase("ru");
    if (track) router.push(mode === "admin" ? `/admin/queue/${encodeURIComponent(track)}` : `/operator/queueNew.tsx/${encodeURIComponent(track)}`);
  }

  const metrics = mode === "admin" ? [
    { title: "Новые обращения", value: String(adminDashboard.new_count), description: "Требуют внимания", icon: "mail" as const, tone: "blue" as const },
    { title: "Срочные", value: String(adminDashboard.urgent_count), description: "Приоритетные", icon: "urgent" as const, tone: "red" as const },
    { title: "Активные", value: String(adminDashboard.active_count), description: "В работе", icon: "clock" as const, tone: "orange" as const },
    { title: "Возвраты", value: String(adminDashboard.returned_count), description: "Повторные", icon: "return" as const, tone: "blue" as const },
  ] : [
    { title: "Новые обращения", value: String(dashboard.new_count), badge: `+${dashboard.new_count} обращений`, description: "Требуют распределения", icon: "mail" as const, tone: "blue" as const },
    { title: "Кризисные", value: String(dashboard.crisis_count), description: "Нуждаются в решении", icon: "urgent" as const, tone: "red" as const },
    { title: "Просроченные", value: String(dashboard.overdue_count), description: "Требуют внимания", icon: "clock" as const, tone: "orange" as const },
    { title: "Возвраты", value: String(dashboard.returned_count), description: "Повторные рассмотрения", icon: "return" as const, tone: "blue" as const },
  ];

  return (
    <div className="min-w-0 flex-1 px-5 py-4 sm:px-6">
      <div className="w-full">
        <header>
          <h1 className="text-[clamp(22px,2vw,30px)] font-bold leading-tight text-[#4562f0]">{mode === "admin" ? "Панель администратора" : "Здравствуйте, Олег!"}</h1>
          <p className="mt-1 text-xs text-[#4f5873]">Спасибо, что помогаете. Ваша работа важна.</p>
        </header>

        <form onSubmit={submitSearch} className="mt-4 flex h-[42px] items-center gap-3 rounded-xl border border-[#7990ff] bg-white px-3 focus-within:ring-3 focus-within:ring-[#4562f0]/15">
          <svg viewBox="0 0 24 24" aria-hidden="true" className="h-5 w-5 shrink-0 text-[#4562f0]" fill="none"><circle cx="10.5" cy="10.5" r="6.5" stroke="currentColor" strokeWidth="2" /><path d="m16 16 5 5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg>
          <span className="sr-only">Поиск обращения по трек-номеру</span>
          <input value={search} onChange={(event) => setSearch(event.target.value)} className="h-full min-w-0 flex-1 border-0 bg-transparent text-[17px] text-[#000828] outline-none placeholder:text-[#9199af]" placeholder="Поиск обращения по треку" />
          <button type="submit" className="cursor-pointer rounded-lg bg-[#4562f0] px-3 py-1.5 text-xs text-white">Открыть</button>
        </form>

        <section aria-label="Сводка обращений" className="mt-[21px] grid grid-cols-2 gap-4 lg:grid-cols-4 lg:gap-10">
          {metrics.map((metric) => <MetricCard key={metric.title} {...metric} />)}
        </section>

        <div className="mt-[30px] grid gap-5 lg:grid-cols-[194px_minmax(0,1fr)] lg:gap-[43px]">
          <section aria-labelledby="processed-title" className="flex min-h-[181px] flex-col items-center justify-center rounded-[13px] border border-[#7990ff] bg-white p-6 text-center">
            <h2 id="processed-title" className="text-base font-medium text-[#000828]">Отработано</h2>
            <strong className="mt-4 text-5xl font-extrabold text-[#4562f0]">67</strong>
            <p className="mt-4 text-xs text-[#4f5873]">обращений</p>
          </section>
          <WeeklyChart />
        </div>
      </div>
    </div>
  );
}
