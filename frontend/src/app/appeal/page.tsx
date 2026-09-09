import { getAppealRole } from "@/features/appeal/roles";
import RoleSelection from "@/widgets/role-selection/role-selection";
import styles from "./appeal.module.css";

export default async function Appeal({ searchParams }: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const selectedRole = getAppealRole((await searchParams).role);
  return <main className={styles.page}><RoleSelection selectedRole={selectedRole.id} /></main>;
}
