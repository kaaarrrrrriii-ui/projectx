import { notFound } from "next/navigation";
import ExpertAppealDetail from "@/features/expert/expert-appeal-detail";
import { getExpertAppeal } from "@/features/expert/expert-data";

export default async function ExpertAppealPage({
  params,
  searchParams,
}: {
  params: Promise<{ track: string }>;
  searchParams: Promise<{ expert?: string | string[] }>;
}) {
  const { track } = await params;
  const appeal = getExpertAppeal(decodeURIComponent(track));
  if (!appeal) notFound();
  const expert = (await searchParams).expert;
  return <ExpertAppealDetail appeal={appeal} initialExecutor={typeof expert === "string" ? expert : ""} />;
}
