import { redirect } from "next/navigation";

export default async function OperatorAssignmentPage({ params }: { params: Promise<{ track: string }> }) {
  const { track } = await params;
  redirect(`/operator/queueNew.tsx/${encodeURIComponent(decodeURIComponent(track))}`);
}
