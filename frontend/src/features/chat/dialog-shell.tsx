"use client";

import Surface from "@/shared/ui/surface";
import type { ReactNode } from "react";
import { useEffect } from "react";

export default function DialogShell({
  children,
  labelledBy,
  onClose,
  showClose = false,
  className = "",
}: {
  children: ReactNode;
  labelledBy: string;
  onClose?: () => void;
  showClose?: boolean;
  className?: string;
}) {
  useEffect(() => {
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape" && onClose) onClose();
    }

    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.body.style.overflow = previousOverflow;
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [onClose]);

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-[#11131a]/55 p-4 backdrop-blur-[1px] max-[499px]:items-end max-[499px]:p-0"
      role="dialog"
      aria-modal="true"
      aria-labelledby={labelledBy}
    >
      <Surface
        as="section"
        labelledBy={labelledBy}
        className={`relative max-h-[calc(100dvh-32px)] w-full overflow-y-auto p-7 [--surface-radius:15px] [--surface-shadow:0_18px_70px_rgba(0,8,40,0.2)] !border-[#dee7fd] !bg-white max-[499px]:max-h-[92dvh] max-[499px]:rounded-b-none max-[499px]:p-5 max-[379px]:p-4 ${className}`}
      >
        {showClose && onClose && (
          <button
            type="button"
            onClick={onClose}
            aria-label="Закрыть окно"
            className="absolute top-5 right-5 flex h-10 w-10 cursor-pointer items-center justify-center rounded-full text-[34px] leading-none font-light text-[var(--color-primary)] transition-colors hover:bg-[#eef1ff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--color-primary)] max-[499px]:top-3 max-[499px]:right-3 max-[379px]:h-9 max-[379px]:w-9"
          >
            ×
          </button>
        )}

        {children}
      </Surface>
    </div>
  );
}
