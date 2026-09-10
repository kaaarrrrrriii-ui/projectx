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
          flex h-full min-h-[360px] w-full flex-col
          rounded-[10px]
          border border-[var(--color-text-muted)]
          bg-white
          px-3 pt-4 pb-5

          transition-[border-color,background-color,box-shadow,color,transform]
          duration-150

          group-hover:-translate-y-0.5
          group-hover:border-[var(--color-primary)]
          group-hover:text-[var(--color-primary)]
          group-hover:shadow-[0_8px_24px_rgb(69_98_240_/_10%)]

          peer-checked:border-[var(--color-primary)]
          peer-checked:bg-[#dfe6ff]
          peer-checked:text-[var(--color-text)]
          peer-checked:shadow-[inset_0_0_0_1px_var(--color-primary)]

          peer-focus-visible:outline
          peer-focus-visible:outline-[3px]
          peer-focus-visible:outline-[var(--color-primary)]
          peer-focus-visible:outline-offset-4

          max-[699px]:min-h-[150px]
          max-[699px]:flex-row
          max-[699px]:items-center
          max-[699px]:px-3
          max-[699px]:py-3

          max-[379px]:min-h-[124px]
          max-[379px]:px-2.5

          motion-reduce:transition-none
          motion-reduce:transform-none
        "
      >
        <span
          className="
            relative mx-auto aspect-square
            w-[calc(100%-20px)]
            shrink-0
            overflow-hidden
            border border-[var(--color-text-muted)]
            bg-white

            max-[699px]:mx-0
            max-[699px]:w-[132px]
            max-[379px]:w-[96px]
          "
        >
          <Image
            src={role.image}
            alt=""
            fill
            sizes="(max-width: 379px) 96px, (max-width: 699px) 132px, (max-width: 1100px) 28vw, 320px"
            className="object-contain object-bottom"
          />
        </span>

        <span className="mt-5 flex min-w-0 flex-col gap-2 px-1 max-[699px]:mt-0 max-[699px]:pl-4 max-[379px]:gap-1 max-[379px]:pl-3">
          <span
            id={`${role.id}-title`}
            className="
              text-[clamp(17px,1.5vw,22px)]
              font-extrabold
              leading-[1.25]
              tracking-[-0.03em]

              max-[699px]:text-lg
              max-[379px]:text-base
            "
          >
            {role.title}
          </span>

          <span
            id={`${role.id}-description`}
            className="
              text-[clamp(11px,1vw,14px)]
              leading-[1.5]

              max-[699px]:text-xs
              max-[379px]:text-[11px]
            "
          >
            {role.description}
          </span>
        </span>
      </span>
    </label>
  );
}
