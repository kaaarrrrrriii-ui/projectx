"use client";

import DialogShell from "@/features/chat/dialog-shell";
import Button from "@/shared/ui/button";
import { useMemo, useState, type FormEvent } from "react";

type Specialist = {
  id: string;
  name: string;
  shortName: string;
  area: string;
  load: number;
  rating: number;
};

const specialists: Specialist[] = [
  { id: "baikova", name: "Байкова Елена Сергеевна", shortName: "Байкова Е.С.", area: "Психология", load: 5, rating: 5 },
  { id: "petrova", name: "Петрова Анна Викторовна", shortName: "Петрова А.В.", area: "Психология", load: 3, rating: 4.9 },
  { id: "sokolova", name: "Соколова Мария Андреевна", shortName: "Соколова М.А.", area: "Конфликтология", load: 7, rating: 4.8 },
  { id: "voronov", name: "Воронов Илья Максимович", shortName: "Воронов И.М.", area: "Юридическое сопровождение", load: 4, rating: 4.7 },
  { id: "belova", name: "Белова Ольга Игоревна", shortName: "Белова О.И.", area: "Социальная педагогика", load: 2, rating: 4.9 },
];

function recommendedArea(category: string) {
  if (category.includes("юридическ")) return "Юридическое сопровождение";
  if (category.includes("конфликт")) return "Конфликтология";
  if (category.includes("давление")) return "Социальная педагогика";
  return "Психология";
}

export function SpecialistAssignmentModal({
  category,
  onAssign,
  onClose,
}: {
  category: string;
  onAssign: (specialist: string) => void;
  onClose: () => void;
}) {
  const [selectedAreas, setSelectedAreas] = useState<string[]>([]);
  const areas = Array.from(new Set(specialists.map((item) => item.area)));
  const filteredSpecialists = useMemo(() => specialists.filter((item) =>
    selectedAreas.length === 0 || selectedAreas.includes(item.area),
  ), [selectedAreas]);
  const recommendation = recommendedArea(category);

  function toggleArea(area: string) {
    setSelectedAreas((current) => current.includes(area)
      ? current.filter((item) => item !== area)
      : [...current, area]);
  }

  return (
    <DialogShell labelledBy="assignment-modal-heading" onClose={onClose} showClose className="max-w-[730px] !p-6 sm:!p-[30px]">
      <h2 id="assignment-modal-heading" className="pr-12 text-[25px] leading-tight font-extrabold text-[#4562f0]">Назначение исполнителя</h2>

      <section aria-label="Фильтр специалистов" className="mt-5">
        <div className="flex flex-wrap items-start gap-3">
          <details className="group relative">
            <summary className="flex h-[32px] min-w-[170px] cursor-pointer list-none items-center justify-between gap-3 rounded-[6px] border border-[#000828] bg-white px-4 text-[11px] text-[#000828] hover:border-[#4562f0] hover:text-[#4562f0] [&::-webkit-details-marker]:hidden">
              Выбрать область специалистов
              <svg viewBox="0 0 12 8" aria-hidden="true" className="h-2 w-3 transition-transform group-open:rotate-180" fill="none"><path d="m1 1 5 5 5-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" /></svg>
            </summary>
            <div className="absolute top-10 left-0 z-10 grid min-w-[260px] gap-1 rounded-xl border border-[#8799f8] bg-white p-3 shadow-[0_14px_35px_rgba(0,8,40,0.14)]">
              {areas.map((area) => (
                <label key={area} className="flex cursor-pointer items-center gap-3 rounded-lg px-2 py-2 text-sm hover:bg-[#eef1ff]">
                  <input type="checkbox" checked={selectedAreas.includes(area)} onChange={() => toggleArea(area)} className="h-4 w-4 accent-[#4562f0]" />
                  {area}
                </label>
              ))}
            </div>
          </details>
          <Button text="Сбросить фильтры" variant="secondary" size="small" disabled={selectedAreas.length === 0} onClick={() => setSelectedAreas([])} className="ml-auto h-[32px] rounded-[6px] px-4 text-[11px] font-normal max-[599px]:ml-0" />
        </div>

        <div className="mt-3 flex min-h-7 flex-wrap gap-2">
          {selectedAreas.map((area) => <button key={area} type="button" onClick={() => toggleArea(area)} className="inline-flex h-7 cursor-pointer items-center gap-1 rounded-full bg-[#dfe6ff] px-3 text-[10px] text-[#4562f0]">{area} <span aria-hidden="true">×</span></button>)}
        </div>
        <p className="mt-1 text-xs text-[#9199af]">Рекомендуемая область: {recommendation.toLocaleLowerCase("ru")}</p>
      </section>

      <section aria-labelledby="available-specialists-heading" className="mt-2">
        <h3 id="available-specialists-heading" className="text-base font-medium text-[#151515]">Доступные специалисты</h3>
        <div className="mt-2 overflow-x-auto rounded-[12px] border border-[#4562f0]">
          <table className="w-full min-w-[630px] table-fixed border-collapse">
            <thead className="bg-[#dfe6ff] text-[#4562f0]">
              <tr className="h-[42px]"><th className="w-[27%] border-r border-[#4562f0] px-3 text-center text-sm font-normal">Специалист</th><th className="w-[27%] border-r border-[#4562f0] px-3 text-center text-sm font-normal">Область</th><th className="w-[31%] border-r border-[#4562f0] px-3 text-center text-sm font-normal">Текущая нагрузка</th><th className="w-[15%]"><span className="sr-only">Выбрать</span></th></tr>
            </thead>
            <tbody>
              {filteredSpecialists.map((specialist) => (
                <tr key={specialist.id} className="h-[40px] border-t border-[#4562f0] hover:bg-[#f7f8ff]">
                  <td className="border-r border-[#4562f0] px-3 text-sm">{specialist.shortName} <span className="rounded bg-[#dfe6ff] px-1 py-0.5 text-[10px] text-[#4562f0]">{specialist.rating.toFixed(1)}</span></td>
                  <td className="border-r border-[#4562f0] px-3 text-center text-sm">{specialist.area}</td>
                  <td className="border-r border-[#4562f0] px-3 text-center text-sm">{specialist.load} чел.</td>
                  <td className="px-2 text-center"><button type="button" onClick={() => onAssign(`${specialist.name} — ${specialist.area.toLocaleLowerCase("ru")}`)} className="cursor-pointer rounded-[6px] bg-[#4562f0] px-3 py-1.5 text-[11px] text-white hover:bg-[#4f71fc] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]">Выбрать</button></td>
                </tr>
              ))}
              {filteredSpecialists.length === 0 && <tr className="h-20 border-t border-[#4562f0]"><td colSpan={4} className="text-center text-sm text-[#646d86]">Специалисты не найдены</td></tr>}
            </tbody>
          </table>
        </div>
      </section>
    </DialogShell>
  );
}

export function CloseTicketModal({
  onSubmit,
  onClose,
}: {
  onSubmit: (reason: string) => void;
  onClose: () => void;
}) {
  const [reason, setReason] = useState("");

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const normalizedReason = reason.trim();
    if (normalizedReason) onSubmit(normalizedReason);
  }

  return (
    <DialogShell labelledBy="close-ticket-heading" onClose={onClose} className="max-w-[650px] !p-[18px]">
      <form onSubmit={handleSubmit}>
        <button type="button" onClick={onClose} aria-label="Закрыть окно" className="absolute top-4 right-5 flex h-8 w-8 cursor-pointer items-center justify-center rounded-full text-2xl text-[#e5141b] hover:bg-[#fff1f1] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#e5141b]">×</button>
        <h2 id="close-ticket-heading" className="pr-12 text-[21px] leading-tight font-extrabold text-[#e5141b]">Укажите причину закрытия обращения</h2>
        <label className="mt-4 block">
          <span className="sr-only">Причина закрытия обращения</span>
          <textarea value={reason} onChange={(event) => setReason(event.target.value)} required autoFocus className="block min-h-[148px] w-full resize-y rounded-[12px] border border-[#333] bg-[#fcfdff] p-3 text-sm leading-5 outline-none focus:border-[#e5141b] focus:ring-2 focus:ring-[#e5141b]/15" />
        </label>
        <button type="submit" disabled={!reason.trim()} className="mt-4 h-10 w-full cursor-pointer rounded-[8px] border border-[#e5141b] bg-[#e5141b] text-sm font-medium text-white transition-colors hover:bg-[#c81017] focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-[#e5141b] disabled:cursor-not-allowed disabled:opacity-45">Отправить</button>
      </form>
    </DialogShell>
  );
}
