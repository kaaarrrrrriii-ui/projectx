import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import StatusChecker from "@/features/status/status-checker";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function Status({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);
  const formal = isFormalAppealRole(role.id);

  return (
    <main className="app-page-background flex min-h-dvh min-w-[320px] flex-col overflow-x-clip">
      <SiteHeader formal={formal} />
      <StatusChecker roleId={role.id} formal={formal} />
    </main>
  );
}
