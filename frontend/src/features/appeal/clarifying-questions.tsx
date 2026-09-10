"use client";

import { useEffect, useState } from "react";
import AppealNavigation from "@/shared/ui/appeal-navigation";
import { getPublicQuestions, type PublicQuestion } from "@/shared/api/public-api";
import { getAppealRoute } from "./routes";
import { readAppealDraft, type DraftAnswer, updateAppealDraft } from "./appeal-draft";

export default function ClarifyingQuestions({ role, topic = "", formal = false }: { role: string; topic?: string; formal?: boolean }) {
  const [initialDraft] = useState(readAppealDraft);
  const [questions, setQuestions] = useState<PublicQuestion[]>([]);
  const [answers, setAnswers] = useState<Record<string, DraftAnswer>>(initialDraft.answers ?? {});
  const [skipped, setSkipped] = useState<Record<string, boolean>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const routeParams = { role, topic };

  useEffect(() => {
    if (!initialDraft.categoryId) {
      queueMicrotask(() => {
        setError("Сначала выберите тему обращения.");
        setLoading(false);
      });
      return;
    }
    const controller = new AbortController();
    getPublicQuestions(initialDraft.categoryId, controller.signal)
      .then(({ questions: loaded }) => setQuestions(loaded))
      .catch((reason: unknown) => {
        if ((reason as Error).name !== "AbortError") setError("Не удалось загрузить уточняющие вопросы.");
      })
      .finally(() => setLoading(false));
    return () => controller.abort();
  }, [initialDraft.categoryId]);

  function setAnswer(questionID: number, answerID: number, text: string) {
    setAnswers((current) => ({ ...current, [questionID]: { answerId: answerID, text } }));
    setSkipped((current) => ({ ...current, [questionID]: false }));
  }

  function toggleSkipped(questionID: number) {
    const key = String(questionID);
    const next = !skipped[key];
    setSkipped((current) => ({ ...current, [key]: next }));
    if (next) {
      setAnswers((current) => {
        const copy = { ...current };
        delete copy[key];
        return copy;
      });
    }
  }

  function saveAnswers() {
    updateAppealDraft({ answers, role });
  }

  return (
    <section
      className="flex flex-1 flex-col px-[4.5%] pt-9 pb-[38px] max-[699px]:px-5 max-[699px]:pt-8 max-[699px]:pb-7 max-[379px]:px-3 max-[379px]:pt-6"
      aria-labelledby="details-heading"
    >
      <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col">
        <form className="flex flex-1 flex-col">
          <div className="rounded-[15px] border border-[var(--color-primary)] bg-[var(--color-background)] px-[30px] py-[50px] max-[699px]:px-5 max-[699px]:py-7 max-[379px]:px-4">
            <header>
              <h1 id="details-heading" className="m-0 text-[28px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[#4562f0]">{formal ? "Несколько уточнений" : "Пару уточнений"}</h1>
              <p className="mt-px mb-0 text-[14px] leading-5 font-normal text-[#17191f]">Эти вопросы необязательные, но помогут лучше понять ситуацию.</p>
            </header>
            {loading && <p className="mt-5 text-sm text-[#646d86]">Загружаем вопросы…</p>}
            {error && <p role="alert" className="mt-5 text-sm text-[#b42318]">{error}</p>}
            <div className="mt-4 w-full max-w-[670px]">
            {questions.map((question, questionIndex) => {
              const questionKey = String(question.id);
              const isSkipped = skipped[questionKey];

              return (
                <fieldset
                  className={`m-0 min-w-0 border-0 p-0 ${questionIndex > 0 ? "mt-5" : ""}`}
                  key={questionKey}
                >
                  <legend
                    className={`block w-full p-0 text-[16px] leading-5 font-medium ${
                      isSkipped ? "text-[#b1b3ba]" : "text-[#11131a]"
                    }`}
                  >
                    {question.text}
                  </legend>

                  <div
                    className="mt-[7px] grid w-full grid-cols-1 gap-2 sm:grid-cols-3 sm:gap-x-3"
                  >
                    {question.answers.map((option) => {
                      const inputId = `${questionKey}-${option.id}`;

                      return (
                        <label className="min-w-0 cursor-pointer" htmlFor={inputId} key={option.id}>
                          <input
                            checked={answers[questionKey]?.answerId === option.id}
                            className="peer sr-only"
                            disabled={isSkipped}
                            id={inputId}
                            name={questionKey}
                            onChange={() => setAnswer(question.id, option.id, option.text)}
                            type="radio"
                            value={option.id}
                          />
                          <span className="flex min-h-[30px] w-full items-center justify-center rounded-[12px] border border-[#18223f] bg-white/20 px-3 py-1 text-center text-[12px] leading-4 font-normal text-[#18213e] transition-[border-color,background-color,color] duration-150 hover:border-[#4562f0] hover:bg-[#eef1ff] peer-checked:border-[#4562f0] peer-checked:bg-[#4562f0] peer-checked:text-white peer-disabled:cursor-default peer-disabled:border-[#b9becb] peer-disabled:bg-transparent peer-disabled:text-[#b8bbc5] peer-focus-visible:outline-2 peer-focus-visible:outline-offset-3 peer-focus-visible:outline-[#4562f0] motion-reduce:transition-none">
                            {option.text}
                          </span>
                        </label>
                      );
                    })}
                  </div>

                  <label
                    className={`relative mt-1 flex w-fit cursor-pointer items-center gap-1.5 text-[12px] leading-[14px] font-normal ${
                      isSkipped ? "text-[#5470ff]" : "text-[#8e929d]"
                    }`}
                  >
                    <input
                      checked={isSkipped}
                      className="peer m-0 grid size-[14px] shrink-0 cursor-pointer appearance-none place-content-center rounded-[3px] border border-[#a8adb8] bg-transparent checked:border-[#5b75ff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]"
                      onChange={() => toggleSkipped(question.id)}
                      type="checkbox"
                    />
                    <svg
                      aria-hidden="true"
                      className="pointer-events-none absolute top-[3px] left-[3px] hidden h-[7px] w-2 text-[#4562f0] peer-checked:block"
                      viewBox="0 0 10 7"
                    >
                      <path d="M1 3.5 3.6 6 9 1" fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.4" />
                    </svg>
                    <span>Не отвечать на вопрос</span>
                  </label>
                </fieldset>
              );
            })}
            </div>
          </div>
          <AppealNavigation backHref={getAppealRoute("description", routeParams)} skipHref={getAppealRoute("attachments", routeParams)} onSkip={() => updateAppealDraft({ answers: {} })} primaryText="Продолжить" primaryHref={getAppealRoute("attachments", routeParams)} primaryOnClick={saveAnswers} />
        </form>
      </div>
    </section>
  );
}
