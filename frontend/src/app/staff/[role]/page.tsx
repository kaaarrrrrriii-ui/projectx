import Link from "next/link";
import SiteHeader from "@/widgets/site-header/site-header";
import { notFound } from "next/navigation";

const roleNames = { admin: "администратора", expert: "эксперта" } as const;

export default async function StaffRolePage({ params }: { params: Promise<{ role: string }> }) {
  const { role } = await params;
  if (!(role in roleNames)) notFound();
  const roleName = roleNames[role as keyof typeof roleNames];

  return (
    <main className="app-page-background min-h-dvh">
      <SiteHeader />
      <section className="mx-auto flex max-w-xl flex-col items-center px-5 py-20 text-center">
        <span className="flex h-16 w-16 items-center justify-center rounded-2xl bg-[#eef1ff] text-3xl text-[#4562f0]">✦</span>
        <h1 className="mt-6 text-3xl font-bold text-[#4562f0]">Кабинет {roleName}</h1>
        <p className="mt-3 leading-6 text-[#646d86]">Раздел готовится к запуску. Сейчас доступен кабинет оператора.</p>
        <div className="mt-8 flex flex-wrap justify-center gap-3">
          <Link href="/staff" className="inline-flex h-12 items-center rounded-xl border border-[#4562f0] px-6 font-medium text-[#4562f0]">Выбрать другую роль</Link>
          <Link href="/operator" className="inline-flex h-12 items-center rounded-xl bg-[#4562f0] px-6 font-medium text-white">Войти как оператор</Link>
        </div>
      </section>
    </main>
  );
}
