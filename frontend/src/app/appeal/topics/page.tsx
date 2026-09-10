import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import TopicSelection from "@/features/appeal/topic-selection";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function AppealTopics({ searchParams }: { searchParams: Promise<{ role?: string | string[]; topic?: string | string[] }> }) {
  const params = await searchParams;
  const role = getAppealRole(params.role);
  const formal = isFormalAppealRole(role.id);
  const initialTopic = typeof params.topic === "string" ? params.topic : "";

  return (
    <main className="app-page-background flex min-h-dvh min-w-[320px] flex-col overflow-x-clip">
      <SiteHeader formal={formal} />
      <section aria-labelledby="topics-heading" className="relative isolate flex flex-1 flex-col overflow-hidden px-[4.5%] pt-9 pb-[38px] max-[699px]:px-5 max-[699px]:pt-8 max-[699px]:pb-7 max-[379px]:px-3 max-[379px]:pt-6">
        <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col">
          <TopicSelection roleId={role.id} initialTopic={initialTopic} formal={formal} />
        </div>
      </section>
    </main>
  );
}
