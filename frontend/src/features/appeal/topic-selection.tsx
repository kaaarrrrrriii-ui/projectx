"use client";

import { useState } from "react";
import Button from "@/shared/ui/button";
import Input from "@/shared/ui/input";
import AppealNavigation from "@/shared/ui/appeal-navigation";
import { appealRoutes, getAppealRoute } from "./routes";
import TopicOption from "./topic-option";
import { appealCategories as topics } from "./categories";

function ActionIcon({ name }: { name: "edit" | "delete" }) {
  if (name === "delete") {
    return (
      <svg viewBox="0 0 20 20" aria-hidden="true" className="h-4 w-4" fill="none">
        <path d="M3.5 5.5h13M8 3h4l1 2.5H7L8 3Zm-2 2.5.7 11h6.6l.7-11M8.5 8v5.5m3-5.5v5.5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    );
  }

  return (
    <svg viewBox="0 0 20 20" aria-hidden="true" className="h-4 w-4" fill="none">
      <path d="m4 13.8-.7 3 3-.7L15.4 7 13 4.6 4 13.8Z" stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round" />
      <path d="m11.8 5.8 2.4 2.4" stroke="currentColor" strokeWidth="1.5" />
    </svg>
  );
}

export default function TopicSelection({ roleId, initialTopic = "", formal = false }: { roleId: string; initialTopic?: string; formal?: boolean }) {
  const initialCustomTopic = initialTopic && !topics.includes(initialTopic as (typeof topics)[number]) ? initialTopic : "";
  const [selectedTopic, setSelectedTopic] = useState(initialTopic);
  const [customTopic, setCustomTopic] = useState(initialCustomTopic);
  const [draft, setDraft] = useState(initialCustomTopic);
  const [isEditing, setIsEditing] = useState(false);
  const [error, setError] = useState("");

  function openEditor() {
    setDraft(customTopic);
    setError("");
    setIsEditing(true);
  }

  function saveCustomTopic() {
    const nextTopic = draft.trim().replace(/\s+/g, " ");
    if (!nextTopic) {
      setError(formal ? "Напишите, как вы хотите назвать тему" : "Напиши, как ты хочешь назвать тему");
      return;
    }

    setCustomTopic(nextTopic);
    setSelectedTopic(nextTopic);
    setError("");
    setIsEditing(false);
  }

  function deleteCustomTopic() {
    if (selectedTopic === customTopic) setSelectedTopic("");
    setCustomTopic("");
    setDraft("");
    setError("");
    setIsEditing(false);
  }

  return (
    <form action={appealRoutes.description} method="get" className="flex flex-1 flex-col">
      <input type="hidden" name="role" value={roleId} />

      <div className="rounded-[15px] border border-[var(--color-primary)] bg-[var(--color-background)] px-[30px] py-[50px]">
        <h1 id="topics-heading" className="text-[28px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[#4562f0]">
          С чем это связано?
        </h1>
        <p id="topics-description" className="mt-[5px] text-[12px] leading-4 text-[#151515]">
          {formal ? "Можно выбрать одну или несколько тем, которые ближе всего к вашей ситуации." : "Можно выбрать одну или несколько тем, которые ближе всего к твоей ситуации."}
        </p>

        <fieldset aria-describedby="topics-description" className="mt-10 min-w-0 sm:mt-[62px]">
        <legend className="sr-only">Выбор одной темы обращения</legend>
        <div className="grid grid-cols-1 gap-x-3 gap-y-4 sm:grid-cols-2 lg:grid-cols-3">
          {topics.map((topic) => (
            <TopicOption
              key={topic}
              value={topic}
              label={topic}
              checked={selectedTopic === topic}
              onChange={setSelectedTopic}
              required
              variant={topic === "не знаю, как это назвать" ? "secondary" : "default"}
            />
          ))}

          {customTopic && !isEditing && (
            <TopicOption
              value={customTopic}
              label={customTopic}
              checked={selectedTopic === customTopic}
              onChange={setSelectedTopic}
              required
            >
              <span className="absolute inset-y-0 right-2 flex items-center gap-1">
                <button type="button" onClick={openEditor} aria-label="Редактировать свою тему" className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-[#4562f0] hover:bg-white/70 focus-visible:outline-2 focus-visible:outline-[#4562f0]">
                  <ActionIcon name="edit" />
                </button>
                <button type="button" onClick={deleteCustomTopic} aria-label="Удалить свою тему" className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-[#d92d20] hover:bg-white/70 focus-visible:outline-2 focus-visible:outline-[#d92d20]">
                  <ActionIcon name="delete" />
                </button>
              </span>
            </TopicOption>
          )}
        </div>
        </fieldset>

        {isEditing ? (
          <div className="mt-6 max-w-[620px] rounded-2xl border border-[#cbd3f5] bg-white/70 p-4">
          <label htmlFor="custom-topic" className="mb-2 block text-sm font-medium text-[#000828]">
            {customTopic
              ? (formal ? "Измените свою тему" : "Измени свою тему")
              : (formal ? "Добавьте свою тему" : "Добавь свою тему")}
          </label>
          <div className="flex flex-col gap-3 sm:flex-row">
            <Input
              id="custom-topic"
              type="text"
              value={draft}
              onChange={(event) => { setDraft(event.target.value); setError(""); }}
              onKeyDown={(event) => {
                if (event.key === "Enter") {
                  event.preventDefault();
                  saveCustomTopic();
                }
              }}
              maxLength={120}
              autoFocus
              placeholder="Например: трудно адаптироваться"
              aria-invalid={Boolean(error)}
              aria-describedby={error ? "custom-topic-error" : undefined}
              className="flex-1 text-left"
            />
            <div className="flex gap-2">
              <Button text="Сохранить" variant="primary" size="small" onClick={saveCustomTopic} />
              <Button text="Отмена" variant="secondary" size="small" onClick={() => { setIsEditing(false); setError(""); }} />
            </div>
          </div>
          {error && <p id="custom-topic-error" role="alert" className="mt-2 text-sm text-[#b42318]">{error}</p>}
          </div>
        ) : !customTopic ? (
          <button type="button" onClick={openEditor} className="mt-4 flex w-fit cursor-pointer items-center gap-3 rounded-sm py-1 text-[12px] text-[#151515] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#4562f0]">
            <span aria-hidden="true" className="text-xl leading-4 font-extralight">+</span>
            {formal ? "Добавить свою тему" : "Добавить свою тему"}
          </button>
        ) : null}
      </div>

      <AppealNavigation
        backHref={getAppealRoute("role", { role: roleId })}
        skipHref={getAppealRoute("description", { role: roleId })}
        primaryText="Продолжить"
        primaryType="submit"
      />
    </form>
  );
}
