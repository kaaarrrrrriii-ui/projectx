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
        px-[4.65%]
        pt-[42px]
        pb-[38px]

        [--heading-size:clamp(26px,2.8vw,40px)]
        [--heading-color:var(--color-primary)]

        max-[699px]:px-5
        max-[699px]:pt-[30px]
        max-[699px]:pb-[max(32px,env(safe-area-inset-bottom))]
        max-[699px]:[--heading-size:28px]
        max-[379px]:px-3
        max-[379px]:pt-6
        max-[379px]:[--heading-size:26px]
      "
    >
      <div className="mx-auto flex w-full max-w-[1440px] flex-1 flex-col gap-9 max-[699px]:max-w-[540px] max-[699px]:gap-6">
        <PageHeading
          id="role-heading"
          title="Кто вы?"
        />

        <form
          action={appealRoutes.topics}
          method="get"
          className="
            flex flex-col items-center gap-9
            max-[699px]:gap-7
          "
        >
        <fieldset
          aria-labelledby="role-heading"
          className="
            grid w-full max-w-[1080px] min-w-0
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
          className="
            min-w-[108px]
            max-[699px]:h-12
            max-[699px]:w-full
            max-[699px]:rounded-xl
            max-[699px]:px-8
            max-[699px]:py-3.5
            max-[699px]:text-base
          "
        />
        </form>

        <Link
          href="/"
          className="
            mt-auto w-fit rounded-[3px]
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
