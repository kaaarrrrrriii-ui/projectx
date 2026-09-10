"use client";

import { useEffect, useState } from "react";
import { readAppealDraft, updateAppealDraft } from "./appeal-draft";

export default function DescriptionInput({ formal }: { formal: boolean }) {
  const [description, setDescription] = useState("");

  useEffect(() => {
    const timer = window.setTimeout(() => setDescription(readAppealDraft().description ?? ""), 0);
    return () => window.clearTimeout(timer);
  }, []);

  return (
    <textarea
      name="description"
      value={description}
      maxLength={20_000}
      onChange={(event) => {
        setDescription(event.target.value);
        updateAppealDraft({ description: event.target.value });
      }}
      aria-labelledby="description-heading"
      aria-describedby="description-intro description-hint"
      placeholder={formal ? "Здесь можно написать всё, что вас беспокоит..." : "Здесь можно написать всё, что тебя беспокоит..."}
      className="mt-7 block min-h-[218px] w-full resize-y rounded-[15px] border border-[#333] bg-[#FCFDFF] px-[13px] py-4 text-[13px] leading-5 text-[#151515] placeholder:text-[#85899f] placeholder:opacity-100 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#4562F0] max-[699px]:mt-6 max-[699px]:text-[16px]"
    />
  );
}
