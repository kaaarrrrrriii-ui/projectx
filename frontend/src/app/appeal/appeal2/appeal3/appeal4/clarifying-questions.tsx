"use client";

import Button from "@/shared/ui/button";
import Link from "next/link";
import { useState } from "react";

type QuestionKey = "place" | "duration" | "askedForHelp";

type Question = {
  key: QuestionKey;
  title: string;
  options: string[];
};

const questions: Question[] = [
  {
    key: "place",
    title: "Где это происходит?",
    options: ["В школе", "Онлайн", "В другом месте"],
  },
  {
    key: "duration",
    title: "Как давно это длится?",
    options: ["Недавно", "Давно", "Пару дней"],
  },
  {
    key: "askedForHelp",
    title: "Ты уже обращался за помощью?",
    options: ["Нет", "Да", "Не уверен"],
  },
];

export default function ClarifyingQuestions({ role }: { role: string }) {
  const [answers, setAnswers] = useState<Record<QuestionKey, string>>({
    place: "",
    duration: "",
    askedForHelp: "",
  });
  const [skipped, setSkipped] = useState<Record<QuestionKey, boolean>>({
    place: true,
    duration: false,
    askedForHelp: false,
  });

  const setAnswer = (question: QuestionKey, answer: string) => {
    setAnswers((current) => ({ ...current, [question]: answer }));
  };

  const toggleSkipped = (question: QuestionKey) => {
    const nextValue = !skipped[question];

    setSkipped((current) => ({ ...current, [question]: nextValue }));

    if (nextValue) {
      setAnswers((current) => ({ ...current, [question]: "" }));
    }
  };

  return (
    <section
      className="flex-1 bg-[#f7f9fe] px-[61px] pt-[55px] pb-[39px]"
      aria-labelledby="details-heading"
    >
      <div className="mx-auto w-full max-w-[1440px]">
        <header>
          <h1
            id="details-heading"
            className="m-0 text-[34px] leading-[41px] font-extrabold tracking-[-0.035em] text-[#4562f0]"
          >
            Пару уточнений
          </h1>
          <p className="mt-px mb-0 text-[14px] leading-5 font-normal text-[#17191f]">
            Эти вопросы необязательные, но помогут лучше понять ситуацию.
          </p>
        </header>

        <form className="mt-5">
          <div className="w-[772px]">
            {questions.map((question, questionIndex) => {
              const isSkipped = skipped[question.key];

              return (
                <fieldset
                  className={`m-0 min-w-0 border-0 p-0 ${questionIndex > 0 ? "mt-[34px]" : ""}`}
                  key={question.key}
                >
                  <legend
                    className={`block w-full p-0 text-[19px] leading-6 font-medium ${
                      isSkipped ? "text-[#b1b3ba]" : "text-[#11131a]"
                    }`}
                  >
                    {question.title}
                  </legend>

                  <div
                    className={`mt-[9px] grid grid-cols-[repeat(3,246px)] ${
                      questionIndex === 0 ? "w-[772px] gap-x-[17px]" : "w-[755px] gap-x-[8.5px]"
                    }`}
                  >
                    {question.options.map((option, optionIndex) => {
                      const inputId = `${question.key}-${optionIndex}`;

                      return (
                        <label className="min-w-0 cursor-pointer" htmlFor={inputId} key={option}>
                          <input
                            checked={answers[question.key] === option}
                            className="peer sr-only"
                            disabled={isSkipped}
                            id={inputId}
                            name={question.key}
                            onChange={() => setAnswer(question.key, option)}
                            type="radio"
                            value={option}
                          />
                          <span className="flex h-[38px] w-[246px] items-center justify-center rounded-[14px] border border-[#18223f] bg-white/20 text-[14px] leading-5 font-normal text-[#18213e] transition-[border-color,background-color,color] duration-150 hover:border-[#4562f0] hover:bg-[#eef1ff] peer-checked:border-[#4562f0] peer-checked:bg-[#4562f0] peer-checked:text-white peer-disabled:cursor-default peer-disabled:border-[#b9becb] peer-disabled:bg-transparent peer-disabled:text-[#b8bbc5] peer-focus-visible:outline-2 peer-focus-visible:outline-offset-3 peer-focus-visible:outline-[#4562f0] motion-reduce:transition-none">
                            {option}
                          </span>
                        </label>
                      );
                    })}
                  </div>

                  <label
                    className={`relative mt-1.5 flex w-fit cursor-pointer items-center gap-[7px] text-[14px] leading-[18px] font-normal ${
                      isSkipped ? "text-[#5470ff]" : "text-[#8e929d]"
                    }`}
                  >
                    <input
                      checked={isSkipped}
                      className="peer m-0 grid size-[18px] shrink-0 cursor-pointer appearance-none place-content-center rounded-[5px] border border-[#a8adb8] bg-transparent checked:border-[#5b75ff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]"
                      onChange={() => toggleSkipped(question.key)}
                      type="checkbox"
                    />
                    <svg
                      aria-hidden="true"
                      className="pointer-events-none absolute top-[5px] left-[4px] hidden h-[7px] w-[10px] text-[#4562f0] peer-checked:block"
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

          <nav
            className="relative mt-[17px] flex w-full flex-col items-center"
            aria-label="Навигация по обращению"
          >
            <Button text="Продолжить" variant="primary" size="default" link={`/appeal/appeal2/appeal3/appeal4/mediaAdd?role=${role}`} />
            <Link
              className="mt-2.5 rounded-[3px] text-[14px] leading-5 font-normal text-[#9296a4] no-underline transition-colors duration-150 hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[#4562f0] motion-reduce:transition-none"
              href="/"
            >
              Пропустить
            </Link>
            <Link
              className="absolute bottom-0 left-0 rounded-[3px] text-[14px] leading-5 font-normal text-[#9296a4] no-underline transition-colors duration-150 hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[#4562f0] motion-reduce:transition-none"
              href={`/appeal/appeal2/appeal3?role=${role}`}
            >
              Вернуться назад
            </Link>
          </nav>
        </form>
      </div>
    </section>
  );
}
