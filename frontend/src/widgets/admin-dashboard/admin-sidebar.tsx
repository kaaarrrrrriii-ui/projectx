"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { StaffIdentity, StaffLogout } from "@/features/staff/staff-session";

const primaryNavigation = [
  { label: "Главная", href: "/admin" },
  { label: "Очередь новых", href: "/admin/queue" },
  { label: "Сотрудники", href: "/admin/employees" },
  { label: "Категории", href: "/admin/categories" },
] as const;

const secondaryNavigation = [
  { label: "Аналитика", href: "/admin/analytics" },
] as const;

function NavigationGroup({
  items,
  pathname,
}: {
  items: ReadonlyArray<{ label: string; href: string }>;
  pathname: string;
}) {
  return (
    <ul className="flex flex-col gap-1 max-[799px]:contents">
      {items.map((item) => {
        const active = item.href === "/admin"
          ? pathname === item.href
          : pathname.startsWith(item.href);

        return (
          <li key={item.href} className="max-[799px]:shrink-0">
            <Link
              href={item.href}
              aria-current={active ? "page" : undefined}
              className={`flex min-h-10 items-center rounded-xl px-3 text-[13px] font-medium transition-colors ${
                active
                  ? "bg-[#e7ecff] text-[#4562f0]"
                  : "text-[#4f5873] hover:bg-[#f1f3fa] hover:text-[#000828]"
              }`}
            >
              <span className="whitespace-nowrap">{item.label}</span>
            </Link>
          </li>
        );
      })}
    </ul>
  );
}

export default function AdminSidebar() {
  const pathname = usePathname();

  return (
    <aside className="flex w-[224px] shrink-0 flex-col border-r border-[#a9b3ff] bg-[#f7f9fe] px-4 py-5 min-[800px]:sticky min-[800px]:top-[100px] min-[800px]:h-[calc(100dvh-100px)] min-[800px]:self-start max-[799px]:w-full max-[799px]:border-r-0 max-[799px]:border-b max-[799px]:px-4 max-[799px]:py-3">
      <StaffIdentity role="admin" label="Администратор" fallback="Администратор" />

      <nav
        aria-label="Навигация администратора"
        className="flex flex-col gap-6 max-[799px]:flex-row max-[799px]:gap-1 max-[799px]:overflow-x-auto max-[799px]:pb-1"
      >
        <NavigationGroup items={primaryNavigation} pathname={pathname} />
        <NavigationGroup items={secondaryNavigation} pathname={pathname} />
      </nav>

      <StaffLogout className="max-[799px]:hidden" />
    </aside>
  );
}
