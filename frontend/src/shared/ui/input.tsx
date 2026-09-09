import type { InputHTMLAttributes } from "react";

interface Props extends InputHTMLAttributes<HTMLInputElement> {
  className?: string;
}

export default function Input({
  className = "",
  ...props
}: Props) {
  return (
    <input
      {...props}
      className={`
        box-border
        h-[42px]
        w-[279px]
        max-w-full
        rounded-[15px]
        border
        border-[#000828]
        bg-[#FCFDFF]
        px-[10px]

        text-center
        text-[16px]
        font-normal
        leading-[20px]
        text-[#000828]

        placeholder:text-[#000828]
        placeholder:opacity-100

        transition-[background-color,border-color,color]
        duration-150

        hover:border-[#4562F0]
        hover:text-[#4562F0]
        hover:placeholder:text-[#4562F0]

        focus:border-[#4562F0]
        focus:bg-[#DEE7FD]
        focus:text-[#4562F0]
        focus:outline-none
        focus:placeholder:text-[#4562F0]

        disabled:cursor-not-allowed
        disabled:opacity-50

        ${className}
      `}
    />
  );
}
