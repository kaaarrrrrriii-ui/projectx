import ExpertAppealList from "@/features/expert/expert-appeal-list";

export default async function ExpertQueuePage({ searchParams }: { searchParams: Promise<{ search?: string | string[] }> }) {
  const value = (await searchParams).search;
  return <ExpertAppealList mode="queue" initialSearch={typeof value === "string" ? value : ""} />;
}
