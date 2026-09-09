import styles from "./page-heading.module.css";

export default function PageHeading({ id, title, description, descriptionId, className = "" }: {
  id: string;
  title: string;
  description?: string;
  descriptionId?: string;
  className?: string;
}) {
  return (
    <header className={`${styles.heading} ${className}`}>
      <h1 id={id}>{title}</h1>
      {description && <p id={descriptionId}>{description}</p>}
    </header>
  );
}
