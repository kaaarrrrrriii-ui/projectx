import Button from "@/shared/ui/button";
import Checkbox from "@/shared/ui/checkbox";
import Input from "@/shared/ui/input";
import Link from "next/link";
import { appealCategories } from "@/features/appeal/categories";
import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import { getAppealRoute } from "@/features/appeal/routes";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function AppealTopics({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);
  const formal = isFormalAppealRole(role.id);

  return (
    <main className="flex min-h-dvh flex-col bg-[var(--color-background)]">
      <SiteHeader formal={formal} />

      <section
        aria-labelledby="topics-heading"
        className="
          relative isolate flex flex-1 flex-col overflow-hidden
          bg-[radial-gradient(ellipse_at_0%_65%,#eaf0ff_0%,transparent_42%),radial-gradient(ellipse_at_100%_100%,#edf1ff_0%,transparent_40%),radial-gradient(ellipse_at_95%_15%,#faf8f5_0%,transparent_25%)]
          px-5 pt-9 pb-8
          sm:px-[4.65%]
          sm:pt-[51px]
          sm:pb-[38px]
        "
      >
        <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col">
          <h1
            id="topics-heading"
            className="
              text-[28px]
              font-extrabold
              leading-[1.2]
              tracking-[-0.025em]
              text-[#4562F0]
              sm:text-[34px]
            "
          >
            С чем это связано?
          </h1>

          <p
            id="topics-description"
            className="mt-[5px] text-[14px] leading-[22px] text-[#151515]"
          >
            {formal ? "Можете" : "Можешь"} {" "} выбрать одну или несколько тем, которые ближе всего к {formal ? "вашей" : "твоей"}
            {" "}ситуации.
          </p>

          <fieldset
            aria-describedby="topics-description"
            className="mt-10 min-w-0 sm:mt-[84px]"
          >
            <legend className="sr-only">
              Темы обращения
            </legend>

            <div
              className="
                grid grid-cols-1
                gap-x-4 gap-y-[22px]
                sm:grid-cols-2
                lg:grid-cols-3
              "
            >
              {appealCategories.map((topic) => (
                <Checkbox
                  key={topic}
                  variant="chip"
                  name="topics"
                  value={topic}
                  label={topic}
                />
              ))}
            </div>
          </fieldset>

          <details className="group mt-[21px] text-[14px] text-[#151515]">
            <summary
              className="
                flex w-fit cursor-pointer list-none
                items-center gap-[14px]
                rounded-sm py-1

                focus-visible:outline-2
                focus-visible:outline-offset-4
                focus-visible:outline-[#4562F0]

                [&::-webkit-details-marker]:hidden
              "
            >
              <span
                aria-hidden="true"
                className="
                  text-[28px]
                  leading-5
                  font-extralight
                  group-open:rotate-45
                "
              >
                +
              </span>

              Добавить свой вариант
            </summary>

            <label className="mt-3 block max-w-[550px]">
              <span className="sr-only">
                Свой вариант темы обращения
              </span>

              <Input
                name="customTopic"
                type="text"
                placeholder={formal ? "Расскажите, с чем это связано" : "Расскажи, с чем это связано"}
                className="w-full"
              />
            </label>
          </details>

          <nav
            aria-label="Навигация по обращению"
            className="
              relative mt-auto flex flex-col
              items-center gap-4 pt-[84px]
              sm:gap-0
            "
          >
            <div className="flex flex-col items-center gap-[9px]">
              <Button
                text="Продолжить"
                variant="primary"
                size="default"
                link={getAppealRoute("description", role.id)}
              />

              <Link
                href={getAppealRoute("description", role.id)}
                className="
                  rounded-sm
                  text-[14px]
                  leading-[19px]
                  text-[#85899b]

                  hover:text-[#4562F0]

                  focus-visible:outline-2
                  focus-visible:outline-offset-4
                  focus-visible:outline-[#4562F0]
                "
              >
                Пропустить
              </Link>
            </div>

            <Link
              href={getAppealRoute("role", role.id)}
              className="
                rounded-sm
                text-[14px]
                leading-[19px]
                text-[#85899b]

                hover:text-[#4562F0]

                focus-visible:outline-2
                focus-visible:outline-offset-4
                focus-visible:outline-[#4562F0]

                sm:absolute
                sm:bottom-0
                sm:left-0
              "
            >
              Вернуться назад
            </Link>
          </nav>
        </div>
      </section>
    </main>
  );
}
