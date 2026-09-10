import EmployeesDashboard from "@/features/admin/employees-dashboard";
import AdminSidebar from "@/widgets/admin-dashboard/admin-sidebar";
import SiteHeader from "@/widgets/site-header/site-header";

export default function AdminEmployeesPage() {
  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <AdminSidebar active="employees" />
        <EmployeesDashboard />
      </div>
    </main>
  );
}
