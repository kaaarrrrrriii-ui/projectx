import OperatorTicketLoader from "@/features/operator/operator-ticket-loader";
import OperatorSidebar from "@/widgets/operator-dashboard/operator-sidebar";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function OperatorTicketPage({ params }: { params: Promise<{ track: string }> }) {
  const { track } = await params;
  let decodedTrack = track;
  try {
    decodedTrack = decodeURIComponent(track);
  } catch {
    // Next.js usually provides a decoded segment; keep the original if it is already decoded.
  }
  return (
    <main className="flex min-h-dvh flex-col bg-[var(--color-background)]">
      <SiteHeader sticky />
      <div className="flex flex-1 items-stretch max-[799px]:flex-col">
        <OperatorSidebar active="queue" />
        <OperatorTicketLoader track={decodedTrack} />
      </div>
    </main>
  );
}
