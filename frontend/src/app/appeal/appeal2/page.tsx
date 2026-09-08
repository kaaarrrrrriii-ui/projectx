import Button from "@/shared/ui/button";
export default async function appeal2() {
  return (
    <main className="min-h-screen flex items-center justify-center bg-white ">
      <div className="flex flex-col items-center justify-center gap-[30px]">
        <div className="">С чем это связано?</div>
        <div className="">Можно подобрать одну или несколько тем, которые ближе всего к твоей ситуации</div>
        <div className="">
            <div className=""></div>
        </div>
        <Button text="Продолжить" color="var(--primary-color)" padLarge fill textColor="#FFFFFF" link="/" />
      </div>
    </main>
  );
}
