import { getHealth } from "@/shared/api/health";

export const dynamic = "force-dynamic";

export default async function Home() {
  const health = await getHealth();

  return (
    <main>
      <section className="card">
        <p className="eyebrow">GERMAN MONOREPO</p>
        <h1>Проект готов к разработке</h1>
        <p className="description">
          Next.js frontend и Go backend запущены как независимые приложения в одном
          репозитории.
        </p>
        <div className="status">
          <span className={health ? "dot dotOnline" : "dot"} aria-hidden="true" />
          Backend: {health ? `${health.status}, uptime ${health.uptime}` : "недоступен"}
        </div>
      </section>
    </main>
  );
}

