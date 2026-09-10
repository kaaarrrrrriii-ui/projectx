import TicketDetail from "@/features/operator/ticket-detail";
import { getOperatorTicket } from "@/features/operator/tickets";
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

  return <TicketDetail ticket={ticket} initialExpert={initialExpert} returnBasePath="/admin/queue" />;
}
