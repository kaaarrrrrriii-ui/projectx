"use client";

import Button from "@/shared/ui/button";
import type { ReactNode } from "react";
import { useState } from "react";
import CloseAppealFlow from "./close-appeal-flow";

export default function ChatWorkspace({
  children,
  formal,
  trackNumber,
  status,
}: {
  children: ReactNode;
  formal: boolean;
  trackNumber: string;
  status: string;
}) {
  const [activeDialog, setActiveDialog] = useState<"helped" | "not-helped" | null>(null);

  return (
    <>
      <div className="mx-auto flex min-h-[calc(100dvh-100px)] w-full max-w-[1440px] flex-col px-5 pt-2 pb-16 max-[699px]:min-h-[calc(100dvh-82px)] max-[699px]:px-3 max-[699px]:pt-3 max-[699px]:pb-6">
        <div aria-live="polite">{children}</div>

        {status === "answer_ready" && <section
          aria-labelledby="helpfulness-heading"
          className="mx-[18px] mt-14 flex min-h-[164px] flex-col items-center justify-center rounded-[15px] border border-[#4562f0] bg-white/90 px-5 py-6 max-[699px]:mx-0 max-[699px]:mt-7"
        >
          <h2
            id="helpfulness-heading"
            className="text-center text-[clamp(18px,2vw,24px)] leading-[1.25] font-extrabold tracking-[-0.025em] text-[#4562f0]"
          >
            {formal
              ? "Мы помогли вам решить вашу проблему?"
              : "Мы помогли тебе решить твою проблему?"}
          </h2>
          <div className="mt-4 grid w-full max-w-[744px] grid-cols-2 gap-2.5 max-[479px]:mt-5">
            <Button
              text="Нет"
              variant="secondary"
              size="small"
              onClick={() => setActiveDialog("not-helped")}
              className="h-[38px] w-full rounded-[10px] font-normal"
            />
            <Button
              text="Да"
              variant="secondary"
              size="small"
              onClick={() => setActiveDialog("helped")}
              className="h-[38px] w-full rounded-[10px] font-normal"
            />
          </div>
        </section>}
      </div>

      {activeDialog && (
        <CloseAppealFlow
          formal={formal}
          trackNumber={trackNumber}
          outcome={activeDialog}
          onCancel={() => setActiveDialog(null)}
        />
      )}
    </>
  );
}
