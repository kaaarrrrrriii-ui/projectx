import Button from "@/shared/ui/button";
import Image from "next/image";
import InfoCard from "@/shared/ui/info-card";
import PageHeading from "@/shared/ui/page-heading";

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
    max-[699px]:h-12
    max-[699px]:rounded-xl
    max-[699px]:px-8
    max-[699px]:py-3.5
    max-[699px]:text-base
  `;

  return (
    <div
      className="
        mt-4 flex flex-col items-center gap-2.5
        min-[1200px]:mt-6
        max-[699px]:mt-0
        max-[699px]:gap-3.5
      "
    >
      <div
        className="
          flex justify-center gap-5
          max-[699px]:w-full
          max-[699px]:gap-3
          max-[699px]:[&>a]:flex-1
          max-[419px]:flex-col
        "
      >
        <Button
          text="Подать обращение"
          variant="primary"
          size="default"
          link="/appeal"
          className={mobileButtonClasses}
        />

        <Button
          text="Проверить статус"
          variant="secondary"
          size="default"
          link="/status"
          className={mobileButtonClasses}
        />
      </div>

      <span
        className="
          py-0.5 text-center
          text-[clamp(11px,1.2vw,14px)]
          text-[var(--color-text-subtle)]
          max-[699px]:text-xs
        "
      >
        Вход для сотрудников
      </span>
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

        max-[699px]:max-w-[500px]
        max-[699px]:gap-6
        max-[699px]:[--heading-size:clamp(24px,6vw,30px)]
      "
    >
      <header className="text-center">
        <PageHeading
          id="important-heading"
          title="Что важно знать?"
        />
      </header>

      <div
        className="
          ml-auto grid w-[54.2%] gap-3
          min-[1200px]:gap-[18px]
          max-[699px]:w-full
        "
      >
        {information.map(({ icon, ...card }) => (
          <InfoCard
            key={icon}
            {...card}
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

      <HomeActions />
    </section>
  );
}