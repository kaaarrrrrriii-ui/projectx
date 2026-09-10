"use client";

import { appealCategories } from "@/features/appeal/categories";
import Button from "@/shared/ui/button";
import Link from "next/link";
import { useMemo, useState } from "react";
import { returnedTickets, type TicketPriority } from "./tickets";

type ComplaintFilter = "priority" | "category" | "expert";

const priorityOptions = [
  { value: "urgent", label: "Срочное" },
  { value: "standard", label: "Стандартное" },
  { value: "low", label: "Низкое" },
];

const expertOptions = ["Байкова Е.С.", "Петрова А.В.", "Соколова М.А."];

const filterOptions: Record<ComplaintFilter, Array<{ value: string; label: string }>> = {
  priority: priorityOptions,
  category: appealCategories.map((category) => ({
    value: category,
    label: category.charAt(0).toUpperCase() + category.slice(1),
  })),
  expert: expertOptions.map((expert) => ({ value: expert, label: expert })),
};

const filterLabels: Record<ComplaintFilter, string> = {
  priority: "Выбор приоритетов",
  category: "Выбор категорий",
  expert: "Выбор эксперта",
};

const priorityStyles: Record<TicketPriority, { text: string; badge: string }> = {
  urgent: { text: "text-[#e5141b]", badge: "bg-[#f7b6b8] text-[#d70d14]" },
  standard: { text: "text-[#4562f0]", badge: "bg-[#dfe6ff] text-[#4562f0]" },
  low: { text: "text-[#087f1a]", badge: "bg-[#dff2e0] text-[#087f1a]" },
};

const complaints = returnedTickets.map((ticket, index) => ({
  ...ticket,
  expert: expertOptions[index % expertOptions.length],
}));

function FilterMenu({
  filter,
  selected,
  onToggle,
}: {
  filter: ComplaintFilter;
  selected: string[];
  onToggle: (value: string) => void;
}) {
  return (
    <details className="group relative">
      <summary className="flex h-[30px] min-w-[145px] cursor-pointer list-none items-center justify-between gap-3 rounded-[6px] border border-[#000828] bg-white px-4 text-[11px] text-[#000828] hover:border-[#4562f0] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] [&::-webkit-details-marker]:hidden">
        {filterLabels[filter]}
        <svg viewBox="0 0 12 8" aria-hidden="true" className="h-2 w-3 transition-transform group-open:rotate-180" fill="none"><path d="m1 1 5 5 5-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" /></svg>
      </summary>
      <div className="absolute top-9 left-0 z-40 grid max-h-72 min-w-[250px] gap-1 overflow-y-auto rounded-xl border border-[#8799f8] bg-white p-3 shadow-[0_14px_35px_rgba(0,8,40,0.14)]">
        {filterOptions[filter].map((option) => (
          <label key={option.value} className="flex cursor-pointer items-center gap-3 rounded-lg px-2 py-2 text-sm hover:bg-[#eef1ff]">
            <input type="checkbox" checked={selected.includes(option.value)} onChange={() => onToggle(option.value)} className="h-4 w-4 accent-[#4562f0]" />
            <span>{option.label}</span>
          </label>
        ))}
      </div>
    </details>
  );
}

export default function ReturnsTickets() {
  const [filters, setFilters] = useState<Record<ComplaintFilter, string[]>>({
    priority: [],
    category: [],
    expert: [],
  });

  const filteredComplaints = useMemo(() => complaints.filter((ticket) =>
    (filters.priority.length === 0 || filters.priority.includes(ticket.priority)) &&
    (filters.category.length === 0 || filters.category.includes(ticket.category)) &&
    (filters.expert.length === 0 || filters.expert.includes(ticket.expert)),
  ), [filters]);

  const chips = (Object.keys(filters) as ComplaintFilter[]).flatMap((filter) =>
    filters[filter].map((value) => ({ filter, value })),
  );

  function toggleFilter(filter: ComplaintFilter, value: string) {
    setFilters((current) => ({
      ...current,
      [filter]: current[filter].includes(value)
        ? current[filter].filter((item) => item !== value)
        : [...current[filter], value],
    }));
  }

  function resetFilters() {
    setFilters({ priority: [], category: [], expert: [] });
  }

  function optionLabel(filter: ComplaintFilter, value: string) {
    return filterOptions[filter].find((option) => option.value === value)?.label ?? value;
  }

  return (
    <section className="min-w-0 flex-1 px-5 pt-6 pb-0 sm:px-[34px]" aria-labelledby="complaints-heading">
      <div className="w-full max-w-[828px]">
        <header>
          <h1 id="complaints-heading" className="text-[24px] leading-tight font-extrabold tracking-[-0.02em] text-[#4562f0]">Жалобы</h1>
          <p className="mt-1 text-base text-[#151515]">Сначала самые ранние</p>
        </header>

        <section aria-label="Фильтры жалоб" className="mt-5">
          <div className="flex flex-wrap items-center gap-[5px]">
            {(Object.keys(filterOptions) as ComplaintFilter[]).map((filter) => (
              <FilterMenu key={filter} filter={filter} selected={filters[filter]} onToggle={(value) => toggleFilter(filter, value)} />
            ))}
            <Button text="Сбросить фильтры" variant="secondary" size="small" onClick={resetFilters} disabled={chips.length === 0} className="ml-auto h-[30px] rounded-[6px] px-4 text-[11px] font-normal max-[699px]:ml-0" />
          </div>

          <div className="mt-4 flex min-h-[28px] flex-wrap gap-2" aria-live="polite">
            {chips.map(({ filter, value }) => (
              <button key={filter + value} type="button" onClick={() => toggleFilter(filter, value)} className="inline-flex h-7 cursor-pointer items-center gap-1.5 rounded-full bg-[#dfe6ff] px-3 text-[10px] text-[#4562f0] hover:bg-[#d4ddff]">
                {optionLabel(filter, value)} <span aria-hidden="true">×</span>
              </button>
            ))}
          </div>
        </section>

        <section aria-label="Список жалоб" className="mt-[18px] overflow-x-auto rounded-[13px] border border-[#4562f0] bg-white/85">
          <table className="w-full min-w-[720px] table-fixed border-collapse">
            <thead className="bg-[#dfe6ff] text-[#4562f0]">
              <tr className="h-[49px]">
                <th className="w-[30%] border-r border-[#4562f0] px-4 text-center font-normal">Трек-номер</th>
                <th className="w-[25%] border-r border-[#4562f0] px-4 text-center font-normal">Категория</th>
                <th className="w-[23%] border-r border-[#4562f0] px-4 text-center font-normal">Эксперт</th>
                <th className="w-[22%] px-4"><span className="sr-only">Действия</span></th>
              </tr>
            </thead>
            <tbody>
              {filteredComplaints.map((ticket) => {
                const priority = priorityStyles[ticket.priority];
                return (
                  <tr key={ticket.track} className="h-[49px] border-t border-[#4562f0] hover:bg-[#f7f8ff]">
                    <td className="border-r border-[#4562f0] px-3">
                      <div className="flex items-center gap-3">
                        <Link href={`/operator/queueNew.tsx/${encodeURIComponent(ticket.track)}`} className={`rounded-sm text-sm font-medium hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] ${priority.text}`}>{ticket.track}</Link>
                        <span className={`rounded-full px-3 py-1 text-[10px] ${priority.badge}`}>{optionLabel("priority", ticket.priority)}</span>
                      </div>
                    </td>
                    <td className="border-r border-[#4562f0] px-3 text-center text-sm">{ticket.category}</td>
                    <td className="border-r border-[#4562f0] px-3 text-center text-sm">{ticket.expert}</td>
                    <td className="px-3 text-center"><Link href={`/operator/queueNew.tsx/${encodeURIComponent(ticket.track)}`} className="inline-flex min-w-[78px] items-center justify-center rounded-[7px] bg-[#4562f0] px-4 py-2 text-xs text-white transition-colors hover:bg-[#4f71fc] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]">Открыть</Link></td>
                  </tr>
                );
              })}

              {filteredComplaints.length === 0 && <tr className="h-24 border-t border-[#4562f0]"><td colSpan={4} className="text-center text-sm text-[#646d86]">Жалобы не найдены</td></tr>}

              {filteredComplaints.length > 0 && Array.from({ length: Math.max(0, 8 - filteredComplaints.length) }).map((_, index) => (
                <tr key={`empty-${index}`} aria-hidden="true" className="h-[49px] border-t border-[#4562f0]"><td className="border-r border-[#4562f0]" /><td className="border-r border-[#4562f0]" /><td className="border-r border-[#4562f0]" /><td /></tr>
              ))}
            </tbody>
          </table>
        </section>
      </div>
    </section>
  );
}
