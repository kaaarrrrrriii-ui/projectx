"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";
import Button from "@/shared/ui/button";
import AppealNavigation from "@/shared/ui/appeal-navigation";
import { getPublicCategories, type PublicCategory } from "@/shared/api/public-api";
import { appealRoutes, getAppealRoute } from "./routes";
import TopicOption from "./topic-option";
import { readAppealDraft, startAppealDraft, updateAppealDraft } from "./appeal-draft";

function ActionIcon({ name }: { name: "edit" | "delete" }) {
  if (name === "delete") {
    return <svg viewBox="0 0 20 20" aria-hidden="true" className="h-4 w-4" fill="none"><path d="M3.5 5.5h13M8 3h4l1 2.5H7L8 3Zm-2 2.5.7 11h6.6l.7-11M8.5 8v5.5m3-5.5v5.5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" /></svg>;
  }
  return <svg viewBox="0 0 20 20" aria-hidden="true" className="h-4 w-4" fill="none"><path d="m4 13.8-.7 3 3-.7L15.4 7 13 4.6 4 13.8Z" stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round" /><path d="m11.8 5.8 2.4 2.4" stroke="currentColor" strokeWidth="1.5" /></svg>;
}

function isUnknownCategory(category: PublicCategory) {
  return category.name.toLocaleLowerCase("ru").startsWith("не знаю");
}

export default function TopicSelection({ roleId, initialTopic = "", formal = false }: { roleId: string; initialTopic?: string; formal?: boolean }) {
  const [categories, setCategories] = useState<PublicCategory[]>([]);
  const [selectedChoice, setSelectedChoice] = useState(initialTopic);
  const [customTopic, setCustomTopic] = useState("");
  const [draft, setDraft] = useState("");
  const [isEditing, setIsEditing] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const unknownCategory = useMemo(() => categories.find(isUnknownCategory), [categories]);

  useEffect(() => {
    startAppealDraft(roleId);
    const saved = readAppealDraft();
    const controller = new AbortController();
    getPublicCategories(controller.signal)
      .then(({ categories: loaded }) => {
        setCustomTopic(saved.customTopic ?? "");
        setDraft(saved.customTopic ?? "");
        setCategories(loaded);
        const category = loaded.find((item) => item.id === saved.categoryId);
        if (saved.customTopic && category && isUnknownCategory(category)) setSelectedChoice("custom");
        else if (category) setSelectedChoice(`category:${category.id}`);
        else if (initialTopic.startsWith("category:")) setSelectedChoice(initialTopic);
        setError(loaded.length ? "" : "Администратор пока не добавил категории.");
      })
      .catch((reason: unknown) => {
        if ((reason as Error).name !== "AbortError") setError("Не удалось загрузить категории. Попробуйте обновить страницу.");
      })
      .finally(() => setLoading(false));
    return () => controller.abort();
  }, [initialTopic, roleId]);

  function selectCategory(value: string) {
    setSelectedChoice(value);
    setError("");
    if (value !== "custom") setCustomTopic("");
  }

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
    if (!unknownCategory) {
      setError("Категория «Не знаю, как это назвать» пока недоступна.");
      return;
    }
    setCustomTopic(nextTopic);
    setSelectedChoice("custom");
    setError("");
    setIsEditing(false);
  }

  function deleteCustomTopic() {
    setCustomTopic("");
    setDraft("");
    setSelectedChoice("");
    setError("");
    setIsEditing(false);
  }

  function selectedCategory() {
    if (selectedChoice === "custom") return unknownCategory;
    const id = Number(selectedChoice.replace("category:", ""));
    return categories.find((category) => category.id === id);
  }

  function persistSelection() {
    const category = selectedCategory() ?? unknownCategory;
    if (!category) return false;
    updateAppealDraft({ role: roleId, categoryId: category.id, categoryName: category.name, customTopic: selectedChoice === "custom" ? customTopic : "", answers: {} });
    return true;
  }

  function submit(event: FormEvent<HTMLFormElement>) {
    if (!persistSelection()) {
      event.preventDefault();
      setError("Выберите тему обращения.");
    }
  }

  return (
    <form action={appealRoutes.description} method="get" onSubmit={submit} className="flex flex-1 flex-col">
      <input type="hidden" name="role" value={roleId} />

      <div className="rounded-[15px] border border-[var(--color-primary)] bg-[var(--color-background)] px-[30px] py-[50px] max-[699px]:px-5 max-[699px]:py-7 max-[379px]:px-4">
        <h1 id="topics-heading" className="text-[28px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[#4562f0] max-[379px]:text-[25px]">
          С чем это связано?
        </h1>
        <p id="topics-description" className="mt-[5px] text-[12px] leading-4 text-[#151515]">
          {formal ? "Можно выбрать одну или несколько тем, которые ближе всего к вашей ситуации." : "Можно выбрать одну или несколько тем, которые ближе всего к твоей ситуации."}
        </p>

        <fieldset aria-describedby="topics-description" className="mt-8 min-w-0">
          <legend className="sr-only">Выбор одной темы обращения</legend>
          <div className="grid grid-cols-1 gap-x-3 gap-y-4 sm:grid-cols-2 lg:grid-cols-3">
            {categories.map((category) => <TopicOption key={category.id} value={`category:${category.id}`} label={category.name} checked={selectedChoice === `category:${category.id}`} onChange={selectCategory} required variant={isUnknownCategory(category) ? "secondary" : "default"} />)}
            {customTopic && !isEditing && (
              <TopicOption value="custom" label={customTopic} checked={selectedChoice === "custom"} onChange={selectCategory} required>
                <span className="absolute inset-y-0 right-2 flex items-center gap-1">
                  <button type="button" onClick={openEditor} aria-label="Редактировать свою тему" className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-[#4562f0] hover:bg-white/70 focus-visible:outline-2 focus-visible:outline-[#4562f0]"><ActionIcon name="edit" /></button>
                  <button type="button" onClick={deleteCustomTopic} aria-label="Удалить свою тему" className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-[#d92d20] hover:bg-white/70 focus-visible:outline-2 focus-visible:outline-[#d92d20]"><ActionIcon name="delete" /></button>
                </span>
              </TopicOption>
            )}
          </div>
        </fieldset>
        {loading && <p className="mt-5 text-sm text-[#646d86]">Загружаем категории…</p>}
        {isEditing ? (
          <div className="mt-8 w-fit max-w-full rounded-2xl border border-[#808393] bg-transparent p-4">
            <label htmlFor="custom-topic" className="mb-2 block text-[12px] leading-4 font-medium text-[#808393]">{customTopic ? "Измените свою тему" : "Добавьте свою тему"}</label>
            <div className="flex flex-col gap-3 sm:flex-row">
              <input id="custom-topic" type="text" value={draft} onChange={(event) => { setDraft(event.target.value); setError(""); }} onKeyDown={(event) => { if (event.key === "Enter") { event.preventDefault(); saveCustomTopic(); } }} maxLength={120} autoFocus className="h-7 w-[279px] max-w-full rounded-[15px] border border-[#808393] bg-transparent px-[10px] text-[14px] leading-5 text-[#000828] outline-none" />
              <div className="flex gap-1.5"><Button text="Сохранить" variant="primary" size="small" className="!h-6 !min-h-0 !rounded-md !px-2 !py-0.5 text-[10px]" onClick={saveCustomTopic} /><Button text="Отмена" variant="secondary" size="small" className="!h-6 !min-h-0 !rounded-md !px-2 !py-0.5 text-[10px]" onClick={() => { setIsEditing(false); setError(""); }} /></div>
            </div>
          </div>
        ) : !customTopic ? (
          <button type="button" onClick={openEditor} className="mt-8 flex w-fit cursor-pointer items-center gap-1.5 rounded-sm py-1 text-[10px] leading-4 text-[#808393] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#4562f0]"><span aria-hidden="true">+</span><span>Добавить свою тему</span></button>
        ) : null}
        {error && <p role="alert" className="mt-3 text-sm text-[#b42318]">{error}</p>}
      </div>
      <AppealNavigation backHref={getAppealRoute("role", { role: roleId })} skipHref={getAppealRoute("description", { role: roleId })} onSkip={persistSelection} primaryText="Продолжить" primaryType="submit" primaryDisabled={loading || categories.length === 0} />
    </form>
  );
}
