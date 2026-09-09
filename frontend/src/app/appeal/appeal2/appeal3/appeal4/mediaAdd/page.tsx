import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import SiteHeader from "@/widgets/site-header/site-header";
import MediaUpload from "./media-upload";

export default async function MediaAddPage({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[]; topic?: string | string[] }>;
}) {
  const params = await searchParams;
  const role = getAppealRole(params.role);
  const topic = typeof params.topic === "string" ? params.topic : "";

  return (
    <main className="app-page-background flex min-h-dvh min-w-[320px] flex-col text-[var(--color-text)]">
      <SiteHeader />

      <MediaUpload role={role.id} topic={topic} />
    </main>
  );
}
