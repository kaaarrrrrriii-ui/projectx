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
          py-5

          max-[699px]:flex-col
          max-[699px]:px-5
          max-[699px]:pt-[30px]
          max-[699px]:pb-0
        "
      >
        <Image
          src="/images/main-page-girl-hq.png"
          alt=""
          width={1402}
          height={1122}
          priority
          sizes="(max-width: 699px) 320px, 42vw"
          className="
            absolute top-1/2 left-[clamp(20px,3.1vw,58px)] z-0
            h-auto w-[clamp(400px,42.2vw,607px)]
            -translate-y-1/2
            transition-[filter] duration-300 ease-out
            hover:drop-shadow-[0_0_18px_rgba(69,98,240,0.28)]
            motion-reduce:transition-none

            max-[699px]:static
            max-[699px]:order-1
            max-[699px]:mx-auto
            max-[699px]:mt-5
            max-[699px]:w-[320px]
            max-[699px]:max-w-full
            max-[699px]:translate-y-0
          "
        />

        <HomeInformation />
      </div>
    </main>
  );
}
