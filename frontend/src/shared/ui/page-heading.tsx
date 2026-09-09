interface Props {
  id?: string;
  title: string;
  description?: string;
}

export default function PageHeading({
  id,
  title,
  description,
}: Props) {
  return (
    <div>
      <h1
        id={id}
        className="
          text-[length:var(--heading-size,30px)]
          font-extrabold
          leading-[1.2]
          tracking-[-0.025em]
          text-[var(--heading-color,var(--color-text))]
        "
      >
        {title}
      </h1>

      {description && (
        <p
          className="
            mt-3
            max-w-[var(--heading-description-width,370px)]
            text-[length:var(--heading-description-size,15px)]
            leading-[1.55]
            text-[var(--color-text-muted)]
          "
        >
          {description}
        </p>
      )}
    </div>
  );
}