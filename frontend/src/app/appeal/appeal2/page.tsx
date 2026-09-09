import Button from "@/shared/ui/button";
import Link from "next/link";
import { getAppealRole } from "@/features/appeal/roles";

export default async function AppealTopics({ searchParams }: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);
  return (
    <main className="min-h-screen flex items-center justify-center bg-[var(--color-surface)]">
      <div className="flex max-w-[600px] flex-col items-center justify-center gap-[30px] p-6 text-center">
        <Link href={`/appeal?role=${role.id}`} className="text-[var(--primary-color)] underline underline-offset-4">
          Назад к выбору роли
        </Link>
        <p className="rounded-xl bg-[var(--secondary-color)] px-4 py-2">{role.title}</p>
        <div className="">С чем это связано?</div>
        <div className="">Можно подобрать одну или несколько тем, которые ближе всего к твоей ситуации</div>
        <div className="">
            <div className=""></div>
        </div>
        <Button text="Продолжить" fill link="/" />
      </div>
    </main>
  );
}
