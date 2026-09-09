import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import SiteHeader from "@/widgets/site-header/site-header";
import SubmissionSuccess from "./submission-success";

export default async function SubmissionSuccessPage({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);
  const formal = isFormalAppealRole(role.id);

  return (
    <main className="flex min-h-dvh min-w-[320px] flex-col overflow-hidden bg-[var(--color-background)] text-[var(--color-text)]">
      <SiteHeader formal={formal} />

      <SubmissionSuccess
        trackNumber="НАШК-УАЫВ-АВАМ-ВАФВ"
        formal={formal}
      />
    </main>
  );
}
