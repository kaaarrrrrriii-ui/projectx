import type { ReactNode } from "react";
import styles from "./surface.module.css";

export default function Surface({ children, className = "", as: Tag = "div", labelledBy }: {
  children: ReactNode;
  className?: string;
  as?: "div" | "section";
  labelledBy?: string;
}) {
  return <Tag className={`${styles.surface} ${className}`} aria-labelledby={labelledBy}>{children}</Tag>;
}
