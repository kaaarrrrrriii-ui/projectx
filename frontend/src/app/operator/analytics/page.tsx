import AnalyticsDashboard from "@/features/operator/analytics-dashboard";
import OperatorSidebar from "@/widgets/operator-dashboard/operator-sidebar";
import SiteHeader from "@/widgets/site-header/site-header";

export default function AnalyticsPage() {
  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader compact />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <OperatorSidebar active="analytics" />
        <AnalyticsDashboard />
      </div>
    </main>
  );
}
