import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import SubmissionSuccess from "@/features/appeal/submission-success";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function SubmissionSuccessPage({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);
  const formal = isFormalAppealRole(role.id);

  return (
    <main className="app-page-background flex min-h-dvh min-w-[320px] flex-col overflow-hidden text-[var(--color-text)]">
      <SiteHeader formal={formal} />

      <SubmissionSuccess
        trackNumber="НАШК-УАЫВ-АВАМ-ВАФВ"
        formal={formal}
      />
    </main>
  );
}
