import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import MediaUpload from "@/features/appeal/media-upload";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function AppealAttachmentsPage({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[]; topic?: string | string[] }>;
}) {
  const params = await searchParams;
  const role = getAppealRole(params.role);
  const formal = isFormalAppealRole(role.id);
  const topic = typeof params.topic === "string" ? params.topic : "";

  return (
    <main className="app-page-background flex min-h-dvh min-w-[320px] flex-col overflow-x-clip text-[var(--color-text)]">
      <SiteHeader formal={formal} />

      <MediaUpload role={role.id} topic={topic} formal={formal} />
    </main>
  );
}
