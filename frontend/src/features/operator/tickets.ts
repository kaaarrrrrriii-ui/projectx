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

export const returnedTickets: OperatorTicket[] = [
  {
    track: "ОТКЛ-6767-0676",
    status: "returned",
    category: "конфликт с родителями",
    applicant: "schoolchild",
    waiting: "1 час",
    priority: "urgent",
    submittedAt: "09.09.2026, 12:44",
    description:
      "После первого ответа ситуация дома не изменилась. Мне всё ещё трудно разговаривать с родителями без ссор, поэтому я прошу помочь ещё раз.",
    clarifications: [
      { question: "Что не помогло в первом ответе?", answer: "Не получилось начать спокойный разговор" },
      { question: "Нужна ли повторная консультация?", answer: "Да" },
    ],
    attachments: [],
  },
  {
    track: "ОТКЛ-4821-3095",
    status: "returned",
    category: "кибербуллинг",
    applicant: "student",
    waiting: "2 часа 18 минут",
    priority: "standard",
    submittedAt: "09.09.2026, 11:26",
    description:
      "Оскорбительные сообщения продолжают приходить с новых аккаунтов. Предыдущих рекомендаций оказалось недостаточно, нужна дополнительная помощь.",
    clarifications: [
      { question: "Сохранились ли доказательства?", answer: "Да, есть новые скриншоты" },
      { question: "Обращались ли к администрации площадки?", answer: "Да" },
    ],
    attachments: ["новые-сообщения.png"],
  },
  {
    track: "ОТКЛ-7314-5520",
    status: "returned",
    category: "давление и угрозы",
    applicant: "parent",
    waiting: "3 часа 05 минут",
    priority: "low",
    submittedAt: "09.09.2026, 10:39",
    description:
      "После консультации появились новые обстоятельства. Хочу уточнить, как безопасно действовать дальше и к кому ещё можно обратиться.",
    clarifications: [
      { question: "Появились ли новые угрозы?", answer: "Нет, но давление продолжается" },
      { question: "Требуется ли срочная помощь?", answer: "Нет" },
    ],
    attachments: [],
  },
];

export const applicantLabels: Record<string, string> = {
  schoolchild: "Школьник",
  parent: "Родитель",
  student: "Студент",
};

export function getOperatorTicket(track: string) {
  return [...operatorTickets, ...returnedTickets].find((ticket) => ticket.track === track);
}
