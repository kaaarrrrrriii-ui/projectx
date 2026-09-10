import Link from "next/link";
import Button from "@/shared/ui/button";

const navigation = [
  { id: "home", label: "Главная", href: "/operator" },
  { id: "queue", label: "Очередь новых", href: "/operator/queueNew.tsx" },
  { id: "assigned", label: "Запросы от экспертов", href: "/operator/assigned" },
  { id: "returns", label: "Жалобы", href: "/operator/returns" },
] as const;

const secondaryNavigation = [
  { id: "analytics", label: "Аналитика", href: "/operator/analytics" },
] as const;

type OperatorSection = (typeof navigation)[number]["id"] | (typeof secondaryNavigation)[number]["id"] | "exports";

function NavigationGroup({
  items,
  active,
}: {
  items: ReadonlyArray<{ id: string; label: string; href: string }>;
  active: OperatorSection;
}) {
  return (
    <ul className="flex flex-col gap-1 max-[799px]:contents">
      {items.map((item) => (
        <li key={item.label} className="max-[799px]:shrink-0">
          <Link
            href={item.href}
            aria-current={item.id === active ? "page" : undefined}
            className={`flex min-h-[33px] items-center rounded-xl px-4 text-[13px] font-medium transition-colors ${item.id === active ? "bg-[#dfe6ff] text-[#4562f0]" : "text-[#30384f] hover:bg-[#f1f3fa] hover:text-[#000828]"}`}
          >
            <span className="whitespace-nowrap">{item.label}</span>
          </Link>
        </li>
      ))}
    </ul>
  );
}

export default function OperatorSidebar({ active = "home" }: { active?: OperatorSection }) {
  return (
    <aside className="flex w-[184px] shrink-0 flex-col border-r border-[#7990ff] bg-[#f7f9fe] px-[5px] py-5 max-[799px]:w-full max-[799px]:border-r-0 max-[799px]:border-b max-[799px]:px-4 max-[799px]:py-3">
      <div className="mb-7 px-2 max-[799px]:mb-3">
        <p className="font-semibold text-[#4562f0]">Олег Зетник</p>
        <p className="mt-0.5 text-xs text-[#646d86]">Оператор</p>
      </div>

      <nav aria-label="Навигация оператора" className="flex flex-col gap-7 max-[799px]:flex-row max-[799px]:gap-1 max-[799px]:overflow-x-auto max-[799px]:pb-1">
        <NavigationGroup items={navigation} active={active} />
        <NavigationGroup items={secondaryNavigation} active={active} />
      </nav>

      <Button
        text="Выйти"
        link="/"
        variant="primary"
        size="default"
        className="mx-2 mt-auto w-[calc(100%_-_16px)] max-[799px]:mt-3 max-[799px]:w-fit"
      />
    </aside>
  );
}
