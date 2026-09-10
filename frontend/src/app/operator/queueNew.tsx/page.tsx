import SiteHeader from "@/widgets/site-header/site-header";
import OperatorSidebar from "@/widgets/operator-dashboard/operator-sidebar";
import QueueNew from "./queue-new";

export default function QueueNewPage() {
  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader compact />
      <div className="flex flex-1 items-stretch  max-[799px]:flex-col">
        <OperatorSidebar active="queue" />
        <QueueNew />
      </div>
    </main>
  );
}
