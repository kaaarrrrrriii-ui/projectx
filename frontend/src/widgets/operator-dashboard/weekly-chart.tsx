const week = [
  { day: "Понедельник", shortDay: "Пн", value: 58 },
  { day: "Вторник", shortDay: "Вт", value: 76 },
  { day: "Среда", shortDay: "Ср", value: 48 },
  { day: "Четверг", shortDay: "Чт", value: 88 },
  { day: "Пятница", shortDay: "Пт", value: 68 },
  { day: "Суббота", shortDay: "Сб", value: 38 },
  { day: "Воскресенье", shortDay: "Вс", value: 52 },
];

export default function WeeklyChart() {
  return (
    <section aria-labelledby="weekly-chart-title" className="flex min-h-[181px] flex-col rounded-[13px] border border-[#7990ff] bg-white p-6">
      <h2 id="weekly-chart-title" className="text-base font-medium text-[#000828]">Мои отработки за неделю</h2>

      <div className="mt-5 grid min-h-[112px] flex-1 grid-cols-7 items-end gap-3 border-b border-[#c9d0e5] px-2 max-[599px]:gap-1.5" role="img" aria-label="Гистограмма отработок с понедельника по воскресенье">
        {week.map(({ day, shortDay, value }) => (
          <div key={day} className="flex h-full min-w-0 flex-col items-center justify-end gap-1.5">
            <div
              className="w-full max-w-[46px] rounded-t-[6px] bg-[#4562f0]"
              style={{ height: `${value}%` }}
              title={`${day}: ${value}`}
            />
            <span className="pb-1 text-[11px] leading-4 text-[#646d86]" aria-label={day}>{shortDay}</span>
          </div>
        ))}
      </div>
    </section>
  );
}
