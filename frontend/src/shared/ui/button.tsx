"use client";

import Link from "next/link";

interface Props {
  text: string;
  variant?: "primary" | "secondary";
  size?: "default" | "small";

  fill?: boolean;

  link?: string;
  onClick?: () => void;
  type?: "button" | "submit" | "reset";
  disabled?: boolean;
  target?: "_blank" | "_self";
  ariaLabel?: string;
  className?: string;
}

export default function Button({
  text,
  variant,
  size = "default",
  fill,
  link,
  onClick,
  type = "button",
  disabled = false,
  target = "_self",
  ariaLabel,
  className = "",
}: Props) {
  // Старый fill=true считаем primary.
  // Без fill старая кнопка будет secondary.
  const resolvedVariant =
    variant ?? (fill ? "primary" : "secondary");

  const baseClasses = `
    inline-flex shrink-0
    items-center justify-center
    border
    font-medium
    leading-none
    whitespace-nowrap
    cursor-pointer
    select-none

    transition-[background-color,border-color,color]
    duration-150

    focus-visible:outline
    focus-visible:outline-[3px]
    focus-visible:outline-[#4562F0]
    focus-visible:outline-offset-[3px]

    disabled:pointer-events-none
    disabled:cursor-not-allowed

    motion-reduce:transition-none
  `;

  const sizeClasses =
    size === "small"
      ? `
          h-[38px]
          px-5
          py-2.5
          rounded-lg
          text-sm
        `
      : `
          h-12
          px-8
          py-3.5
          rounded-xl
          text-base
        `;

  const variantClasses =
    resolvedVariant === "primary"
      ? `
          border-[#4562F0]
          bg-[#4562F0]
          text-white

          hover:border-[#4F71FC]
          hover:bg-[#4F71FC]

          active:border-[#374ECC]
          active:bg-[#374ECC]

          disabled:border-[#4562F0]/40
          disabled:bg-[#4562F0]/40
          disabled:text-white
        `
      : `
          border-[#4562F0]
          bg-transparent
          text-[#4562F0]

          hover:border-[#4F71FC]
          hover:bg-[#4562F0]/5
          hover:text-[#4F71FC]

          active:border-[#374ECC]
          active:bg-[#4562F0]/15
          active:text-[#374ECC]

          disabled:border-[#4562F0]/40
          disabled:bg-transparent
    disabled:text-[#4562F0]/40
        `;

  const classes = `
    ${baseClasses}
    ${sizeClasses}
    ${variantClasses}
    max-[419px]:h-auto
    max-[419px]:min-h-12
    max-[419px]:whitespace-normal
    max-[419px]:text-center
    max-[419px]:leading-[1.25]
    ${className}
  `;

  if (link && !disabled) {
    return (
      <Link
        href={link}
        target={target}
        aria-label={ariaLabel}
        onClick={onClick}
        className={classes}
      >
        {text}
      </Link>
    );
  }

  return (
    <button
      type={type}
      disabled={disabled}
      aria-label={ariaLabel}
      onClick={onClick}
      className={classes}
    >
      {text}
    </button>
  );
}
