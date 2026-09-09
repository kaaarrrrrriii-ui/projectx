"use client";

import Link from "next/link";
import type { CSSProperties } from "react";

interface Props {
  text: string;
  color?: string;
  padSmall?: boolean;
  padMedium?: boolean;
  padLarge?: boolean;
  fill?: boolean;
  textColor?: string;
  link?: string;
  onClick?: () => void;
  type?: "button" | "submit" | "reset";
  disabled?: boolean;
  target?: "_blank" | "_self";
  ariaLabel?: string;
  className?: string;
  style?: CSSProperties;
}

export default function Button({
  text,
  color,
  padSmall,
  padMedium,
  padLarge,
  fill,
  textColor,
  link,
  onClick,
  type = "button",
  disabled = false,
  target = "_self",
  ariaLabel,
  className = "",
  style,
}: Props) {
  const customStyle = {
    ...(color && { "--button-accent": color }),
    ...(textColor && { "--button-text": textColor }),
    ...(padSmall && { padding: "5px 24px" }),
    ...(padMedium && { padding: "10px 32px" }),
    ...(padLarge && { padding: "20px 40px" }),
    ...style,
  } as CSSProperties;

  const baseClasses = `
    inline-flex min-w-0
    items-center justify-center

    min-h-[var(--button-height,54px)]
    px-[20px] py-[14px]
    [padding:var(--button-padding,14px_20px)]

    rounded-[var(--button-radius,var(--radius-control))]
    border border-transparent

    text-center
    text-[length:var(--button-font-size,17px)]
    font-[var(--button-weight,600)]
    leading-[1.3]
    no-underline

    cursor-pointer

    transition-[background-color,border-color]
    duration-150

    focus-visible:outline
    focus-visible:outline-3
    focus-visible:outline-[var(--color-primary)]
    focus-visible:outline-offset-4

    disabled:cursor-not-allowed
    disabled:border-[var(--color-border)]
    disabled:bg-[var(--color-primary-disabled)]
    disabled:text-[var(--color-on-primary)]

    motion-reduce:transition-none
  `;

  const variantClasses = fill
    ? `
        bg-[var(--button-accent,var(--color-primary))]
        text-[var(--button-text,var(--color-on-primary))]

        hover:bg-[var(--color-primary-hover)]
        active:bg-[var(--color-primary-active)]
      `
    : `
        border-[var(--button-accent,var(--color-primary))]
        bg-[var(--color-surface)]
        text-[var(--button-text,var(--button-accent,var(--color-primary)))]

        hover:bg-[var(--color-primary-soft-hover)]
        active:bg-[var(--color-primary-soft-active)]
      `;

  const classes = `${baseClasses} ${variantClasses} ${className}`;

  if (link && !disabled) {
    return (
      <Link
        href={link}
        className={classes}
        style={customStyle}
        target={target}
        aria-label={ariaLabel}
        onClick={onClick}
      >
        {text}
      </Link>
    );
  }

  return (
    <button
      type={type}
      className={classes}
      style={customStyle}
      onClick={onClick}
      disabled={disabled}
      aria-label={ariaLabel}
    >
      {text}
    </button>
  );
}