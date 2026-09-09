import Button from "@/shared/ui/button";
import Icon, { type IconName } from "@/shared/ui/icon";
import InfoCard from "@/shared/ui/info-card";
import PageHeading from "@/shared/ui/page-heading";
import Surface from "@/shared/ui/surface";
import styles from "./home-information.module.css";

const information: { title: string; description: string; icon: IconName }[] = [
  {
    title: "Ты остаёшься анонимным",
    description: "Мы не спрашиваем твоё имя и не передаём информацию третьим лицам.",
    icon: "lock",
  },
  {
    title: "Тебе помогут специалисты",
    description: "Твоё обращение увидят психологи, педагоги и другие эксперты.",
    icon: "support",
  },
  {
    title: "Ты можешь отслеживать статус",
    description: "После отправки ты получишь уникальный трек-номер и сможешь проверить, на каком этапе находится обращение.",
    icon: "bell",
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
    <Surface as="section" labelledBy="important-heading" className={styles.panel}>
      <PageHeading id="important-heading" title="Что важно знать?" />
      <div className={styles.cards}>
        {information.map(({ icon, ...card }) => <InfoCard key={icon} {...card} icon={<Icon name={icon} />} />)}
      </div>
      <HomeActions />
    </Surface>
  );
}
