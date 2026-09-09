import Button from "@/shared/ui/button";
import Link from "next/link";
import { getAppealRole } from "@/features/appeal/roles";
import SiteHeader from "@/widgets/site-header/site-header";

const topics = [
  "травля и оскорбления",
  "кибербуллинг",
  "конфликт с родителями",
  "конфликт с одноклассниками",
  "конфликт с учителем",
  "давление и угрозы",
  "юридический вопрос",
  "конфликт с сестрой/братом",
  "я не знаю, как это назвать",
];

export default async function AppealTopics({ searchParams }: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);

  return (
    <main className="flex min-h-dvh flex-col bg-[var(--color-background)] [--color-primary:#465fff]">
      <div className="">
        <SiteHeader />
      </div>

      <section
        aria-labelledby="topics-heading"
        className="relative isolate flex flex-1 flex-col overflow-hidden bg-[radial-gradient(ellipse_at_0%_65%,#eaf0ff_0%,transparent_42%),radial-gradient(ellipse_at_100%_100%,#edf1ff_0%,transparent_40%),radial-gradient(ellipse_at_95%_15%,#faf8f5_0%,transparent_25%)] px-5 pt-9 pb-8 sm:px-[4.65%] sm:pt-[51px] sm:pb-[38px]"
      >
        <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col">
          <h1 id="topics-heading" className="text-[28px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[var(--color-primary)] sm:text-[34px]">
            С чем это связано?
          </h1>
          <p id="topics-description" className="mt-[5px] text-[14px] leading-[22px] text-[#151515]">
            Можно выбрать одну или несколько тем, которые ближе всего к твоей ситуации.
          </p>

          <fieldset aria-describedby="topics-description" className="mt-10 min-w-0 sm:mt-[84px]">
            <legend className="sr-only">Темы обращения</legend>
            <div className="grid grid-cols-1 gap-x-4 gap-y-[22px] sm:grid-cols-2 lg:grid-cols-3">
              {topics.map((topic, index) => (
                <label key={topic} className="relative cursor-pointer">
                  <input type="checkbox" name="topics" value={topic} className="peer sr-only" />
                  <span className={`flex min-h-9 items-center justify-center rounded-[13px] border bg-white/65 px-3 py-[7px] text-center text-[14px] leading-5 transition-colors hover:border-[var(--color-primary)] hover:bg-[#eef1ff] peer-checked:border-[var(--color-primary)] peer-checked:bg-[#dfe5ff] peer-checked:text-[var(--color-primary)] peer-focus-visible:outline-2 peer-focus-visible:outline-offset-4 peer-focus-visible:outline-[var(--color-primary)] motion-reduce:transition-none ${index === 8 ? "border-[var(--color-primary)] text-[var(--color-primary)]" : "border-[#222c50] text-[var(--color-text)]"}`}>
                    {topic}
                  </span>
                </label>
              ))}
            </div>
          </fieldset>

          <details className="group mt-[21px] text-[14px] text-[#151515]">
            <summary className="flex w-fit cursor-pointer list-none items-center gap-[14px] rounded-sm py-1 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[var(--color-primary)] [&::-webkit-details-marker]:hidden">
              <span aria-hidden="true" className="text-[28px] leading-5 font-extralight group-open:rotate-45">+</span>
              Добавить свой вариант
            </summary>
            <label className="mt-3 block max-w-[550px]">
              <span className="sr-only">Свой вариант темы обращения</span>
              <input name="customTopic" type="text" placeholder="Расскажи, с чем это связано" className="w-full rounded-[13px] border border-[#222c50] bg-white/75 px-4 py-3 text-[14px] outline-offset-4 focus-visible:outline-[var(--color-primary)]" />
            </label>
          </details>

          <nav aria-label="Навигация по обращению" className="relative mt-auto flex flex-col items-center gap-4 pt-[84px] sm:gap-0">
            <div className="flex flex-col items-center gap-[9px]">
              <Button text="Продолжить" fill link="/appeal/appeal2/appeal3" className="w-[138px] [--button-height:41px] [--button-padding:10px_20px] [--button-radius:12px] [--button-font-size:14px] [--button-weight:500]" />
              <Link href="/" className="rounded-sm text-[14px] leading-[19px] text-[#85899b] hover:text-[var(--color-primary)] focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[var(--color-primary)]">
                Пропустить
              </Link>
            </div>
            <Link href={`/appeal?role=${role.id}`} className="rounded-sm text-[14px] leading-[19px] text-[#85899b] hover:text-[var(--color-primary)] focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[var(--color-primary)] sm:absolute sm:bottom-0 sm:left-0">
              Вернуться назад
            </Link>
          </nav>
        </div>
      </section>
    </main>
  );
}
