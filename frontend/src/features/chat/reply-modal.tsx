"use client";

import Button from "@/shared/ui/button";
import Image from "next/image";
import {
  ChangeEvent,
  DragEvent,
  FormEvent,
  KeyboardEvent,
  useEffect,
  useRef,
  useState,
} from "react";
import DialogShell from "./dialog-shell";

const MAX_FILES = 5;
const MAX_FILE_SIZE = 10 * 1024 * 1024;

type ReplyFile = {
  file: File;
  previewUrl?: string;
};

function isSupportedFile(file: File) {
  return (
    file.type.startsWith("image/") ||
    file.type === "application/pdf" ||
    file.name.toLowerCase().endsWith(".pdf")
  );
}

export default function ReplyModal({
  onClose,
  onSend,
}: {
  onClose: () => void;
  onSend: (message: string, fileNames: string[]) => void;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const previewUrls = useRef(new Set<string>());
  const [message, setMessage] = useState("");
  const [items, setItems] = useState<ReplyFile[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    const urls = previewUrls.current;
    return () => urls.forEach((url) => URL.revokeObjectURL(url));
  }, []);

  function openFilePicker() {
    inputRef.current?.click();
  }

  function addFiles(fileList: FileList | File[]) {
    const files = Array.from(fileList);
    const supported = files.filter(isSupportedFile);
    const valid = supported.filter((file) => file.size <= MAX_FILE_SIZE);

    if (supported.length !== files.length) {
      setError("Можно добавить только изображения или PDF.");
    } else if (valid.length !== supported.length) {
      setError("Размер каждого файла не должен превышать 10 МБ.");
    } else if (items.length + valid.length > MAX_FILES) {
      setError("Можно добавить не больше 5 файлов.");
    } else {
      setError("");
    }

    const availableSlots = MAX_FILES - items.length;
    const additions = valid.slice(0, availableSlots).map((file) => {
      if (!file.type.startsWith("image/")) return { file };
      const previewUrl = URL.createObjectURL(file);
      previewUrls.current.add(previewUrl);
      return { file, previewUrl };
    });

    setItems((current) => [...current, ...additions]);
  }

  function removeFile(index: number) {
    setItems((current) => {
      const item = current[index];
      if (item.previewUrl) {
        URL.revokeObjectURL(item.previewUrl);
        previewUrls.current.delete(item.previewUrl);
      }
      return current.filter((_, itemIndex) => itemIndex !== index);
    });
    setError("");
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

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!message.trim() && items.length === 0) return;
    onSend(message.trim(), items.map(({ file }) => file.name));
  }

  return (
    <DialogShell
      labelledBy="reply-dialog-heading"
      onClose={onClose}
      showClose
      className="max-w-[738px]"
    >
      <form onSubmit={handleSubmit}>
        <h2
          id="reply-dialog-heading"
          className="pr-12 text-[28px] leading-[1.2] font-extrabold tracking-[-0.025em] text-[var(--color-primary)] max-[599px]:text-[22px]"
        >
          Добавить ответное сообщение специалисту
        </h2>

        <textarea
          value={message}
          onChange={(event) => setMessage(event.target.value)}
          placeholder="Здесь можно написать всё, что тебя беспокоит..."
          aria-label="Сообщение специалисту"
          className="mt-4 block min-h-[168px] w-full resize-y rounded-[12px] border border-[#333] bg-[#fcfdff] px-3 py-2.5 text-[15px] leading-6 outline-none placeholder:text-[#9196a7] focus:border-[var(--color-primary)] focus:ring-2 focus:ring-[#4562f0]/20 max-[499px]:min-h-[130px]"
        />

        <input
          ref={inputRef}
          className="sr-only"
          type="file"
          accept="image/*,.pdf,application/pdf"
          multiple
          onChange={(event: ChangeEvent<HTMLInputElement>) => {
            if (event.target.files) addFiles(event.target.files);
            event.target.value = "";
          }}
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
          onDragLeave={() => setIsDragging(false)}
          onDrop={handleDrop}
          className={`mt-[18px] flex min-h-[140px] cursor-pointer flex-col items-center justify-center rounded-[12px] border px-4 py-4 text-center transition-colors focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[var(--color-primary)] ${
            isDragging
              ? "border-[var(--color-primary)] bg-[#eef1ff]"
              : "border-[#333] bg-[#fcfdff] hover:border-[var(--color-primary)]"
          }`}
        >
          <Image
            src="/images/mediaAddIcon.svg"
            alt=""
            width={58}
            height={58}
            className="h-[58px] w-[58px]"
          />
          <p className="mt-2 text-[13px] leading-5">Нажми, чтобы выбрать файлы</p>
          <p className="text-[12px] leading-4 text-[#9196a7]">
            Можно добавить до 5 файлов (фото, скриншоты, PDF), размер до 10 МБ каждый.
          </p>
        </div>

        <div className="mt-3 flex min-h-[74px] flex-wrap gap-2" aria-live="polite">
          {items.map((item, index) => (
            <div
              key={`${item.file.name}-${item.file.lastModified}`}
              className="group relative h-[74px] w-[74px] overflow-hidden rounded-[12px] border border-[#d7dbea] bg-[#eef1ff]"
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
                <span className="flex h-full items-center justify-center px-1 text-center text-[10px] text-[var(--color-primary)]">
                  PDF<br />{item.file.name}
                </span>
              )}
              <button
                type="button"
                onClick={() => removeFile(index)}
                aria-label={`Удалить файл ${item.file.name}`}
                className="absolute top-1 right-1 flex h-5 w-5 cursor-pointer items-center justify-center rounded-full bg-[#000828]/75 text-sm text-white opacity-0 group-hover:opacity-100 focus-visible:opacity-100"
              >
                ×
              </button>
            </div>
          ))}

          {items.length < MAX_FILES && (
            <button
              type="button"
              onClick={openFilePicker}
              aria-label="Добавить файлы"
              className="flex h-[74px] w-[74px] cursor-pointer items-center justify-center rounded-[12px] border border-[var(--color-primary)] bg-[#e4e9ff] text-[32px] font-light text-[var(--color-primary)] hover:bg-[#dce3ff]"
            >
              +
            </button>
          )}
        </div>

        {error && <p className="mt-2 text-[13px] text-[#b42318]" role="alert">{error}</p>}

        <Button
          text="Отправить"
          type="submit"
          variant="primary"
          size="small"
          disabled={!message.trim() && items.length === 0}
          className="mt-[18px] w-full"
        />
      </form>
    </DialogShell>
  );
}
