"use client";

import Button from "@/shared/ui/button";
import Image from "next/image";
import { FormEvent, useState } from "react";
import DialogShell from "./dialog-shell";

type CloseStage = "feedback" | "closed" | "complaint";

export default function CloseAppealFlow({
  formal,
  trackNumber,
  outcome,
  onCancel,
}: {
  formal: boolean;
  trackNumber: string;
  outcome: "helped" | "not-helped";
  onCancel: () => void;
}) {
  const [stage, setStage] = useState<CloseStage>(
    outcome === "not-helped" ? "feedback" : "closed",
  );
  const [feedback, setFeedback] = useState("");
  const [complaint, setComplaint] = useState("");
  const [rating, setRating] = useState(3);

  function submitFeedback(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setStage("closed");
  }

  function submitComplaint(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    onCancel();
  }

  if (stage === "feedback") {
    return (
      <DialogShell labelledBy="feedback-heading" className="max-w-[988px]">
        <form onSubmit={submitFeedback}>
          <h2
            id="feedback-heading"
            className="text-[25px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[var(--color-primary)] max-[599px]:text-[22px]"
          >
            {formal ? "Расскажите, что не помогло?" : "Расскажи, что не помогло?"}
          </h2>
          <p className="mt-2 text-[13px] leading-[1.4] text-[#151515]">
            {formal
              ? "Спасибо, что делитесь этим. Ваш ответ поможет нам лучше понять ситуацию и подобрать другие рекомендации."
              : "Спасибо, что делишься этим. Твой ответ поможет нам лучше понять ситуацию и подобрать другие рекомендации."}
          </p>

          <div className="relative mt-4">
            <textarea
              value={feedback}
              onChange={(event) => setFeedback(event.target.value)}
              placeholder={
                formal
                  ? "Что конкретно вам не помогло?"
                  : "Что конкретно тебе не помогло?"
              }
              aria-label="Что не помогло"
              className="block min-h-[220px] w-full resize-y rounded-[12px] border border-[#333] bg-[#fcfdff] px-3 py-2.5 pr-20 text-[14px] leading-5 outline-none placeholder:text-[#9196a7] focus:border-[var(--color-primary)] focus:ring-2 focus:ring-[#4562f0]/20 max-[499px]:min-h-[170px] max-[379px]:pr-3"
            />
            <Image
              src="/images/roles/ticher.png"
              alt="Специалист"
              width={58}
              height={58}
              className="absolute top-[-28px] right-10 h-[58px] w-[58px] rounded-full bg-[#11131a] object-cover shadow-[0_4px_10px_rgba(0,0,0,0.35)] max-[499px]:right-4 max-[499px]:h-12 max-[499px]:w-12"
            />
          </div>

          <div className="mt-4 flex items-center gap-3 rounded-[11px] border border-[var(--color-primary)] bg-[#dee7fd] px-3 py-2.5">
            <Image
              src="/images/LightbulbFilament.svg"
              alt=""
              width={22}
              height={22}
              className="shrink-0"
            />
            <p className="text-[12px] leading-[1.35]">
              {formal ? "Ваше" : "Твоё"} обращение вернётся специалисту на
              дополнительное рассмотрение. Мы постараемся ответить как можно
              скорее.
            </p>
          </div>

          <Button
            text="Отправить"
            type="submit"
            variant="primary"
            size="small"
            className="mt-[18px] w-full"
          />
          <button
            type="button"
            onClick={() => setStage("closed")}
            className="mx-auto mt-3 block cursor-pointer rounded-sm px-3 py-1 text-[13px] text-[#9196a7] hover:text-[var(--color-primary)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--color-primary)]"
          >
            Пропустить
          </button>
        </form>
      </DialogShell>
    );
  }

  if (stage === "complaint") {
    return (
      <DialogShell
        labelledBy="complaint-heading"
        onClose={() => setStage("closed")}
        showClose
        className="max-w-[988px]"
      >
        <form onSubmit={submitComplaint}>
          <h2
            id="complaint-heading"
            className="pr-12 text-[25px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[var(--color-primary)] max-[599px]:text-[22px]"
          >
            Создание жалобы
          </h2>
          <p className="mt-7 text-[16px] leading-[1.4] text-[#151515] max-[599px]:mt-5 max-[599px]:text-[14px]">
            {formal
              ? "Опишите, что пошло не так?"
              : "Опиши, что пошло не так?"}
          </p>
          <textarea
            value={complaint}
            onChange={(event) => setComplaint(event.target.value)}
            aria-label="Текст жалобы"
            className="mt-7 block min-h-[224px] w-full resize-y rounded-[12px] border border-[#333] bg-[#fcfdff] px-3 py-2.5 text-[14px] leading-5 outline-none focus:border-[var(--color-primary)] focus:ring-2 focus:ring-[#4562f0]/20 max-[599px]:mt-4 max-[499px]:min-h-[180px]"
          />
          <Button
            text="Отправить"
            type="submit"
            variant="primary"
            size="small"
            className="mt-7 w-full max-[599px]:mt-5"
          />
        </form>
      </DialogShell>
    );
  }

  return (
    <DialogShell labelledBy="closed-heading" onClose={onCancel} showClose className="max-w-[988px]">
      <h2
        id="closed-heading"
        className="text-[25px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[var(--color-primary)] max-[599px]:text-[22px]"
      >
        Обращение закрыто
      </h2>
      <p className="mt-2 text-[13px] leading-[1.4] text-[#151515]">
        {formal
          ? "Мы рады, что смогли быть рядом. Если вам снова понадобится поддержка — вы всегда можете написать нам."
          : "Мы рады, что смогли быть рядом. Если тебе снова понадобится поддержка — ты всегда можешь написать нам."}
      </p>

      <section className="mt-7 rounded-[13px] border border-[var(--color-primary)] bg-white px-5 py-5 max-[379px]:px-3 max-[379px]:py-4" aria-labelledby="rating-heading">
        <h3 id="rating-heading" className="text-[16px] font-medium text-[#000828]">
          {formal ? "Оцените работу специалиста" : "Оцени работу специалиста"}
        </h3>
        <div className="mt-5 flex justify-center gap-5 max-[379px]:gap-2" aria-label="Оценка специалиста">
          {[1, 2, 3, 4, 5].map((value) => (
            <button
              key={value}
              type="button"
              onClick={() => setRating(value)}
              aria-label={`${value} из 5`}
              aria-pressed={rating === value}
              className="flex h-11 w-11 cursor-pointer items-center justify-center text-[var(--color-primary)] transition-transform hover:scale-110 focus-visible:rounded focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--color-primary)] max-[499px]:h-9 max-[499px]:w-9 max-[379px]:h-8 max-[379px]:w-8"
            >
              <svg viewBox="0 0 24 24" aria-hidden="true" className="h-full w-full">
                <path
                  d="m12 2.8 2.76 5.59 6.17.9-4.47 4.35 1.06 6.15L12 16.89l-5.52 2.9 1.06-6.15-4.47-4.35 6.17-.9L12 2.8Z"
                  fill={value <= rating ? "currentColor" : "none"}
                  stroke="currentColor"
                  strokeWidth="1.7"
                  strokeLinejoin="round"
                />
              </svg>
            </button>
          ))}
        </div>
        <div className="mx-auto mt-5 grid max-w-[496px] grid-cols-2 gap-8 max-[499px]:grid-cols-1 max-[499px]:gap-2">
          <Button
            text="Оставить жалобу"
            variant="secondary"
            size="small"
            onClick={() => setStage("complaint")}
            className="w-full"
          />
          <Button
            text="Отправить"
            variant="primary"
            size="small"
            onClick={onCancel}
            className="w-full"
          />
        </div>
      </section>

      <section className="mt-6 rounded-[13px] border border-[var(--color-primary)] bg-[#dee7fd] p-2.5" aria-labelledby="save-track-heading">
        <div className="flex items-start gap-3 px-1">
          <Image
            src="/images/LightbulbFilament.svg"
            alt=""
            width={24}
            height={24}
            className="mt-0.5 shrink-0"
          />
          <div>
            <h3 id="save-track-heading" className="text-[16px] font-medium">
              {formal ? "Сохраните номер обращения" : "Сохрани номер обращения"}
            </h3>
            <p className="mt-0.5 text-[12px] leading-[1.35]">
              {formal
                ? "Он поможет быстрее найти вашу историю, если вы вернётесь."
                : "Он поможет быстрее найти твою историю, если ты вернёшься."}
            </p>
          </div>
        </div>
        <output className="mt-2 flex min-h-8 items-center justify-center rounded-[9px] border border-[var(--color-primary)] bg-white/80 px-3 text-center text-[13px] font-medium text-[var(--color-primary)]">
          {trackNumber}
        </output>
      </section>

      <Button
        text="Вернуться на главный экран"
        variant="primary"
        size="small"
        link="/"
        className="mt-[18px] w-full"
      />
    </DialogShell>
  );
}
