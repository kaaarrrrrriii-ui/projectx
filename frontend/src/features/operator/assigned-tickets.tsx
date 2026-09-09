"use client";

import { appealCategories } from "@/features/appeal/categories";
import Button from "@/shared/ui/button";
import Input from "@/shared/ui/input";
import Link from "next/link";
import { useMemo, useState } from "react";
import { applicantLabels, operatorTickets, type TicketPriority } from "./tickets";

type FilterKey = "priority" | "category" | "applicant";
type FilterOption = { value: string; label: string };

const priorityOptions: FilterOption[] = [
  { value: "urgent", label: "Срочное" },
  { value: "standard", label: "Стандартное" },
  { value: "low", label: "Низкое" },
];

const filterOptions: Record<FilterKey, FilterOption[]> = {
  priority: priorityOptions,
  category: appealCategories.map((category) => ({
    value: category,
    label: category.charAt(0).toUpperCase() + category.slice(1),
  })),
  applicant: Object.entries(applicantLabels).map(([value, label]) => ({ value, label })),
};

const filterLabels: Record<FilterKey, string> = {
  priority: "Выбор приоритетов",
  category: "Выбор категорий",
  applicant: "Выбор типа заявителя",
};

const priorityStyles: Record<TicketPriority, { text: string; badge: string }> = {
  urgent: { text: "text-[#d70d14]", badge: "bg-[#ffdfe0] text-[#d70d14]" },
  standard: { text: "text-[#4562f0]", badge: "bg-[#dfe6ff] text-[#4562f0]" },
  low: { text: "text-[#087f1a]", badge: "bg-[#dff2e0] text-[#087f1a]" },
};

const assignedTickets = operatorTickets.slice(0, 2);

function FilterMenu({
  filter,
  selected,
  onToggle,
}: {
  filter: FilterKey;
  selected: string[];
  onToggle: (value: string) => void;
}) {
  return (
    <details className="group relative">
      <summary className="flex h-[33px] min-w-[157px] cursor-pointer list-none items-center justify-between gap-3 rounded-[7px] border border-[#000828] bg-white px-4 text-xs text-[#000828] hover:border-[#4562f0] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] [&::-webkit-details-marker]:hidden">
        {filterLabels[filter]}
        <svg viewBox="0 0 12 8" aria-hidden="true" className="h-2 w-3 transition-transform group-open:rotate-180" fill="none">
          <path d="m1 1 5 5 5-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </summary>
      <div className="absolute top-10 left-0 z-40 grid max-h-72 min-w-[260px] gap-1 overflow-y-auto rounded-xl border border-[#8799f8] bg-white p-3 shadow-[0_14px_35px_rgba(0,8,40,0.14)]">
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

export default function AssignedTickets() {
  const [search, setSearch] = useState("");
  const [filters, setFilters] = useState<Record<FilterKey, string[]>>({
    priority: [],
    category: [],
    applicant: [],
  });

  const filteredTickets = useMemo(() => {
    const query = search.trim().toLocaleLowerCase("ru");
    return assignedTickets.filter((ticket) =>
      (!query || ticket.track.toLocaleLowerCase("ru").includes(query)) &&
      (filters.priority.length === 0 || filters.priority.includes(ticket.priority)) &&
      (filters.category.length === 0 || filters.category.includes(ticket.category)) &&
      (filters.applicant.length === 0 || filters.applicant.includes(ticket.applicant)),
    );
  }, [search, filters]);

  function toggleFilter(filter: FilterKey, value: string) {
    setFilters((current) => ({
      ...current,
      [filter]: current[filter].includes(value)
        ? current[filter].filter((item) => item !== value)
        : [...current[filter], value],
    }));
  }

  function resetFilters() {
    setSearch("");
    setFilters({ priority: [], category: [], applicant: [] });
  }

  const chips = (Object.keys(filters) as FilterKey[]).flatMap((filter) =>
    filters[filter].map((value) => ({ filter, value })),
  );
  const hasFilters = Boolean(search || chips.length);

  function optionLabel(filter: FilterKey, value: string) {
    return filterOptions[filter].find((option) => option.value === value)?.label ?? value;
  }

  return (
    <section className="min-w-0 flex-1 px-5 pt-5 pb-0 sm:px-[26px]" aria-labelledby="assigned-heading">
      <div className="mx-auto w-full max-w-[1180px]">
        <label className="relative block">
          <svg viewBox="0 0 24 24" aria-hidden="true" className="pointer-events-none absolute top-1/2 left-3 z-10 h-5 w-5 -translate-y-1/2 text-[#4562f0]" fill="none">
            <circle cx="10.5" cy="10.5" r="6.5" stroke="currentColor" strokeWidth="2" />
            <path d="m16 16 5 5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
          </svg>
          <Input
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            aria-label="Поиск распределённого обращения по трек-номеру"
            placeholder="Поиск обращения по треку"
            className="!h-[42px] !w-full !rounded-[12px] !border-[#4562f0] !bg-white !pr-3 !pl-[42px] !text-left !text-base placeholder:!text-[#9196a7] focus:!bg-white focus:!text-[#000828]"
          />
        </label>

        <header className="mt-4">
          <h1 id="assigned-heading" className="text-[26px] leading-tight font-extrabold tracking-[-0.02em] text-[#4562f0]">
            Распределённые обращения
          </h1>
          <p className="mt-1 text-base text-[#151515]">Сначала самые ранние</p>
        </header>

        <section aria-label="Фильтры распределённых обращений" className="mt-5">
          <div className="flex flex-wrap items-center gap-2">
            {(Object.keys(filterOptions) as FilterKey[]).map((filter) => (
              <FilterMenu key={filter} filter={filter} selected={filters[filter]} onToggle={(value) => toggleFilter(filter, value)} />
            ))}
            <Button text="Сбросить фильтры" variant="secondary" size="small" onClick={resetFilters} disabled={!hasFilters} className="ml-auto h-[33px] rounded-[7px] px-4 text-xs font-normal max-[699px]:ml-0" />
          </div>

          <div className="mt-3 flex min-h-8 flex-wrap gap-2" aria-live="polite">
            {chips.map(({ filter, value }) => (
              <button key={filter + value} type="button" onClick={() => toggleFilter(filter, value)} className="inline-flex h-7 cursor-pointer items-center gap-1.5 rounded-full bg-[#dfe6ff] px-3 text-[11px] text-[#4562f0] hover:bg-[#d4ddff]">
                {optionLabel(filter, value)} <span aria-hidden="true">×</span>
              </button>
            ))}
          </div>
        </section>

        <section aria-label="Список распределённых обращений" className="mt-3 overflow-x-auto rounded-[13px] border border-[#4562f0] bg-white/85">
          <table className="w-full min-w-[760px] table-fixed border-collapse">
            <thead className="bg-[#dfe6ff] text-[#4562f0]">
              <tr className="h-[52px]">
                <th className="w-[31%] border-r border-[#4562f0] px-4 text-center font-normal">Трек-номер</th>
                <th className="w-[29%] border-r border-[#4562f0] px-4 text-center font-normal">Категория</th>
                <th className="w-[23%] border-r border-[#4562f0] px-4 text-center font-normal">Тип заявителя</th>
                <th className="w-[17%] px-4 text-center font-normal">Ожидает</th>
              </tr>
            </thead>
            <tbody>
              {filteredTickets.map((ticket) => {
                const priority = priorityStyles[ticket.priority];
                return (
                  <tr key={ticket.track} className="h-[52px] border-t border-[#4562f0] hover:bg-[#f7f8ff]">
                    <td className="border-r border-[#4562f0] px-3">
                      <div className="flex items-center gap-3">
                        <Link href={"/operator/queueNew.tsx/" + encodeURIComponent(ticket.track)} className={"rounded-sm text-sm font-medium hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] " + priority.text}>
                          {ticket.track}
                        </Link>
                        <span className={"rounded-full px-3 py-1 text-[10px] " + priority.badge}>
                          {optionLabel("priority", ticket.priority)}
                        </span>
                      </div>
                    </td>
                    <td className="border-r border-[#4562f0] px-3 text-center text-sm">{optionLabel("category", ticket.category)}</td>
                    <td className="border-r border-[#4562f0] px-3 text-center text-sm">{optionLabel("applicant", ticket.applicant)}</td>
                    <td className="px-3 text-center text-sm">{ticket.waiting}</td>
                  </tr>
                );
              })}

              {filteredTickets.length === 0 && (
                <tr className="h-24 border-t border-[#4562f0]">
                  <td colSpan={4} className="text-center text-sm text-[#646d86]">Обращения не найдены</td>
                </tr>
              )}

            </tbody>
          </table>
        </section>
      </div>
    </section>
  );
}
