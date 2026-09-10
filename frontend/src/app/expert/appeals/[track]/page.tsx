import ExpertAppealDetail from "@/features/expert/expert-appeal-detail";

export default async function ExpertAppealPage({
  params,
}: {
  params: Promise<{ track: string }>;
}) {
  const { track } = await params;
  return <ExpertAppealDetail trackID={decodeURIComponent(track)} />;
}
