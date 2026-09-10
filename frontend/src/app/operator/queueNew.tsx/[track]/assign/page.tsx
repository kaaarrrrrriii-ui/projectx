import { getOperatorTicket } from "@/features/operator/tickets";
import { notFound, redirect } from "next/navigation";

export default async function SpecialistAssignmentPage({ params }: { params: Promise<{ track: string }> }) {
  const { track } = await params;
  let decodedTrack = track;
  try {
    decodedTrack = decodeURIComponent(track);
  } catch {
    // The segment may already be decoded by Next.js.
  }
  const ticket = getOperatorTicket(decodedTrack);
  if (!ticket) notFound();

  redirect(`/operator/queueNew.tsx/${encodeURIComponent(ticket.track)}`);
}
