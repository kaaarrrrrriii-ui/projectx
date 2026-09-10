"use client";

import Button from "@/shared/ui/button";
import { useMemo, useState } from "react";
import AnalyticsDonut, { type ChartSegment } from "./analytics-donut";
import { saveAnalyticsExport } from "./analytics-export-storage";
import { downloadAnalyticsPdf } from "./analytics-pdf";

const statusSegments: ChartSegment[] = [
  { label: "Закрыты", value: 52, color: "#4562f0" },
  { label: "В работе", value: 30, color: "#ff7a45" },
  { label: "Возвращены", value: 18, color: "#e5141b" },
];

const applicantSegments: ChartSegment[] = [
  { label: "Школьники", value: 48, color: "#4562f0" },
  { label: "Родители", value: 29, color: "#ff8a3d" },
  { label: "Студенты", value: 23, color: "#36b5a2" },
];

const categorySegments: ChartSegment[] = [
  { label: "Кибербуллинг", value: 31, color: "#4562f0" },
  { label: "Семейные конфликты", value: 27, color: "#8b5cf6" },
  { label: "Травля", value: 24, color: "#ff7a45" },
  { label: "Другое", value: 18, color: "#35b7a0" },
];

function MetricCard({
  title,
  value,
  unit,
  description,
  tone = "blue",
  className = "",
}: {
  title: string;
  value: string;
  unit?: string;
  description?: string;
  tone?: "blue" | "red";
  className?: string;
}) {
  const valueColor = tone === "red" ? "text-[#e5141b]" : "text-[#4562f0]";

  return (
    <article className={`flex min-h-[102px] flex-col items-center rounded-[12px] border border-[#7990ff] bg-white/85 px-3 py-3 text-center ${className}`}>
      <h2 className="max-w-[230px] text-[15px] leading-[18px] font-medium text-[#000828]">{title}</h2>
      <div className="mt-auto pt-2">
        <strong className={`block text-[35px] leading-none font-extrabold ${valueColor}`}>
          {value}{unit === "%" ? unit : ""}
        </strong>
        {unit && unit !== "%" && <span className="mt-1 block text-[11px] text-[#30384f]">{unit}</span>}
        {description && <span className="mt-1 block text-[10px] leading-4 text-[#30384f]">{description}</span>}
      </div>
    </article>
  );
}

function inclusiveDays(dateFrom: string, dateTo: string) {
  const start = new Date(`${dateFrom}T00:00:00`).getTime();
  const end = new Date(`${dateTo}T00:00:00`).getTime();
  if (!Number.isFinite(start) || !Number.isFinite(end) || end < start) return 0;
  return Math.floor((end - start) / 86400000) + 1;
}

function formatShortDate(value: string) {
  const [year, month, day] = value.split("-");
  return year && month && day ? `${day}.${month}.${year}` : "—";
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
  const urgentCount = Math.round((processed * urgentShare) / 100);
  const returnedCount = Math.round((processed * returnedShare) / 100);

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
    <section className="min-w-0 flex-1 px-5 pt-8 pb-5 sm:px-[22px]" aria-labelledby="analytics-heading">
      <div className="w-full">
        <header>
          <h1 id="analytics-heading" className="text-[24px] leading-tight font-extrabold tracking-[-0.02em] text-[#4562f0]">Аналитика</h1>
        </header>

        <section aria-labelledby="period-heading" className="mt-4">
          <h2 id="period-heading" className="text-[12px] font-medium text-[#30384f]">Выбор даты</h2>
          <details className="group relative mt-1.5 w-fit">
            <summary className="flex h-[28px] min-w-[132px] cursor-pointer list-none items-center justify-between gap-3 rounded-[6px] border border-[#808393] bg-[#f7f9fe] px-3 text-[10px] text-[#000828] outline-none hover:border-[#4562f0] focus-visible:ring-2 focus-visible:ring-[#4562f0]/20 [&::-webkit-details-marker]:hidden">
              {formatShortDate(dateFrom)} – {formatShortDate(dateTo)}
              <svg viewBox="0 0 12 8" aria-hidden="true" className="h-1.5 w-2.5 transition-transform group-open:rotate-180" fill="none">
                <path d="m1 1 5 5 5-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            </summary>
            <div className="absolute top-9 left-0 z-30 flex min-w-[295px] gap-3 rounded-[12px] border border-[#7990ff] bg-white p-3 shadow-[0_14px_35px_rgba(0,8,40,0.14)] max-[419px]:min-w-[calc(100vw-40px)] max-[419px]:flex-col">
              <label className="grid flex-1 gap-1 text-[11px] text-[#646d86]">
                С
                <input type="date" value={dateFrom} max={dateTo || undefined} onChange={(event) => setDateFrom(event.target.value)} className="h-9 rounded-[8px] border border-[#000828] bg-white px-2 text-xs text-[#000828] outline-none focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/15" />
              </label>
              <label className="grid flex-1 gap-1 text-[11px] text-[#646d86]">
                По
                <input type="date" value={dateTo} min={dateFrom || undefined} onChange={(event) => setDateTo(event.target.value)} className="h-9 rounded-[8px] border border-[#000828] bg-white px-2 text-xs text-[#000828] outline-none focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/15" />
              </label>
            </div>
          </details>
          {!days && <p role="alert" className="mt-2 text-xs text-[#d70d14]">Укажите корректный период</p>}
        </section>

        <section aria-label="Основные показатели" className="mt-8 grid gap-[13px] md:grid-cols-[112px_minmax(0,1fr)_minmax(0,1fr)]">
          <MetricCard title="Отработано" value={String(processed)} unit="обращений" />
          <MetricCard title="Доля срочных обращений" value={String(urgentShare)} unit="%" tone="red" description={`${urgentCount} из ${processed} обращений`} />
          <MetricCard title="Доля возвращённых обращений" value={String(returnedShare)} unit="%" description={`${returnedCount} из ${processed} обращений`} />
        </section>

        <section aria-label="Среднее время обработки" className="mt-5 grid gap-[13px] md:grid-cols-3">
          <MetricCard title="Среднее время до принятия оператором" value="14" unit="мин" className="min-h-[118px]" />
          <MetricCard title="Среднее время до первого ответа" value="36" unit="мин" className="min-h-[118px]" />
          <MetricCard title="Среднее время до закрытия" value="67" unit="мин" className="min-h-[118px]" />
        </section>

        <section aria-label="Круговые диаграммы" className="mt-5 grid gap-[13px] md:grid-cols-3">
          <AnalyticsDonut title="Статусы" segments={statusSegments} />
          <AnalyticsDonut title="Тип заявителя" segments={applicantSegments} />
          <AnalyticsDonut title="Категории" segments={categorySegments} />
        </section>

        <div className="mt-14 flex justify-center md:mt-[88px]">
          <Button
            text={isGenerating ? "Формируем выгрузку…" : "Скачать полную выгрузку"}
            variant="primary"
            size="small"
            disabled={isGenerating || !days}
            onClick={() => void downloadReport()}
            className="h-[34px] min-w-[190px] rounded-[8px] px-5 text-[11px] font-normal"
          />
        </div>
      </div>
    </section>
  );
}
