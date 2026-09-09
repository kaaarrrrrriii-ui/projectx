import SiteHeader from "@/widgets/site-header/site-header";
import SubmissionSuccess from "./submission-success";

export default function SubmissionSuccessPage() {
  return (
    <main className="app-page-background flex min-h-dvh min-w-[320px] flex-col overflow-hidden text-[var(--color-text)]">
      <SiteHeader />

      <SubmissionSuccess trackNumber="НАШК-УАЫВ-АВАМ-ВАФВ" />
    </main>
  );
}
