import MetricCard from "./metric-card";
import WeeklyChart from "./weekly-chart";

const metrics = [
  { title: "Новые обращения", value: "10", badge: "+10 обращений", description: "Требуют распределения", icon: "mail", tone: "blue" },
  { title: "Срочные", value: "5", description: "Нуждаются в решении", icon: "urgent", tone: "red" },
  { title: "Просроченные", value: "3", description: "Более 24 часов", icon: "clock", tone: "orange" },
  { title: "Возвраты", value: "3", description: "Повторные рассмотрения", icon: "return", tone: "blue" },
] as const;

export default function OperatorDashboard() {
  return (
    <div className="w-full min-w-0 px-5 py-6 sm:px-7 lg:px-9 lg:py-7">
      <div className="mx-auto w-full max-w-[1180px]">
        <header>
          <h1 className="text-[clamp(22px,2vw,30px)] font-bold leading-tight text-[#4562f0]">Здравствуйте, Олег!</h1>
          <p className="mt-1 text-xs text-[#4f5873]">Спасибо, что помогаете. Ваша работа важна.</p>
        </header>

        <label className="mt-5 flex h-12 items-center gap-3 rounded-xl border border-[#7990ff] bg-white px-4 focus-within:ring-3 focus-within:ring-[#4562f0]/15">
          <svg viewBox="0 0 24 24" aria-hidden="true" className="h-5 w-5 shrink-0 text-[#4562f0]" fill="none"><circle cx="10.5" cy="10.5" r="6.5" stroke="currentColor" strokeWidth="2" /><path d="m16 16 5 5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg>
          <span className="sr-only">Поиск обращения по трек-номеру</span>
          <input className="h-full min-w-0 flex-1 border-0 bg-transparent text-sm text-[#000828] outline-none placeholder:text-[#9199af]" placeholder="Поиск обращения по треку" />
        </label>

        <section aria-label="Сводка обращений" className="mt-5 grid grid-cols-2 gap-4 xl:grid-cols-4">
          {metrics.map((metric) => <MetricCard key={metric.title} {...metric} />)}
        </section>

        <div className="mt-5 grid gap-5 lg:grid-cols-[240px_minmax(0,1fr)]">
          <section aria-labelledby="processed-title" className="flex min-h-[276px] flex-col items-center justify-center rounded-2xl border border-[#7990ff] bg-white p-6 text-center">
            <h2 id="processed-title" className="text-base font-medium text-[#000828]">Отработано</h2>
            <strong className="mt-4 text-5xl font-extrabold text-[#4562f0]">67</strong>
            <p className="mt-4 text-xs text-[#4f5873]">обращений</p>
          </section>
          <WeeklyChart />
        </div>
      </div>
    </div>
  );
}
