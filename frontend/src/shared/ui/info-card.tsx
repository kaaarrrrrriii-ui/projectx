import type { ReactNode } from "react";

export default function InfoCard({
  title,
  description,
  icon,
  className = "",
}: {
  title: string;
  description: string;
  icon: ReactNode;
  className?: string;
}) {
  return (
    <article
      className={`
        flex min-h-[73px] items-center gap-4
        rounded-[10px]
        border border-[var(--color-primary)]
        bg-[rgb(255_255_255_/_65%)]
        px-4 py-2.5

        min-[1200px]:min-h-24
        min-[1200px]:gap-[22px]
        min-[1200px]:px-[22px]
        min-[1200px]:py-4

        max-[699px]:min-h-24
        max-[699px]:gap-3
        max-[699px]:px-3.5
        max-[699px]:py-4
        ${className}
      `}
    >
      <div
        aria-hidden="true"
        className="
          flex w-[30px] shrink-0
          items-center justify-center
          text-[var(--color-primary)]
        "
      >
        {icon}
      </div>

      <div className="min-w-0">
        <h2
          className="
            text-[clamp(14px,1.5vw,20px)]
            font-medium
            leading-[1.3]
            text-[var(--color-primary)]

            max-[699px]:text-[15px]
          "
        >
          {title}
        </h2>

        <p
          className="
            mt-[5px]
            text-[clamp(11px,1.16vw,15px)]
            leading-[1.3]
            text-[var(--color-text)]

            max-[699px]:text-xs
            max-[699px]:leading-[1.5]
          "
        >
          {description}
        </p>
      </div>
    </article>
  );
}
