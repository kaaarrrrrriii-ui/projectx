import TicketDetail from "@/features/operator/ticket-detail";
import { getOperatorTicket } from "@/features/operator/tickets";
import AdminSidebar from "@/widgets/admin-dashboard/admin-sidebar";
import SiteHeader from "@/widgets/site-header/site-header";
import { notFound } from "next/navigation";

export default async function AdminTicketPage({
  params,
  searchParams,
}: {
  params: Promise<{ track: string }>;
  searchParams: Promise<{ expert?: string | string[] }>;
}) {
  const { track } = await params;
  let decodedTrack = track;
  try {
    decodedTrack = decodeURIComponent(track);
  } catch {
    // The segment may already be decoded by Next.js.
  }

  const ticket = getOperatorTicket(decodedTrack);
  if (!ticket) notFound();

  const expertParam = (await searchParams).expert;
  const initialExpert = typeof expertParam === "string" ? expertParam : "";

  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <AdminSidebar active="queue" />
        <TicketDetail ticket={ticket} initialExpert={initialExpert} basePath="/admin/queue" />
      </div>
    </main>
  );
}
