import AppealNavigation from "@/shared/ui/appeal-navigation";
import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import SiteHeader from "@/widgets/site-header/site-header";
import Image from "next/image";

export default async function AppealDescription({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[]; topic?: string | string[] }>;
}) {
  const params = await searchParams;
  const role = getAppealRole(params.role);
  const formal = isFormalAppealRole(role.id);
  const topic = typeof params.topic === "string" ? params.topic : "";
  const topicQuery = topic ? `&topic=${encodeURIComponent(topic)}` : "";

  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader formal={formal} />

      <section
        aria-labelledby="description-heading"
        className="
          flex flex-1 flex-col
          px-[4.88%] pt-14 pb-[38px]

          max-[699px]:px-5
          max-[699px]:py-8
        "
      >
        <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col">
          <h1
            id="description-heading"
            className="
              text-[36px]
              leading-[1.2]
              font-black
              tracking-[-0.025em]
              text-[#4562F0]

              max-[699px]:text-[28px]
            "
          >
            {formal ? "Расскажите, что происходит" : "Расскажи, что происходит"}
          </h1>

          <p
            id="description-intro"
            className="
              mt-1
              text-[15px]
              leading-[22px]
              text-[#151515]

              max-[699px]:mt-2.5
              max-[699px]:text-[14px]
            "
          >
            {formal
              ? "Опишите ситуацию своими словами. Чем больше деталей, тем проще нам будет помочь. Если не знаете, с чего начать — просто напишите, что чувствуете."
              : "Опиши ситуацию своими словами. Чем больше деталей, тем проще нам будет помочь. Если не знаешь, с чего начать — просто напиши, что чувствуешь."}
          </p>

          <textarea
            name="description"
            aria-labelledby="description-heading"
            aria-describedby="description-intro description-hint"
            placeholder={
              formal
                ? "Здесь можно написать всё, что вас беспокоит..."
                : "Здесь можно написать всё, что тебя беспокоит..."
            }
            className="
              mt-7 block
              min-h-[260px] w-full
              resize-y
              rounded-[15px]
              border border-[#333]
              bg-[#FCFDFF]
              px-[13px] py-4

              text-[15px]
              leading-6
              text-[#151515]

              placeholder:text-[#85899f]
              placeholder:opacity-100

              focus-visible:outline-2
              focus-visible:outline-offset-4
              focus-visible:outline-[#4562F0]

              max-[699px]:mt-6
              max-[699px]:text-[16px]
            "
          />

          <aside
            id="description-hint"
            className="
              mt-[38px]
              flex min-h-[49px]
              items-center gap-2.5
              rounded-[15px]
              border border-[#677aff]
              bg-[#dfe6ff]
              px-2.5 py-[9px]

              text-[15px]
              leading-[22px]
              text-[#151515]

              max-[699px]:mt-6
              max-[699px]:items-start
              max-[699px]:text-[14px]
            "
          >
            <Image
              src="/images/LightbulbFilament.svg"
              alt=""
              width={28}
              height={28}
              className="shrink-0"
            />

            <p>
              {formal
                ? "Нет правильных или неправильных слов. Пишите так, как вам удобно. Мы внимательно читаем каждое обращение."
                : "Нет правильных или неправильных слов. Пиши так, как тебе удобно. Мы внимательно читаем каждое обращение."}
            </p>
          </aside>

          <AppealNavigation
            backHref={`/appeal/appeal2?role=${role.id}${topicQuery}`}
            skipHref={`/appeal/appeal2/appeal3/appeal4?role=${role.id}${topicQuery}`}
            primaryText="Продолжить"
            primaryHref={`/appeal/appeal2/appeal3/appeal4?role=${role.id}${topicQuery}`}
          />
        </div>
      </section>
    </main>
  );
}
