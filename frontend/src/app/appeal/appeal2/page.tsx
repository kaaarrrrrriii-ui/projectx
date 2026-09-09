import { getAppealRole } from "@/features/appeal/roles";
import TopicSelection from "@/features/appeal/topic-selection";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function AppealTopics({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[]; topic?: string | string[] }>;
}) {
  const params = await searchParams;
  const role = getAppealRole(params.role);
  const initialTopic = typeof params.topic === "string" ? params.topic : "";

  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader />

      <section
        aria-labelledby="topics-heading"
        className="
          relative isolate flex flex-1 flex-col overflow-hidden
          px-5 pt-9 pb-8
          sm:px-[4.65%]
          sm:pt-[51px]
          sm:pb-[38px]
        "
      >
        <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col">
          <h1
            id="topics-heading"
            className="
              text-[28px]
              font-extrabold
              leading-[1.2]
              tracking-[-0.025em]
              text-[#4562F0]
              sm:text-[34px]
            "
          >
            Выбери тему обращения
          </h1>

          <p
            id="topics-description"
            className="mt-[5px] text-[14px] leading-[22px] text-[#151515]"
          >
            Выбери одну тему, которая точнее всего описывает твою ситуацию.
            Если подходящей темы нет, добавь свою.
          </p>
          <TopicSelection roleId={role.id} initialTopic={initialTopic} />
        </div>
      </section>
    </main>
  );
}
