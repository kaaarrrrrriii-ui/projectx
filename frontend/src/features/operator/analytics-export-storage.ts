import type { PdfReportData } from "./analytics-pdf";

export type AnalyticsExportRecord = {
  id: string;
  createdAt: string;
  fileName: string;
  data: PdfReportData;
};

export const analyticsExportsStorageKey = "otklik-analytics-exports";

export const defaultAnalyticsExport: AnalyticsExportRecord = {
  id: "initial-analytics-report",
  createdAt: "2026-09-09T12:00:00+05:00",
  fileName: "Аналитический отчёт.pdf",
  data: {
    dateFrom: "2026-09-01",
    dateTo: "2026-09-09",
    processed: 67,
    urgentShare: 18,
    returnedShare: 9,
    averageAcceptance: "14 мин",
    averageResolution: "2 ч 36 мин",
    categories: [
      { label: "Кибербуллинг", value: 31, color: "#4562f0" },
      { label: "Семейные конфликты", value: 25, color: "#8b5cf6" },
      { label: "Травля и оскорбления", value: 22, color: "#ff7a45" },
      { label: "Давление и угрозы", value: 14, color: "#ef476f" },
      { label: "Другие темы", value: 8, color: "#35b7a0" },
    ],
    applicants: [
      { label: "Школьники", value: 48, color: "#4562f0" },
      { label: "Родители", value: 29, color: "#ff8a3d" },
      { label: "Студенты", value: 23, color: "#36b5a2" },
    ],
  },
};

export function readAnalyticsExports(serializedValue?: string) {
  if (typeof window === "undefined") return [defaultAnalyticsExport];

  try {
    const saved = JSON.parse(
      serializedValue ?? window.localStorage.getItem(analyticsExportsStorageKey) ?? "[]",
    ) as AnalyticsExportRecord[];
    return saved.length > 0 ? saved : [defaultAnalyticsExport];
  } catch {
    return [defaultAnalyticsExport];
  }
}

export function saveAnalyticsExport(data: PdfReportData) {
  if (typeof window === "undefined") return;

  const createdAt = new Date().toISOString();
  const record: AnalyticsExportRecord = {
    id: createdAt,
    createdAt,
    fileName: `Аналитика ${data.dateFrom} — ${data.dateTo}.pdf`,
    data,
  };

  const current = readAnalyticsExports().filter((item) => item.id !== defaultAnalyticsExport.id);
  window.localStorage.setItem(
    analyticsExportsStorageKey,
    JSON.stringify([record, ...current].slice(0, 10)),
  );
}
