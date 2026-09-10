import Image from "next/image";
import Link from "next/link";

export default function SiteHeader({ formal = false, compact = false }: { formal?: boolean; compact?: boolean }) {
  return (
    <header
      className={`
        sticky top-0 z-50 box-border h-[100px] w-full shrink-0
        rounded-b-[30px] border border-[#4562f0] bg-white
        px-5 pt-[30px] pb-[15px]

        ${compact ? "min-[700px]:h-[80px] min-[700px]:rounded-none min-[700px]:px-[15px] min-[700px]:pt-[18px] min-[700px]:pb-[10px]" : ""}

        max-[699px]:h-[82px]
        max-[699px]:rounded-b-[24px]
        max-[699px]:px-5
        max-[699px]:pt-4
        max-[699px]:pb-2.5
      `}
    >
      <div
        className="
          flex h-full w-full items-center justify-between gap-[50px]

          max-[699px]:gap-5

          max-[379px]:gap-3
        "
      >
        <Link
          href="/"
          aria-label="Отклик — на главную"
          className="
            shrink-0 rounded
            focus-visible:outline-2
            focus-visible:outline-[var(--color-primary)]
            focus-visible:outline-offset-4
          "
        >
          <Image
            src="/images/logo.svg"
            alt="Отклик — ты не один, мы рядом"
            width={180}
            height={53}
            priority
            className="
              h-auto w-[clamp(120px,13vw,180px)]
              max-[699px]:w-28
              max-[379px]:w-[100px]
            "
          />
        </Link>

        <p
          className={`
            min-w-0 flex-1 text-right
            text-[clamp(14px,1.4vw,18px)]
            font-semibold
            leading-[1.5]
            text-[var(--color-primary)]

            max-[699px]:max-w-[260px]
            max-[699px]:text-xs

            max-[379px]:text-[9px]
            ${compact ? "min-[700px]:text-[13px]" : ""}
          `}
        >
          {formal
            ? "Безопасный способ рассказать о том, что вас беспокоит, и получить помощь от профессионалов."
            : "Безопасный способ рассказать о том, что тебя беспокоит, и получить помощь от профессионалов."}
        </p>
      </div>
    </header>
  );
}
