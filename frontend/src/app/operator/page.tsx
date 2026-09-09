import SiteHeader from "@/widgets/site-header/site-header";
import OperatorDashboard from "@/widgets/operator-dashboard/operator-dashboard";
import OperatorSidebar from "@/widgets/operator-dashboard/operator-sidebar";

export default function OperatorPage() {
  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <OperatorSidebar />
        <OperatorDashboard />
      </div>
    </main>
  );
}
