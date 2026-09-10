import Link from "next/link";
import type { MouseEvent } from "react";
import Button from "./button";

type AppealNavigationProps = {
  backHref: string;
  skipHref: string;
  primaryText: string;
  primaryHref?: string;
  primaryType?: "button" | "submit";
  primaryVariant?: "primary" | "secondary";
  primaryClassName?: string;
  primaryOnClick?: () => void;
  primaryDisabled?: boolean;
  onSkip?: (event: MouseEvent<HTMLAnchorElement>) => void;
};

const navigationLinkClassName = `
  rounded-[3px] text-sm leading-5 text-[#85899b]
  transition-colors hover:text-[#4562f0]
  focus-visible:outline-2 focus-visible:outline-offset-4
  focus-visible:outline-[#4562f0]
`;

export default function AppealNavigation({
  backHref,
  skipHref,
  primaryText,
  primaryHref,
  primaryType = "button",
  primaryVariant = "primary",
  primaryClassName = "",
  primaryOnClick,
  primaryDisabled = false,
  onSkip,
}: AppealNavigationProps) {
  return (
    <nav
      aria-label="Навигация по обращению"
      className="relative mt-auto min-h-[110px] w-full shrink-0 pt-6"
    >
      <div className="flex justify-center">
        <Button
          text={primaryText}
          link={primaryHref}
          type={primaryType}
          onClick={primaryOnClick}
          disabled={primaryDisabled}
          variant={primaryVariant}
          size="default"
          className={`!w-full ${primaryClassName}`}
        />
      </div>

      <Link href={backHref} className={`${navigationLinkClassName} absolute bottom-0 left-0`}>
        Вернуться назад
      </Link>

      <Link href={skipHref} onClick={onSkip} className={`${navigationLinkClassName} absolute bottom-0 left-1/2 -translate-x-1/2 max-[699px]:right-0 max-[699px]:left-auto max-[699px]:translate-x-0`}>
        Пропустить
      </Link>
    </nav>
  );
}
