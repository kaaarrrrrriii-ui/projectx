"use client";

import { useCallback, useEffect, useState } from "react";
import Button from "@/shared/ui/button";
import {
  getTicketChat,
  postTicketMessage,
  ticketAttachmentURL,
  type TicketChat,
} from "@/shared/api/public-api";
import ChatWorkspace from "./chat-workspace";
import ReplyModal from "./reply-modal";

const replyStatuses = new Set(["new", "assigned", "in_progress", "needs_clarification", "answer_ready"]);

function formatDate(value: string) {
  return new Intl.DateTimeFormat("ru-RU", { dateStyle: "short", timeStyle: "short" }).format(new Date(value));
}

export default function ChatContent({ trackNumber, formal }: { trackNumber: string; formal: boolean }) {
  const [chat, setChat] = useState<TicketChat | null>(null);
  const [error, setError] = useState("");
  const [replyOpen, setReplyOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const loadChat = useCallback(async (signal?: AbortSignal) => {
    if (!trackNumber) return;
    try {
      setChat(await getTicketChat(trackNumber, signal));
      setError("");
    } catch (reason) {
      if ((reason as Error).name !== "AbortError") {
        setError(reason instanceof Error ? reason.message : "Не удалось загрузить переписку.");
      }
    }
  }, [trackNumber]);

  useEffect(() => {
    const controller = new AbortController();
    if (trackNumber) {
      getTicketChat(trackNumber, controller.signal)
        .then((loaded) => {
          setChat(loaded);
          setError("");
        })
        .catch((reason: unknown) => {
          if ((reason as Error).name !== "AbortError") {
            setError(reason instanceof Error ? reason.message : "Не удалось загрузить переписку.");
          }
        });
    }
    return () => controller.abort();
  }, [trackNumber]);

  async function sendReply(text: string, files: File[]) {
    setSubmitting(true);
    setError("");
    try {
      await postTicketMessage(trackNumber, text, files);
      setReplyOpen(false);
      await loadChat();
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Не удалось отправить сообщение.");
    } finally {
      setSubmitting(false);
    }
  }

  if (!chat) {
    const loadingError = !trackNumber ? "Не указан номер обращения. Вернитесь к проверке статуса." : error;
    return (
      <div className="mx-auto flex w-full max-w-[900px] flex-1 flex-col items-center justify-center px-5 text-center">
        <p className={loadingError ? "text-[#b42318]" : "text-[#646d86]"} role={loadingError ? "alert" : undefined}>
          {loadingError || "Загружаем переписку…"}
        </p>
        {loadingError && <Button text="Проверить статус" link="/status" variant="primary" className="mt-5" />}
      </div>
    );
  }

  const specialistTitle = chat.specialist?.label || chat.specialist?.expert_group || "Специалист";
  return (
    <ChatWorkspace formal={formal} trackNumber={chat.track_id} status={chat.status}>
      <div className="space-y-4">
        {chat.messages.length === 0 && (
          <div className="rounded-[15px] border border-[#dee7fd] bg-white px-6 py-7 text-[#646d86]">
            Сообщений пока нет. Ответ специалиста появится здесь.
          </div>
        )}
        {chat.messages.map((message) => {
          const specialist = message.type === "specialist";
          return (
            <article
              key={message.id}
              aria-label={specialist ? "Сообщение специалиста" : "Ваше сообщение"}
              className={`overflow-hidden rounded-[15px] border ${specialist ? "border-[#dee7fd] bg-[#dee7fd]" : "ml-auto max-w-[85%] border-[#4562f0] bg-white"}`}
            >
              <header className="flex items-center justify-between gap-4 bg-white px-5 py-3">
                <h2 className="font-semibold text-[var(--color-primary)]">{specialist ? specialistTitle : "Вы"}</h2>
                <time className="text-xs text-[#777d91]" dateTime={message.created_at}>{formatDate(message.created_at)}</time>
              </header>
              <div className="whitespace-pre-wrap px-5 py-4 text-[15px] leading-6 text-[#11131a]">{message.text}</div>
            </article>
          );
        })}

        {chat.attachments.length > 0 && (
          <section className="rounded-[15px] border border-[#dee7fd] bg-white px-5 py-4" aria-labelledby="chat-attachments-heading">
            <h2 id="chat-attachments-heading" className="font-semibold text-[var(--color-primary)]">Вложения обращения</h2>
            <ul className="mt-2 space-y-1 text-sm">
              {chat.attachments.map((attachment) => (
                <li key={attachment.id}>
                  <a className="text-[#4562f0] underline underline-offset-2" href={ticketAttachmentURL(chat.track_id, attachment.id)} target="_blank" rel="noreferrer">
                    {attachment.name}
                  </a>
                </li>
              ))}
            </ul>
          </section>
        )}

        {replyStatuses.has(chat.status) && (
          <div className="flex justify-end">
            <Button text="Добавить сообщение" variant="secondary" onClick={() => setReplyOpen(true)} />
          </div>
        )}
        {error && <p role="alert" className="text-sm text-[#b42318]">{error}</p>}
      </div>

      {replyOpen && <ReplyModal formal={formal} submitting={submitting} onClose={() => setReplyOpen(false)} onSend={(text, files) => void sendReply(text, files)} />}
    </ChatWorkspace>
  );
}
