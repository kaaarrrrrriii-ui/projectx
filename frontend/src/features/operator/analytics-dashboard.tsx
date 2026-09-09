"use client";

import Button from "@/shared/ui/button";
import { useMemo, useState } from "react";
import AnalyticsDonut, { type ChartSegment } from "./analytics-donut";
import { saveAnalyticsExport } from "./analytics-export-storage";
import { downloadAnalyticsPdf } from "./analytics-pdf";

const categorySegments: ChartSegment[] = [
  { label: "Кибербуллинг", value: 31, color: "#4562f0" },
  { label: "Семейные конфликты", value: 25, color: "#8b5cf6" },
  { label: "Травля и оскорбления", value: 22, color: "#ff7a45" },
  { label: "Давление и угрозы", value: 14, color: "#ef476f" },
  { label: "Другие темы", value: 8, color: "#35b7a0" },
];

const applicantSegments: ChartSegment[] = [
  { label: "Школьники", value: 48, color: "#4562f0" },
  { label: "Родители", value: 29, color: "#ff8a3d" },
  { label: "Студенты", value: 23, color: "#36b5a2" },
];

function MetricCard({
  title,
  value,
  suffix,
  tone = "blue",
  className = "",
}: {
  title: string;
  value: string;
  suffix: string;
  tone?: "blue" | "red" | "orange";
  className?: string;
}) {
  const valueColor = tone === "red" ? "text-[#e5141b]" : tone === "orange" ? "text-[#ef7b19]" : "text-[#4562f0]";
  return (
    <article className={"flex min-h-[150px] flex-col rounded-[15px] border border-[#7990ff] bg-white/85 p-5 " + className}>
      <h2 className="text-sm leading-5 font-medium text-[#000828]">{title}</h2>
      <div className="mt-auto text-center">
        <strong className={"text-[34px] leading-none font-extrabold " + valueColor}>{value}</strong>
        <span className={"ml-2 text-xs font-medium " + valueColor}>{suffix}</span>
      </div>
    </article>
  );
}

function inclusiveDays(dateFrom: string, dateTo: string) {
  const start = new Date(dateFrom + "T00:00:00").getTime();
  const end = new Date(dateTo + "T00:00:00").getTime();
  if (!Number.isFinite(start) || !Number.isFinite(end) || end < start) return 0;
  return Math.floor((end - start) / 86400000) + 1;
}

export default function AnalyticsDashboard() {
  const [dateFrom, setDateFrom] = useState("2026-09-01");
  const [dateTo, setDateTo] = useState("2026-09-09");
  const [isGenerating, setIsGenerating] = useState(false);
  const days = inclusiveDays(dateFrom, dateTo);
  const processed = useMemo(() => Math.max(0, Math.round((67 * days) / 9)), [days]);
  const urgentShare = 18;
  const returnedShare = 9;
  const averageAcceptance = "14 мин";
  const averageResolution = "2 ч 36 мин";

  async function downloadReport() {
    if (!days) return;
    setIsGenerating(true);
    try {
      const reportData = {
        dateFrom,
        dateTo,
        processed,
        urgentShare,
        returnedShare,
        averageAcceptance,
        averageResolution,
        categories: categorySegments,
        applicants: applicantSegments,
      };
      await downloadAnalyticsPdf(reportData);
      saveAnalyticsExport(reportData);
    } finally {
      setIsGenerating(false);
    }
  }

  return (
    <section className="min-w-0 flex-1 px-5 pt-5 pb-7 sm:px-[28px]" aria-labelledby="analytics-heading">
      <div className="mx-auto w-full max-w-[1180px]">
        <label className="relative block">
          <svg viewBox="0 0 24 24" aria-hidden="true" className="pointer-events-none absolute top-1/2 left-3 h-5 w-5 -translate-y-1/2 text-[#4562f0]" fill="none">
            <circle cx="10.5" cy="10.5" r="6.5" stroke="currentColor" strokeWidth="2" />
            <path d="m16 16 5 5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
          </svg>
          <input type="search" aria-label="Поиск обращения по трек-номеру" placeholder="Поиск обращения по треку" className="h-[42px] w-full rounded-[12px] border border-[#7990ff] bg-white px-4 pl-10 text-sm text-[#000828] outline-none placeholder:text-[#9199af] focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/15" />
        </label>

        <header className="mt-5">
          <h1 id="analytics-heading" className="text-[26px] leading-tight font-extrabold text-[#4562f0]">Аналитика</h1>
        </header>

        <section aria-labelledby="period-heading" className="mt-4">
          <h2 id="period-heading" className="text-sm font-medium text-[#30384f]">Период отчёта</h2>
          <div className="mt-2 flex flex-wrap items-end gap-3">
            <label className="grid gap-1 text-xs text-[#646d86]">
              С
              <input type="date" value={dateFrom} max={dateTo || undefined} onChange={(event) => setDateFrom(event.target.value)} className="h-10 rounded-[9px] border border-[#000828] bg-white px-3 text-sm text-[#000828] outline-none focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/15" />
            </label>
            <label className="grid gap-1 text-xs text-[#646d86]">
              По
              <input type="date" value={dateTo} min={dateFrom || undefined} onChange={(event) => setDateTo(event.target.value)} className="h-10 rounded-[9px] border border-[#000828] bg-white px-3 text-sm text-[#000828] outline-none focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/15" />
            </label>
            {!days && <p role="alert" className="pb-2 text-xs text-[#d70d14]">Укажите корректный период</p>}
          </div>
        </section>

        <section aria-label="Основные показатели" className="mt-5 grid gap-4 sm:grid-cols-2 xl:grid-cols-6">
          <MetricCard title="Обработано" value={String(processed)} suffix="обращений" className="xl:col-span-2" />
          <MetricCard title="Доля срочных обращений" value={String(urgentShare)} suffix="%" tone="red" className="xl:col-span-2" />
          <MetricCard title="Доля возвращённых обращений" value={String(returnedShare)} suffix="%" tone="orange" className="xl:col-span-2" />
          <MetricCard title="Среднее время до принятия оператором" value="14" suffix="мин" className="xl:col-span-3" />
          <MetricCard title="Среднее время до решения обращения" value="2:36" suffix="ч" className="xl:col-span-3" />
        </section>

        <section aria-label="Круговые диаграммы" className="mt-4 grid gap-4 xl:grid-cols-2">
          <AnalyticsDonut title="Доля каждой категории" segments={categorySegments} centerValue={String(processed)} centerLabel="обращений" />
          <AnalyticsDonut title="Доля каждого типа заявителя" segments={applicantSegments} centerValue={String(processed)} centerLabel="заявителей" />
        </section>

        <div className="mt-6 flex justify-center">
          <Button
            text={isGenerating ? "Формируем PDF…" : "Скачать отчёт в PDF"}
            variant="primary"
            size="default"
            disabled={isGenerating || !days}
            onClick={() => void downloadReport()}
          />
        </div>
      </div>
    </section>
  );
}
