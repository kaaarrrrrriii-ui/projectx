"use client";

import Button from "@/shared/ui/button";
import AppealNavigation from "@/shared/ui/appeal-navigation";
import { getAppealRoute } from "./routes";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { createPublicTicket } from "@/shared/api/public-api";
import { clearAppealDraft, draftLooksCritical, readAppealDraft, saveSubmissionResult } from "./appeal-draft";
import {
  ChangeEvent,
  DragEvent,
  KeyboardEvent,
  useEffect,
  useRef,
  useState,
} from "react";

const MAX_FILES = 5;
const MAX_FILE_SIZE = 10 * 1024 * 1024;

type UploadItem = {
  file: File;
  previewUrl?: string;
};

function isSupportedFile(file: File) {
  return file.type === "image/jpeg" || file.type === "image/png" || /\.(jpe?g|png)$/i.test(file.name);
}

export default function MediaUpload({ role, topic = "", formal = false }: { role: string; topic?: string; formal?: boolean }) {
  const router = useRouter();
  const inputRef = useRef<HTMLInputElement>(null);
  const previewUrls = useRef(new Set<string>());
  const [items, setItems] = useState<UploadItem[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const [error, setError] = useState("");
  const [crisisContact, setCrisisContact] = useState("");
  const [showCrisisHelp, setShowCrisisHelp] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const routeParams = { role, topic };

  useEffect(() => {
    const urls = previewUrls.current;
    const timer = window.setTimeout(() => setShowCrisisHelp(draftLooksCritical(readAppealDraft())), 0);

    return () => {
      window.clearTimeout(timer);
      urls.forEach((url) => URL.revokeObjectURL(url));
    };
  }, []);

  function openFilePicker() {
    inputRef.current?.click();
  }

  function addFiles(fileList: FileList | File[]) {
    const files = Array.from(fileList);
    const supportedFiles = files.filter(isSupportedFile);
    const validFiles = supportedFiles.filter(
      (file) => file.size <= MAX_FILE_SIZE,
    );

    if (supportedFiles.length !== files.length) {
      setError("Можно добавить только изображения JPEG или PNG.");
    } else if (validFiles.length !== supportedFiles.length) {
      setError("Размер каждого файла не должен превышать 10 МБ.");
    } else if (items.length + validFiles.length > MAX_FILES) {
      setError("Можно добавить не больше 5 файлов.");
    } else {
      setError("");
    }

    const availableSlots = MAX_FILES - items.length;
    const newItems = validFiles.slice(0, availableSlots).map((file) => {
      if (!file.type.startsWith("image/")) {
        return { file };
      }

      const previewUrl = URL.createObjectURL(file);
      previewUrls.current.add(previewUrl);
      return { file, previewUrl };
    });

    setItems((currentItems) => [...currentItems, ...newItems]);
  }

  function handleInputChange(event: ChangeEvent<HTMLInputElement>) {
    if (event.target.files) {
      addFiles(event.target.files);
    }

    event.target.value = "";
  }

  function handleDrop(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    setIsDragging(false);
    addFiles(event.dataTransfer.files);
  }

  function handleDropzoneKeyDown(event: KeyboardEvent<HTMLDivElement>) {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      openFilePicker();
    }
  }

  function removeItem(index: number) {
    setItems((currentItems) => {
      const item = currentItems[index];
      if (item.previewUrl) {
        URL.revokeObjectURL(item.previewUrl);
        previewUrls.current.delete(item.previewUrl);
      }

      return currentItems.filter((_, itemIndex) => itemIndex !== index);
    });
    setError("");
  }

  async function submitTicket() {
    if (submitting) return;
    const draft = readAppealDraft();
    if (!draft.categoryId) {
      setError("Не удалось определить категорию. Вернитесь к выбору темы.");
      return;
    }
    setSubmitting(true);
    setError("");
    const form = new FormData();
    form.set("applicant_type", draft.role || role);
    form.set("category_id", String(draft.categoryId));
    form.set("description", draft.description ?? "");
    form.set("custom_topic", draft.customTopic ?? "");
    if (showCrisisHelp) form.set("crisis_contact", crisisContact);
    for (const [questionID, answer] of Object.entries(draft.answers ?? {})) {
      form.set(`answers[${questionID}]`, String(answer.answerId));
    }
    for (const item of items) form.append("attachments[]", item.file);

    try {
      const result = await createPublicTicket(form);
      saveSubmissionResult(result);
      clearAppealDraft();
      router.push(getAppealRoute("success", { role }));
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Не удалось отправить обращение. Попробуйте ещё раз.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section
      aria-labelledby="media-heading"
      className="flex flex-1 flex-col px-[4.5%] pt-9 pb-[38px] max-[699px]:px-5 max-[699px]:pt-8 max-[699px]:pb-7"
    >
      <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col">
        <div className="flex flex-col rounded-[15px] border border-[var(--color-primary)] bg-[var(--color-background)] px-[30px] py-[50px]">
          <h1
            id="media-heading"
            className="text-[28px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[var(--color-primary)]"
          >
            Можно добавить фото
          </h1>

          <p className="mt-1 text-[15px] leading-[22px] text-[#151515] max-[699px]:mt-2 max-[699px]:text-[14px]">
            Это также необязательно, но фото могут помочь лучше понять ситуацию.
          </p>

          <input
            ref={inputRef}
            className="sr-only"
            type="file"
            accept="image/jpeg,image/png,.jpg,.jpeg,.png"
            multiple
            onChange={handleInputChange}
            aria-label="Выбрать файлы"
          />

          <div
            role="button"
            tabIndex={0}
            aria-label="Выбрать файлы или перетащить их сюда"
            onClick={openFilePicker}
            onKeyDown={handleDropzoneKeyDown}
            onDragEnter={(event) => {
              event.preventDefault();
              setIsDragging(true);
            }}
            onDragOver={(event) => event.preventDefault()}
            onDragLeave={(event) => {
              if (!event.currentTarget.contains(event.relatedTarget as Node)) {
                setIsDragging(false);
              }
            }}
            onDrop={handleDrop}
            className={`
              group mt-[9px] flex min-h-[270px] w-full cursor-pointer flex-col items-center justify-center
              rounded-[15px] border bg-[#fcfdff] px-5 pt-4 pb-3 text-center
              transition-[border-color,background-color] duration-150
              focus-visible:border-[var(--color-primary)] focus-visible:bg-[#dfe6ff]
              focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[var(--color-primary)]
              motion-reduce:transition-none
              ${
                isDragging
                  ? "border-[var(--color-primary)] bg-[#dfe6ff]"
                  : "border-[#333] hover:border-[var(--color-primary)] hover:bg-[#dfe6ff]"
              }
            `}
          >
            <div className="relative top-2 flex flex-col items-center">
              <svg
                aria-hidden="true"
                viewBox="0 0 80 80"
                fill="none"
                className={`mb-[23px] h-[72px] w-[72px] transition-colors duration-150 group-hover:text-[var(--color-primary)] group-focus-visible:text-[var(--color-primary)] motion-reduce:transition-none max-[699px]:mb-4 max-[699px]:h-14 max-[699px]:w-14 ${
                  isDragging
                    ? "text-[var(--color-primary)]"
                    : "text-[var(--color-text)]"
                }`}
              >
                <path
                  d="M40 49V3M40 3 22 21M40 3l18 18M4 49v27h72V49"
                  stroke="currentColor"
                  strokeWidth="5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>

              <span className="text-[15px] leading-[22px] text-[#000828] max-[699px]:text-[14px]">
                {formal
                  ? "Нажмите, чтобы выбрать файлы, или перетащите их сюда"
                  : "Нажми, чтобы выбрать файлы, или перетащи их сюда"}
              </span>
              <span className="text-[15px] leading-[22px] text-[#9196a7] max-[699px]:text-[13px]">
                 Можно добавить до 5 файлов (JPEG или PNG). Размер до 10
                МБ каждый.
              </span>
            </div>
          </div>

          <div className="mt-[13px] min-h-[89px]">
            <div className="flex flex-wrap gap-2.5" aria-live="polite">
              {items.map((item, index) => (
                <div
                  key={`${item.file.name}-${item.file.lastModified}`}
                  className="group relative h-[88px] w-[86px] overflow-hidden rounded-[15px] border border-[#d7dbea] bg-[#eef1ff]"
                >
                  {item.previewUrl ? (
                    <Image
                      src={item.previewUrl}
                      alt={`Превью файла ${item.file.name}`}
                      fill
                      unoptimized
                      className="object-cover"
                    />
                  ) : (
                    <div className="flex h-full flex-col items-center justify-center px-2 text-center text-[11px] leading-4 text-[#4562f0]">
                      <span className="text-lg font-bold">IMG</span>
                      <span className="w-full truncate">{item.file.name}</span>
                    </div>
                  )}

                  <button
                    type="button"
                    onClick={(event) => {
                      event.stopPropagation();
                      removeItem(index);
                    }}
                    aria-label={`Удалить файл ${item.file.name}`}
                    className="absolute top-1 right-1 flex h-6 w-6 cursor-pointer items-center justify-center rounded-full bg-[#000828]/75 text-lg leading-none text-white opacity-0 transition-opacity group-hover:opacity-100 focus-visible:opacity-100 focus-visible:outline-2 focus-visible:outline-white"
                  >
                    ×
                  </button>
                </div>
              ))}

              {items.length < MAX_FILES && (
                <button
                  type="button"
                  onClick={openFilePicker}
                  aria-label="Добавить ещё файлы"
                  className="group flex h-[88px] w-[86px] cursor-pointer items-center justify-center rounded-[15px] border border-[var(--color-primary)] bg-[#e4e9ff] text-[var(--color-primary)] transition-colors hover:bg-[#dce3ff] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[var(--color-primary)]"
                >
                  <span className="relative h-7 w-7" aria-hidden="true">
                    <span className="absolute top-1/2 left-0 h-px w-full -translate-y-1/2 bg-current" />
                    <span className="absolute top-0 left-1/2 h-full w-px -translate-x-1/2 bg-current" />
                  </span>
                </button>
              )}
            </div>

            {error && (
              <p className="mt-2 text-sm text-[#b42318]" role="alert">
                {error}
              </p>
            )}
          </div>

          {showCrisisHelp && (
            <aside className="mt-5 rounded-[15px] border border-[#d92d20] bg-[#fff1f1] p-4" aria-labelledby="crisis-help-heading">
              <h2 id="crisis-help-heading" className="text-[17px] font-semibold text-[#8f1111]">Если опасность непосредственная — позвоните 112</h2>
              <p className="mt-1 text-[13px] leading-5 text-[#451313]">Можно также бесплатно и анонимно позвонить на детский телефон доверия: 124 или 8 800 2000-122. Экстренная психологическая помощь МЧС: +7 (495) 989-50-50. Отправка обращения останется доступной.</p>
              <label className="mt-4 block text-[13px] text-[#451313]">
                Способ связи, если вы хотите его оставить (необязательно)
                <input value={crisisContact} onChange={(event) => setCrisisContact(event.target.value)} maxLength={500} className="mt-2 block h-10 w-full rounded-[10px] border border-[#b56b6b] bg-white px-3 text-[#151515] outline-none focus:border-[#4562f0]" placeholder="Телефон, почта или другой удобный способ" />
              </label>
              <p className="mt-2 text-[12px] text-[#6d3333]">Контакт сохранится отдельно от текста и будет доступен только сотрудникам, работающим с кризисным обращением.</p>
            </aside>
          )}
        </div>

          <AppealNavigation
            backHref={getAppealRoute("details", routeParams)}
            skipHref="#"
            onSkip={(event) => { event.preventDefault(); void submitTicket(); }}
            primaryText={submitting ? "Отправляем…" : "Отправить обращение"}
            primaryOnClick={() => void submitTicket()}
            primaryDisabled={submitting}
          primaryVariant="secondary"
          primaryClassName="h-[44px] w-[234px] rounded-[11px] px-5 text-[15px] font-normal"
        />
      </div>
    </section>
  );
}
