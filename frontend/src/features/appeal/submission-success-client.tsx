"use client";

import { useEffect, useState } from "react";
import Button from "@/shared/ui/button";
import type { CreatedTicket } from "@/shared/api/public-api";
import { readSubmissionResult } from "./appeal-draft";
import SubmissionSuccess from "./submission-success";

export default function SubmissionSuccessClient({ formal }: { formal: boolean }) {
  const [result, setResult] = useState<CreatedTicket | null>(null);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      setResult(readSubmissionResult());
      setLoaded(true);
    }, 0);
    return () => window.clearTimeout(timer);
  }, []);

  if (!loaded) return <div className="flex flex-1 items-center justify-center text-[#646d86]">Загружаем обращение…</div>;
  if (!result) {
    return (
      <section className="mx-auto flex w-full max-w-[720px] flex-1 flex-col items-center justify-center px-5 text-center">
        <h1 className="text-2xl font-bold text-[#4562f0]">Трек-номер не найден</h1>
        <p className="mt-3 text-[#646d86]">Вернитесь к форме и отправьте обращение ещё раз.</p>
        <Button text="Подать обращение" link="/appeal" variant="primary" className="mt-6" />
      </section>
    );
  }
  return <SubmissionSuccess trackNumber={result.track_id} formal={formal} crisisContacts={result.crisis_contacts} />;
}
