import { appealRoles } from "@/features/appeal/roles";
import RoleCard from "@/features/appeal/role-card";
import Button from "@/shared/ui/button";
import PageHeading from "@/shared/ui/page-heading";
import styles from "./role-selection.module.css";

export default function RoleSelection({ selectedRole }: { selectedRole: (typeof appealRoles)[number]["id"] }) {
  return (
    <section className={styles.panel} aria-labelledby="role-heading">
      <PageHeading id="role-heading" title="Кто вы?" />
      <form action="/appeal/appeal2" method="get" className={styles.form}>
        <fieldset className={styles.roles} aria-labelledby="role-heading">
          <legend className="sr-only">Выбери свою роль</legend>
          {appealRoles.map(role => <RoleCard key={role.id} role={role} selected={role.id === selectedRole} />)}
        </fieldset>
        <Button text="Продолжить" type="submit" fill className={styles.continue} />
      </form>
    </section>
  );
}
