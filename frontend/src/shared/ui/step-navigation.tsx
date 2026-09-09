import Link from "next/link";
import Icon from "./icon";
import styles from "./step-navigation.module.css";

export default function StepNavigation({ backHref, backLabel, stepLabel }: {
  backHref: string; backLabel: string; stepLabel: string;
}) {
  return (
    <nav className={styles.navigation} aria-label="Навигация по обращению">
      <Link href={backHref} className={styles.back} aria-label={backLabel}><Icon name="chevron-left" /></Link>
      <div className={styles.progress}><span className="sr-only">{stepLabel}</span></div>
    </nav>
  );
}
