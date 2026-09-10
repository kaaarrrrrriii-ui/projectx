"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import Button from "@/shared/ui/button";
import { staffRequest } from "@/shared/api/staff-api";
import {
  expertAppeals,
  expertRequests,
  type ExpertAppeal,
  type ExpertPriority,
  type ExpertRequest,
} from "./expert-data";

type ListMode = "queue" | "assigned" | "returns" | "requests";
type FilterKey = "category" | "applicant";

const modeContent: Record<ListMode, { title: string; subtitle: string }> = {
  queue: { title: "Очередь новых обращений", subtitle: "Сначала самые ранние" },
  assigned: { title: "Распределённые обращения", subtitle: "Обращения в вашей работе" },
  returns: { title: "Возвраты", subtitle: "Сначала самые ранние" },
  requests: { title: "Запросы", subtitle: "Сначала самые ранние" },
};

const priorityStyles: Record<ExpertPriority, { label: string; text: string; badge: string }> = {
  urgent: { label: "Срочное", text: "text-[#e5141b]", badge: "bg-[#ffd7d9] text-[#d70d14]" },
  standard: { label: "Стандартный", text: "text-[#4562f0]", badge: "bg-[#dfe6ff] text-[#4562f0]" },
  low: { label: "Низкий", text: "text-[#087f1a]", badge: "bg-[#dff2e0] text-[#087f1a]" },
};

function FilterMenu({
  label,
  options,
  selected,
  onToggle,
}: {
  label: string;
  options: string[];
  selected: string[];
  onToggle: (value: string) => void;
}) {
  return (
    <details className="group relative">
      <summary className="flex h-[38px] min-w-[166px] cursor-pointer list-none items-center justify-between gap-3 rounded-lg border border-[#808393] bg-[#f7f9fe] px-4 text-xs text-[#000828] transition-colors hover:border-[#4562f0] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] [&::-webkit-details-marker]:hidden">
        {label}
        <svg viewBox="0 0 12 8" aria-hidden="true" className="h-2 w-3 transition-transform group-open:rotate-180" fill="none"><path d="m1 1 5 5 5-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" /></svg>
      </summary>
      <div className="absolute top-11 left-0 z-40 grid min-w-[240px] gap-1 rounded-xl border border-[#9aa8f5] bg-white p-2 shadow-[0_18px_42px_rgba(0,8,40,0.16)]">
        {options.map((option) => (
          <label key={option} className="flex cursor-pointer items-center gap-3 rounded-lg px-3 py-2 text-sm text-[#30384f] hover:bg-[#eef1ff]">
            <input type="checkbox" checked={selected.includes(option)} onChange={() => onToggle(option)} className="h-4 w-4 accent-[#4562f0]" />
            <span>{option}</span>
          </label>
        ))}
      </div>
    </details>
  );
}

function PriorityCell({ row }: { row: ExpertAppeal | ExpertRequest }) {
  const priority = priorityStyles[row.priority];
  return (
    <div className="flex flex-wrap items-center gap-3">
      <Link href={`/expert/appeals/${row.track}`} className={`rounded-sm text-sm font-semibold underline-offset-4 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] ${priority.text}`}>
        {row.track}
      </Link>
      <span className={`rounded-full px-3 py-1 text-[10px] font-medium ${priority.badge}`}>{priority.label}</span>
    </div>
  );
}

export default function ExpertAppealList({ mode, initialSearch = "" }: { mode: ListMode; initialSearch?: string }) {
  const isRequests = mode === "requests";
  const fallbackRows: Array<ExpertAppeal | ExpertRequest> = isRequests ? expertRequests : expertAppeals.filter((appeal) => appeal.group === (mode === "returns" ? "return" : mode));
  const [remoteRows, setRemoteRows] = useState<Array<ExpertAppeal | ExpertRequest> | null>(null);
  const sourceRows = remoteRows ?? fallbackRows;
  useEffect(() => {
    const path = isRequests ? "/api/expert/requests?limit=100" : `/api/expert/tickets?queue=${mode === "returns" ? "returned" : mode}&limit=100`;
    staffRequest<{ items: Array<{ track_id: string; category: { name: string }; applicant_type: string; priority: ExpertPriority; status: string; created_at?: string; waiting_seconds?: number }> }>(path)
      .then(({ items }) => setRemoteRows(items.map((item) => isRequests ? ({ track: item.track_id, category: item.category?.name ?? "", applicant: item.applicant_type ?? "", priority: item.priority ?? "standard", status: item.status }) : ({ track: item.track_id, category: item.category.name, applicant: item.applicant_type, priority: item.priority, status: item.status, waiting: `${Math.floor((item.waiting_seconds ?? 0) / 60)} мин.`, submittedAt: new Date(item.created_at ?? Date.now()).toLocaleString("ru-RU"), description: "", clarifications: [], attachments: [], group: mode === "returns" ? "return" : mode as "queue" | "assigned" }))));
  }, [isRequests, mode]);
  const categories = Array.from(new Set(sourceRows.map((row) => row.category)));
  const applicants = Array.from(new Set(sourceRows.map((row) => row.applicant)));
  const [filters, setFilters] = useState<Record<FilterKey, string[]>>({ category: [], applicant: [] });
  const [search, setSearch] = useState(initialSearch);

  const rows = useMemo(() => {
    const query = search.trim().toLocaleLowerCase("ru");
    return sourceRows.filter((row) =>
      (!query || row.track.toLocaleLowerCase("ru").includes(query)) &&
      (filters.category.length === 0 || filters.category.includes(row.category)) &&
      (filters.applicant.length === 0 || filters.applicant.includes(row.applicant)),
    );
  }, [filters, search, sourceRows]);

  function toggleFilter(filter: FilterKey, value: string) {
    setFilters((current) => ({
      ...current,
      [filter]: current[filter].includes(value)
        ? current[filter].filter((item) => item !== value)
        : [...current[filter], value],
    }));
  }

  const chips = (Object.keys(filters) as FilterKey[]).flatMap((filter) => filters[filter].map((value) => ({ filter, value })));
  const hasFilters = Boolean(search.trim() || chips.length);
  const content = modeContent[mode];

  function resetFilters() {
    setFilters({ category: [], applicant: [] });
    setSearch("");
  }

  return (
    <section className="min-w-0 flex-1 px-5 py-6 sm:px-7 lg:px-9" aria-labelledby={`${mode}-heading`}>
      <div className="w-full">
        <header>
          <h1 id={`${mode}-heading`} className="text-[clamp(24px,2.4vw,32px)] font-extrabold leading-tight tracking-[-0.025em] text-[#4562f0]">{content.title}</h1>
          <p className="mt-1 text-base text-[#151515]">{content.subtitle}</p>
        </header>

        <section aria-label="Фильтры обращений" className="mt-6">
          <div className="flex flex-wrap items-center gap-2">
            <FilterMenu label="Выбор категорий" options={categories} selected={filters.category} onToggle={(value) => toggleFilter("category", value)} />
            <FilterMenu label="Выбор типа заявителя" options={applicants} selected={filters.applicant} onToggle={(value) => toggleFilter("applicant", value)} />
            <Button text="Сбросить фильтры" variant="secondary" size="small" onClick={resetFilters} disabled={!hasFilters} className="ml-auto h-[38px] rounded-lg px-4 text-xs font-normal max-[699px]:ml-0" />
          </div>

          {(chips.length > 0 || search) && (
            <div className="mt-4 flex flex-wrap gap-2" aria-live="polite">
              {search && (
                <button type="button" onClick={() => setSearch("")} className="inline-flex h-8 cursor-pointer items-center gap-2 rounded-full bg-[#dfe6ff] px-3 text-xs text-[#4562f0] hover:bg-[#d4ddff]">
                  Поиск: {search} <span aria-hidden="true">×</span>
                </button>
              )}
              {chips.map(({ filter, value }) => (
                <button key={`${filter}-${value}`} type="button" onClick={() => toggleFilter(filter, value)} className="inline-flex h-8 cursor-pointer items-center gap-2 rounded-full bg-[#dfe6ff] px-3 text-xs text-[#4562f0] hover:bg-[#d4ddff]">
                  {value} <span aria-hidden="true">×</span>
                </button>
              ))}
            </div>
          )}
        </section>

        {rows.length > 0 ? (
          <section aria-label="Список обращений" className="mt-7 overflow-x-auto rounded-[15px] border border-[#7990ff] bg-white/90 shadow-[0_12px_30px_rgba(69,98,240,0.07)]">
            <table className="w-full min-w-[760px] table-fixed border-collapse">
              <caption className="sr-only">{content.title}</caption>
              <thead className="bg-[#dfe6ff] text-[#4562f0]">
                <tr className="h-[52px]">
                  <th scope="col" className="w-[29%] border-r border-[#7990ff] px-4 text-center font-normal">Трек-номер</th>
                  <th scope="col" className="w-[25%] border-r border-[#7990ff] px-4 text-center font-normal">Категория</th>
                  {isRequests ? (
                    <th scope="col" className="w-[28%] border-r border-[#7990ff] px-4 text-center font-normal">Статус</th>
                  ) : (
                    <th scope="col" className="w-[22%] border-r border-[#7990ff] px-4 text-center font-normal">Тип заявителя</th>
                  )}
                  <th scope="col" className="px-4 text-center font-normal">{isRequests ? "Тип заявителя" : ""}<span className="sr-only">{isRequests ? "" : "Действия"}</span></th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => (
                  <tr key={row.track} className="min-h-[58px] border-t border-[#7990ff] text-[#20263a] transition-colors hover:bg-[#f5f7ff]">
                    <td className="border-r border-[#7990ff] px-4 py-3"><PriorityCell row={row} /></td>
                    <td className="border-r border-[#7990ff] px-4 py-3 text-center text-sm">{row.category}</td>
                    <td className="border-r border-[#7990ff] px-4 py-3 text-center text-sm">{isRequests ? row.status : row.applicant}</td>
                    <td className="px-4 py-3 text-center">
                      {isRequests ? (
                        <span className="text-sm">{row.applicant}</span>
                      ) : (
                        <Link href={`/expert/appeals/${row.track}`} className="inline-flex min-h-8 min-w-[110px] items-center justify-center rounded-lg bg-[#4562f0] px-4 text-xs font-medium text-white transition-all hover:bg-[#374ecc] hover:shadow-[0_7px_16px_rgba(69,98,240,0.22)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]">Открыть</Link>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </section>
        ) : (
          <div className="mt-7 rounded-[15px] border border-[#b8c2fa] bg-white/85 px-6 py-12 text-center text-sm text-[#646d86]">
            По выбранным фильтрам обращений не найдено.
          </div>
        )}
      </div>
    </section>
  );
}
