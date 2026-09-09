import SiteHeader from "@/widgets/site-header/site-header";
import SubmissionSuccess from "./submission-success";

export default function SubmissionSuccessPage() {
  return (
    <main className="flex min-h-dvh min-w-[320px] flex-col overflow-hidden bg-[var(--color-background)] text-[var(--color-text)]">
      <SiteHeader />

      <SubmissionSuccess trackNumber="НАШК-УАЫВ-АВАМ-ВАФВ" />
    </main>
  );
}
