export type ChartSegment = {
  label: string;
  value: number;
  color: string;
};

export default function AnalyticsDonut({
  title,
  segments,
  centerValue,
  centerLabel,
}: {
  title: string;
  segments: ChartSegment[];
  centerValue: string;
  centerLabel: string;
}) {
  const chartSegments = segments.map((segment, index) => ({
    ...segment,
    offset: segments
      .slice(0, index)
      .reduce((total, previousSegment) => total + previousSegment.value, 0),
  }));

  return (
    <section className="rounded-[15px] border border-[#7990ff] bg-white/85 p-5" aria-label={title}>
      <h2 className="text-base font-medium text-[#000828]">{title}</h2>
      <div className="mt-5 grid items-center gap-6 sm:grid-cols-[190px_minmax(0,1fr)]">
        <div className="relative mx-auto h-[190px] w-[190px]">
          <svg viewBox="0 0 120 120" role="img" aria-label={title} className="h-full w-full -rotate-90">
            <circle cx="60" cy="60" r="45" fill="none" stroke="#edf0fa" strokeWidth="18" />
            {chartSegments.map((segment) => (
                <circle
                  key={segment.label}
                  cx="60"
                  cy="60"
                  r="45"
                  fill="none"
                  stroke={segment.color}
                  strokeWidth="18"
                  pathLength="100"
                  strokeDasharray={String(segment.value) + " " + String(100 - segment.value)}
                  strokeDashoffset={-segment.offset}
                  strokeLinecap="butt"
                />
            ))}
          </svg>
          <div className="pointer-events-none absolute inset-0 flex flex-col items-center justify-center text-center">
            <strong className="text-3xl font-extrabold text-[#4562f0]">{centerValue}</strong>
            <span className="mt-0.5 max-w-20 text-[11px] leading-4 text-[#7b849b]">{centerLabel}</span>
          </div>
        </div>

        <ul className="grid gap-2.5">
          {segments.map((segment) => (
            <li key={segment.label} className="flex items-center justify-between gap-4 text-xs">
              <span className="flex min-w-0 items-center gap-2 text-[#30384f]">
                <span className="h-2.5 w-2.5 shrink-0 rounded-full" style={{ backgroundColor: segment.color }} />
                <span className="truncate">{segment.label}</span>
              </span>
              <strong className="shrink-0 text-[#000828]">{segment.value}%</strong>
            </li>
          ))}
        </ul>
      </div>
    </section>
  );
}
