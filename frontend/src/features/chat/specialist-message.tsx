import type { ReactNode } from "react";

export default function SpecialistMessage({ children }: { children: ReactNode }) {
  return (
    <article
      aria-label="Сообщение специалиста"
      className="max-w-[980px] rounded-[15px] bg-[#dee7fd] px-6 py-4 text-[17px] leading-[1.18] text-[#11131a] [&>p]:m-0 max-[799px]:px-5 max-[799px]:text-[15px] max-[799px]:leading-[1.4] max-[479px]:rounded-[13px] max-[479px]:px-4 max-[479px]:py-4 max-[479px]:text-[14px]"
    >
      {children}
    </article>
  );
}
