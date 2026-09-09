"use client";

import Button from "@/shared/ui/button";
import Image from "next/image";
import { FormEvent, useState } from "react";
import DialogShell from "./dialog-shell";

type CloseStage = "confirm" | "result" | "feedback" | "closed";

export default function CloseAppealFlow({
  trackNumber,
  onCancel,
}: {
  trackNumber: string;
  onCancel: () => void;
}) {
  const [stage, setStage] = useState<CloseStage>("confirm");
  const [feedback, setFeedback] = useState("");

  function submitFeedback(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setStage("closed");
  }

  if (stage === "confirm") {
    return (
      <DialogShell labelledBy="close-confirm-heading" onClose={onCancel} className="max-w-[884px]">
        <h2
          id="close-confirm-heading"
          className="text-center text-[30px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[var(--color-primary)] max-[599px]:text-[23px]"
        >
          Ты уверен, что хочешь закрыть обращение?
        </h2>
        <div className="mx-auto mt-10 grid max-w-[666px] grid-cols-2 gap-2.5 max-[499px]:mt-7 max-[499px]:grid-cols-1">
          <Button
            text="Да"
            variant="secondary"
            size="small"
            onClick={() => setStage("result")}
            className="w-full"
          />
          <Button
            text="Нет"
            variant="secondary"
            size="small"
            onClick={onCancel}
            className="w-full"
          />
        </div>
      </DialogShell>
    );
  }

  if (stage === "result") {
    return (
      <DialogShell labelledBy="result-heading" className="max-w-[884px]">
        <h2
          id="result-heading"
          className="text-center text-[30px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[var(--color-primary)] max-[599px]:text-[23px]"
        >
          Мы помогли тебе решить твою проблему?
        </h2>
        <div className="mx-auto mt-10 grid max-w-[666px] grid-cols-2 gap-2.5 max-[499px]:mt-7 max-[499px]:grid-cols-1">
          <Button
            text="Нет"
            variant="secondary"
            size="small"
            onClick={() => setStage("feedback")}
            className="w-full"
          />
          <Button
            text="Да"
            variant="secondary"
            size="small"
            onClick={() => setStage("closed")}
            className="w-full"
          />
        </div>
      </DialogShell>
    );
  }

  if (stage === "feedback") {
    return (
      <DialogShell labelledBy="feedback-heading" className="max-w-[704px]">
        <form onSubmit={submitFeedback}>
          <h2
            id="feedback-heading"
            className="text-[25px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[var(--color-primary)] max-[599px]:text-[22px]"
          >
            Расскажи, что не помогло?
          </h2>
          <p className="mt-2 text-[13px] leading-[1.4] text-[#151515]">
            Спасибо, что делишься этим. Твой ответ поможет нам лучше понять ситуацию и подобрать другие рекомендации.
          </p>

          <textarea
            value={feedback}
            onChange={(event) => setFeedback(event.target.value)}
            placeholder="Что конкретно тебе не помогло?"
            aria-label="Что не помогло"
            className="mt-4 block min-h-[160px] w-full resize-y rounded-[12px] border border-[#333] bg-[#fcfdff] px-3 py-2.5 text-[14px] leading-5 outline-none placeholder:text-[#9196a7] focus:border-[var(--color-primary)] focus:ring-2 focus:ring-[#4562f0]/20"
          />

          <div className="mt-4 flex items-center gap-3 rounded-[11px] border border-[var(--color-primary)] bg-[#dee7fd] px-3 py-2.5">
            <Image
              src="/images/LightbulbFilament.svg"
              alt=""
              width={22}
              height={22}
              className="shrink-0"
            />
            <p className="text-[12px] leading-[1.35]">
              Твоё обращение вернётся специалисту на дополнительное рассмотрение. Мы постараемся ответить как можно скорее.
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

  return (
    <DialogShell labelledBy="closed-heading" className="max-w-[704px]">
      <h2
        id="closed-heading"
        className="text-[25px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[var(--color-primary)] max-[599px]:text-[22px]"
      >
        Обращение закрыто
      </h2>
      <p className="mt-2 text-[13px] leading-[1.4] text-[#151515]">
        Мы рады, что смогли быть рядом. Если тебе снова понадобится поддержка — ты всегда можешь написать нам.
      </p>

      <section className="mt-4 rounded-[13px] border border-[var(--color-primary)] bg-[#dee7fd] p-2.5" aria-labelledby="save-track-heading">
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
              Сохрани номер обращения
            </h3>
            <p className="mt-0.5 text-[12px] leading-[1.35]">
              Он поможет быстрее найти твою историю, если ты вернёшься.
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
