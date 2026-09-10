import ExpertTicketLoader from "@/features/expert/expert-ticket-loader";

export default async function ExpertAppealPage({
  params,
}: {
  params: Promise<{ track: string }>;
}) {
  const { track } = await params;
  return <ExpertTicketLoader track={decodeURIComponent(track)} />;
}
