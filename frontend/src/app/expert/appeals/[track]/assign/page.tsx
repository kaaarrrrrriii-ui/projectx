import { redirect } from "next/navigation";

export default async function ExpertAssignmentPage({ params }: { params: Promise<{ track: string }> }) {
  const { track } = await params;
  redirect(`/expert/appeals/${encodeURIComponent(decodeURIComponent(track))}`);
}
