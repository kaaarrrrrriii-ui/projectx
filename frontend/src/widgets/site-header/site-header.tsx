import Image from "next/image";
import styles from "./site-header.module.css";

export default function SiteHeader() {
  return (
    <header className={styles.header}>
      <div className={styles.content}>
        <Image src="/images/logo.svg" alt="Отклик — ты не один, мы рядом"
          width={180} height={53} priority className={styles.logo} />
        <p>Безопасный способ рассказать о том, что тебя беспокоит, и получить помощь от профессионалов.</p>
      </div>
    </header>
  );
}
