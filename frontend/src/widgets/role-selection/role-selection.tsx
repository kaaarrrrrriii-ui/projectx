import { appealRoles } from "@/features/appeal/roles";
import RoleCard from "@/features/appeal/role-card";
import Button from "@/shared/ui/button";
import PageHeading from "@/shared/ui/page-heading";

export default function RoleSelection({
  selectedRole,
}: {
  selectedRole: (typeof appealRoles)[number]["id"];
}) {
  return (
    <section
      aria-labelledby="role-heading"
      className="
        mx-auto flex w-full max-w-[1440px] flex-col gap-9
        px-[clamp(24px,5vw,72px)]
        pt-[42px]
        pb-[max(40px,env(safe-area-inset-bottom))]

        [--heading-size:clamp(26px,2.8vw,40px)]
        [--heading-color:var(--color-primary)]

        [--button-height:36px]
        [--button-font-size:clamp(11px,1.2vw,15px)]
        [--button-padding:8px_22px]
        [--button-radius:8px]
        [--button-weight:500]

        max-[699px]:max-w-[540px]
        max-[699px]:gap-6
        max-[699px]:px-5
        max-[699px]:pt-[30px]
        max-[699px]:pb-[max(28px,env(safe-area-inset-bottom))]
        max-[699px]:[--heading-size:28px]
        max-[699px]:[--button-height:48px]
        max-[699px]:[--button-font-size:14px]
      "
    >
      <PageHeading id="role-heading" title="Кто вы?" />

      <form
        action="/appeal/appeal2"
        method="get"
        className="
          flex flex-col items-center gap-9
          max-[699px]:gap-7
        "
      >
        <fieldset
          aria-labelledby="role-heading"
          className="
            grid w-full min-w-0
            grid-cols-3 items-stretch
            gap-3 border-0 p-0

            max-[699px]:grid-cols-1
            max-[699px]:gap-3.5
          "
        >
          <legend className="sr-only">Выбери свою роль</legend>

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
          fill
          className="
            min-w-[108px]
            max-[699px]:w-full
          "
        />
      </form>
    </section>
  );
}

