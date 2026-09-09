import HomeInformation from "@/widgets/home-information/home-information";
import Image from "next/image";
import SiteHeader from "@/widgets/site-header/site-header";
import styles from "./home.module.css";

export default function Home() {
  return (
    <main className={styles.page}>
      <SiteHeader />
      <div className={styles.layout}>
        <Image src="/images/welcomeIMG.png" alt="" width={407} height={636}
          priority sizes="(max-width: 699px) 240px, 35vw" className={styles.illustration} />
        <HomeInformation />
      </div>
    </main>
  );
}
