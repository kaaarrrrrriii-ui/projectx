import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import ChatContent from "@/features/chat/chat-content";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function ChatPage({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[]; track?: string | string[] }>;
}) {
  const params = await searchParams;
  const role = getAppealRole(params.role);
  const formal = isFormalAppealRole(role.id);
  const trackNumber = Array.isArray(params.track) ? params.track[0] : params.track ?? "";

  return (
    <main className="app-page-background flex min-h-dvh min-w-[320px] flex-col text-[var(--color-text)]">
      <SiteHeader formal={formal} />

      <section
        aria-labelledby="specialist-answer-heading"
        className="flex flex-1"
      >
        <h1 id="specialist-answer-heading" className="sr-only">
          Ответ специалиста
        </h1>

        <ChatContent formal={formal} trackNumber={trackNumber} />
      </section>
    </main>
  );
}
