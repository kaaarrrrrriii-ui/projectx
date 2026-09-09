import Button from "@/shared/ui/button";
import Image from "next/image";
import InfoCard from "@/shared/ui/info-card";
import PageHeading from "@/shared/ui/page-heading";
import styles from "./home-information.module.css";

const information = [
  {
    title: "Ты остаёшься анонимным",
    description: "Мы не спрашиваем твоё имя и не передаём информацию третьим лицам.",
    icon: "lock-key",
  },
  {
    title: "Тебе помогут специалисты",
    description: "Твоё обращение увидят психологи, педагоги и другие эксперты.",
    icon: "user-sound",
  },
  {
    title: "Ты можешь отслеживать статус",
    description: "После отправки ты получишь уникальный трек-номер и сможешь проверить, на каком этапе находится обращение.",
    icon: "bell-ringing",
  },
];

function HomeActions() {
  return (
    <div className={styles.actions}>
      <div className={styles.buttons}>
        <Button text="Подать обращение" fill link="/appeal" />
        <Button text="Проверить статус" link="/status" />
      </div>
      <span className={styles.staff}>Вход для сотрудников</span>
    </div>
  );
}

export default function HomeInformation() {
  return (
    <section aria-labelledby="important-heading" className={styles.panel}>
      <PageHeading id="important-heading" title="Что важно знать?" />
      <div className={styles.cards}>
        {information.map(({ icon, ...card }) => <InfoCard key={icon} {...card}
          icon={<Image src={`/images/${icon}.svg`} alt="" width={28} height={28} />} />)}
      </div>
      <HomeActions />
    </section>
  );
}
