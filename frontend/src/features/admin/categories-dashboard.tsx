"use client";

import {
  clarifyingQuestions as initialQuestions,
  type ClarifyingQuestion,
} from "@/features/appeal/clarifying-question-data";
import DialogShell from "@/features/chat/dialog-shell";
import Button from "@/shared/ui/button";
import Input from "@/shared/ui/input";
import Link from "next/link";
import { useMemo, useState, type FormEvent } from "react";
import {
  adminCategories,
  specialistOptions,
  type SpecialistOption,
} from "./admin-options";

type CategoryConfig = {
  id: number;
  category: string;
  specialist: SpecialistOption;
  questionIds: string[];
};

const controlClass = "h-[38px] w-full rounded-[9px] border border-[#000828] bg-white px-3 text-sm text-[#000828] outline-none transition-colors hover:border-[#4562f0] focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/20";

function defaultSpecialist(category: string): SpecialistOption {
  if (category.includes("юридическ")) return "Юрист";
  if (category.includes("конфликт")) return "Конфликтолог";
  if (category.includes("давление")) return "Социальный педагог";
  return "Психолог";
}

const initialConfigs: CategoryConfig[] = adminCategories.slice(0, 4).map((category, index) => ({
  id: index + 1,
  category,
  specialist: defaultSpecialist(category),
  questionIds: initialQuestions.map((question) => question.id),
}));

function QuestionSelectionModal({
  category,
  configs,
  questions,
  onCategoryChange,
  onToggle,
  onEdit,
  onAddQuestion,
  onClose,
}: {
  category: string;
  configs: CategoryConfig[];
  questions: ClarifyingQuestion[];
  onCategoryChange: (category: string) => void;
  onToggle: (category: string, questionId: string) => void;
  onEdit: (questionId: string) => void;
  onAddQuestion: () => void;
  onClose: () => void;
}) {
  const selectedIds = configs.find((item) => item.category === category)?.questionIds ?? [];

  return (
    <DialogShell labelledBy="questions-modal-heading" onClose={onClose} showClose className="max-w-[720px]">
      <h2 id="questions-modal-heading" className="pr-12 text-[27px] leading-tight font-extrabold text-[#4562f0] max-[499px]:text-[22px]">Выбор уточняющих вопросов</h2>

      <label className="mt-6 grid gap-2 text-base font-medium text-[#151515]">
        Категория
        <select value={category} onChange={(event) => onCategoryChange(event.target.value)} className={controlClass}>
          {adminCategories.map((item) => <option key={item} value={item}>{item.charAt(0).toUpperCase() + item.slice(1)}</option>)}
        </select>
      </label>

      <section aria-labelledby="questions-list-heading" className="mt-5">
        <h3 id="questions-list-heading" className="text-base font-medium text-[#151515]">Список вопросов</h3>
        <div className="mt-2 overflow-hidden rounded-[12px] border border-[#000828] bg-[#fcfdff]">
          <button type="button" onClick={onAddQuestion} className="flex min-h-11 w-full cursor-pointer items-center gap-3 border-b border-[#bdc2cf] px-4 text-left text-sm text-[#30384f] transition-colors hover:bg-[#eef1ff] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-inset focus-visible:outline-[#4562f0]">
            <span className="text-xl leading-none">+</span>Добавить новый вопрос
          </button>

          <ul className="max-h-[330px] overflow-y-auto">
            {questions.map((question) => (
              <li key={question.id} className="flex min-h-[58px] items-center gap-3 border-b border-[#e1e4ec] px-4 last:border-b-0">
                <label className="flex min-w-0 flex-1 cursor-pointer items-center gap-3 text-sm text-[#30384f]">
                  <input type="checkbox" checked={selectedIds.includes(question.id)} onChange={() => onToggle(category, question.id)} className="h-4 w-4 shrink-0 accent-[#4562f0]" />
                  <span className="truncate">{question.title}</span>
                </label>
                <button type="button" onClick={() => onEdit(question.id)} className="shrink-0 cursor-pointer rounded-[8px] bg-[#4562f0] px-4 py-2 text-xs text-white transition-colors hover:bg-[#4f71fc] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]">Редактировать вопрос</button>
              </li>
            ))}
          </ul>
        </div>
      </section>

      <Button text={`Готово · выбрано ${selectedIds.length}`} variant="primary" onClick={onClose} className="mt-6 w-full" />
    </DialogShell>
  );
}

function EditQuestionModal({
  question,
  onSave,
  onClose,
}: {
  question?: ClarifyingQuestion;
  onSave: (question: Omit<ClarifyingQuestion, "id">) => void;
  onClose: () => void;
}) {
  const [title, setTitle] = useState(question?.title ?? "");
  const [options, setOptions] = useState(question?.options.length ? question.options : [""]);
  const [error, setError] = useState("");

  function updateOption(index: number, value: string) {
    setOptions((current) => current.map((item, itemIndex) => itemIndex === index ? value : item));
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const normalizedTitle = title.trim();
    const normalizedOptions = options.map((option) => option.trim()).filter(Boolean);
    if (!normalizedTitle) {
      setError("Введите текст вопроса.");
      return;
    }
    if (normalizedOptions.length === 0) {
      setError("Добавьте хотя бы один вариант ответа.");
      return;
    }
    onSave({ title: normalizedTitle, options: normalizedOptions, key: question?.key });
  }

  return (
    <DialogShell labelledBy="edit-question-heading" onClose={onClose} showClose className="max-w-[670px]">
      <form onSubmit={handleSubmit} noValidate>
        <h2 id="edit-question-heading" className="pr-12 text-[27px] leading-tight font-extrabold text-[#4562f0] max-[499px]:text-[22px]">{question ? "Редактировать вопрос и ответы" : "Новый уточняющий вопрос"}</h2>

        <label className="mt-6 grid gap-2 text-base font-medium text-[#151515]">
          Вопрос
          <Input value={title} onChange={(event) => setTitle(event.target.value)} required autoFocus className="!h-[42px] !w-full !rounded-[12px] !bg-white !px-3 !text-left !text-sm focus:!bg-white focus:!text-[#000828]" />
        </label>

        <section aria-labelledby="answers-heading" className="mt-5">
          <div className="flex items-center justify-between gap-4">
            <h3 id="answers-heading" className="text-base font-medium text-[#151515]">Варианты ответа</h3>
            <button type="button" onClick={() => setOptions((current) => [...current, ""])} className="cursor-pointer rounded-md px-2 py-1 text-sm text-[#4562f0] hover:bg-[#eef1ff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]">+ Добавить ответ</button>
          </div>

          <div className="mt-2 grid max-h-[330px] gap-3 overflow-y-auto rounded-[12px] border border-[#000828] bg-[#fcfdff] p-3">
            {options.map((option, index) => (
              <div key={index} className="flex items-center gap-2">
                <Input value={option} onChange={(event) => updateOption(index, event.target.value)} aria-label={`Вариант ответа ${index + 1}`} placeholder={`Вариант ${index + 1}`} className="!h-[38px] !w-full !rounded-[9px] !bg-white !px-3 !text-left !text-sm focus:!bg-white focus:!text-[#000828]" />
                <button type="button" onClick={() => setOptions((current) => current.filter((_, itemIndex) => itemIndex !== index))} disabled={options.length === 1} aria-label={`Удалить вариант ответа ${index + 1}`} className="flex h-9 w-9 shrink-0 cursor-pointer items-center justify-center rounded-lg text-xl text-[#646d86] hover:bg-[#fff1f1] hover:text-[#d70d14] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] disabled:pointer-events-none disabled:opacity-30">×</button>
              </div>
            ))}
          </div>
        </section>

        {error && <p role="alert" className="mt-4 text-sm text-[#d70d14]">{error}</p>}
        <Button text={question ? "Сохранить вопрос" : "Добавить вопрос"} type="submit" variant="primary" className="mt-6 w-full" />
      </form>
    </DialogShell>
  );
}

export default function CategoriesDashboard() {
  const [configs, setConfigs] = useState<CategoryConfig[]>(initialConfigs);
  const [questions, setQuestions] = useState<ClarifyingQuestion[]>(initialQuestions);
  const [questionModalCategory, setQuestionModalCategory] = useState<string | null>(null);
  const [editingQuestionId, setEditingQuestionId] = useState<string | null | undefined>(undefined);
  const [saved, setSaved] = useState(false);

  const usedCategories = useMemo(() => new Set(configs.map((item) => item.category)), [configs]);
  const nextCategory = adminCategories.find((category) => !usedCategories.has(category));
  const editingQuestion = questions.find((question) => question.id === editingQuestionId);

  function updateConfig(id: number, patch: Partial<CategoryConfig>) {
    setConfigs((current) => current.map((item) => item.id === id ? { ...item, ...patch } : item));
    setSaved(false);
  }

  function ensureCategoryConfig(category: string) {
    setConfigs((current) => current.some((item) => item.category === category) ? current : [...current, { id: Math.max(0, ...current.map((item) => item.id)) + 1, category, specialist: defaultSpecialist(category), questionIds: [] }]);
  }

  function toggleQuestion(category: string, questionId: string) {
    setConfigs((current) => {
      const existing = current.find((item) => item.category === category);
      if (!existing) return [...current, { id: Math.max(0, ...current.map((item) => item.id)) + 1, category, specialist: defaultSpecialist(category), questionIds: [questionId] }];
      return current.map((item) => item.id === existing.id ? { ...item, questionIds: item.questionIds.includes(questionId) ? item.questionIds.filter((id) => id !== questionId) : [...item.questionIds, questionId] } : item);
    });
    setSaved(false);
  }

  function saveQuestion(question: Omit<ClarifyingQuestion, "id">) {
    if (editingQuestionId) {
      setQuestions((current) => current.map((item) => item.id === editingQuestionId ? { ...item, ...question } : item));
    } else if (questionModalCategory) {
      const id = `custom-${Date.now()}`;
      setQuestions((current) => [...current, { ...question, id }]);
      setConfigs((current) => current.map((item) => item.category === questionModalCategory ? { ...item, questionIds: [...item.questionIds, id] } : item));
    }
    setEditingQuestionId(undefined);
    setSaved(false);
  }

  return (
    <section className="min-w-0 flex-1 px-5 pt-5 pb-8 sm:px-[26px]" aria-labelledby="categories-heading">
      <div className="w-full">
        <Link href="/admin" className="inline-flex rounded-sm text-sm text-[#85899b] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[#4562f0]">Вернуться назад</Link>
        <h1 id="categories-heading" className="mt-7 text-[26px] leading-tight font-extrabold tracking-[-0.02em] text-[#4562f0]">Категории</h1>

        <section aria-label="Настройка категорий" className="mt-7 overflow-x-auto rounded-[13px] border border-[#4562f0] bg-white/85">
          <table className="w-full min-w-[760px] table-fixed border-collapse">
            <thead className="bg-[#dfe6ff] text-[#4562f0]"><tr className="h-[52px]"><th scope="col" className="w-[34%] border-r border-[#4562f0] px-4 text-center font-normal">Категория</th><th scope="col" className="w-[32%] border-r border-[#4562f0] px-4 text-center font-normal">Специализация</th><th scope="col" className="w-[34%] px-4 text-center font-normal">Список уточняющих вопросов</th></tr></thead>
            <tbody>
              {configs.map((config) => (
                <tr key={config.id} className="h-[62px] border-t border-[#4562f0] hover:bg-[#f7f8ff]">
                  <td className="border-r border-[#4562f0] px-3">
                    <select value={config.category} onChange={(event) => updateConfig(config.id, { category: event.target.value, specialist: defaultSpecialist(event.target.value) })} aria-label={`Категория строки ${config.id}`} className={controlClass}>
                      {adminCategories.map((category) => <option key={category} value={category} disabled={category !== config.category && usedCategories.has(category)}>{category.charAt(0).toUpperCase() + category.slice(1)}</option>)}
                    </select>
                  </td>
                  <td className="border-r border-[#4562f0] px-3">
                    <select value={config.specialist} onChange={(event) => updateConfig(config.id, { specialist: event.target.value as SpecialistOption })} aria-label={`Специализация для категории ${config.category}`} className={controlClass}>
                      {specialistOptions.map((specialist) => <option key={specialist} value={specialist}>{specialist}</option>)}
                    </select>
                  </td>
                  <td className="px-4 text-center"><button type="button" onClick={() => setQuestionModalCategory(config.category)} className="inline-flex min-w-[128px] cursor-pointer items-center justify-center rounded-[8px] bg-[#4562f0] px-4 py-2 text-xs text-white transition-colors hover:bg-[#4f71fc] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]">Выбрать ({config.questionIds.length})</button></td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>

        <button type="button" onClick={() => { if (nextCategory) ensureCategoryConfig(nextCategory); setSaved(false); }} disabled={!nextCategory} className="mt-4 cursor-pointer rounded-md px-1 py-1 text-sm text-[#30384f] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] disabled:cursor-default disabled:text-[#a0a6b7]">+ {nextCategory ? "Добавить ещё категорию" : "Все категории добавлены"}</button>

        <div className="mt-12 flex flex-col items-center gap-3">
          <Button text="Сохранить изменения" variant="primary" onClick={() => setSaved(true)} />
          {saved && <p role="status" className="text-sm text-[#087f1a]">Изменения сохранены</p>}
        </div>
      </div>

      {questionModalCategory && editingQuestionId === undefined && (
        <QuestionSelectionModal category={questionModalCategory} configs={configs} questions={questions} onCategoryChange={(category) => { ensureCategoryConfig(category); setQuestionModalCategory(category); }} onToggle={toggleQuestion} onEdit={(questionId) => setEditingQuestionId(questionId)} onAddQuestion={() => setEditingQuestionId(null)} onClose={() => setQuestionModalCategory(null)} />
      )}
      {questionModalCategory && editingQuestionId !== undefined && (
        <EditQuestionModal question={editingQuestion} onSave={saveQuestion} onClose={() => setEditingQuestionId(undefined)} />
      )}
    </section>
  );
}
