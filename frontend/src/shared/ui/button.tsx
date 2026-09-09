"use client";

import Link from "next/link";
import type { CSSProperties } from "react";
import styles from "./button.module.css";

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
  text, color, padSmall, padMedium, padLarge, fill, textColor,
  link, onClick, type = "button", disabled = false, target = "_self",
  ariaLabel, className = "", style,
}: Props) {
  const customStyle = {
    ...(color && { "--button-accent": color }),
    ...(textColor && { "--button-text": textColor }),
    ...(padSmall && { padding: "5px 24px" }),
    ...(padMedium && { padding: "10px 32px" }),
    ...(padLarge && { padding: "20px 40px" }),
    ...style,
  } as CSSProperties;
  const classes = [styles.button, fill ? styles.primary : styles.secondary, className].join(" ");

  if (link && !disabled) {
    return (
      <Link href={link} className={classes} style={customStyle} target={target}
        aria-label={ariaLabel} onClick={onClick}>
        {text}
      </Link>
    );
  }

  return (
    <button type={type} className={classes} style={customStyle} onClick={onClick}
      disabled={disabled} aria-label={ariaLabel}>
      {text}
    </button>
  );
}
