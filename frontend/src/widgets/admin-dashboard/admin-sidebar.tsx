import Link from "next/link";

const primaryNavigation = [
  { id: "home", label: "Главная", href: "/admin" },
  { id: "queue", label: "Очередь новых", href: "/admin/queue" },
  { id: "employees", label: "Сотрудники", href: "/admin/employees" },
  { id: "categories", label: "Категории", href: "/admin/categories" },
] as const;

const secondaryNavigation = [
  { id: "analytics", label: "Аналитика", href: "/admin/analytics" },
] as const;

type AdminSection =
  | (typeof primaryNavigation)[number]["id"]
  | (typeof secondaryNavigation)[number]["id"];

function NavigationGroup({
  items,
  active,
}: {
  items: ReadonlyArray<{ id: string; label: string; href: string }>;
  active: AdminSection;
}) {
  return (
    <ul className="flex flex-col gap-1 max-[799px]:contents">
      {items.map((item) => (
        <li key={item.id} className="max-[799px]:shrink-0">
          <Link
            href={item.href}
            aria-current={item.id === active ? "page" : undefined}
            className={`flex min-h-10 items-center rounded-xl px-3 text-[13px] font-medium transition-colors ${
              item.id === active
                ? "bg-[#e7ecff] text-[#4562f0]"
                : "text-[#4f5873] hover:bg-[#f1f3fa] hover:text-[#000828]"
            }`}
          >
            <span className="whitespace-nowrap">{item.label}</span>
          </Link>
        </li>
      ))}
    </ul>
  );
}

export default function AdminSidebar({ active = "home" }: { active?: AdminSection }) {
  return (
    <aside className="flex w-[224px] shrink-0 flex-col border-r border-[#a9b3ff] bg-[#f7f9fe] px-4 py-5 max-[799px]:w-full max-[799px]:border-r-0 max-[799px]:border-b max-[799px]:px-4 max-[799px]:py-3">
      <div className="mb-7 px-2 max-[799px]:mb-3">
        <p className="font-semibold text-[#4562f0]">Олег Зетник</p>
        <p className="mt-0.5 text-xs text-[#646d86]">Администратор</p>
      </div>

      <nav
        aria-label="Навигация администратора"
        className="flex flex-col gap-6 max-[799px]:flex-row max-[799px]:gap-1 max-[799px]:overflow-x-auto max-[799px]:pb-1"
      >
        <NavigationGroup items={primaryNavigation} active={active} />
        <NavigationGroup items={secondaryNavigation} active={active} />
      </nav>
    </aside>
  );
}
