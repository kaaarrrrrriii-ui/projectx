import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import TopicSelection from "@/features/appeal/topic-selection";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function AppealTopics({ searchParams }: { searchParams: Promise<{ role?: string | string[]; topic?: string | string[] }> }) {
  const params = await searchParams;
  const role = getAppealRole(params.role);
  const formal = isFormalAppealRole(role.id);
  const initialTopic = typeof params.topic === "string" ? params.topic : "";

  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader formal={formal} />
      <section aria-labelledby="topics-heading" className="relative isolate flex flex-1 flex-col overflow-hidden px-5 pt-9 pb-8 sm:px-[4.65%] sm:pt-9 sm:pb-[38px]">
        <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col">
          <TopicSelection roleId={role.id} initialTopic={initialTopic} formal={formal} />
        </div>
      </section>
    </main>
  );
}
