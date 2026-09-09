import Image from "next/image";
import Link from "next/link";

export default function SiteHeader() {
  return (
    <header className="border-b border-[#a9b3ff] bg-[var(--color-surface)]">
      <div
        className="
          mx-auto flex min-h-16 max-w-[1600px]
          items-center justify-between gap-8
          px-3 py-2.5

          max-[699px]:min-h-[72px]
          max-[699px]:gap-5
          max-[699px]:px-5
          max-[699px]:py-3

          max-[379px]:gap-3
          max-[379px]:px-4
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
          className="
            min-w-0 flex-1 text-right
            text-[clamp(14px,1.4vw,18px)]
            font-semibold
            leading-[1.5]
            text-[var(--color-primary)]

            max-[699px]:max-w-[260px]
            max-[699px]:text-xs

            max-[379px]:text-[9px]
          "
        >
          Безопасный способ рассказать о том, что тебя беспокоит, и получить
          помощь от профессионалов.
        </p>
      </div>
    </header>
  );
}
