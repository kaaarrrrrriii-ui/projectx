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
    <main className="min-h-dvh bg-[var(--color-background)]">
      <SiteHeader />
      <RoleSelection selectedRole={selectedRole.id} />
    </main>
  );
}