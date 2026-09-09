import { getAppealRole } from "@/features/appeal/roles";
import SiteHeader from "@/widgets/site-header/site-header";
import ClarifyingQuestions from "./clarifying-questions";

export default async function AppealDetails({ searchParams }: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);

  return (
    <main className="flex min-h-dvh flex-col bg-[#f7f9fe] text-[#11131a]">
      <SiteHeader />

      <ClarifyingQuestions role={role.id} />
    </main>
  );
}
