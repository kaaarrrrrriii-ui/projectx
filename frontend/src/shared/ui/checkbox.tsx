"use client";

import type { InputHTMLAttributes, ReactNode } from "react";

interface Props extends Omit<InputHTMLAttributes<HTMLInputElement>, "type"> {
  label: ReactNode;
  variant?: "default" | "chip";
  className?: string;
}

export default function Checkbox({
  label,
  variant = "default",
  className = "",
  ...props
}: Props) {
  if (variant === "chip") {
    return (
      <label className={`relative block cursor-pointer ${className}`}>
        <input
          {...props}
          type="checkbox"
          className="peer sr-only"
        />

        <span
          className="
            flex min-h-[42px] w-full
            items-center justify-center
            rounded-[15px]
            border border-[#000828]
            bg-[#FCFDFF]
            px-[10px]

            text-center
            text-[16px]
            font-normal
            leading-5
            text-[#000828]

            transition-[background-color,border-color,color]
            duration-150

            hover:border-[#4562F0]
            hover:text-[#4562F0]

            peer-checked:border-[#4562F0]
            peer-checked:bg-[#DEE7FD]
            peer-checked:text-[#4562F0]

            peer-focus-visible:outline
            peer-focus-visible:outline-[3px]
            peer-focus-visible:outline-[#4562F0]
            peer-focus-visible:outline-offset-[3px]

            peer-disabled:cursor-not-allowed
            peer-disabled:opacity-40

            motion-reduce:transition-none
          "
        >
          {label}
        </span>
      </label>
    );
  }

  return (
    <label className={`inline-flex cursor-pointer items-center gap-2 ${className}`}>
      <input
        {...props}
        type="checkbox"
        className="peer sr-only"
      />

      <span
        className="
          relative flex h-5 w-5 shrink-0
          items-center justify-center
          rounded-[5px]
          border border-black

          peer-checked:border-[#4562F0]
          peer-checked:bg-[#DEE7FD]

          after:absolute
          after:h-[7px]
          after:w-[12px]
          after:-translate-y-[1px]
          after:rotate-[-45deg]
          after:border-b-2
          after:border-l-2
          after:border-[#4562F0]
          after:opacity-0
          after:content-['']

          peer-checked:after:opacity-100

          peer-focus-visible:outline
          peer-focus-visible:outline-[3px]
          peer-focus-visible:outline-[#4562F0]
          peer-focus-visible:outline-offset-[3px]
        "
      />

      <span>{label}</span>
    </label>
  );
}