import { getOperatorTicket } from "@/features/operator/tickets";
import TicketDetail from "@/features/operator/ticket-detail";
import OperatorSidebar from "@/widgets/operator-dashboard/operator-sidebar";
import SiteHeader from "@/widgets/site-header/site-header";
import { notFound } from "next/navigation";

export default async function OperatorTicketPage({ params, searchParams }: { params: Promise<{ track: string }>; searchParams: Promise<{ expert?: string | string[] }> }) {
  const { track } = await params;
  let decodedTrack = track;
  try {
    decodedTrack = decodeURIComponent(track);
  } catch {
    // Next.js usually provides a decoded segment; keep the original if it is already decoded.
  }
  const ticket = getOperatorTicket(decodedTrack);
  if (!ticket) notFound();
  const expertParam = (await searchParams).expert;
  const initialExpert = typeof expertParam === "string" ? expertParam : "";

  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader compact />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <OperatorSidebar active="queue" />
        <TicketDetail ticket={ticket} initialExpert={initialExpert} />
      </div>
    </main>
  );
}
