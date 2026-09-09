import HomeInformation from "@/widgets/home-information/home-information";
import Image from "next/image";
import SiteHeader from "@/widgets/site-header/site-header";

export default function Home() {
  return (
    <main className="app-page-background flex min-h-dvh flex-col">
      <SiteHeader />

      <div
        className="
          relative isolate mx-auto flex w-full max-w-[1600px] flex-1
          items-center
          px-[clamp(32px,4.7vw,72px)]
          pt-[44px] pb-5

          max-[699px]:flex-col
          max-[699px]:px-5
          max-[699px]:pt-[30px]
          max-[699px]:pb-0
        "
      >
        <Image
          src="/images/welcomeIMG.png"
          alt=""
          width={407}
          height={636}
          priority
          sizes="(max-width: 699px) 240px, 35vw"
          className="
            pointer-events-none
            absolute bottom-0 left-0 -z-10
            h-auto w-[34.3%]

            max-[699px]:static
            max-[699px]:order-1
            max-[699px]:mt-[-50px]
            max-[699px]:w-[240px]
            max-[699px]:max-w-[80%]
          "
        />

        <HomeInformation />
      </div>
    </main>
  );
}
