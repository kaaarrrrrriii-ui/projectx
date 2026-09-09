import AssignedTickets from "@/features/operator/assigned-tickets";
import OperatorSidebar from "@/widgets/operator-dashboard/operator-sidebar";
import SiteHeader from "@/widgets/site-header/site-header";

export default function AssignedTicketsPage() {
  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <OperatorSidebar active="assigned" />
        <AssignedTickets />
      </div>
    </main>
  );
}
