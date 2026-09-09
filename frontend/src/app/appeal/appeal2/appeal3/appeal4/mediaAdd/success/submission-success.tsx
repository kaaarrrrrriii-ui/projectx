"use client";

import Button from "@/shared/ui/button";
import InfoCard from "@/shared/ui/info-card";
import Surface from "@/shared/ui/surface";
import Image from "next/image";
import { useRef, useState } from "react";

export default function SubmissionSuccess({
  trackNumber,
  formal,
}: {
  trackNumber: string;
  formal: boolean;
}) {
  const [isCopied, setIsCopied] = useState(false);
  const resetTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  async function copyTrackNumber() {
    try {
      await navigator.clipboard.writeText(trackNumber);
    } catch {
      const textarea = document.createElement("textarea");
      textarea.value = trackNumber;
      textarea.style.position = "fixed";
      textarea.style.opacity = "0";
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand("copy");
      textarea.remove();
    }

    setIsCopied(true);
    if (resetTimer.current) clearTimeout(resetTimer.current);
    resetTimer.current = setTimeout(() => setIsCopied(false), 1800);
  }

  return (
    <section
      aria-labelledby="success-heading"
      className="
        relative isolate flex flex-1 overflow-hidden
        bg-[radial-gradient(ellipse_at_0%_12%,rgba(222,230,255,0.9)_0%,transparent_19%),radial-gradient(ellipse_at_100%_19%,rgba(247,243,238,0.92)_0%,transparent_24%),radial-gradient(ellipse_at_13%_84%,rgba(227,234,255,0.75)_0%,transparent_25%),radial-gradient(ellipse_at_93%_84%,rgba(232,237,255,0.82)_0%,transparent_28%)]
        px-[4.6vw] pt-[58px] pb-[18px]

        max-[699px]:overflow-visible
        max-[699px]:px-5
        max-[699px]:pt-8
        max-[699px]:pb-7
      "
    >
      <div className="relative mx-auto flex w-full max-w-[1440px] flex-col">
        <header className="relative z-20 w-[59.2%] max-w-[690px] max-[699px]:w-full max-[699px]:max-w-none">
          <h1
            id="success-heading"
            className="text-[32px] leading-[1.2] font-black tracking-[-0.025em] text-[var(--color-primary)] max-[699px]:text-[27px]"
          >
            {formal ? "Спасибо, что поделились!" : "Спасибо за твоё обращение!"}
          </h1>
          <p className="mt-1 text-[13px] leading-[1.55] text-[#151515] max-[699px]:mt-2 max-[699px]:text-[13px]">
            {formal ? "Ваше" : "Твоё"} обращение успешно отправлено. Мы рядом
            и уже работаем над тем, чтобы помочь.
          </p>
        </header>

        <div
          className="
            relative z-20 mt-[28px] w-[59.2%] max-w-[690px]
            [&_article]:min-h-[79px]
            [&_article]:gap-[13px]
            [&_article]:rounded-[12px]
            [&_article]:bg-[#dfe6ff]/90
            [&_article]:px-3
            [&_article]:py-2
            [&_article_h2]:text-[16px]
            [&_article_h2]:font-medium
            [&_article_h2]:text-[#11131a]
            [&_article_p]:mt-[4px]
            [&_article_p]:text-[13px]
            [&_article_p]:leading-[1.35]
            [&_article_p]:text-[#151515]
            [&_article>div:first-child]:w-4
            [&_article>div:first-child]:self-center

            max-[699px]:w-full
            max-[699px]:[&_article]:min-h-[102px]
            max-[699px]:[&_article_h2]:text-[15px]
            max-[699px]:[&_article_p]:text-xs
          "
        >
          <InfoCard
            title="Что дальше?"
            description={
              formal
                ? "Наш специалист изучит ваше обращение. Вы сможете в любое время проверить статус по трек-номеру."
                : "Наш специалист изучит твоё обращение. Ты сможешь в любое время проверить статус по трек-номеру."
            }
            icon={ ""
            }
          />
        </div>

        <Image
          src="/images/girl.png"
          alt=""
          width={372}
          height={372}
          priority
          sizes="(max-width: 699px) 260px, 29vw"
          className="
            pointer-events-none absolute top-[-75px] right-[2%] z-10
            h-auto w-[clamp(290px,29vw,360px)]

            max-[899px]:right-0
            max-[699px]:static
            max-[699px]:mx-auto
            max-[699px]:mt-5
            max-[699px]:w-[260px]
          "
        />

        <Surface
          as="section"
          labelledBy="track-number-heading"
          className="
            relative z-30 mt-[102px] flex min-h-[186px] w-full flex-col items-center
            justify-center px-6 py-[25px]
            [--surface-radius:12px]
            [--surface-shadow:none]
            !border-[#333]

            max-[699px]:mt-5
            max-[699px]:min-h-[190px]
            max-[699px]:px-4
          "
        >
          <h2
            id="track-number-heading"
            className="text-center text-[17px] leading-[1.4] font-medium text-[#11131a]"
          >
            {formal ? "Ваш трек-номер" : "Твой трек-номер"}
          </h2>

          <output
            aria-label="Трек-номер обращения"
            className="mt-2 flex h-[51px] w-full max-w-[479px] items-center justify-center rounded-[12px] border border-[var(--color-primary)] bg-white/70 px-4 text-center text-[24px] leading-none font-extrabold tracking-[-0.02em] text-[var(--color-primary)] max-[699px]:text-[clamp(16px,5vw,21px)]"
          >
            {trackNumber}
          </output>

          <Button
            text={isCopied ? "Номер скопирован" : "Скопировать номер"}
            ariaLabel="Скопировать трек-номер"
            variant="primary"
            size="default"
            onClick={() => void copyTrackNumber()}
            className="mt-[9px] h-[39px] w-full !max-w-[479px] rounded-[9px] px-5 text-[13px] font-normal"
          />
          <span className="sr-only" aria-live="polite">
            {isCopied ? "Трек-номер скопирован" : ""}
          </span>
        </Surface>

        <div className="relative z-30 mt-10 flex justify-center max-[699px]:mt-6">
          <Button
            text="Вернуться на главную страницу"
            variant="secondary"
            size="default"
            link="/"
            className="h-[40px] w-[272px] rounded-[10px] px-5 text-[13px] font-normal"
          />
        </div>
      </div>
    </section>
  );
}
