type TopicOptionProps = {
  value: string;
  label: string;
  checked: boolean;
  onChange: (value: string) => void;
  required?: boolean;
  children?: React.ReactNode;
};

export default function TopicOption({
  value,
  label,
  checked,
  onChange,
  required = false,
  children,
}: TopicOptionProps) {
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
          className={`flex min-h-[42px] w-full items-center justify-center rounded-[15px] border border-[#000828] bg-[#fcfdff] px-10 text-center text-base leading-5 text-[#000828] transition-[background-color,border-color,color] duration-150 hover:border-[#4562f0] hover:text-[#4562f0] peer-checked:border-[#4562f0] peer-checked:bg-[#dee7fd] peer-checked:text-[#4562f0] peer-focus-visible:outline-3 peer-focus-visible:outline-offset-3 peer-focus-visible:outline-[#4562f0] ${children ? "pr-[82px]" : ""}`}
        >
          {label}
        </span>
      </label>
      {children}
    </div>
  );
}
