import Image from "next/image";
import Icon from "@/shared/ui/icon";
import type { appealRoles } from "./roles";
import styles from "./role-card.module.css";

export default function RoleCard({ role, selected }: {
  role: (typeof appealRoles)[number]; selected: boolean;
}) {
  return (
    <label className={styles.option}>
      <input className={styles.radio} type="radio" name="role" value={role.id}
        defaultChecked={selected} aria-labelledby={`${role.id}-title`}
        aria-describedby={`${role.id}-description`} required />
      <span className={styles.card}>
        <span className={styles.portrait}>
          <Image src={role.image} alt="" fill unoptimized loading="eager" className={styles.image} />
        </span>
        <span className={styles.copy}>
          <span id={`${role.id}-title`} className={styles.title}>{role.title}</span>
          <span id={`${role.id}-description`} className={styles.description}>{role.description}</span>
        </span>
        <Icon name="chevron-right" className={styles.chevron} width={20} height={20} />
      </span>
    </label>
  );
}
