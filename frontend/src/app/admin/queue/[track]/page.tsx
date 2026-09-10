import TicketDetail from "@/features/operator/ticket-detail";

export default async function AdminTicketPage({
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

  return <TicketDetail trackID={decodedTrack} returnBasePath="/admin/queue" mode="admin" />;
}
