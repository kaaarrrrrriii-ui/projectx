export type ClarifyingQuestionKey = "place" | "duration" | "askedForHelp";

export type ClarifyingQuestion = {
  id: string;
  key?: ClarifyingQuestionKey;
  title: string;
  options: string[];
};

export const clarifyingQuestions: ClarifyingQuestion[] = [
  {
    id: "place",
    key: "place",
    title: "Где это происходит?",
    options: ["В школе", "Онлайн", "В другом месте"],
  },
  {
    id: "duration",
    key: "duration",
    title: "Как давно это длится?",
    options: ["Недавно", "Давно", "Пару дней"],
  },
  {
    id: "askedForHelp",
    key: "askedForHelp",
    title: "Ты уже обращался за помощью?",
    options: ["Нет", "Да", "Не уверен"],
  },
];
