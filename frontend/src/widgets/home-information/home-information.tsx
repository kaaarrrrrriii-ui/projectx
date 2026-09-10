import Button from "@/shared/ui/button";
import Image from "next/image";
import InfoCard from "@/shared/ui/info-card";
import PageHeading from "@/shared/ui/page-heading";
import Link from "next/link";

const information = [
  {
    title: "Ты остаёшься анонимным",
    description:
      "Мы не спрашиваем твоё имя и не передаём информацию третьим лицам.",
    icon: "lock-key",
  },
  {
    title: "Тебе помогут специалисты",
    description:
      "Твоё обращение увидят психологи, педагоги и другие эксперты.",
    icon: "user-sound",
  },
  {
    title: "Ты можешь отслеживать статус",
    description:
      "После отправки ты получишь уникальный трек-номер и сможешь проверить, на каком этапе находится обращение.",
    icon: "bell-ringing",
  },
];

function HomeActions() {
  const mobileButtonClasses = `
    text-[13px]
    max-[699px]:h-12
    max-[699px]:rounded-xl
    max-[699px]:px-8
    max-[699px]:py-3.5
    max-[699px]:text-base
  `;

  return (
    <div
      className="
        mt-[22px] flex flex-col items-center gap-2.5
        min-[1200px]:mt-6
        max-[699px]:mt-0
        max-[699px]:gap-3.5
      "
    >
      <div
        className="
          flex justify-center gap-[30px]
          max-[699px]:w-full
          max-[699px]:gap-3
          max-[699px]:[&>a]:flex-1
          max-[419px]:flex-col
        "
      >
        <Button
          text="Проверить статус"
          variant="secondary"
          size="default"
          link="/status"
          className={mobileButtonClasses}
        />

        <Button
          text="Подать обращение"
          variant="primary"
          size="default"
          link="/appeal"
          className={mobileButtonClasses}
        />
      </div>

      <Link
        href="/staff"
        className="
          rounded-md py-0.5 text-center underline-offset-4
          text-[clamp(11px,1.2vw,14px)]
          text-[var(--color-text-subtle)]
          transition-colors hover:text-[var(--color-primary)] hover:underline
          focus-visible:outline-2 focus-visible:outline-[var(--color-primary)]
          focus-visible:outline-offset-4
          max-[699px]:text-xs
        "
      >
        Вход для сотрудников
      </Link>
    </div>
  );
}

export default function HomeInformation() {
  return (
    <section
      aria-labelledby="important-heading"
      className="
        flex w-full flex-col gap-[26px]

        [--heading-size:clamp(26px,2.8vw,40px)]
        [--heading-color:var(--color-primary)]

        min-[1200px]:gap-9

        max-[899px]:max-w-[680px]
        max-[899px]:gap-6
        max-[899px]:[--heading-size:clamp(24px,6vw,30px)]
      "
    >
      <header className="relative left-[-13px] text-center max-[899px]:left-0">
        <PageHeading
          id="important-heading"
          title="Что важно знать?"
        />
      </header>

      <div
        className="
          mt-3 ml-auto grid w-[54.2%] gap-[14px]
          min-[1200px]:gap-[18px]
          max-[899px]:w-full
        "
      >
        {information.map(({ icon, ...card }) => (
          <InfoCard
            key={icon}
            {...card}
            className="
              min-h-[85px] gap-[22px] px-[22px]
              transition-[transform,box-shadow,background-color]
              duration-200 ease-out
              hover:-translate-y-0.5
              hover:bg-white/80
              hover:shadow-[4px_6px_14px_rgba(69,98,240,0.16)]
              motion-reduce:transform-none
              motion-reduce:transition-none
            "
            icon={
              <Image
                src={`/images/${icon}.svg`}
                alt=""
                width={28}
                height={28}
              />
            }
          />
        ))}
      </div>

      <div className="relative top-9 left-[-13px] max-[899px]:top-0 max-[899px]:left-0">
        <HomeActions />
      </div>
    </section>
  );
}
