export type ExpertPriority = "urgent" | "standard" | "low";

export type ExpertAppeal = {
  track: string;
  category: string;
  applicant: string;
  priority: ExpertPriority;
  status: string;
  waiting: string;
  submittedAt: string;
  description: string;
  clarifications: Array<{ question: string; answer: string }>;
  attachments: string[];
  group: "queue" | "assigned" | "return";
};

export type ExpertRequest = {
  track: string;
  category: string;
  applicant: string;
  priority: ExpertPriority;
  status: string;
};

export const expertAppeals: ExpertAppeal[] = [
  {
    track: "6767-6767-6767",
    category: "Семейные проблемы",
    applicant: "Ученик",
    priority: "urgent",
    status: "Новое",
    waiting: "12 минут",
    submittedAt: "13.03.2023, 16:39",
    description: "Меня били в детстве. Много раз. Отчим. Пять раз за день.",
    clarifications: [
      { question: "Как давно это длится?", answer: "Недавно" },
      { question: "Где это было?", answer: "В школе" },
    ],
    attachments: ["bvjdsfieo.png"],
    group: "queue",
  },
  {
    track: "2471-9382-1104",
    category: "Кибербуллинг",
    applicant: "Студент",
    priority: "standard",
    status: "Новое",
    waiting: "34 минуты",
    submittedAt: "10.09.2026, 12:14",
    description: "В общем чате публикуют мои фотографии и пишут обидные сообщения. Не знаю, как это остановить.",
    clarifications: [
      { question: "Где это происходит?", answer: "Онлайн" },
      { question: "Как давно это длится?", answer: "Пару недель" },
    ],
    attachments: ["скриншот-чата.png"],
    group: "queue",
  },
  {
    track: "5138-2047-8821",
    category: "Конфликт с родителями",
    applicant: "Родитель",
    priority: "standard",
    status: "В работе",
    waiting: "1 час 18 минут",
    submittedAt: "10.09.2026, 11:42",
    description: "Мы с ребёнком перестали понимать друг друга. Любой разговор заканчивается ссорой.",
    clarifications: [
      { question: "Как давно это длится?", answer: "Несколько месяцев" },
      { question: "Обращались ли вы за помощью?", answer: "Нет" },
    ],
    attachments: [],
    group: "assigned",
  },
  {
    track: "9054-6713-4302",
    category: "Давление и угрозы",
    applicant: "Студент",
    priority: "urgent",
    status: "Нужен повторный ответ",
    waiting: "2 часа 05 минут",
    submittedAt: "10.09.2026, 10:55",
    description: "Одногруппник продолжает угрожать мне после первого обращения. Мне нужна дополнительная помощь.",
    clarifications: [
      { question: "Появились ли новые угрозы?", answer: "Да" },
      { question: "Нужна ли срочная помощь?", answer: "Да" },
    ],
    attachments: ["новые-сообщения.jpg"],
    group: "return",
  },
  {
    track: "7314-5520-1846",
    category: "Конфликт с учителем",
    applicant: "Ученик",
    priority: "low",
    status: "В работе",
    waiting: "3 часа 21 минута",
    submittedAt: "10.09.2026, 09:39",
    description: "Мне трудно общаться с учителем, и я не понимаю, как спокойно обсудить ситуацию.",
    clarifications: [{ question: "Где это происходит?", answer: "В школе" }],
    attachments: [],
    group: "assigned",
  },
];

export const expertRequests: ExpertRequest[] = [
  {
    track: "6767-6767-6767",
    category: "Семейные проблемы",
    applicant: "Ученик",
    priority: "urgent",
    status: "Отправлен оператору",
  },
  {
    track: "2471-9382-1104",
    category: "Кибербуллинг",
    applicant: "Студент",
    priority: "standard",
    status: "Ожидает ответа",
  },
];

export function getExpertAppeal(track: string) {
  return expertAppeals.find((appeal) => appeal.track === track);
}
