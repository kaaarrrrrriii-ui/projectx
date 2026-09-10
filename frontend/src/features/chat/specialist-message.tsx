import type { ReactNode } from "react";

export default function SpecialistMessage({ children }: { children: ReactNode }) {
  return (
    <article
      aria-label="Сообщение специалиста"
      className="overflow-hidden rounded-[15px] border border-[#dee7fd] bg-white shadow-[0_4px_12px_rgba(69,98,240,0.05)] max-[479px]:rounded-[13px]"
    >
      <h2 className="flex min-h-[72px] items-center px-6 text-[24px] leading-tight font-extrabold text-[var(--color-primary)] max-[699px]:min-h-[56px] max-[699px]:px-4 max-[699px]:text-[19px]">
        Специалист
      </h2>
      <div className="min-h-[310px] bg-[#dee7fd] px-8 py-6 text-[17px] leading-[1.18] text-[#11131a] [&>p]:m-0 [&>p]:max-w-[980px] max-[799px]:min-h-0 max-[799px]:px-5 max-[799px]:text-[15px] max-[799px]:leading-[1.4] max-[479px]:px-4 max-[479px]:py-4 max-[479px]:text-[14px]">
        {children}
      </div>
    </article>
  );
}
