import Image from "next/image";
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
          <Image src={role.image} alt="" fill sizes="(max-width: 699px) 80px, (max-width: 1100px) 116px, 160px" className={styles.image} />
        </span>
        <span className={styles.copy}>
          <span id={`${role.id}-title`} className={styles.title}>{role.title}</span>
          <span id={`${role.id}-description`} className={styles.description}>{role.description}</span>
        </span>
      </span>
    </label>
  );
}
