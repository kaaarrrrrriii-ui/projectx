"use client";

import Input from "@/shared/ui/input";
import { useMemo, useState, useSyncExternalStore } from "react";
import {
  analyticsExportsStorageKey,
  defaultAnalyticsExport,
  readAnalyticsExports,
  type AnalyticsExportRecord,
} from "./analytics-export-storage";
import { downloadAnalyticsPdf } from "./analytics-pdf";

function formatExportDate(value: string) {
  return new Intl.DateTimeFormat("ru-RU", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value));
}

function subscribeToExports(onStoreChange: () => void) {
  function handleStorage(event: StorageEvent) {
    if (event.key === analyticsExportsStorageKey) onStoreChange();
  }

  window.addEventListener("storage", handleStorage);
  return () => window.removeEventListener("storage", handleStorage);
}

function getExportsSnapshot() {
  return window.localStorage.getItem(analyticsExportsStorageKey) ?? "";
}

export default function ExportsDashboard() {
  const [search, setSearch] = useState("");
  const serializedRecords = useSyncExternalStore(subscribeToExports, getExportsSnapshot, () => "");
  const records = useMemo<AnalyticsExportRecord[]>(
    () => serializedRecords ? readAnalyticsExports(serializedRecords) : [defaultAnalyticsExport],
    [serializedRecords],
  );
  const [downloadingId, setDownloadingId] = useState<string | null>(null);

  const filteredRecords = useMemo(() => {
    const query = search.trim().toLocaleLowerCase("ru");
    if (!query) return records;

    return records.filter((record) =>
      `${record.fileName} ${formatExportDate(record.createdAt)}`
        .toLocaleLowerCase("ru")
        .includes(query),
    );
  }, [records, search]);

  async function downloadRecord(record: AnalyticsExportRecord) {
    setDownloadingId(record.id);
    try {
      await downloadAnalyticsPdf(record.data);
    } finally {
      setDownloadingId(null);
    }
  }

  return (
    <section className="min-w-0 flex-1 px-5 pt-5 pb-7 sm:px-[28px]" aria-labelledby="exports-heading">
      <div className="w-full">
        <label className="relative block">
          <svg viewBox="0 0 24 24" aria-hidden="true" className="pointer-events-none absolute top-1/2 left-3 z-10 h-5 w-5 -translate-y-1/2 text-[#4562f0]" fill="none">
            <circle cx="10.5" cy="10.5" r="6.5" stroke="currentColor" strokeWidth="2" />
            <path d="m16 16 5 5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
          </svg>
          <Input
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            aria-label="Поиск выгрузки"
            placeholder="Поиск выгрузки по названию или дате"
            className="!h-[42px] !w-full !rounded-[12px] !border-[#4562f0] !bg-white !pr-3 !pl-[42px] !text-left !text-base placeholder:!text-[#9196a7] focus:!bg-white focus:!text-[#000828]"
          />
        </label>

        <header className="mt-5">
          <h1 id="exports-heading" className="text-[26px] leading-tight font-extrabold tracking-[-0.02em] text-[#4562f0]">
            Последние выгрузки
          </h1>
        </header>

        <section aria-label="Список сформированных отчётов" className="mt-3 overflow-x-auto rounded-[13px] border border-[#4562f0] bg-white/85">
          <table className="w-full min-w-[560px] table-fixed border-collapse">
            <thead className="bg-[#dfe6ff] text-[#4562f0]">
              <tr className="h-[52px]">
                <th scope="col" className="w-[38%] border-r border-[#4562f0] px-4 text-center font-normal">Дата</th>
                <th scope="col" className="px-4 text-center font-normal">Файл</th>
              </tr>
            </thead>
            <tbody>
              {filteredRecords.map((record) => (
                <tr key={record.id} className="h-[58px] border-t border-[#4562f0] transition-colors hover:bg-[#f7f8ff]">
                  <td className="border-r border-[#4562f0] px-4 text-center text-sm text-[#30384f]">
                    {formatExportDate(record.createdAt)}
                  </td>
                  <td className="px-4 text-center">
                    <button
                      type="button"
                      disabled={downloadingId !== null}
                      onClick={() => void downloadRecord(record)}
                      className="cursor-pointer rounded-sm text-sm font-medium text-[#4562f0] underline-offset-4 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] disabled:cursor-wait disabled:opacity-50"
                    >
                      {downloadingId === record.id ? "Формируем PDF…" : `Скачать «${record.fileName}»`}
                    </button>
                  </td>
                </tr>
              ))}

              {filteredRecords.length === 0 && (
                <tr className="h-24 border-t border-[#4562f0]">
                  <td colSpan={2} className="px-6 text-center text-sm text-[#646d86]">
                    Выгрузки не найдены
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </section>
      </div>
    </section>
  );
}
