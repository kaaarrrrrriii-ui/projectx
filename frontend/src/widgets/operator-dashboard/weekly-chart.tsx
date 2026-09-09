const days = [
  { label: "Пн", value: 8 },
  { label: "Вт", value: 12 },
  { label: "Ср", value: 9 },
  { label: "Чт", value: 15 },
  { label: "Пт", value: 11 },
  { label: "Сб", value: 7 },
  { label: "Вс", value: 5 },
] as const;

export default function WeeklyChart() {
  const maxValue = Math.max(...days.map((day) => day.value));

  return (
    <section aria-labelledby="weekly-chart-title" className="min-h-[276px] rounded-2xl border border-[#7990ff] bg-white p-5 sm:p-6">
      <div className="flex items-start justify-between gap-5">
        <div>
          <h2 id="weekly-chart-title" className="text-base font-medium text-[#000828]">Мои отработки за неделю</h2>
          <p className="mt-1 text-xs text-[#7a8298]">Обработанные обращения</p>
        </div>
        <div className="shrink-0 text-right">
          <strong className="text-2xl font-bold text-[#4562f0]">67</strong>
          <span className="ml-1 text-xs text-[#7a8298]">всего</span>
        </div>
      </div>

      <div className="mt-6 flex h-[150px] items-end gap-[clamp(8px,2vw,22px)] border-b border-[#e5e9f5] px-1 sm:px-3">
        {days.map((day) => (
          <div key={day.label} className="group relative flex h-full min-w-0 flex-1 flex-col items-center justify-end">
            <span className="mb-1 text-[11px] font-semibold text-[#4562f0] opacity-0 transition-opacity group-hover:opacity-100 group-focus-within:opacity-100">{day.value}</span>
            <div
              tabIndex={0}
              aria-label={`${day.label}: ${day.value} обращений`}
              className="w-full max-w-[46px] rounded-t-[7px] bg-[#4562f0] outline-none transition-[filter,transform] hover:brightness-110 focus-visible:ring-3 focus-visible:ring-[#4562f0]/20"
              style={{ height: `${Math.max(18, (day.value / maxValue) * 116)}px` }}
            />
            <span className="absolute translate-y-6 text-xs text-[#646d86]">{day.label}</span>
          </div>
        ))}
      </div>
    </section>
  );
}
