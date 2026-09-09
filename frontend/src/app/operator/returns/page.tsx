import ReturnsTickets from "@/features/operator/returns-tickets";
import OperatorSidebar from "@/widgets/operator-dashboard/operator-sidebar";
import SiteHeader from "@/widgets/site-header/site-header";

export default function ReturnsTicketsPage() {
  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <OperatorSidebar active="returns" />
        <ReturnsTickets />
      </div>
    </main>
  );
}
