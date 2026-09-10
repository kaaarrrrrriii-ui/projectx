import SpecialistAssignment from "@/features/operator/specialist-assignment";
import { getOperatorTicket } from "@/features/operator/tickets";
import { notFound } from "next/navigation";

export default async function AdminSpecialistAssignmentPage({
  params,
}: {
  params: Promise<{ track: string }>;
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

  return <SpecialistAssignment ticket={ticket} returnBasePath="/admin/queue" />;
}
