"use client";

import Button from "@/shared/ui/button";
import Input from "@/shared/ui/input";
import CloseAppealFlow from "@/features/chat/close-appeal-flow";
import { FormEvent, useState } from "react";

export default function StatusChecker({
  roleId,
  formal,
}: {
  roleId: string;
  formal: boolean;
}) {
  const [trackNumber, setTrackNumber] = useState("");
  const [checkedTrackNumber, setCheckedTrackNumber] = useState("");
  const [isCloseDialogOpen, setIsCloseDialogOpen] = useState(false);

  function checkStatus(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const normalizedTrackNumber = trackNumber.trim();
    if (normalizedTrackNumber) setCheckedTrackNumber(normalizedTrackNumber);
  }

  return (
    <section className="flex flex-1 flex-col justify-start px-5 pt-10 pb-5 sm:justify-center sm:px-[5%] sm:py-2 max-[379px]:px-3 max-[379px]:pt-7">
      <div className="mx-auto flex w-full max-w-[1440px] flex-col">
        <section aria-labelledby="status-check-heading" className="rounded-[12px] border border-[#4562f0] bg-white/80 px-[18px] py-4 sm:px-[26px] sm:py-[18px]">
          <h1 id="status-check-heading" className="text-[32px] leading-[1.3] font-extrabold text-[#4562f0] max-[699px]:text-xl">
            Проверка статуса
          </h1>
          <form onSubmit={checkStatus} className="mt-4 grid grid-cols-1 items-end gap-5 sm:grid-cols-[1fr_1fr] sm:gap-8">
            <label className="flex min-w-0 flex-col gap-2">
              <span className="text-xl font-medium text-[#000828]">Номер обращения</span>
              <Input name="trackCode" type="text" value={trackNumber} onChange={(event) => { setTrackNumber(event.target.value); setCheckedTrackNumber(""); }} placeholder="Введите номер обращения" className="!w-full" />
            </label>
            <Button text="Проверить" type="submit" variant="primary" size="small" disabled={!trackNumber.trim()} className="w-full" />
          </form>
        </section>

        {checkedTrackNumber && (
          <section aria-labelledby="appeal-status-heading" aria-live="polite" className="mt-4 rounded-[12px] border border-[#4562f0] bg-white/80 px-[18px] py-5 sm:px-[26px] sm:py-[18px]">
            <header className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
              <h2 id="appeal-status-heading" className="text-[32px] font-extrabold text-[#4562f0] max-[699px]:text-xl">Моё обращение</h2>
              <span className="text-xl font-medium text-[#4562f0]">{checkedTrackNumber}</span>
            </header>
            <div className="mt-7">
              <p className="text-2xl leading-[1.4] font-extrabold text-black max-[699px]:text-lg">Статус: Ответ специалиста готов</p>
              <div className="mt-7 flex flex-col gap-1 text-xl leading-[1.5] font-medium text-[#151515] max-[699px]:text-base">
                <p>Дата отправки: 09.09.2026</p>
                <p>Категория: Кибербуллинг</p>
              </div>
            </div>
            <div className="mt-7 grid grid-cols-1 gap-3 sm:grid-cols-2 sm:gap-[46px]">
              <Button
                text="Закрыть обращение"
                variant="secondary"
                size="default"
                onClick={() => setIsCloseDialogOpen(true)}
                className="w-full"
              />
              <Button text="Перейти к ответу специалиста" variant="primary" size="default" link={`/chat?role=${roleId}`} className="w-full" />
            </div>
          </section>
        )}

        <div className="flex justify-center py-4">
          <Button text="Вернуться на главную страницу" variant="primary" size="default" link="/" className="w-full" />
        </div>
      </div>

      {isCloseDialogOpen && (
        <CloseAppealFlow
          formal={formal}
          trackNumber={checkedTrackNumber}
          outcome="helped"
          onCancel={() => setIsCloseDialogOpen(false)}
        />
      )}
    </section>
  );
}
