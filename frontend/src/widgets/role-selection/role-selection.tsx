import { appealRoles } from "@/features/appeal/roles";
import RoleCard from "@/features/appeal/role-card";
import Button from "@/shared/ui/button";
import PageHeading from "@/shared/ui/page-heading";
import StepNavigation from "@/shared/ui/step-navigation";
import Surface from "@/shared/ui/surface";
import styles from "./role-selection.module.css";

export default function RoleSelection({ selectedRole }: { selectedRole: (typeof appealRoles)[number]["id"] }) {
  return (
    <Surface className={styles.panel}>
      <StepNavigation backHref="/" backLabel="На главную" stepLabel="Первый шаг: выбор роли" />
      <PageHeading id="role-heading" title="Кто ты?" descriptionId="role-description"
        description="Это поможет нам подобрать нужного специалиста и говорить с тобой на одном языке." />
      <form action="/appeal/appeal2" method="get" className={styles.form}>
        <fieldset className={styles.roles} aria-labelledby="role-heading" aria-describedby="role-description">
          <legend className="sr-only">Выбери свою роль</legend>
          {appealRoles.map(role => <RoleCard key={role.id} role={role} selected={role.id === selectedRole} />)}
        </fieldset>
        <Button text="Продолжить" type="submit" fill className={styles.continue} />
      </form>
    </Surface>
  );
}
