import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import SiteHeader from "@/widgets/site-header/site-header";
import MediaUpload from "./media-upload";

export default async function MediaAddPage({
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
