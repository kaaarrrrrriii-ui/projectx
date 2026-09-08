import Button from "@/shared/ui/button";
import Image from "next/image";

const cards = [
  {
    title: "Ты остаёшься анонимным",
    description:
      "Мы не спрашиваем твоё имя и не передаём информацию третьим лицам.",
    icon: "/images/lock-key.svg",
  },
  {
    title: "Тебе помогут специалисты",
    description: "Твоё обращение увидят психологи, педагоги и другие эксперты.",
    icon: "/images/user-sound.svg",
  },
  {
    title: "Ты можешь отслеживать статус",
    description:
      "После отправки ты получишь уникальный трек-номер и сможешь проверить, на каком этапе находится обращение.",
    icon: "/images/bell-ringing.svg",
  },
];

const buttonStyle = {
  height: 55,
  padding: "0 35px",
  borderRadius: 15,
  fontSize: 20,
};

const buttonClass =
  "w-full max-w-[250px] whitespace-nowrap font-medium sm:w-auto";

function InfoCard({ title, description, icon }: (typeof cards)[number]) {
  return (
    <article className="flex min-h-[100px] items-center gap-[20px] rounded-[15px] border border-[var(--primary-color)] bg-[var(--text-color-secondary)] p-[20px] 2xl:h-[130px] 2xl:gap-[30px] 2xl:px-[40px]">
      <div className="relative h-[35px] w-[35px] shrink-0 2xl:h-[50px] 2xl:w-[50px]">
        <Image
          src={icon}
          alt=""
          fill
          sizes="(min-width: 1500px) 50px, 35px"
          className="object-contain"
        />
      </div>

      <div className="flex min-w-0 flex-col justify-center gap-[10px] 2xl:gap-[15px]">
        <h2 className="text-[20px] font-medium leading-tight text-[var(--primary-color)] 2xl:text-[25px]">
          {title}
        </h2>
        <p className="text-[15px] leading-snug 2xl:max-w-[695px] 2xl:text-[20px]">
          {description}
        </p>
      </div>
    </article>
  );
}

export default async function Home() {
  return (
    <main className="relative min-h-screen overflow-hidden bg-[var(--background-color)] text-[var(--text-color)] 2xl:h-screen">
      <header className="relative z-20 border border-[var(--primary-color)] bg-[var(--default-color-wborder)]">
        <div className="mx-auto flex min-h-[70px] max-w-[1920px] items-center justify-around px-5 py-[15px] 2xl:h-[90px] 2xl:px-[25px] 2xl:py-[10px]">
          <Image
            src="/images/logo.svg"
            alt="Отклик — ты не один, мы рядом"
            width={180}
            height={55}
            priority
            className="w-[115px] shrink-0 2xl:w-[215px]"
          />
          <p className="hidden whitespace-nowrap text-[20px] font-medium text-[var(--primary-color)] md:block">
            Безопасный способ рассказать о том, что тебя беспокоит, и получить помощь от профессионалов.
          </p>
        </div>
      </header>

      <Image
        src="/images/imageBG.png"
        alt="Девушка с ноутбуком"
        width={500}
        height={770}
        priority
        className="pointer-events-none absolute left-0 top-[10px] hidden w-[605px] select-none 2xl:block"
      />

      <section className="relative z-10 mx-auto max-w-[1920px] px-5 pt-[25px] 2xl:h-[985px] 2xl:px-0 2xl:pt-0">
        <h1 className="text-center font-['Dela_Gothic_One'] text-[25px] leading-tight text-[var(--primary-color)] 2xl:absolute 2xl:left-[715px] 2xl:top-[105px] 2xl:whitespace-nowrap 2xl:text-[45px] 2xl:leading-none font-black">
          Что важно знать?
        </h1>

        <div className="mx-auto mt-[25px] flex max-w-[840px] flex-col gap-[15px] 2xl:absolute 2xl:left-[900px] 2xl:top-[210px] 2xl:mt-0 2xl:gap-[20px]">
          {cards.map((card) => (
            <InfoCard key={card.title} {...card} />
          ))}
        </div>

        <div className="mx-auto mt-[25px] flex max-w-[560px] flex-col items-center gap-[15px] 2xl:absolute 2xl:left-[660px] 2xl:top-[730px] 2xl:mt-0">
          <div className="flex w-full flex-col items-center gap-[15px] sm:flex-row sm:justify-center 2xl:gap-[35px]">
            <Button
              text="Подать обращение"
              fill
              link="/appeal"
              className={buttonClass}
              style={buttonStyle}
            />
            <Button
              text="Проверить статус"
              link="/status"
              className={buttonClass}
              style={buttonStyle}
            />
          </div>
          <span className="text-center text-[15px] opacity-50 2xl:text-[20px]">
            Вход для сотрудников
          </span>
        </div>
      </section>
    </main>
  );
}
