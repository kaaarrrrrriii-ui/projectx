import { getAppealRole } from "@/features/appeal/roles";
import RoleSelection from "@/widgets/role-selection/role-selection";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function Appeal({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const selectedRole = getAppealRole((await searchParams).role);

  return (
    <main className="app-page-background flex min-h-dvh flex-col text-[clamp(14px,1.4vw,18px)]">
      <SiteHeader />
      <RoleSelection selectedRole={selectedRole.id} />
    </main>
  );
}
