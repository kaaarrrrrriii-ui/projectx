import { getAppealRole } from "@/features/appeal/roles";
import SiteHeader from "@/widgets/site-header/site-header";
import MediaUpload from "./media-upload";

export default async function MediaAddPage({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);

  return (
    <main className="flex min-h-dvh min-w-[320px] flex-col bg-[var(--color-background)] text-[var(--color-text)]">
      <div
        className="
          [&_header>div]:min-h-[89px]
          [&_header>div]:max-w-none
          [&_header>div]:px-[18px]
          [&_header>div]:py-4
          [&_header_img]:w-[162px]
          [&_header_p]:text-[15px]
          [&_header_p]:leading-[1.45]
          [&_header_p]:font-medium

          max-[699px]:[&_header>div]:min-h-[72px]
          max-[699px]:[&_header>div]:gap-4
          max-[699px]:[&_header>div]:px-5
          max-[699px]:[&_header>div]:py-3
          max-[699px]:[&_header_img]:w-28
          max-[699px]:[&_header_p]:text-[11px]
        "
      >
        <SiteHeader />
      </div>

      <MediaUpload role={role.id} />
    </main>
  );
}
