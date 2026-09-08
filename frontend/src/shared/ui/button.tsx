"use client";
import Link from "next/link";

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
    style?: React.CSSProperties;
}

export default function Button({
    text,
    color = "var(--primary-color)",
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
    const padding = padSmall
        ? '5px 45px'
        : padMedium
            ? '10px 50px'
            : padLarge
                ? '20px 80px'
                : '10px 50px';
    const stateClass = disabled
        ? "bg-(--disabled-color) cursor-not-allowed"
        : fill
            ? "hover:bg-(--hoover-color) active:bg-(--pressed-color)"
            : "hover:bg-(--hoover-color-wborder) active:bg-(--pressed-color-wborder)";
    const commonClass = `inline-flex cursor-pointer items-center justify-center rounded-3xl text-xl transition-colors ${stateClass} ${className}`;
    const commonStyle: React.CSSProperties = {
        padding,
        backgroundColor: fill ? color : 'transparent',
        border: fill ? 'none' : `1px solid ${color}`,
        color: textColor ?? (fill ? 'var(--text-color-secondary)' : color),
        textDecoration: 'none',
        ...style,
    };

    if (link) {
        return (
            <Link
                href={link}
                className={commonClass}
                style={commonStyle}
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
            className={commonClass}
            style={commonStyle}
            onClick={onClick}
            disabled={disabled}
            aria-label={ariaLabel}
        >
            {text}
        </button>
    );
}
