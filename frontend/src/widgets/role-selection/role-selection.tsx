import { appealRoles } from "@/features/appeal/roles";
import { appealRoutes } from "@/features/appeal/routes";
import RoleCard from "@/features/appeal/role-card";
import Button from "@/shared/ui/button";
import PageHeading from "@/shared/ui/page-heading";
import Link from "next/link";

export default function RoleSelection({
  selectedRole,
}: {
  selectedRole: (typeof appealRoles)[number]["id"];
}) {
  return (
    <section
      aria-labelledby="role-heading"
      className="
        flex w-full flex-1 flex-col
        px-[4.5%]
        pt-9
        pb-[38px]

        [--heading-size:28px]
        [--heading-color:var(--color-primary)]

        max-[699px]:px-5
        max-[699px]:pt-8
        max-[699px]:pb-[max(32px,env(safe-area-inset-bottom))]
        max-[699px]:[--heading-size:28px]
      "
    >
      <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col">
        <div className="rounded-[15px] border border-[var(--color-primary)] bg-[var(--color-background)] px-[30px] py-[50px]">
          <PageHeading
            id="role-heading"
            title="Кто вы?"
          />

          <form
            action={appealRoutes.topics}
            method="get"
            className="mt-9 flex w-full flex-col items-center gap-9 max-[699px]:mt-6 max-[699px]:gap-7"
          >
            <fieldset
              aria-labelledby="role-heading"
              className="
                grid w-full max-w-[900px] min-w-0
                grid-cols-3 items-stretch
                gap-5 border-0 p-0

                max-[699px]:grid-cols-1
                max-[699px]:gap-3.5
              "
            >
              <legend className="sr-only">
                Выбери свою роль
              </legend>

              {appealRoles.map((role) => (
                <RoleCard
                  key={role.id}
                  role={role}
                  selected={role.id === selectedRole}
                />
              ))}
            </fieldset>

            <Button
              text="Продолжить"
              type="submit"
              variant="primary"
              size="default"
              className="w-full max-[699px]:h-12 max-[699px]:rounded-xl max-[699px]:px-8 max-[699px]:py-3.5 max-[699px]:text-base"
            />
          </form>
        </div>

        <Link
          href="/"
          className="
            mt-auto w-fit rounded-[3px] pt-9
            text-sm leading-5 text-[#85899b]
            transition-colors hover:text-[#4562f0]
            focus-visible:outline-2 focus-visible:outline-offset-4
            focus-visible:outline-[#4562f0]
          "
        >
          Вернуться назад
        </Link>
      </div>
    </section>
  );
}
