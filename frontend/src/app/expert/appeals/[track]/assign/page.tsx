import { notFound } from "next/navigation";
import SpecialistAssignment from "@/features/operator/specialist-assignment";
import { getExpertAppeal } from "@/features/expert/expert-data";

export default async function ExpertAssignmentPage({ params }: { params: Promise<{ track: string }> }) {
  const { track } = await params;
  const appeal = getExpertAppeal(decodeURIComponent(track));
  if (!appeal) notFound();

  return <SpecialistAssignment ticket={appeal} returnBasePath="/expert/appeals" />;
}
