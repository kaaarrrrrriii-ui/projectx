import { redirect } from "next/navigation";

export default async function AdminAssignmentPage({ params }: { params: Promise<{ track: string }> }) {
  const { track } = await params;
  redirect(`/admin/queue/${encodeURIComponent(decodeURIComponent(track))}`);
}
