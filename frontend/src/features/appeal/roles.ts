export const appealRoles = [
  { id: "student", title: "Я школьник", description: "Мне нужна помощь или совет", image: "/images/roles/student.png" },
  { id: "parent", title: "Я родитель", description: "Хочу получить совет о ситуации с ребёнком", image: "/images/roles/parent.png" },
  { id: "teacher", title: "Я педагог", description: "Нужна консультация по рабочей ситуации", image: "/images/roles/ticher.png" },
] as const;

export type AppealRoleId = (typeof appealRoles)[number]["id"];

export function getAppealRole(value: string | string[] | undefined) {
  return appealRoles.find((role) => role.id === value) ?? appealRoles[0];
}

export function isFormalAppealRole(role: AppealRoleId | string) {
  return role !== "student";
}
