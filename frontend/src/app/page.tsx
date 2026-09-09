import HomeInformation from "@/widgets/home-information/home-information";
import SiteHeader from "@/widgets/site-header/site-header";
import styles from "./home.module.css";

export default function Home() {
  return (
    <main className={styles.page}>
      <SiteHeader />
      <div className={styles.layout}>
        <div className={styles.information}>
          <HomeInformation />
        </div>
      </div>
    </main>
  );
}
