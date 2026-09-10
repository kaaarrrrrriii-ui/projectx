import AdminSidebar from "@/widgets/admin-dashboard/admin-sidebar";
import OperatorDashboard from "@/widgets/operator-dashboard/operator-dashboard";
import SiteHeader from "@/widgets/site-header/site-header";

export default function AdminPage() {
  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <AdminSidebar />
        <OperatorDashboard />
      </div>
    </main>
  );
}
