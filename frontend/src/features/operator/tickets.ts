export type TicketPriority = "urgent" | "standard" | "low";

export type OperatorTicket = {
  track: string;
  status: string;
  category: string;
  applicant: string;
  waiting: string;
  priority: TicketPriority;
  submittedAt: string;
  description: string;
  clarifications: Array<{ question: string; answer: string }>;
  attachments: string[];
};

export const operatorTickets: OperatorTicket[] = [
  {
    track: "ОТКЛ-2471-9382",
    status: "new",
    category: "кибербуллинг",
    applicant: "schoolchild",
    waiting: "12 минут",
    priority: "urgent",
    submittedAt: "09.09.2026, 14:48",
    description:
      "Мне пишут обидные сообщения в чате класса и выкладывают мои фотографии. Я боюсь идти в школу и не знаю, как это остановить.",
    clarifications: [
      { question: "Где это происходит?", answer: "Онлайн" },
      { question: "Как давно это длится?", answer: "Недавно" },
      { question: "Уже обращался за помощью?", answer: "Нет" },
    ],
    attachments: ["скриншот-чата-1.png", "скриншот-чата-2.png"],
  },
  {
    track: "ОТКЛ-5138-2047",
    status: "assigned",
    category: "конфликт с родителями",
    applicant: "parent",
    waiting: "34 минуты",
    priority: "standard",
    submittedAt: "09.09.2026, 14:26",
    description:
      "Мы с ребёнком перестали понимать друг друга. Любой разговор заканчивается ссорой. Хочу понять, как вернуть доверие.",
    clarifications: [
      { question: "Как давно это длится?", answer: "Давно" },
      { question: "Уже обращались за помощью?", answer: "Нет" },
    ],
    attachments: [],
  },
  {
    track: "ОТКЛ-9054-6713",
    status: "clarification",
    category: "давление и угрозы",
    applicant: "student",
    waiting: "1 час 08 минут",
    priority: "low",
    submittedAt: "09.09.2026, 13:52",
    description:
      "Одногруппник постоянно давит на меня и угрожает испортить мои вещи. Пока не решаюсь рассказать об этом преподавателю.",
    clarifications: [
      { question: "Где это происходит?", answer: "В учебном заведении" },
      { question: "Как давно это длится?", answer: "Пару дней" },
    ],
    attachments: ["фото-сообщения.jpg"],
  },
];

export const applicantLabels: Record<string, string> = {
  schoolchild: "Школьник",
  parent: "Родитель",
  student: "Студент",
};

export function getOperatorTicket(track: string) {
  return operatorTickets.find((ticket) => ticket.track === track);
}
