import MetricCard from "@/widgets/operator-dashboard/metric-card";
import WeeklyChart from "@/widgets/operator-dashboard/weekly-chart";

const metrics = [
  { title: "Новые обращения", value: "10", badge: "+10 обращений", description: "Требуют рассмотрения", icon: "mail", tone: "blue" },
  { title: "Срочные", value: "5", description: "Нуждаются в решении", icon: "urgent", tone: "red" },
  { title: "Просроченные", value: "3", description: "Более 24 часов", icon: "clock", tone: "orange" },
  { title: "Возвраты", value: "3", description: "Повторные рассмотрения", icon: "return", tone: "blue" },
] as const;

export default function ExpertDashboard() {
  return (
    <section className="min-w-0 flex-1 px-5 py-6 sm:px-7 lg:px-9" aria-labelledby="expert-dashboard-heading">
      <div className="w-full">
        <header>
          <h1 id="expert-dashboard-heading" className="text-[clamp(24px,2.4vw,32px)] font-extrabold leading-tight tracking-[-0.025em] text-[#4562f0]">Здравствуйте, Олег!</h1>
          <p className="mt-1 text-sm text-[#30384f]">Спасибо, что помогаете. Ваша работа важна.</p>
        </header>

        <form action="/expert/queue" className="mt-5">
          <label className="flex h-12 items-center gap-3 rounded-xl border border-[#7990ff] bg-white/90 px-4 shadow-[0_8px_24px_rgba(69,98,240,0.06)] focus-within:ring-3 focus-within:ring-[#4562f0]/15">
            <svg viewBox="0 0 24 24" aria-hidden="true" className="h-5 w-5 shrink-0 text-[#4562f0]" fill="none"><circle cx="10.5" cy="10.5" r="6.5" stroke="currentColor" strokeWidth="2" /><path d="m16 16 5 5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg>
            <span className="sr-only">Поиск обращения по трек-номеру</span>
            <input name="search" className="h-full min-w-0 flex-1 border-0 bg-transparent text-base text-[#000828] outline-none placeholder:text-[#9199af]" placeholder="Поиск обращения по треку" />
          </label>
        </form>

        <section aria-label="Сводка обращений" className="mt-5 grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4 [&>article]:transition-[transform,box-shadow] [&>article]:duration-200 [&>article]:hover:-translate-y-1 [&>article]:hover:shadow-[0_14px_30px_rgba(49,67,145,0.13)]">
          {metrics.map((metric) => <MetricCard key={metric.title} {...metric} />)}
        </section>

        <div className="mt-5 grid gap-5 lg:grid-cols-[240px_minmax(0,1fr)]">
          <section aria-labelledby="processed-title" className="flex min-h-[276px] flex-col items-center justify-center rounded-2xl border border-[#7990ff] bg-white/90 p-6 text-center shadow-[0_10px_26px_rgba(69,98,240,0.06)]">
            <h2 id="processed-title" className="text-base font-medium text-[#000828]">Отработано</h2>
            <strong className="mt-4 text-5xl font-extrabold text-[#4562f0]">67</strong>
            <p className="mt-4 text-xs text-[#4f5873]">обращений</p>
          </section>
          <WeeklyChart />
        </div>
      </div>
    </section>
  );
}
