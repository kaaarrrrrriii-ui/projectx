import type { SVGProps } from "react";

export type IconName = "lock" | "support" | "bell" | "chevron-left" | "chevron-right";

export default function Icon({ name, ...props }: SVGProps<SVGSVGElement> & { name: IconName }) {
  return (
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor"
      strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true" {...props}>
      {name === "lock" && <>
        <rect x="5" y="10" width="14" height="11" rx="3" />
        <path d="M8 10V7a4 4 0 0 1 8 0v3M12 14v3" />
      </>}
      {name === "support" && <>
        <circle cx="10" cy="8" r="4" />
        <path d="M3 21v-2a7 7 0 0 1 14 0v2M18 5a7 7 0 0 1 0 8M21 3a10 10 0 0 1 0 12" />
      </>}
      {name === "bell" && <>
        <path d="M18 8a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9M10 21h4M2 7a10 10 0 0 1 2-4M22 7a10 10 0 0 0-2-4" />
      </>}
      {name === "chevron-left" && <path d="m14 6-6 6 6 6" />}
      {name === "chevron-right" && <path d="m9 6 6 6-6 6" />}
    </svg>
  );
}
