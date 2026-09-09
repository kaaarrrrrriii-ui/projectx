import Button from "@/shared/ui/button";
import { getAppealRole, isFormalAppealRole } from "@/features/appeal/roles";
import Input from "@/shared/ui/input";
import SiteHeader from "@/widgets/site-header/site-header";

export default async function Status({
  searchParams,
}: {
  searchParams: Promise<{ role?: string | string[] }>;
}) {
  const role = getAppealRole((await searchParams).role);
  const formal = isFormalAppealRole(role.id);
  const trackNumber: string = "сиксевен-6767-6767";
  const status: string = "послан нахуй";
  const sentDate: string = "09/09/09";
  const category: string = " кибербулинг олега";

  return (
    <main className="flex min-h-dvh flex-col bg-[var(--color-background)]">
      <SiteHeader formal={formal} />

      <section
        className="
          flex flex-1 flex-col justify-center
          bg-[radial-gradient(ellipse_at_0%_30%,#eef3ff_0%,transparent_32%),radial-gradient(ellipse_at_100%_85%,#eef2ff_0%,transparent_38%)]
          px-5 py-2
          sm:px-[5%]
        "
      >
        <div className="mx-auto flex w-full max-w-[1440px] flex-col">
          <section
            aria-labelledby="status-check-heading"
            className="
              rounded-[12px]
              border border-[#4562F0]
              bg-white/80
              px-[18px] py-4
              sm:px-[26px]
              sm:py-[18px]
            "
          >
            <h1
              id="status-check-heading"
              className="
                text-[32px]
                font-extrabold
                leading-[1.3]
                text-[#4562F0]
                max-[699px]:text-[20px]
              "
            >
              Проверка статуса
            </h1>

            <form
              className="
                mt-4
                grid grid-cols-1
                items-end gap-5
                sm:grid-cols-[1fr_1fr]
                sm:gap-8
              "
            >
              <input type="hidden" name="role" value={role.id} />

              <label className="flex min-w-0 flex-col gap-2">
                <span className="text-[20px] font-medium text-[#000828]">
                  Номер обращения
                </span>

                <Input
                  name="trackCode"
                  type="text"
                  placeholder="Введите номер обращения"
                  className="!w-full"
                />
              </label>

              <Button
                text="Проверить"
                type="submit"
                variant="primary"
                size="small"
                className="w-full"
              />
            </form>
          </section>

          <section
            aria-labelledby="appeal-status-heading"
            className="
              mt-4
              rounded-[12px]
              border border-[#4562F0]
              bg-white/80
              px-[18px] py-5
              sm:px-[26px]
              sm:py-[18px]
            "
          >
            <header
              className="
                flex flex-col gap-2
                sm:flex-row
                sm:items-center
                sm:justify-between
              "
            >
              <h2
                id="appeal-status-heading"
                className="
                  text-[32px]
                  font-extrabold
                  text-[#4562F0]
                  max-[699px]:text-[20px]
                "
              >
                Моё обращение
              </h2>

              <span className="text-[20px] font-medium text-[#4562F0]">
                {trackNumber}
              </span>
            </header>

            <div className="mt-7">
              <p
                className="
                  text-[24px]
                  font-extrabold
                  leading-[1.4]
                  text-black
                  max-[699px]:text-[18px]
                "
              >
                Статус: {status}
              </p>

              <div
                className="
                  mt-7 flex flex-col gap-1
                  text-[20px]
                  font-medium
                  leading-[1.5]
                  text-[#151515]
                "
              >
                <p>
                  Дата отправки: {sentDate}
                </p>

                <p>
                  Категория: {category}
                </p>
              </div>
            </div>

            <div
              className="
                mt-7
                grid grid-cols-1 gap-3
                sm:grid-cols-2
                sm:gap-[46px]
              "
            >
              <Button
                text="Закрыть обращение"
                variant="secondary"
                size="default"
                className="w-full"
              />

              <Button
                text="Перейти к ответу специалиста"
                variant="primary"
                size="default"
                link={`/chat?role=${role.id}`}
                className="w-full"
              />
            </div>
          </section>

          <div className="flex justify-center py-4">
            <Button
              text="Вернуться на главную страницу"
              variant="primary"
              size="default"
              link="/"
            />
          </div>
        </div>
      </section>
    </main>
  );
}
