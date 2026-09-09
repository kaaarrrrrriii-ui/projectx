import DashboardIcon from "./dashboard-icons";

type MetricCardProps = {
  title: string;
  value: string;
  description: string;
  icon: "mail" | "urgent" | "clock" | "return";
  tone: "blue" | "red" | "orange";
  badge?: string;
};

const tones = {
  blue: { border: "border-[#7990ff]", text: "text-[#4562f0]", badge: "bg-[#4562f0]" },
  red: { border: "border-[#ff6d6d]", text: "text-[#e5141b]", badge: "bg-[#e5141b]" },
  orange: { border: "border-[#ff9a4d]", text: "text-[#ff7000]", badge: "bg-[#ff7000]" },
};

export default function MetricCard({ title, value, description, icon, tone, badge }: MetricCardProps) {
  const colors = tones[tone];

  return (
    <article className={`flex min-h-[218px] flex-col rounded-2xl border bg-white p-6 ${colors.border}`}>
      <DashboardIcon name={icon} className={`h-10 w-10 ${colors.text}`} />
      <h2 className="mt-4 min-h-12 text-base font-medium leading-6 text-[#000828]">{title}</h2>
      <div className="mt-auto">
        {badge ? (
          <span className={`inline-flex rounded-full px-3 py-1 text-sm font-semibold text-white ${colors.badge}`}>{badge}</span>
        ) : (
          <strong className={`block text-2xl font-extrabold ${colors.text}`}>{value}</strong>
        )}
        <p className="mt-2 text-xs leading-4 text-[#4f5873]">{description}</p>
      </div>
    </article>
  );
}
