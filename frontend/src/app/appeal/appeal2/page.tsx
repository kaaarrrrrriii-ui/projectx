import { getAppealRole } from "@/features/appeal/roles";
import TopicSelection from "@/features/appeal/topic-selection";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function AppealTopics({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[]; topic?: string | string[] }>;
}) {

  const params = await searchParams;
  const initialTopic = typeof params.topic === "string" ? params.topic : "";

  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader />
  const role = getAppealRole((await searchParams).role);
  const formal = isFormalAppealRole(role.id);
      <SiteHeader formal={formal} />


      <section
        aria-labelledby="topics-heading"
        className="
          relative isolate flex flex-1 flex-col overflow-hidden
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
            Выбери тему обращения
          </h1>

          <p
            id="topics-description"
            className="mt-[5px] text-[14px] leading-[22px] text-[#151515]"
          >
            Можно выбрать одну или несколько тем, которые ближе всего к {formal ? "вашей" : "твоей"}
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
                link={`/appeal/appeal2/appeal3?role=${role.id}`}
              />

              <Link
                href={`/appeal/appeal2/appeal3?role=${role.id}`}
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
              href={`/appeal?role=${role.id}`}
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
