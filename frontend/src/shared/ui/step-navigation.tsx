import Link from "next/link";
import Icon from "./icon";

export default function StepNavigation({
  backHref,
  backLabel,
  stepLabel,
}: {
  backHref: string;
  backLabel: string;
  stepLabel: string;
}) {
  return (
    <nav
      aria-label="Навигация по обращению"
      className="
        grid grid-cols-[44px_1fr_44px]
        items-center
      "
    >
      <Link
        href={backHref}
        aria-label={backLabel}
        className="
          -ml-3 flex h-11 w-11
          items-center justify-center
          rounded-xl
          text-[var(--color-primary)]

          hover:bg-[var(--color-primary-soft)]

          focus-visible:outline
          focus-visible:outline-3
          focus-visible:outline-[var(--color-primary)]
          focus-visible:outline-offset-4
        "
      >
        <Icon name="chevron-left" />
      </Link>

      <div
        className="
          relative h-[7px] w-28
          justify-self-center
          overflow-hidden
          rounded-full
          bg-[var(--color-progress-track)]

          before:absolute
          before:inset-y-0
          before:left-0
          before:w-[28%]
          before:rounded-[inherit]
          before:bg-[var(--color-primary)]
          before:content-['']
        "
      >
        <span className="sr-only">{stepLabel}</span>
      </div>
    </nav>
  );
}