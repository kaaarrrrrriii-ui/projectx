type TopicOptionProps = {
  value: string;
  label: string;
  checked: boolean;
  onChange: (value: string) => void;
  required?: boolean;
  variant?: "default" | "secondary";
  children?: React.ReactNode;
};

export default function TopicOption({
  value,
  label,
  checked,
  onChange,
  required = false,
  variant = "default",
  children,
}: TopicOptionProps) {
  const colorClasses =
    variant === "secondary"
      ? "border-[#4562f0] bg-transparent text-[#4562f0] hover:border-[#4562f0] hover:text-[#4562f0] peer-checked:border-[#4562f0] peer-checked:bg-[#dee7fd] peer-checked:text-[#4562f0]"
      : "border-[#808393] bg-transparent text-[#808393] hover:border-[#4562f0] hover:text-[#4562f0] peer-checked:border-[#4562f0] peer-checked:bg-[#dee7fd] peer-checked:text-[#4562f0]";

  return (
    <div className="relative min-w-0">
      <label className="group block cursor-pointer">
        <input
          type="radio"
          name="topic"
          value={value}
          checked={checked}
          onChange={() => onChange(value)}
          required={required}
          className="peer sr-only"
        />
        <span
          className={`flex min-h-[40px] w-full items-center justify-center rounded-[12px] border px-6 py-2 text-center text-[14px] leading-5 transition-[background-color,border-color,color] duration-150 peer-focus-visible:outline-3 peer-focus-visible:outline-offset-3 peer-focus-visible:outline-[#4562f0] ${colorClasses} ${children ? "pr-[82px]" : ""}`}
        >
          {label}
        </span>
      </label>
      {children}
    </div>
  );
}
