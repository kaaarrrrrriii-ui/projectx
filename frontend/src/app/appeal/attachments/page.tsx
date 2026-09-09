import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import MediaUpload from "@/features/appeal/media-upload";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function AppealAttachmentsPage({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);
  const formal = isFormalAppealRole(role.id);

  return (
    <main className="flex min-h-dvh min-w-[320px] flex-col bg-[var(--color-background)] text-[var(--color-text)]">
      <SiteHeader formal={formal} />

      <MediaUpload role={role.id} formal={formal} />
    </main>
  );
}
