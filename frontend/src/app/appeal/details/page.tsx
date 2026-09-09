import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import ClarifyingQuestions from "@/features/appeal/clarifying-questions";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function AppealDetails({ searchParams }: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);
  const formal = isFormalAppealRole(role.id);

  return (
    <main className="flex min-h-dvh flex-col bg-[#f7f9fe] text-[#11131a]">
      <SiteHeader formal={formal} />

      <ClarifyingQuestions role={role.id} formal={formal} />
    </main>
  );
}
