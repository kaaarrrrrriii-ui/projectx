import type { ReactNode } from "react";
import styles from "./info-card.module.css";

export default function InfoCard({ title, description, icon }: {
  title: string;
  description: string;
  icon: ReactNode;
}) {
  return (
    <article className={styles.card}>
      <div className={styles.icon} aria-hidden="true">{icon}</div>
      <div className={styles.copy}>
        <h2>{title}</h2>
        <p>{description}</p>
      </div>
    </article>
  );
}
