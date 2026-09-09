import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import SiteHeader from "@/widgets/site-header/site-header";
import ClarifyingQuestions from "./clarifying-questions";

export default async function AppealDetails({ searchParams }: {
  searchParams: Promise<{ role?: string | string[]; topic?: string | string[] }>;
}) {
  const params = await searchParams;
  const role = getAppealRole(params.role);
  const topic = typeof params.topic === "string" ? params.topic : "";

  return (
    <main className="app-page-background flex min-h-dvh flex-col text-[#11131a]">
      <SiteHeader />

      <ClarifyingQuestions role={role.id} topic={topic} />
    </main>
  );
}
