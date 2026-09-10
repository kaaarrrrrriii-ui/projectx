import EmployeeLogin from "@/features/staff/employee-login";

export default function StaffPage() {
  return (
    <main className="app-page-background flex min-h-dvh min-w-[320px] items-center justify-center px-6 py-10 text-[var(--color-text)] max-[599px]:items-end max-[599px]:px-0 max-[599px]:py-0">
      <EmployeeLogin />
    </main>
  );
}
