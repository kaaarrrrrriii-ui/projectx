"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

const roles = [
  {
    id: "operator",
    title: "Оператор",
    description: "Принимает и распределяет новые обращения",
    icon: "headset",
  },
  {
    id: "admin",
    title: "Администратор",
    description: "Управляет командой, доступами и настройками",
    icon: "shield",
  },
  {
    id: "expert",
    title: "Эксперт",
    description: "Работает с направленными ему обращениями",
    icon: "expert",
  },
] as const;

type RoleId = (typeof roles)[number]["id"];

function RoleIcon({ name }: { name: (typeof roles)[number]["icon"] }) {
  if (name === "shield") {
    return (
      <svg viewBox="0 0 32 32" aria-hidden="true" className="h-8 w-8" fill="none">
        <path d="M16 3.5 27 8v7.2c0 6.6-4.6 11.1-11 13.3C9.6 26.3 5 21.8 5 15.2V8l11-4.5Z" stroke="currentColor" strokeWidth="2.2" />
        <path d="m11.5 16 3 3 6.5-7" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    );
  }

  if (name === "expert") {
    return (
      <svg viewBox="0 0 32 32" aria-hidden="true" className="h-8 w-8" fill="none">
        <circle cx="16" cy="11" r="5" stroke="currentColor" strokeWidth="2.2" />
        <path d="M6 27c.8-5.4 4.2-8 10-8s9.2 2.6 10 8" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" />
        <path d="m23 7 1 1 2-2" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    );
  }

  return (
    <svg viewBox="0 0 32 32" aria-hidden="true" className="h-8 w-8" fill="none">
      <path d="M7 18v-3a9 9 0 0 1 18 0v3" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" />
      <path d="M5 18.5A2.5 2.5 0 0 1 7.5 16H10v8H7.5A2.5 2.5 0 0 1 5 21.5v-3Zm22 0a2.5 2.5 0 0 0-2.5-2.5H22v8h2.5a2.5 2.5 0 0 0 2.5-2.5v-3Z" fill="currentColor" />
      <path d="M22 25c-1 2-2.9 3-5.5 3" stroke="currentColor" strokeWidth="2.2" strokeLinecap="round" />
    </svg>
  );
}

export default function EmployeeRoleSelection() {
  const router = useRouter();
  const [selectedRole, setSelectedRole] = useState<RoleId>("operator");

  function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (selectedRole === "operator") {
      router.push("/operator");
      return;
    }
    router.push(selectedRole === "expert" ? "/expert" : "/admin");
  }

  return (
    <form onSubmit={submit} className="w-full">
      <fieldset>
        <legend className="sr-only">Выберите роль сотрудника</legend>
        <div className="grid gap-4 md:grid-cols-3">
          {roles.map((role) => {
            const selected = selectedRole === role.id;
            return (
              <label
                key={role.id}
                className={`group relative min-h-52 cursor-pointer rounded-[20px] border bg-white p-6 transition-[border-color,box-shadow,transform] focus-within:ring-3 focus-within:ring-[#4562f0]/20 hover:-translate-y-0.5 hover:border-[#4562f0] ${selected ? "border-[#4562f0] shadow-[0_12px_30px_rgba(69,98,240,0.12)]" : "border-[#e0e5f5]"}`}
              >
                <input
                  type="radio"
                  name="employee-role"
                  value={role.id}
                  checked={selected}
                  onChange={() => setSelectedRole(role.id)}
                  className="sr-only"
                />
                <span className={`mb-5 flex h-14 w-14 items-center justify-center rounded-2xl ${selected ? "bg-[#4562f0] text-white" : "bg-[#eef1ff] text-[#4562f0] group-hover:bg-[#e4e8ff]"}`}>
                  <RoleIcon name={role.icon} />
                </span>
                <span className="block text-lg font-semibold text-[#000828]">{role.title}</span>
                <span className="mt-2 block text-sm leading-5 text-[#646d86]">{role.description}</span>
                <span className={`absolute right-5 top-5 flex h-5 w-5 items-center justify-center rounded-full border ${selected ? "border-[#4562f0]" : "border-[#cbd2e8]"}`}>
                  {selected && <span className="h-2.5 w-2.5 rounded-full bg-[#4562f0]" />}
                </span>
              </label>
            );
          })}
        </div>
      </fieldset>

      <div className="mt-8 flex flex-wrap items-center justify-center gap-4">
        <Link href="/" className="inline-flex h-12 items-center justify-center rounded-xl border border-[#4562f0] px-8 font-medium text-[#4562f0] transition-colors hover:bg-[#eef1ff]">
          Назад
        </Link>
        <button type="submit" className="inline-flex h-12 cursor-pointer items-center justify-center rounded-xl bg-[#4562f0] px-8 font-medium text-white transition-colors hover:bg-[#4f71fc] focus-visible:outline-3 focus-visible:outline-offset-3 focus-visible:outline-[#4562f0]">
          Войти
        </button>
      </div>
    </form>
  );
}
