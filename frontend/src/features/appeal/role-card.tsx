import Image from "next/image";
import type { appealRoles } from "./roles";

export default function RoleCard({
  role,
  selected,
}: {
  role: (typeof appealRoles)[number];
  selected: boolean;
}) {
  return (
    <label className="group relative flex min-w-0 cursor-pointer">
      <input
        className="peer sr-only"
        type="radio"
        name="role"
        value={role.id}
        defaultChecked={selected}
        aria-labelledby={`${role.id}-title`}
        aria-describedby={`${role.id}-description`}
        required
      />

      <span
        className="
          flex min-h-[232px] w-full flex-col gap-7
          rounded-[11px]
          border border-[var(--color-text-muted)]
          bg-[rgb(255_255_255_/_65%)]
          px-[13px] pt-5 pb-[18px]

          transition-[border-color,background-color,box-shadow]
          duration-150

          group-hover:border-[var(--color-primary)]
          group-hover:bg-[var(--color-primary-soft)]

          peer-checked:border-[var(--color-primary)]
          peer-checked:bg-[#dfe6ff]
          peer-checked:text-[var(--color-primary)]
          peer-checked:shadow-[inset_0_0_0_1px_var(--color-primary)]

          peer-focus-visible:outline
          peer-focus-visible:outline-[3px]
          peer-focus-visible:outline-[var(--color-primary)]
          peer-focus-visible:outline-offset-4

          min-[1100px]:min-h-[300px]
          min-[1100px]:px-5
          min-[1100px]:pt-7
          min-[1100px]:pb-6

          max-[699px]:min-h-32
          max-[699px]:flex-row
          max-[699px]:items-center
          max-[699px]:gap-4
          max-[699px]:px-4
          max-[699px]:py-[18px]

          max-[359px]:gap-3
          max-[359px]:px-3
          max-[359px]:py-4

          motion-reduce:transition-none
        "
      >
        <span
          className="
            relative mx-auto
            aspect-square
            w-[clamp(116px,12.3vw,160px)]
            shrink-0
            overflow-hidden
            rounded-full
            bg-[var(--color-primary-soft)]

            max-[699px]:m-0
            max-[699px]:w-20

            max-[359px]:w-16
          "
        >
          <Image
            src={role.image}
            alt=""
            fill
            sizes="(max-width: 699px) 80px, (max-width: 1100px) 116px, 160px"
            className="object-cover object-top"
          />
        </span>

        <span className="flex min-w-0 flex-col gap-2">
          <span
            id={`${role.id}-title`}
            className="
              text-[clamp(20px,2.15vw,28px)]
              font-extrabold
              leading-[1.25]
              tracking-[-0.03em]

              max-[699px]:text-xl
              max-[359px]:text-lg
            "
          >
            {role.title}
          </span>

          <span
            id={`${role.id}-description`}
            className="
              text-[clamp(11px,1.16vw,15px)]
              leading-[1.5]

              max-[699px]:text-xs
            "
          >
            {role.description}
          </span>
        </span>
      </span>
    </label>
  );
}