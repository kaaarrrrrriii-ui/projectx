"use client";

import Button from "@/shared/ui/button";
import Surface from "@/shared/ui/surface";
import Link from "next/link";
import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { loginStaff } from "@/shared/api/staff-api";

export default function EmployeeLogin() {
  const router = useRouter();
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = new FormData(event.currentTarget);
    setLoading(true);
    setError("");
    try {
      const user = await loginStaff(String(form.get("username") ?? ""), String(form.get("password") ?? ""));
      router.push(`/${user.role}`);
    } catch (reason) {
      setError(reason instanceof Error && reason.message !== "invalid username or password" ? reason.message : "Неверный логин или пароль");
    } finally { setLoading(false); }
  }

  return (
    <Surface
      as="section"
      labelledBy="employee-login-heading"
      className="relative w-full max-w-[800px] px-9 pt-7 pb-11 [--surface-radius:15px] [--surface-shadow:0_18px_70px_rgba(0,8,40,0.12)] !border-[#dee7fd] !bg-white max-[599px]:rounded-b-none max-[599px]:px-5 max-[599px]:pt-6 max-[599px]:pb-8"
    >
      <Link
        href="/"
        aria-label="Закрыть страницу входа"
        className="absolute top-4 right-5 flex h-10 w-10 items-center justify-center rounded-full text-[34px] leading-none font-light text-[var(--color-primary)] transition-colors hover:bg-[#eef1ff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--color-primary)]"
      >
        ×
      </Link>

      <h1
        id="employee-login-heading"
        className="pr-12 text-[28px] leading-tight font-extrabold tracking-[-0.025em] text-[var(--color-primary)]"
      >
        Вход
      </h1>

      <form onSubmit={submit} className="mt-8">
        <label className="block">
          <span className="text-[20px] leading-6 font-medium text-[#151515]">
            Логин
          </span>
          <input
            type="text"
            name="username"
            required
            autoComplete="username"
            placeholder="Логин"
            className="mt-2 block h-[44px] w-full rounded-[14px] border border-[#000828] bg-[#fcfdff] px-4 text-center text-[15px] text-[#000828] outline-none transition-[border-color,box-shadow,background-color] placeholder:text-[#646d86] hover:border-[var(--color-primary)] focus:border-[var(--color-primary)] focus:bg-white focus:ring-2 focus:ring-[#4562f0]/20"
          />
        </label>

        <label className="mt-8 block">
          <span className="text-[20px] leading-6 font-medium text-[#151515]">
            Пароль
          </span>
          <input
            type="password"
            name="password"
            required
            autoComplete="current-password"
            placeholder="Пароль"
            className="mt-2 block h-[44px] w-full rounded-[14px] border border-[#000828] bg-[#fcfdff] px-4 text-center text-[15px] text-[#000828] outline-none transition-[border-color,box-shadow,background-color] placeholder:text-[#646d86] hover:border-[var(--color-primary)] focus:border-[var(--color-primary)] focus:bg-white focus:ring-2 focus:ring-[#4562f0]/20"
          />
        </label>

        <div className="h-8" aria-hidden="true" />
        {error && <p role="alert" className="mb-3 text-sm text-[#d70d14]">{error}</p>}

        <Button
          text="Войти"
          type="submit"
          disabled={loading}
          variant="primary"
          size="small"
          className="w-full !h-[45px] !rounded-[10px] font-normal"
        />
      </form>
    </Surface>
  );
}
