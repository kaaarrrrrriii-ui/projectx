"use client";

import Button from "@/shared/ui/button";
import Checkbox from "@/shared/ui/checkbox";
import Input from "@/shared/ui/input";
import { appealCategories } from "@/features/appeal/categories";
import { operatorTickets, type OperatorTicket } from "@/features/operator/tickets";
import Link from "next/link";
import { useEffect, useMemo, useRef, useState } from "react";

type FilterKey = "priority" | "status" | "category" | "applicant";

type FilterOption = {
  value: string;
  label: string;
  color?: "red" | "yellow" | "green";
};

const priorityOptions: FilterOption[] = [
  { value: "urgent", label: "Срочный", color: "red" },
  { value: "standard", label: "Стандартный", color: "yellow" },
  { value: "low", label: "Низкий", color: "green" },
];

const statusOptions: FilterOption[] = [
  { value: "new", label: "Новое" },
  { value: "assigned", label: "Распределено" },
  { value: "in-progress", label: "В работе" },
  { value: "clarification", label: "Нужно уточнение" },
  { value: "answer-ready", label: "Ответ готов" },
  { value: "returned", label: "Возвращено" },
  { value: "rejected", label: "Отклонено" },
  { value: "closed", label: "Закрыто без ответа" },
];

const categoryOptions: FilterOption[] = appealCategories.map((category) => ({
  value: category,
  label: category.charAt(0).toUpperCase() + category.slice(1),
}));

const applicantOptions: FilterOption[] = [
  { value: "schoolchild", label: "Школьник" },
  { value: "parent", label: "Родитель" },
  { value: "student", label: "Студент" },
];

const filterGroups: Record<FilterKey, FilterOption[]> = {
  priority: priorityOptions,
  status: statusOptions,
  category: categoryOptions,
  applicant: applicantOptions,
};

const filterLabels: Record<FilterKey, string> = {
  priority: "Выбор приоритетов",
  status: "Выбор статусов",
  category: "Выбор категорий",
  applicant: "Выбор типа заявителя",
};

const priorityColors = {
  red: "border-[#ed7777] bg-[#f5b1b1] text-[#d70d14]",
  yellow: "border-[#e4bd52] bg-[#fff0b3] text-[#9a6500]",
  green: "border-[#72bf78] bg-[#adddad] text-[#069b1e]",
};

const rowColors: Record<string, string> = {
  urgent: "bg-[#f5b1b1]/90",
  standard: "bg-[#fff0b3]/55",
  low: "bg-[#adddad]/45",
};

function getOptionLabel(filter: FilterKey, value: string) {
  return filterGroups[filter].find((option) => option.value === value)?.label ?? value;
}

function FilterDropdown({
  filter,
  selected,
  open,
  onOpen,
  onToggle,
}: {
  filter: FilterKey;
  selected: string[];
  open: boolean;
  onOpen: () => void;
  onToggle: (value: string) => void;
}) {
  const options = filterGroups[filter];
  const menuId = `${filter}-filter-menu`;

  return (
    <div className="relative" data-filter-root>
      <button
        type="button"
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-controls={menuId}
        onClick={onOpen}
        className="flex h-[33px] min-w-[157px] cursor-pointer items-center justify-between gap-3 rounded-[7px] border border-[#000828] bg-white px-4 text-[12px] leading-4 text-[#000828] transition-colors hover:border-[#4562f0] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]"
      >
        <span>{filterLabels[filter]}</span>
        <svg
          viewBox="0 0 12 8"
          aria-hidden="true"
          className={`h-2 w-3 transition-transform ${open ? "rotate-180" : ""}`}
          fill="none"
        >
          <path d="m1 1 5 5 5-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </button>

      {open && (
        <div
          id={menuId}
          role="listbox"
          aria-multiselectable="true"
          className={`absolute top-[41px] left-0 z-50 grid max-h-[330px] min-w-full gap-2 overflow-y-auto rounded-[13px] border border-[#8799f8] bg-white p-3 shadow-[0_14px_35px_rgba(0,8,40,0.14)] ${
            filter === "priority"
              ? "w-[250px]"
              : filter === "category"
                ? "w-[290px]"
                : "w-[235px]"
          }`}
        >
          {options.map((option) => {
            const isSelected = selected.includes(option.value);
            const priorityClass = option.color
              ? priorityColors[option.color]
              : isSelected
                ? "border-[#4562f0] bg-[#dfe6ff] text-[#4562f0]"
                : "border-[#d8ddea] bg-white text-[#000828]";

            return (
              <button
                key={option.value}
                type="button"
                role="option"
                aria-selected={isSelected}
                onClick={() => onToggle(option.value)}
                className={`flex min-h-[42px] cursor-pointer items-center justify-between gap-3 rounded-[10px] border px-4 text-left text-[14px] font-medium transition-[filter,box-shadow] hover:brightness-[0.97] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] ${priorityClass} ${
                  isSelected ? "shadow-[inset_0_0_0_1px_currentColor]" : ""
                }`}
              >
                <span>{option.label}</span>
                <span aria-hidden="true" className={`text-base ${isSelected ? "opacity-100" : "opacity-0"}`}>
                  ✓
                </span>
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}

export default function QueueNew() {
  const filtersRef = useRef<HTMLDivElement>(null);
  const [search, setSearch] = useState("");
  const [openFilter, setOpenFilter] = useState<FilterKey | null>(null);
  const [selectedFilters, setSelectedFilters] = useState<Record<FilterKey, string[]>>({
    priority: [],
    status: [],
    category: [],
    applicant: [],
  });
  const [ticketItems, setTicketItems] = useState<OperatorTicket[]>(operatorTickets);
  const [selectedTracks, setSelectedTracks] = useState<string[]>([]);

  useEffect(() => {
    function handlePointerDown(event: PointerEvent) {
      if (!filtersRef.current?.contains(event.target as Node)) {
        setOpenFilter(null);
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") setOpenFilter(null);
    }

    document.addEventListener("pointerdown", handlePointerDown);
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, []);

  const filteredTickets = useMemo(() => {
    const normalizedSearch = search.trim().toLocaleLowerCase("ru");

    return ticketItems.filter((ticket) => {
      const matchesSearch =
        !normalizedSearch || ticket.track.toLocaleLowerCase("ru").includes(normalizedSearch);
      const matchesPriority =
        selectedFilters.priority.length === 0 ||
        selectedFilters.priority.includes(ticket.priority);
      const matchesStatus =
        selectedFilters.status.length === 0 ||
        selectedFilters.status.includes(ticket.status);
      const matchesCategory =
        selectedFilters.category.length === 0 ||
        selectedFilters.category.includes(ticket.category);
      const matchesApplicant =
        selectedFilters.applicant.length === 0 ||
        selectedFilters.applicant.includes(ticket.applicant);

      return (
        matchesSearch &&
        matchesPriority &&
        matchesStatus &&
        matchesCategory &&
        matchesApplicant
      );
    });
  }, [search, selectedFilters, ticketItems]);

  const selectedChips = (Object.keys(selectedFilters) as FilterKey[]).flatMap((filter) =>
    selectedFilters[filter].map((value) => ({ filter, value })),
  );

  const hasFilters = search.length > 0 || selectedChips.length > 0;

  function toggleFilter(filter: FilterKey, value: string) {
    setSelectedFilters((current) => ({
      ...current,
      [filter]: current[filter].includes(value)
        ? current[filter].filter((item) => item !== value)
        : [...current[filter], value],
    }));
  }

  function resetFilters() {
    setSearch("");
    setSelectedFilters({ priority: [], status: [], category: [], applicant: [] });
    setOpenFilter(null);
  }

  function toggleTrack(track: string) {
    setSelectedTracks((current) =>
      current.includes(track)
        ? current.filter((item) => item !== track)
        : [...current, track],
    );
  }

  function updateTicket(
    track: string,
    field: "status" | "category",
    value: string,
  ) {
    setTicketItems((current) =>
      current.map((ticket) =>
        ticket.track === track ? { ...ticket, [field]: value } : ticket,
      ),
    );
  }

  return (
      <div className="min-w-0 flex-1 px-5 pt-6 pb-0 sm:px-[30px]">
        <div className="mx-auto w-full max-w-[1180px]">
          <label className="relative block">
            <svg viewBox="0 0 24 24" aria-hidden="true" className="pointer-events-none absolute top-1/2 left-3 z-10 h-[21px] w-[21px] -translate-y-1/2 text-[#4562f0]" fill="none">
              <circle cx="10.5" cy="10.5" r="6.5" stroke="currentColor" strokeWidth="2" />
              <path d="m16 16 5 5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
            </svg>
            <Input
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              aria-label="Поиск обращения по трек-номеру"
              className="!h-[42px] !w-full !rounded-[12px] !border-[#4562f0] !bg-white !pr-3 !pl-[42px] !text-left !text-[16px] placeholder:!text-[#9196a7] focus:!bg-white focus:!text-[#000828] focus:placeholder:!text-[#9196a7]"
              placeholder="Поиск обращения по треку"
            />
          </label>

          <header className="mt-[29px]">
            <h1 className="text-[24px] leading-[1.2] font-extrabold tracking-[-0.02em] text-[#4562f0] sm:text-[26px]">
              Очередь новых обращений
            </h1>
            <p className="mt-[7px] text-[16px] leading-5 text-[#151515]">Сначала самые ранние</p>
          </header>

          <section aria-label="Фильтры обращений" className="mt-[20px]">
            <div ref={filtersRef} className="flex flex-wrap items-center gap-[6px]">
              {(Object.keys(filterGroups) as FilterKey[]).map((filter) => (
                <FilterDropdown
                  key={filter}
                  filter={filter}
                  selected={selectedFilters[filter]}
                  open={openFilter === filter}
                  onOpen={() => setOpenFilter((current) => (current === filter ? null : filter))}
                  onToggle={(value) => toggleFilter(filter, value)}
                />
              ))}

              <Button
                text="Сбросить фильтры"
                variant="secondary"
                size="small"
                onClick={resetFilters}
                disabled={!hasFilters}
                className="ml-auto h-[33px] rounded-[7px] px-4 text-[12px] font-normal max-[699px]:ml-0"
              />
            </div>

            <div className="mt-4 flex min-h-[30px] flex-wrap gap-2" aria-live="polite">
              {selectedChips.map(({ filter, value }) => {
                const option = filterGroups[filter].find((item) => item.value === value);
                const chipColor = option?.color
                  ? priorityColors[option.color]
                  : "border-[#dfe6ff] bg-[#dfe6ff] text-[#4562f0]";

                return (
                  <button
                    key={`${filter}-${value}`}
                    type="button"
                    onClick={() => toggleFilter(filter, value)}
                    aria-label={`Удалить фильтр ${getOptionLabel(filter, value)}`}
                    className={`inline-flex h-[29px] cursor-pointer items-center gap-2 rounded-full border px-3 text-[11px] transition-[filter] hover:brightness-[0.97] ${chipColor}`}
                  >
                    {getOptionLabel(filter, value)}
                    <span aria-hidden="true" className="text-base leading-none">×</span>
                  </button>
                );
              })}
            </div>
          </section>

          <section aria-label="Список обращений" className="mt-[18px] overflow-x-auto rounded-[13px] border border-[#4562f0] bg-white">
            <table className="w-full min-w-[940px] table-fixed border-collapse">
              <thead className="bg-[#dfe6ff] text-[#4562f0]">
                <tr className="h-[53px]">
                  <th scope="col" className="w-[23%] border-r border-[#4562f0] px-4 text-center text-[16px] font-normal">Трек-номер</th>
                  <th scope="col" className="w-[26%] border-r border-[#4562f0] px-4 text-center text-[16px] font-normal">Категория</th>
                  <th scope="col" className="w-[19%] border-r border-[#4562f0] px-4 text-center text-[16px] font-normal">Статус</th>
                  <th scope="col" className="w-[17%] border-r border-[#4562f0] px-4 text-center text-[16px] font-normal">Тип заявителя</th>
                  <th scope="col" className="w-[15%] px-4 text-center text-[16px] font-normal">Ожидает</th>
                </tr>
              </thead>
              <tbody>
                {filteredTickets.map((ticket) => (
                  <tr key={ticket.track} className={`h-[59px] border-t border-[#4562f0] ${rowColors[ticket.priority]}`}>
                    <td className="border-r border-[#4562f0] px-4">
                      <div className="flex items-center gap-3 text-[13px] text-[#000828]">
                        <Checkbox
                          checked={selectedTracks.includes(ticket.track)}
                          onChange={() => toggleTrack(ticket.track)}
                          label={<span className="sr-only">Выбрать обращение {ticket.track}</span>}
                          className="shrink-0 [&>span:first-of-type]:h-[17px] [&>span:first-of-type]:w-[17px] [&>span:first-of-type]:rounded-[4px] [&>span:first-of-type]:border-[#4562f0] [&>span:last-child]:sr-only"
                        />
                        <Link
                          href={`/operator/queueNew.tsx/${encodeURIComponent(ticket.track)}`}
                          className="rounded-sm font-medium underline-offset-4 hover:text-[#4562f0] hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]"
                        >
                          {ticket.track}
                        </Link>
                      </div>
                    </td>
                    <td className="border-r border-[#4562f0] px-2 text-center">
                      <select
                        value={ticket.category}
                        onChange={(event) =>
                          updateTicket(ticket.track, "category", event.target.value)
                        }
                        aria-label={`Изменить категорию обращения ${ticket.track}`}
                        className="h-[33px] w-full cursor-pointer truncate rounded-[8px] border border-[#000828] bg-white/70 px-2 text-[11px] text-[#000828] outline-none transition-colors hover:border-[#4562f0] focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/20"
                      >
                        {categoryOptions.map((option) => (
                          <option key={option.value} value={option.value}>
                            {option.label}
                          </option>
                        ))}
                      </select>
                    </td>
                    <td className="border-r border-[#4562f0] px-2 text-center">
                      <select
                        value={ticket.status}
                        onChange={(event) =>
                          updateTicket(ticket.track, "status", event.target.value)
                        }
                        aria-label={`Изменить статус обращения ${ticket.track}`}
                        className="h-[33px] w-full cursor-pointer truncate rounded-[8px] border border-[#000828] bg-white/70 px-2 text-[11px] text-[#000828] outline-none transition-colors hover:border-[#4562f0] focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/20"
                      >
                        {statusOptions.map((option) => (
                          <option key={option.value} value={option.value}>
                            {option.label}
                          </option>
                        ))}
                      </select>
                    </td>
                    <td className="border-r border-[#4562f0] px-3 text-center text-[13px] text-[#000828]">
                      {getOptionLabel("applicant", ticket.applicant)}
                    </td>
                    <td className="px-3 text-center text-[13px] text-[#000828]">{ticket.waiting}</td>
                  </tr>
                ))}

                {filteredTickets.length === 0 && (
                  <tr className="h-[105px] border-t border-[#4562f0]">
                    <td colSpan={5} className="px-6 text-center text-[14px] text-[#646d86]">
                      Обращения по выбранным фильтрам не найдены
                    </td>
                  </tr>
                )}

                <tr aria-hidden="true" className="h-[327px] border-t border-[#4562f0]">
                  <td className="border-r border-[#4562f0]" />
                  <td className="border-r border-[#4562f0]" />
                  <td className="border-r border-[#4562f0]" />
                  <td className="border-r border-[#4562f0]" />
                  <td />
                </tr>
              </tbody>
            </table>
          </section>
        </div>
      </div>
  );
}
