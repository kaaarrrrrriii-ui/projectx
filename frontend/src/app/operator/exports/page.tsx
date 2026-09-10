import ExportsDashboard from "@/features/operator/exports-dashboard";
import OperatorSidebar from "@/widgets/operator-dashboard/operator-sidebar";
import SiteHeader from "@/widgets/site-header/site-header";

export default function ExportsPage() {
  return (
    <main className="flex min-h-dvh flex-col bg-[var(--color-background)]">
      <SiteHeader compact />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <OperatorSidebar active="exports" />
        <ExportsDashboard />
      </div>
    </main>
  );
}
