"use client";

import Button from "@/shared/ui/button";
import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";

type AssignmentTicket = {
  track: string;
  submittedAt: string;
};

type Specialist = {
  id: string;
  name: string;
  specialization: string;
  load: number;
};

const specialists: Specialist[] = [
  { id: "ba ikova", name: "Байкова Елена Сергеевна", specialization: "Буллинг", load: 5 },
  { id: "petrova", name: "Петрова Анна Викторовна", specialization: "Кибербуллинг", load: 3 },
  { id: "sokolova", name: "Соколова Мария Андреевна", specialization: "Семейные конфликты", load: 7 },
  { id: "voronov", name: "Воронов Илья Максимович", specialization: "Юридическая помощь", load: 4 },
  { id: "belova", name: "Белова Ольга Игоревна", specialization: "Давление и угрозы", load: 2 },
];

function FilterMenu({
  label,
  options,
  selected,
  onToggle,
}: {
  label: string;
  options: string[];
  selected: string[];
  onToggle: (option: string) => void;
}) {
  return (
    <details className="group relative">
      <summary className="flex h-9 min-w-[170px] cursor-pointer list-none items-center justify-between gap-3 rounded-[7px] border border-[#000828] bg-white px-4 text-xs text-[#000828] hover:border-[#4562f0] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] [&::-webkit-details-marker]:hidden">
        {label}
        <svg viewBox="0 0 12 8" aria-hidden="true" className="h-2 w-3 transition-transform group-open:rotate-180" fill="none">
          <path d="m1 1 5 5 5-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </summary>
      <div className="absolute top-11 left-0 z-40 grid max-h-64 min-w-[260px] gap-1.5 overflow-y-auto rounded-xl border border-[#8799f8] bg-white p-3 shadow-[0_14px_35px_rgba(0,8,40,0.14)]">
        {options.map((option) => (
          <label key={option} className="flex cursor-pointer items-center gap-3 rounded-lg px-2 py-2 text-sm hover:bg-[#eef1ff]">
            <input type="checkbox" checked={selected.includes(option)} onChange={() => onToggle(option)} className="h-4 w-4 accent-[#4562f0]" />
            <span>{option}</span>
          </label>
        ))}
      </div>
    </details>
  );
}

export default function SpecialistAssignment({
  ticket,
  returnBasePath = "/operator/queueNew.tsx",
}: {
  ticket: AssignmentTicket;
  returnBasePath?: string;
}) {
  const router = useRouter();
  const [selectedNames, setSelectedNames] = useState<string[]>([]);
  const [selectedSpecializations, setSelectedSpecializations] = useState<string[]>([]);
  const [selectedId, setSelectedId] = useState("");

  const names = specialists.map((item) => item.name);
  const specializations = Array.from(new Set(specialists.map((item) => item.specialization)));
  const selectedSpecialist = specialists.find((item) => item.id === selectedId);

  const filteredSpecialists = useMemo(
    () => specialists.filter((item) =>
      (selectedNames.length === 0 || selectedNames.includes(item.name)) &&
      (selectedSpecializations.length === 0 || selectedSpecializations.includes(item.specialization)),
    ),
    [selectedNames, selectedSpecializations],
  );

  function toggleFilter(value: string, selected: string[], setSelected: (values: string[]) => void) {
    setSelected(selected.includes(value) ? selected.filter((item) => item !== value) : [...selected, value]);
  }

  function resetFilters() {
    setSelectedNames([]);
    setSelectedSpecializations([]);
  }

  function assignSpecialist() {
    if (!selectedSpecialist) return;
    const expert = selectedSpecialist.name + " — " + selectedSpecialist.specialization.toLocaleLowerCase("ru");
    router.push(returnBasePath + "/" + encodeURIComponent(ticket.track) + "?expert=" + encodeURIComponent(expert));
  }

  const selectedChips = [
    ...selectedNames.map((value) => ({ type: "name" as const, value })),
    ...selectedSpecializations.map((value) => ({ type: "specialization" as const, value })),
  ];

  return (
    <section className="min-w-0 flex-1 px-5 py-6 sm:px-[30px]" aria-labelledby="assignment-heading">
      <div className="mx-auto w-full max-w-[1180px]">
        <button type="button" onClick={() => router.push(returnBasePath + "/" + encodeURIComponent(ticket.track))} className="cursor-pointer rounded-sm text-sm text-[#85899b] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[#4562f0]">
          Вернуться назад
        </button>

        <header className="mt-6">
          <h1 id="assignment-heading" className="text-[26px] leading-tight font-extrabold text-[#000828]">
            Назначение исполнителя
          </h1>
          <p className="mt-1 text-base text-[#151515]">Дата: {ticket.submittedAt}</p>
          <p className="mt-1 text-sm text-[#646d86]">Обращение {ticket.track}</p>
        </header>

        <section aria-label="Фильтры специалистов" className="mt-4">
          <div className="flex flex-wrap items-center gap-3">
            <FilterMenu label="Выбрать специалистов" options={names} selected={selectedNames} onToggle={(value) => toggleFilter(value, selectedNames, setSelectedNames)} />
            <FilterMenu label="Выбрать специализацию" options={specializations} selected={selectedSpecializations} onToggle={(value) => toggleFilter(value, selectedSpecializations, setSelectedSpecializations)} />
            <Button text="Сбросить фильтры" variant="secondary" size="small" disabled={selectedChips.length === 0} onClick={resetFilters} className="ml-auto h-9 rounded-[7px] px-4 text-xs font-normal max-[699px]:ml-0" />
          </div>

          <div className="mt-3 flex min-h-8 flex-wrap gap-2" aria-live="polite">
            {selectedChips.map((chip) => (
              <button
                key={chip.type + chip.value}
                type="button"
                onClick={() => chip.type === "name"
                  ? toggleFilter(chip.value, selectedNames, setSelectedNames)
                  : toggleFilter(chip.value, selectedSpecializations, setSelectedSpecializations)}
                className="inline-flex h-7 cursor-pointer items-center gap-1.5 rounded-full bg-[#dfe6ff] px-3 text-[11px] text-[#4562f0] hover:bg-[#d4ddff]"
              >
                {chip.value} <span aria-hidden="true">×</span>
              </button>
            ))}
          </div>
        </section>

        <section aria-labelledby="available-heading" className="mt-2">
          <h2 id="available-heading" className="text-lg font-medium text-[#151515]">Доступные специалисты</h2>
          <div className="mt-2 overflow-x-auto rounded-[13px] border border-[#4562f0] bg-white/80">
            <table className="w-full min-w-[680px] table-fixed border-collapse">
              <thead className="bg-[#dfe6ff] text-[#4562f0]">
                <tr className="h-[52px]">
                  <th className="w-[36%] border-r border-[#4562f0] px-4 text-center font-normal">Специалист</th>
                  <th className="w-[38%] border-r border-[#4562f0] px-4 text-center font-normal">Специализация</th>
                  <th className="w-[26%] px-4 text-center font-normal">Текущая нагрузка</th>
                </tr>
              </thead>
              <tbody>
                {filteredSpecialists.map((specialist) => {
                  const selected = selectedId === specialist.id;
                  return (
                    <tr key={specialist.id} className={"h-[52px] border-t border-[#4562f0] transition-colors " + (selected ? "bg-[#dfe6ff]" : "hover:bg-[#f2f4ff]")}>
                      <td className="border-r border-[#4562f0] px-4">
                        <button type="button" onClick={() => setSelectedId(specialist.id)} className="w-full cursor-pointer rounded-sm text-left text-sm font-medium text-[#000828] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]">
                          <span className="inline-flex items-center gap-2">
                            <span className={"h-3.5 w-3.5 rounded-full border " + (selected ? "border-[#4562f0] bg-[#4562f0] shadow-[inset_0_0_0_3px_white]" : "border-[#9199af]")} />
                            {specialist.name}
                          </span>
                        </button>
                      </td>
                      <td className="border-r border-[#4562f0] px-4 text-center text-sm">{specialist.specialization}</td>
                      <td className="px-4 text-center text-sm">{specialist.load} чел.</td>
                    </tr>
                  );
                })}
                {filteredSpecialists.length === 0 && (
                  <tr className="h-24 border-t border-[#4562f0]">
                    <td colSpan={3} className="text-center text-sm text-[#646d86]">Специалисты не найдены</td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </section>

        <div className="mt-5 flex justify-end">
          <Button text="Назначить выбранного специалиста" variant="primary" size="default" disabled={!selectedSpecialist} onClick={assignSpecialist} />
        </div>
      </div>
    </section>
  );
}
