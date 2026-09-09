export const appealRoles = [
  { id: "student", title: "Я школьник", description: "Мне нужна помощь или совет", image: "/images/roles/student.webp" },
  { id: "parent", title: "Я родитель", description: "Хочу получить совет о ситуации с ребёнком", image: "/images/roles/parent.webp" },
  { id: "teacher", title: "Я педагог", description: "Нужна консультация по рабочей ситуации", image: "/images/roles/teacher.webp" },
] as const;

export function getAppealRole(value: string | string[] | undefined) {
  return appealRoles.find((role) => role.id === value) ?? appealRoles[0];
}
