import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import ChatWorkspace from "@/features/chat/chat-workspace";
import SpecialistMessage from "@/features/chat/specialist-message";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function ChatPage({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);
  const formal = isFormalAppealRole(role.id);
  const copy = formal
    ? {
        intro:
          "Спасибо, что поделились этой историей. Поверьте, то, что вы чувствуете — страх, злость, растерянность или даже стыд — абсолютно нормально в такой ситуации. Кибербуллинг ранит не меньше, чем оскорбления вживую, а иногда даже больнее, потому что кажется, что от него не скрыться.",
        reassurance:
          "Я внимательно прочитал ваш рассказ. Хочу, чтобы вы знали: вы не виноваты в том, что происходит. Никто не имеет права унижать вас, угрожать или высмеивать, даже через экран. Это не «шутки» и не «просто слова» — это агрессия, и вы имеете право защищать себя.",
      }
    : {
        intro:
          "Спасибо, что делишься этой историей. Поверь, то, что ты чувствуешь — страх, злость, растерянность или даже стыд — абсолютно нормально в такой ситуации. Кибербуллинг ранит не меньше, чем оскорбления вживую, а иногда даже больнее, потому что кажется, что от него не скрыться.",
        reassurance:
          "Я внимательно прочитал твой рассказ. Помни: ты не несёшь ответственности за то, что происходит. Никто не имеет права унижать тебя, угрожать или высмеивать, даже через экран. Это не «шутки» и не «просто слова» — это агрессия, и ты имеешь право защищать себя.",
      };

  return (
    <main className="app-page-background flex min-h-dvh min-w-[320px] flex-col overflow-x-clip text-[var(--color-text)]">
      <SiteHeader formal={formal} />

      <section
        aria-labelledby="specialist-answer-heading"
        className="flex flex-1"
      >
        <h1 id="specialist-answer-heading" className="sr-only">
          Ответ специалиста
        </h1>

        <ChatWorkspace formal={formal}>
          <SpecialistMessage>
            <p>Здравствуйте.</p>
            <p>{copy.intro}</p>
            <p>{copy.reassurance}</p>
          </SpecialistMessage>
        </ChatWorkspace>
      </section>
    </main>
  );
}
