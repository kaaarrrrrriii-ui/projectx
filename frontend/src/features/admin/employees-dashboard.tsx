"use client";

import DialogShell from "@/features/chat/dialog-shell";
import Button from "@/shared/ui/button";
import Input from "@/shared/ui/input";
import Link from "next/link";
import { useMemo, useState, type FormEvent } from "react";
import { adminCategories } from "./admin-options";

type EmployeeRole = "operator" | "expert";
type EmployeeFilter = "all" | EmployeeRole;

type Employee = {
  id: number;
  fullName: string;
  role: EmployeeRole;
  category: string | null;
  login: string;
  rating?: number;
};

type NewEmployee = Omit<Employee, "id" | "rating">;

const initialEmployees: Employee[] = [
  { id: 1, fullName: "Байкова Е.С.", role: "expert", category: "кибербуллинг", login: "e.baikova", rating: 5 },
  { id: 2, fullName: "Смирнов А.В.", role: "operator", category: null, login: "a.smirnov" },
  { id: 3, fullName: "Петрова А.В.", role: "expert", category: "травля и оскорбления", login: "a.petrova", rating: 4.9 },
  { id: 4, fullName: "Соколова М.А.", role: "expert", category: "конфликт с родителями", login: "m.sokolova", rating: 4.8 },
  { id: 5, fullName: "Орлов Д.С.", role: "operator", category: null, login: "d.orlov" },
  { id: 6, fullName: "Воронов И.М.", role: "expert", category: "юридический вопрос", login: "i.voronov", rating: 4.7 },
];

const roleLabels: Record<EmployeeRole, string> = {
  operator: "Оператор",
  expert: "Эксперт",
};

const filters: Array<{ value: EmployeeFilter; label: string }> = [
  { value: "all", label: "Все" },
  { value: "operator", label: "Операторы" },
  { value: "expert", label: "Эксперты" },
];

const controlClass =
  "h-[42px] w-full rounded-[12px] border border-[#000828] bg-[#fcfdff] px-3 text-sm text-[#000828] outline-none transition-colors hover:border-[#4562f0] focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/20 disabled:cursor-not-allowed disabled:border-[#d8ddea] disabled:bg-[#f1f3f8] disabled:text-[#9199af]";

function EyeIcon({ hidden }: { hidden: boolean }) {
  return hidden ? (
    <svg viewBox="0 0 24 24" aria-hidden="true" className="h-5 w-5" fill="none">
      <path d="M3 3 21 21" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
      <path d="M10.6 5.2c.5-.1.9-.2 1.4-.2 5.7 0 9 7 9 7a17 17 0 0 1-2.5 3.4M6.2 6.3C4.1 8 3 12 3 12s3.3 7 9 7c1.7 0 3.2-.6 4.5-1.5" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M9.9 9.9A3 3 0 0 0 14.1 14" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
    </svg>
  ) : (
    <svg viewBox="0 0 24 24" aria-hidden="true" className="h-5 w-5" fill="none">
      <path d="M3 12s3.3-7 9-7 9 7 9 7-3.3 7-9 7-9-7-9-7Z" stroke="currentColor" strokeWidth="1.8" strokeLinejoin="round" />
      <circle cx="12" cy="12" r="3" stroke="currentColor" strokeWidth="1.8" />
    </svg>
  );
}

function AddEmployeeModal({
  existingLogins,
  onClose,
  onAdd,
}: {
  existingLogins: string[];
  onClose: () => void;
  onAdd: (employee: NewEmployee) => void;
}) {
  const [fullName, setFullName] = useState("");
  const [role, setRole] = useState<EmployeeRole | "">("");
  const [category, setCategory] = useState("");
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");

  function handleRoleChange(nextRole: EmployeeRole | "") {
    setRole(nextRole);
    if (nextRole !== "expert") setCategory("");
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const normalizedName = fullName.trim();
    const normalizedLogin = login.trim();

    if (!normalizedName || !role || !normalizedLogin || !password.trim()) {
      setError("Заполните все обязательные поля.");
      return;
    }
    if (role === "expert" && !category) {
      setError("Выберите категорию для эксперта.");
      return;
    }
    if (existingLogins.some((item) => item.toLocaleLowerCase("ru") === normalizedLogin.toLocaleLowerCase("ru"))) {
      setError("Сотрудник с таким логином уже существует.");
      return;
    }

    onAdd({
      fullName: normalizedName,
      role,
      category: role === "expert" ? category : null,
      login: normalizedLogin,
    });
  }

  return (
    <DialogShell labelledBy="new-employee-heading" onClose={onClose} showClose className="max-w-[750px]">
      <form onSubmit={handleSubmit} noValidate>
        <h2 id="new-employee-heading" className="pr-12 text-[28px] leading-tight font-extrabold text-[#4562f0] max-[499px]:text-[23px]">Новый сотрудник</h2>

        <div className="mt-6 grid gap-5">
          <label className="grid gap-2 text-base font-medium text-[#151515]">
            Фамилия ИО <span className="sr-only">обязательное поле</span>
            <Input value={fullName} onChange={(event) => setFullName(event.target.value)} required autoFocus autoComplete="name" aria-required="true" className="!h-[42px] !w-full !rounded-[12px] !bg-[#fcfdff] !px-3 !text-left !text-sm focus:!bg-white focus:!text-[#000828]" />
          </label>

          <label className="grid gap-2 text-base font-medium text-[#151515]">
            Роль <span className="sr-only">обязательное поле</span>
            <select value={role} onChange={(event) => handleRoleChange(event.target.value as EmployeeRole | "")} required aria-required="true" className={controlClass}>
              <option value="" disabled>Выберите роль</option>
              <option value="expert">Эксперт</option>
              <option value="operator">Оператор</option>
            </select>
          </label>

          <label className="grid gap-2 text-base font-medium text-[#151515]">
            {role === "expert" ? "Категория *" : "Категория"}
            <select value={category} onChange={(event) => setCategory(event.target.value)} required={role === "expert"} disabled={role !== "expert"} aria-required={role === "expert"} aria-describedby="category-help" className={controlClass}>
              <option value="" disabled>{role === "operator" ? "Операторам категория не назначается" : "Выберите категорию"}</option>
              {adminCategories.map((item) => <option key={item} value={item}>{item.charAt(0).toUpperCase() + item.slice(1)}</option>)}
            </select>
            <span id="category-help" className="text-xs font-normal text-[#646d86]">
              {role === "expert" ? "Категория определяет сферу работы эксперта." : role === "operator" ? "Для роли оператора категория не используется." : "Сначала выберите роль сотрудника."}
            </span>
          </label>

          <div className="border-t border-[#b7bccb] pt-5">
            <label className="grid gap-2 text-base font-medium text-[#151515]">
              Логин <span className="sr-only">обязательное поле</span>
              <Input value={login} onChange={(event) => setLogin(event.target.value)} required autoComplete="username" aria-required="true" className="!h-[42px] !w-full !rounded-[12px] !bg-[#fcfdff] !px-3 !text-left !text-sm focus:!bg-white focus:!text-[#000828]" />
            </label>
          </div>

          <label className="grid gap-2 text-base font-medium text-[#151515]">
            Пароль <span className="sr-only">обязательное поле</span>
            <span className="relative block">
              <Input type={showPassword ? "text" : "password"} value={password} onChange={(event) => setPassword(event.target.value)} required autoComplete="new-password" aria-required="true" className="!h-[42px] !w-full !rounded-[12px] !bg-[#fcfdff] !px-3 !pr-12 !text-left !text-sm focus:!bg-white focus:!text-[#000828]" />
              <button type="button" onClick={() => setShowPassword((current) => !current)} aria-label={showPassword ? "Скрыть пароль" : "Показать пароль"} aria-pressed={showPassword} className="absolute top-1/2 right-2 flex h-9 w-9 -translate-y-1/2 cursor-pointer items-center justify-center rounded-lg text-[#646d86] hover:bg-[#eef1ff] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-[#4562f0]">
                <EyeIcon hidden={!showPassword} />
              </button>
            </span>
          </label>
        </div>

        {error && <p role="alert" className="mt-4 text-sm text-[#d70d14]">{error}</p>}
        <Button text="Создать" type="submit" variant="primary" className="mt-6 w-full" />
      </form>
    </DialogShell>
  );
}

export default function EmployeesDashboard() {
  const [employees, setEmployees] = useState(initialEmployees);
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState<EmployeeFilter>("all");
  const [isModalOpen, setIsModalOpen] = useState(false);

  const filteredEmployees = useMemo(() => {
    const query = search.trim().toLocaleLowerCase("ru");
    return employees.filter((employee) => {
      const matchesRole = filter === "all" || employee.role === filter;
      const matchesSearch = !query || employee.fullName.toLocaleLowerCase("ru").includes(query) || employee.login.toLocaleLowerCase("ru").includes(query) || employee.category?.toLocaleLowerCase("ru").includes(query);
      return matchesRole && matchesSearch;
    });
  }, [employees, filter, search]);

  function addEmployee(employee: NewEmployee) {
    setEmployees((current) => [...current, { ...employee, id: Math.max(0, ...current.map((item) => item.id)) + 1 }]);
    setIsModalOpen(false);
  }

  return (
    <section className="min-w-0 flex-1 px-5 pt-5 pb-8 sm:px-[26px]" aria-labelledby="employees-heading">
      <div className="mx-auto w-full max-w-[1180px]">
        <Link href="/admin" className="inline-flex rounded-sm text-sm text-[#85899b] hover:text-[#4562f0] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-[#4562f0]">Вернуться назад</Link>
        <h1 id="employees-heading" className="mt-7 text-[26px] leading-tight font-extrabold tracking-[-0.02em] text-[#4562f0]">Сотрудники</h1>

        <div className="mt-4 grid gap-4 md:grid-cols-2">
          <label className="relative block">
            <svg viewBox="0 0 24 24" aria-hidden="true" className="pointer-events-none absolute top-1/2 left-3 z-10 h-5 w-5 -translate-y-1/2 text-[#4562f0]" fill="none"><circle cx="10.5" cy="10.5" r="6.5" stroke="currentColor" strokeWidth="2" /><path d="m16 16 5 5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" /></svg>
            <Input type="search" value={search} onChange={(event) => setSearch(event.target.value)} aria-label="Поиск сотрудника" placeholder="Поиск сотрудника" className="!h-[42px] !w-full !rounded-[12px] !border-[#4562f0] !bg-white !pr-3 !pl-[42px] !text-left !text-sm placeholder:!text-[#9196a7] focus:!bg-white focus:!text-[#000828]" />
          </label>
          <Button text="Добавить сотрудника" variant="primary" onClick={() => setIsModalOpen(true)} className="!h-[42px] !w-full !rounded-[9px] !text-sm" />
        </div>

        <div className="mt-5 grid grid-cols-3" role="tablist" aria-label="Фильтр сотрудников по роли">
          {filters.map((item) => (
            <button key={item.value} type="button" role="tab" aria-selected={filter === item.value} onClick={() => setFilter(item.value)} className={`cursor-pointer border-b-2 px-3 py-3 text-sm transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0] ${filter === item.value ? "border-[#4562f0] font-medium text-[#4562f0]" : "border-transparent text-[#30384f] hover:border-[#b5c0f9] hover:text-[#4562f0]"}`}>
              {item.label}
            </button>
          ))}
        </div>

        <section aria-label="Список сотрудников" className="mt-3 overflow-x-auto rounded-[13px] border border-[#4562f0] bg-white/85">
          <table className="w-full min-w-[760px] table-fixed border-collapse">
            <thead className="bg-[#dfe6ff] text-[#4562f0]"><tr className="h-[52px]"><th scope="col" className="w-[27%] border-r border-[#4562f0] px-4 text-center font-normal">Фамилия ИО</th><th scope="col" className="w-[19%] border-r border-[#4562f0] px-4 text-center font-normal">Роль</th><th scope="col" className="w-[29%] border-r border-[#4562f0] px-4 text-center font-normal">Категория</th><th scope="col" className="w-[25%] px-4 text-center font-normal">Ссылка на аналитику</th></tr></thead>
            <tbody>
              {filteredEmployees.map((employee) => (
                <tr key={employee.id} className="h-[58px] border-t border-[#4562f0] hover:bg-[#f7f8ff]">
                  <td className="border-r border-[#4562f0] px-5 text-sm font-medium text-[#000828]">{employee.fullName}{employee.rating && <span className="ml-2 rounded-md bg-[#dfe6ff] px-1.5 py-0.5 text-[11px] font-normal text-[#4562f0]">{employee.rating.toFixed(1)}</span>}</td>
                  <td className="border-r border-[#4562f0] px-4 text-center text-sm text-[#30384f]">{roleLabels[employee.role]}</td>
                  <td className="border-r border-[#4562f0] px-4 text-center text-sm text-[#30384f]">{employee.role === "expert" && employee.category ? employee.category.charAt(0).toUpperCase() + employee.category.slice(1) : "—"}</td>
                  <td className="px-4 text-center"><Link href={`/admin/analytics?employee=${employee.id}`} className="inline-flex min-w-[105px] items-center justify-center rounded-full bg-[#4562f0] px-4 py-1 text-xs text-white transition-colors hover:bg-[#4f71fc] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]">Открыть</Link></td>
                </tr>
              ))}
              {filteredEmployees.length === 0 && <tr className="h-28 border-t border-[#4562f0]"><td colSpan={4} className="px-6 text-center text-sm text-[#646d86]">Сотрудники не найдены</td></tr>}
            </tbody>
          </table>
        </section>
        <p className="mt-3 text-right text-xs text-[#7a8298]" aria-live="polite">Найдено сотрудников: {filteredEmployees.length}</p>
      </div>

      {isModalOpen && <AddEmployeeModal existingLogins={employees.map((employee) => employee.login)} onClose={() => setIsModalOpen(false)} onAdd={addEmployee} />}
    </section>
  );
}
