import type { ReactNode } from "react";

export default function Surface({
  children,
  className = "",
  as: Tag = "div",
  labelledBy,
}: {
  children: ReactNode;
  className?: string;
  as?: "div" | "section";
  labelledBy?: string;
}) {
  return (
    <Tag
      aria-labelledby={labelledBy}
      className={`
        min-w-0
        border-[length:var(--surface-border-width,1px)]
        border-solid
        border-[var(--color-border-panel)]
        rounded-[var(--surface-radius,var(--radius-panel))]
        bg-[var(--color-surface)]
        shadow-[var(--surface-shadow,var(--shadow-panel))]
        ${className}
      `}
    >
      {children}
    </Tag>
  );
}