"use client";

import AppealNavigation from "@/shared/ui/appeal-navigation";
import { getAppealRoute } from "./routes";
import { useState } from "react";
import {
  clarifyingQuestions,
  type ClarifyingQuestionKey,
} from "./clarifying-question-data";

export default function ClarifyingQuestions({ role, topic = "", formal = false }: { role: string; topic?: string; formal?: boolean }) {
  const [answers, setAnswers] = useState<Record<ClarifyingQuestionKey, string>>({
    place: "",
    duration: "",
    askedForHelp: "",
  });
  const [skipped, setSkipped] = useState<Record<ClarifyingQuestionKey, boolean>>({
    place: false,
    duration: false,
    askedForHelp: false,
  });
  const routeParams = { role, topic };

  const setAnswer = (question: ClarifyingQuestionKey, answer: string) => {
    setAnswers((current) => ({ ...current, [question]: answer }));
  };

  const toggleSkipped = (question: ClarifyingQuestionKey) => {
    const nextValue = !skipped[question];

    setSkipped((current) => ({ ...current, [question]: nextValue }));

    if (nextValue) {
      setAnswers((current) => ({ ...current, [question]: "" }));
    }
  };

  return (
    <section
      className="flex flex-1 flex-col px-[4.5%] pt-9 pb-[38px] max-[699px]:px-5 max-[699px]:pt-8 max-[699px]:pb-7"
      aria-labelledby="details-heading"
    >
      <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col">
        <form className="flex flex-1 flex-col">
          <div className="rounded-[15px] border border-[var(--color-primary)] bg-[var(--color-background)] px-[30px] py-[50px]">
            <header>
              <h1
                id="details-heading"
                className="m-0 text-[28px] leading-[1.2] font-extrabold tracking-[-0.035em] text-[#4562f0] sm:text-[30px] sm:leading-9"
              >
                {formal ? "Несколько уточнений" : "Пару уточнений"}
              </h1>
              <p className="mt-px mb-0 text-[14px] leading-5 font-normal text-[#17191f]">
                Эти вопросы необязательные, но помогут лучше понять ситуацию.
              </p>
            </header>

            <div className="mt-4 w-full max-w-[670px]">
            {clarifyingQuestions.map((question, questionIndex) => {
              if (!question.key) return null;
              const questionKey = question.key;
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
                    {formal && questionKey === "askedForHelp"
                      ? "Вы уже обращались за помощью?"
                      : question.title}
                  </legend>

                  <div
                    className="mt-[7px] grid w-full grid-cols-1 gap-2.5 sm:grid-cols-3 sm:gap-x-3"
                  >
                    {question.options.map((option, optionIndex) => {
                      const inputId = `${questionKey}-${optionIndex}`;

                      return (
                        <label className="min-w-0 cursor-pointer" htmlFor={inputId} key={option}>
                          <input
                            checked={answers[questionKey] === option}
                            className="peer sr-only"
                            disabled={isSkipped}
                            id={inputId}
                            name={questionKey}
                            onChange={() => setAnswer(questionKey, option)}
                            type="radio"
                            value={option}
                          />
                          <span className="flex min-h-[30px] w-full items-center justify-center rounded-[12px] border border-[#18223f] bg-white/20 px-3 py-1 text-center text-[12px] leading-4 font-normal text-[#18213e] transition-[border-color,background-color,color] duration-150 hover:border-[#4562f0] hover:bg-[#eef1ff] peer-checked:border-[#4562f0] peer-checked:bg-[#4562f0] peer-checked:text-white peer-disabled:cursor-default peer-disabled:border-[#b9becb] peer-disabled:bg-transparent peer-disabled:text-[#b8bbc5] peer-focus-visible:outline-2 peer-focus-visible:outline-offset-3 peer-focus-visible:outline-[#4562f0] motion-reduce:transition-none">
                            {formal &&
                            questionKey === "askedForHelp" &&
                            option === "Не уверен"
                              ? "Не уверены"
                              : option}
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
                      onChange={() => toggleSkipped(questionKey)}
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

          <AppealNavigation
            backHref={getAppealRoute("description", routeParams)}
            skipHref={getAppealRoute("attachments", routeParams)}
            primaryText="Продолжить"
            primaryHref={getAppealRoute("attachments", routeParams)}
          />
        </form>
      </div>
    </section>
  );
}
