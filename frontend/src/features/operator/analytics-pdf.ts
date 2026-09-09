import type { ChartSegment } from "./analytics-donut";

export type PdfReportData = {
  dateFrom: string;
  dateTo: string;
  processed: number;
  urgentShare: number;
  returnedShare: number;
  averageAcceptance: string;
  averageResolution: string;
  categories: ChartSegment[];
  applicants: ChartSegment[];
};

function formatDate(value: string) {
  if (!value) return "не указана";
  return new Intl.DateTimeFormat("ru-RU").format(new Date(value + "T00:00:00"));
}

function drawRoundedRect(context: CanvasRenderingContext2D, x: number, y: number, width: number, height: number, radius: number) {
  context.beginPath();
  context.roundRect(x, y, width, height, radius);
  context.fill();
  context.stroke();
}

function drawDonut(
  context: CanvasRenderingContext2D,
  centerX: number,
  centerY: number,
  segments: ChartSegment[],
  centerValue: string,
  centerLabel: string,
) {
  const total = segments.reduce((sum, segment) => sum + Math.max(0, segment.value), 0);
  const radius = 105;

  context.save();
  context.beginPath();
  context.arc(centerX, centerY, radius, 0, Math.PI * 2);
  context.strokeStyle = "#e9edfb";
  context.lineWidth = 38;
  context.stroke();

  if (total <= 0) {
    context.restore();
    return;
  }

  let startAngle = -Math.PI / 2;
  segments.forEach((segment) => {
    const endAngle = startAngle + (Math.PI * 2 * Math.max(0, segment.value)) / total;
    context.beginPath();
    context.arc(centerX, centerY, radius, startAngle, endAngle);
    context.strokeStyle = segment.color;
    context.lineWidth = 38;
    context.lineCap = "butt";
    context.stroke();
    startAngle = endAngle;
  });

  context.textAlign = "center";
  context.textBaseline = "middle";
  context.fillStyle = "#4562f0";
  setFittedFont(context, centerValue, 130, 34);
  context.fillText(centerValue, centerX, centerY - 10);
  context.fillStyle = "#7b849b";
  context.font = "400 13px Montserrat, Arial, sans-serif";
  context.fillText(centerLabel, centerX, centerY + 25);
  context.restore();
}

function setFittedFont(
  context: CanvasRenderingContext2D,
  text: string,
  maxWidth: number,
  preferredSize = 36,
) {
  let fontSize = preferredSize;
  do {
    context.font = "800 " + fontSize + "px Montserrat, Arial, sans-serif";
    fontSize -= 1;
  } while (context.measureText(text).width > maxWidth && fontSize >= 22);
}

function createPdfFromJpeg(jpeg: Uint8Array, imageWidth: number, imageHeight: number) {
  const encoder = new TextEncoder();
  const chunks: Uint8Array[] = [];
  const offsets: number[] = [0];
  let byteLength = 0;

  function append(value: string | Uint8Array) {
    const bytes = typeof value === "string" ? encoder.encode(value) : value;
    chunks.push(bytes);
    byteLength += bytes.length;
  }

  function object(id: number, body: string) {
    offsets[id] = byteLength;
    append(id + " 0 obj\n" + body + "\nendobj\n");
  }

  append("%PDF-1.4\n%1234\n");
  object(1, "<< /Type /Catalog /Pages 2 0 R >>");
  object(2, "<< /Type /Pages /Kids [3 0 R] /Count 1 >>");
  object(3, "<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595.28 841.89] /Resources << /XObject << /Im0 5 0 R >> >> /Contents 4 0 R >>");
  const content = "q\n595.28 0 0 841.89 0 0 cm\n/Im0 Do\nQ\n";
  object(4, "<< /Length " + content.length + " >>\nstream\n" + content + "endstream");

  offsets[5] = byteLength;
  append("5 0 obj\n<< /Type /XObject /Subtype /Image /Width " + imageWidth + " /Height " + imageHeight + " /ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /DCTDecode /Length " + jpeg.length + " >>\nstream\n");
  append(jpeg);
  append("\nendstream\nendobj\n");

  const xrefOffset = byteLength;
  append("xref\n0 6\n0000000000 65535 f \n");
  for (let id = 1; id <= 5; id += 1) append(String(offsets[id]).padStart(10, "0") + " 00000 n \n");
  append("trailer\n<< /Size 6 /Root 1 0 R >>\nstartxref\n" + xrefOffset + "\n%%EOF");
  return new Blob(chunks as BlobPart[], { type: "application/pdf" });
}

export async function downloadAnalyticsPdf(data: PdfReportData) {
  await document.fonts.ready;
  const canvas = document.createElement("canvas");
  canvas.width = 1240;
  canvas.height = 1754;
  const context = canvas.getContext("2d");
  if (!context) return;

  context.fillStyle = "#f7f9fe";
  context.fillRect(0, 0, canvas.width, canvas.height);
  context.fillStyle = "#ffffff";
  context.strokeStyle = "#4562f0";
  context.lineWidth = 2;
  drawRoundedRect(context, 55, 55, 1130, 1644, 28);

  context.fillStyle = "#4562f0";
  context.font = "800 46px Montserrat, Arial, sans-serif";
  context.fillText("ОТКЛИК — аналитический отчёт", 100, 135);
  context.fillStyle = "#646d86";
  context.font = "500 22px Montserrat, Arial, sans-serif";
  context.fillText("Период: " + formatDate(data.dateFrom) + " — " + formatDate(data.dateTo), 100, 180);
  context.fillText("Сформирован: " + new Intl.DateTimeFormat("ru-RU", { dateStyle: "long", timeStyle: "short" }).format(new Date()), 100, 215);

  const metrics = [
    { label: "Обработано", value: String(data.processed), suffix: "обращений" },
    { label: "Срочные", value: data.urgentShare + "%", suffix: "от всех обращений" },
    { label: "Возвращённые", value: data.returnedShare + "%", suffix: "от всех обращений" },
    { label: "До принятия", value: data.averageAcceptance, suffix: "в среднем" },
    { label: "До решения", value: data.averageResolution, suffix: "в среднем" },
  ];

  metrics.forEach((metric, index) => {
    const isTimeMetric = index >= 3;
    const width = isTimeMetric ? 510 : 333;
    const x = isTimeMetric ? 100 + (index - 3) * 530 : 100 + index * 353;
    const y = isTimeMetric ? 455 : 265;
    context.fillStyle = "#f9faff";
    context.strokeStyle = "#9cacf8";
    drawRoundedRect(context, x, y, width, 160, 18);
    context.fillStyle = "#30384f";
    context.font = "500 17px Montserrat, Arial, sans-serif";
    context.fillText(metric.label, x + 18, y + 40);
    context.fillStyle = "#4562f0";
    setFittedFont(context, metric.value, width - 36, isTimeMetric ? 40 : 36);
    context.fillText(metric.value, x + 18, y + 99);
    context.fillStyle = "#7b849b";
    context.font = "400 14px Montserrat, Arial, sans-serif";
    context.fillText(metric.suffix, x + 18, y + 132);
  });

  const charts = [
    {
      title: "Распределение по категориям",
      segments: data.categories,
      centerLabel: "обращений",
      x: 100,
    },
    {
      title: "Распределение по заявителям",
      segments: data.applicants,
      centerLabel: "заявителей",
      x: 640,
    },
  ];

  charts.forEach((chart) => {
    context.fillStyle = "#f9faff";
    context.strokeStyle = "#9cacf8";
    drawRoundedRect(context, chart.x, 665, 500, 620, 20);
    context.fillStyle = "#000828";
    context.font = "600 22px Montserrat, Arial, sans-serif";
    context.fillText(chart.title, chart.x + 30, 715);
    drawDonut(
      context,
      chart.x + 250,
      890,
      chart.segments,
      String(data.processed),
      chart.centerLabel,
    );
    chart.segments.forEach((segment, index) => {
      const y = 1080 + index * 38;
      context.fillStyle = segment.color;
      context.beginPath();
      context.arc(chart.x + 42, y - 7, 8, 0, Math.PI * 2);
      context.fill();
      context.fillStyle = "#30384f";
      context.font = "400 16px Montserrat, Arial, sans-serif";
      context.fillText(segment.label, chart.x + 62, y, 310);
      context.fillStyle = "#000828";
      context.font = "700 16px Montserrat, Arial, sans-serif";
      context.textAlign = "right";
      context.fillText(segment.value + "%", chart.x + 462, y);
      context.textAlign = "left";
    });
  });

  context.fillStyle = "#eef1ff";
  context.strokeStyle = "#aab7f8";
  drawRoundedRect(context, 100, 1340, 1040, 210, 20);
  context.fillStyle = "#000828";
  context.font = "700 24px Montserrat, Arial, sans-serif";
  context.fillText("Краткий вывод", 135, 1395);
  context.font = "400 19px Montserrat, Arial, sans-serif";
  context.fillStyle = "#30384f";
  const summary = [
    "• Основная доля обращений относится к кибербуллингу и семейным конфликтам.",
    "• Большинство обращений поступает от школьников.",
    "• Доля срочных обращений составляет " + data.urgentShare + "%.",
    "• Среднее время принятия обращения — " + data.averageAcceptance + ".",
  ];
  summary.forEach((line, index) => context.fillText(line, 135, 1440 + index * 29));

  context.fillStyle = "#7b849b";
  context.font = "400 15px Montserrat, Arial, sans-serif";
  context.fillText("Автоматически сформировано системой «Отклик»", 100, 1635);

  const dataUrl = canvas.toDataURL("image/jpeg", 0.94);
  const binary = atob(dataUrl.split(",")[1]);
  const jpeg = new Uint8Array(binary.length);
  for (let index = 0; index < binary.length; index += 1) jpeg[index] = binary.charCodeAt(index);
  const pdf = createPdfFromJpeg(jpeg, canvas.width, canvas.height);
  const url = URL.createObjectURL(pdf);
  const link = document.createElement("a");
  link.href = url;
  link.download = "otklik-analytics-" + data.dateFrom + "-" + data.dateTo + ".pdf";
  link.click();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}
