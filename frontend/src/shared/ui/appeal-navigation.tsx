import Link from "next/link";
import Button from "./button";

type AppealNavigationProps = {
  backHref: string;
  skipHref: string;
  primaryText: string;
  primaryHref?: string;
  primaryType?: "button" | "submit";
  primaryVariant?: "primary" | "secondary";
  primaryClassName?: string;
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
}: AppealNavigationProps) {
  return (
    <nav
      aria-label="Навигация по обращению"
      className="relative mt-auto min-h-[110px] w-full shrink-0 pt-6 max-[699px]:grid max-[699px]:min-h-0 max-[699px]:grid-cols-2 max-[699px]:gap-x-4 max-[699px]:gap-y-4 max-[699px]:pt-5"
    >
      <div className="flex justify-center max-[699px]:col-span-2 max-[699px]:[&>*]:w-full">
        <Button
          text={primaryText}
          link={primaryHref}
          type={primaryType}
          variant={primaryVariant}
          size="default"
          className={primaryClassName}
        />
      </div>

      <Link href={backHref} className={`${navigationLinkClassName} absolute bottom-0 left-0 max-[699px]:static max-[699px]:justify-self-start`}>
        Вернуться назад
      </Link>

      <Link href={skipHref} className={`${navigationLinkClassName} absolute bottom-0 left-1/2 -translate-x-1/2 max-[699px]:static max-[699px]:justify-self-end max-[699px]:translate-x-0`}>
        Пропустить
      </Link>
    </nav>
  );
}
