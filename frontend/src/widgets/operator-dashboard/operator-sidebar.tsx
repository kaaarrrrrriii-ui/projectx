import Link from "next/link";
import Button from "@/shared/ui/button";

const navigation = [
  { label: "Главная", icon: "home", active: true },
  { label: "Очередь новых", icon: "inbox" },
  { label: "Распределённые", icon: "users" },
  { label: "Возвраты", icon: "return" },
] as const;

const secondaryNavigation = [
  { label: "Аналитика", icon: "chart" },
  { label: "Выгрузки", icon: "download" },
] as const;

function NavigationIcon({ name }: { name: string }) {
  const common = {
    className: "h-[18px] w-[18px] shrink-0",
    viewBox: "0 0 24 24",
    fill: "none",
    "aria-hidden": true,
  } as const;

  switch (name) {
    case "home":
      return <svg {...common}><path d="m3.5 11 8.5-7 8.5 7v8.5a1 1 0 0 1-1 1h-15a1 1 0 0 1-1-1V11Z" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" /><path d="M9 20v-6h6v6" stroke="currentColor" strokeWidth="1.8" /></svg>;
    case "inbox":
      return <svg {...common}><path d="M4 5.5h16v14H4v-14Z" stroke="currentColor" strokeWidth="1.8" /><path d="M4 14h4l2 2h4l2-2h4" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" /></svg>;
    case "users":
      return <svg {...common}><circle cx="9" cy="8" r="3" stroke="currentColor" strokeWidth="1.8" /><circle cx="17" cy="9" r="2" stroke="currentColor" strokeWidth="1.8" /><path d="M3.5 19c.5-4 2.4-6 5.5-6s5 2 5.5 6M15 14c3.2-.7 5 .9 5.5 4" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" /></svg>;
    case "return":
      return <svg {...common}><path d="m8 8-4 4 4 4" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" /><path d="M4 12h10.5a5.5 5.5 0 1 1 0 11" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" /></svg>;
    case "chart":
      return <svg {...common}><path d="M5 20V10m7 10V4m7 16v-7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg>;
    default:
      return <svg {...common}><path d="M12 3v12m0 0 4-4m-4 4-4-4" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" /><path d="M5 19v2h14v-2" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" /></svg>;
  }
}

function NavigationGroup({ items }: { items: ReadonlyArray<{ label: string; icon: string; active?: boolean }> }) {
  return (
    <ul className="flex flex-col gap-1 max-[799px]:contents">
      {items.map((item) => (
        <li key={item.label} className="max-[799px]:shrink-0">
          <Link
            href="#"
            aria-current={item.active ? "page" : undefined}
            className={`flex min-h-10 items-center gap-3 rounded-xl px-3 text-[13px] font-medium transition-colors ${item.active ? "bg-[#e7ecff] text-[#4562f0]" : "text-[#4f5873] hover:bg-[#f1f3fa] hover:text-[#000828]"}`}
          >
            <NavigationIcon name={item.icon} />
            <span className="whitespace-nowrap">{item.label}</span>
          </Link>
        </li>
      ))}
    </ul>
  );
}

export default function OperatorSidebar() {
  return (
    <aside className="flex w-[224px] shrink-0 flex-col border-r border-[#a9b3ff] bg-[#f7f9fe] px-4 py-5 max-[799px]:w-full max-[799px]:border-r-0 max-[799px]:border-b max-[799px]:px-4 max-[799px]:py-3">
      <div className="mb-7 px-2 max-[799px]:mb-3">
        <p className="font-semibold text-[#4562f0]">Олег Зетник</p>
        <p className="mt-0.5 text-xs text-[#646d86]">Оператор</p>
      </div>

      <nav aria-label="Навигация оператора" className="flex flex-col gap-6 max-[799px]:flex-row max-[799px]:gap-1 max-[799px]:overflow-x-auto max-[799px]:pb-1">
        <NavigationGroup items={navigation} />
        <NavigationGroup items={secondaryNavigation} />
      </nav>

      <Button
        text="Выйти"
        link="/"
        variant="primary"
        size="default"
        className="mt-auto w-full max-[799px]:mt-3 max-[799px]:w-fit"
      />
    </aside>
  );
}
