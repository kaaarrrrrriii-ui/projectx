"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

type NavigationItem = {
  label: string;
  href: string;
  icon: "home" | "queue" | "assigned" | "return" | "requests" | "analytics";
};

const primaryNavigation: NavigationItem[] = [
  { label: "Главная", href: "/expert", icon: "home" },
  { label: "Очередь новых", href: "/expert/queue", icon: "queue" },
  { label: "Распределённые", href: "/expert/assigned", icon: "assigned" },
  { label: "Возвраты", href: "/expert/returns", icon: "return" },
  { label: "Запросы", href: "/expert/requests", icon: "requests" },
];

const secondaryNavigation: NavigationItem[] = [
  { label: "Аналитика", href: "/expert/analytics", icon: "analytics" },
];

function SidebarIcon({ name }: { name: NavigationItem["icon"] }) {
  const common = { viewBox: "0 0 24 24", fill: "none", className: "h-[18px] w-[18px] shrink-0", "aria-hidden": true } as const;

  if (name === "home") return <svg {...common}><path d="m3.5 11 8.5-7 8.5 7v8.5a1 1 0 0 1-1 1h-15a1 1 0 0 1-1-1V11Z" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" /><path d="M9 20v-6h6v6" stroke="currentColor" strokeWidth="1.8" /></svg>;
  if (name === "queue") return <svg {...common}><path d="M4 5.5h16v14H4v-14Z" stroke="currentColor" strokeWidth="1.8" /><path d="M4 14h4l2 2h4l2-2h4" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" /></svg>;
  if (name === "assigned") return <svg {...common}><circle cx="9" cy="8" r="3" stroke="currentColor" strokeWidth="1.8" /><circle cx="17" cy="9" r="2" stroke="currentColor" strokeWidth="1.8" /><path d="M3.5 19c.5-4 2.4-6 5.5-6s5 2 5.5 6M15 14c3.2-.7 5 .9 5.5 4" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" /></svg>;
  if (name === "return") return <svg {...common}><path d="m8 8-4 4 4 4" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" /><path d="M4 12h10.5a5.5 5.5 0 1 1 0 11" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" /></svg>;
  if (name === "requests") return <svg {...common}><path d="M5 4h14v16H5V4Z" stroke="currentColor" strokeWidth="1.8" /><path d="M8 9h8M8 13h8M8 17h5" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" /></svg>;
  return <svg {...common}><path d="M5 20V10m7 10V4m7 16v-7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg>;
}

function NavigationGroup({ items, pathname }: { items: NavigationItem[]; pathname: string }) {
  return (
    <ul className="flex flex-col gap-1 max-[799px]:contents">
      {items.map((item) => {
        const active = item.href === "/expert" ? pathname === item.href : pathname.startsWith(item.href) || (item.href === "/expert/queue" && pathname.startsWith("/expert/appeals/"));
        return (
          <li key={item.href} className="max-[799px]:shrink-0">
            <Link
              href={item.href}
              aria-current={active ? "page" : undefined}
              className={`flex min-h-10 items-center gap-3 rounded-xl px-3 text-[13px] font-medium transition-all duration-200 ${active ? "bg-[#dfe6ff] text-[#4562f0] shadow-[0_5px_14px_rgba(69,98,240,0.12)]" : "text-[#4f5873] hover:translate-x-0.5 hover:bg-white hover:text-[#000828]"}`}
            >
              <SidebarIcon name={item.icon} />
              <span className="whitespace-nowrap">{item.label}</span>
            </Link>
          </li>
        );
      })}
    </ul>
  );
}

export default function ExpertSidebar() {
  const pathname = usePathname();

  return (
    <aside className="flex w-[224px] shrink-0 flex-col border-r border-[#a9b3ff] bg-[#f7f9fe]/95 px-4 py-5 min-[800px]:sticky min-[800px]:top-[100px] min-[800px]:h-[calc(100dvh-100px)] min-[800px]:self-start max-[799px]:w-full max-[799px]:border-r-0 max-[799px]:border-b max-[799px]:px-4 max-[799px]:py-3">
      <div className="mb-7 px-2 max-[799px]:mb-3">
        <p className="font-semibold text-[#4562f0]">Олег Зетник</p>
        <p className="mt-0.5 text-xs text-[#30384f]">Эксперт</p>
        <p className="text-xs text-[#8c93a8]">Буллинг</p>
      </div>

      <nav aria-label="Навигация эксперта" className="flex flex-col gap-6 max-[799px]:flex-row max-[799px]:gap-1 max-[799px]:overflow-x-auto max-[799px]:pb-1">
        <NavigationGroup items={primaryNavigation} pathname={pathname} />
        <NavigationGroup items={secondaryNavigation} pathname={pathname} />
      </nav>

      <Link href="/" className="mt-auto flex min-h-10 items-center justify-center rounded-xl border border-[#4562f0] px-4 text-sm font-medium text-[#4562f0] transition-colors hover:bg-[#4562f0] hover:text-white max-[799px]:hidden">
        Выйти
      </Link>
    </aside>
  );
}
