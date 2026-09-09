import EmployeeRoleSelection from "@/features/staff/employee-role-selection";
import SiteHeader from "@/widgets/site-header/site-header";

export default function StaffPage() {
  return (
    <main className="app-page-background min-h-dvh">
      <SiteHeader />
      <section className="mx-auto flex w-full max-w-[1120px] flex-col items-center px-5 py-10 md:px-10 md:py-16">
        <div className="mb-9 text-center">
          <p className="mb-2 text-sm font-medium text-[#646d86]">Вход для сотрудников</p>
          <h1 className="text-[clamp(26px,4vw,40px)] font-bold text-[#4562f0]">Выберите роль</h1>
          <p className="mt-3 text-sm text-[#646d86]">Вы попадёте в рабочее пространство своей роли</p>
        </div>
        <EmployeeRoleSelection />
      </section>
    </main>
  );
}
