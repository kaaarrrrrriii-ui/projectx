"use client";

import Button from "@/shared/ui/button";
import Link from "next/link";
import type { ReactNode } from "react";
import { useState } from "react";
import CloseAppealFlow from "./close-appeal-flow";
import ReplyModal from "./reply-modal";
import SpecialistCard from "./specialist-card";

type SentMessage = {
  id: number;
  text: string;
  fileNames: string[];
};

export default function ChatWorkspace({ children }: { children: ReactNode }) {
  const [activeDialog, setActiveDialog] = useState<"reply" | "close" | null>(null);
  const [messages, setMessages] = useState<SentMessage[]>([]);

  function addMessage(text: string, fileNames: string[] = []) {
    setMessages((current) => [
      ...current,
      { id: Date.now(), text, fileNames },
    ]);
    setActiveDialog(null);
  }

  return (
    <>
      <div className="mx-auto grid w-full max-w-[1600px] grid-cols-[268px_minmax(0,1fr)] max-[799px]:grid-cols-1">
        <aside
          className="flex min-h-[calc(100dvh-64px)] flex-col border-r border-[#dee7fd]  px-[14px] py-4 max-[799px]:min-h-0 max-[799px]:border-r-0 max-[799px]:border-b max-[799px]:px-5 max-[799px]:py-5"
          aria-label="Информация о специалисте"
        >
          <SpecialistCard
            avatarSrc="/images/roles/teacher.webp"
            name="Анна Петрова"
            profession="Психолог"
            experience="Более 8 лет"
            specialization="КПТ, тревожные расстройства"
          />

          <div className="mt-[18px] grid gap-2.5 max-[799px]:mx-auto max-[799px]:w-full max-[799px]:max-w-[420px] max-[479px]:grid-cols-1 min-[480px]:max-[799px]:grid-cols-2">
            <Button
              text="Закрыть обращение"
              variant="primary"
              size="small"
              onClick={() => setActiveDialog("close")}
              className="w-full"
            />
            <Button
              text="Добавить свой ответ"
              variant="secondary"
              size="small"
              onClick={() => setActiveDialog("reply")}
              className="w-full"
            />
          </div>

          <Link
            href="/status"
            className="mt-auto w-fit rounded-sm pt-8 text-[13px] text-[#9196a7] transition-colors hover:text-[var(--color-primary)] focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[var(--color-primary)] max-[799px]:mx-auto max-[799px]:pt-5"
          >
            Вернуться назад
          </Link>
        </aside>

        <div className="min-w-0 px-6 py-7 max-[799px]:px-5 max-[799px]:py-6">
          <div className="space-y-4" aria-live="polite">
            {children}

            {messages.map((message) => (
              <article
                key={message.id}
                aria-label="Ваше сообщение"
                className="ml-auto max-w-[72%] rounded-[15px] bg-[var(--color-primary)] px-5 py-3 text-[15px] leading-[1.45] text-white max-[599px]:max-w-[88%] max-[499px]:px-4 max-[499px]:text-[14px]"
              >
                {message.text && <p className="whitespace-pre-wrap">{message.text}</p>}
                {message.fileNames.length > 0 && (
                  <ul className={`${message.text ? "mt-2" : ""} grid gap-1 text-[12px] text-white/85`}>
                    {message.fileNames.map((fileName) => (
                      <li key={fileName} className="truncate">📎 {fileName}</li>
                    ))}
                  </ul>
                )}
              </article>
            ))}
          </div>
        </div>
      </div>

      {activeDialog === "reply" && (
        <ReplyModal
          onClose={() => setActiveDialog(null)}
          onSend={addMessage}
        />
      )}

      {activeDialog === "close" && (
        <CloseAppealFlow
          trackNumber="НАШК-УАЫВ-АВАМ-ВАФВ"
          onCancel={() => setActiveDialog(null)}
        />
      )}
    </>
  );
}
