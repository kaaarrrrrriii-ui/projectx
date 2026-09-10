import { appealCategories } from "@/features/appeal/categories";

export const adminCategories = appealCategories.filter(
  (category) => !category.toLocaleLowerCase("ru").startsWith("не знаю"),
);

export const specialistOptions = [
  "Психолог",
  "Конфликтолог",
  "Юрист",
  "Социальный педагог",
] as const;

export type SpecialistOption = (typeof specialistOptions)[number];
