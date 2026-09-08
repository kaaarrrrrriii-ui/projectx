import Button from "@/shared/ui/button";
export default async function appeal2() {
  return (
    <main className="min-h-screen flex items-center justify-center bg-white ">
      <div className="flex flex-col items-center justify-center gap-[30px]">
        <div className="">Кто вы?</div>
        <div className="">Это поможет нам подобрать нужного специалиста и говорить с тобой на одном языке</div>
        <div className="">
          <div className="">Я школьник</div>
          <div className="">Я родитель</div>
          <div className="">Я педагог</div>
        </div>
        <Button text="Продолжить" color="var(--primary-color)" padLarge fill textColor="#FFFFFF" link="/appeal/appeal2" />
      </div>
    </main>
  );
}
