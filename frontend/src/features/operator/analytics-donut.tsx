export type ChartSegment = {
  label: string;
  value: number;
  color: string;
};

export default function AnalyticsDonut({
  title,
  segments,
}: {
  title: string;
  segments: ChartSegment[];
}) {
  const total = segments.reduce((sum, segment) => sum + Math.max(0, segment.value), 0);
  let offset = 0;
  const gradientStops = segments.map((segment) => {
    const start = offset;
    offset += total > 0 ? (Math.max(0, segment.value) / total) * 100 : 0;
    return `${segment.color} ${start}% ${offset}%`;
  });
  const chartDescription = segments.map((segment) => `${segment.label}: ${segment.value}%`).join(", ");

  return (
    <section className="min-h-[185px] rounded-[12px] border border-[#7990ff] bg-white/85 px-3 py-3" aria-label={title}>
      <h2 className="text-center text-[15px] leading-5 font-medium text-[#000828]">{title}</h2>
      <div className="mt-4 flex items-center justify-center gap-3">
        <div
          role="img"
          aria-label={`${title}. ${chartDescription}`}
          className="h-[100px] w-[100px] shrink-0 rounded-full border border-white shadow-[0_0_0_1px_rgba(69,98,240,0.15)]"
          style={{
            backgroundImage: total > 0 ? `conic-gradient(${gradientStops.join(", ")})` : "none",
            backgroundColor: "#edf0fa",
          }}
        />

        <ul className="grid min-w-0 gap-1.5">
          {segments.map((segment) => (
            <li key={segment.label} className="flex min-w-0 items-start gap-1.5 text-[9px] leading-[11px] text-[#30384f]">
              <span className="mt-0.5 h-2 w-2 shrink-0 rounded-[2px]" style={{ backgroundColor: segment.color }} />
              <span className="min-w-0">
                {segment.label} <strong className="font-semibold text-[#000828]">{segment.value}%</strong>
              </span>
            </li>
          ))}
        </ul>
      </div>
    </section>
  );
}
