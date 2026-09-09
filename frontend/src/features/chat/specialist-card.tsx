import Surface from "@/shared/ui/surface";
import Image from "next/image";

type SpecialistCardProps = {
  avatarSrc: string;
  name: string;
  profession: string;
  experience: string;
  specialization: string;
};

export default function SpecialistCard({
  avatarSrc,
  name,
  profession,
  experience,
  specialization,
}: SpecialistCardProps) {
  return (
    <Surface
      as="section"
      labelledBy="specialist-name"
      className="flex min-h-[148px] w-full flex-col items-start gap-3 p-4 [--surface-radius:16px] [--surface-shadow:0_4px_12px_rgba(69,98,240,0.05)] !border-[#dee7fd] !bg-[#fcfdff] max-[799px]:mx-auto max-[799px]:max-w-[420px]"
    >
      <div className="flex items-center gap-2.5">
        <Image
          src={avatarSrc}
          alt=""
          width={48}
          height={48}
          sizes="48px"
          className="h-12 w-12 shrink-0 rounded-full object-cover"
        />

        <div className="min-w-0">
          <p className="text-[11px] leading-4 text-[#777d91] uppercase">
            {profession}
          </p>
          <h2
            id="specialist-name"
            className="truncate text-[14px] leading-5 font-semibold text-[var(--color-primary)]"
          >
            {name}
          </h2>
        </div>
      </div>

      <dl className="grid gap-1 text-[11px] leading-4 text-[#7c8295]">
        <div className="flex items-start gap-2">
          <dt aria-label="Опыт работы" className="w-3 shrink-0 text-[#f2b500]">
            ★
          </dt>
          <dd>Опыт работы: {experience}</dd>
        </div>
        <div className="flex items-start gap-2">
          <dt aria-label="Специализация" className="w-3 shrink-0 text-[var(--color-primary)]">
            ◇
          </dt>
          <dd>{specialization}</dd>
        </div>
      </dl>
    </Surface>
  );
}
