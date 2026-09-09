import { getAppealRole } from "@/features/appeal/roles";
import SiteHeader from "@/widgets/site-header/site-header";
import ClarifyingQuestions from "./clarifying-questions";

export default async function AppealDetails({ searchParams }: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);

  return (
    <main className="flex min-h-dvh min-w-[900px] flex-col bg-[#f7f9fe] text-[#11131a]">
      <div className="[&_header]:border-x [&_header]:border-t [&_header]:border-[#9dacf9] [&_header>div]:min-h-[88px] [&_header>div]:max-w-none [&_header>div]:px-[18px] [&_header>div]:py-4 [&_header_img]:w-[160px] [&_header_p]:text-[15px] [&_header_p]:leading-[1.45] [&_header_p]:font-medium">
        <SiteHeader />
      </div>

      <ClarifyingQuestions role={role.id} />
    </main>
  );
}
