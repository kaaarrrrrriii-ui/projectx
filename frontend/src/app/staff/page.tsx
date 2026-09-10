import Link from "next/link";

export default function StaffPage() {
  return (
    <main className="flex min-h-dvh min-w-[320px] items-center justify-center bg-[#202020] px-6 py-10 max-[599px]:items-end max-[599px]:px-0 max-[599px]:py-0">
      <section
        aria-labelledby="employee-login-heading"
        className="relative w-full max-w-[800px] rounded-[14px] bg-white px-9 pt-7 pb-9 text-[#11131a] shadow-[0_24px_80px_rgba(0,0,0,0.22)] max-[599px]:rounded-b-none max-[599px]:px-5 max-[599px]:pt-6 max-[599px]:pb-[max(28px,env(safe-area-inset-bottom))]"
      >
        <Link
          href="/"
          aria-label="Закрыть окно входа"
          className="absolute top-4 right-5 flex size-9 items-center justify-center rounded-full text-[28px] leading-none font-light text-[#4562f0] transition-colors hover:bg-[#eef1ff] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#4562f0]"
        >
          ×
        </Link>

        <h1
          id="employee-login-heading"
          className="pr-12 text-[32px] leading-tight font-extrabold tracking-[-0.03em] text-[#4562f0] max-[599px]:text-[28px]"
        >
          Вход
        </h1>

        <form className="mt-8 space-y-7">
          <label className="block">
            <span className="mb-2 block text-[20px] font-medium max-[599px]:text-base">Логин</span>
            <input
              type="text"
              name="login"
              autoComplete="username"
              placeholder="текст текст"
              className="h-[43px] w-full rounded-[15px] border border-[#18223f] bg-[#fcfdff] px-4 text-center text-[15px] outline-none placeholder:text-[#30364a] focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/20 max-[599px]:h-12 max-[599px]:text-base"
            />
          </label>

          <label className="block">
            <span className="mb-2 block text-[20px] font-medium max-[599px]:text-base">Пароль</span>
            <input
              type="password"
              name="password"
              autoComplete="current-password"
              placeholder="текст текст"
              className="h-[43px] w-full rounded-[15px] border border-[#18223f] bg-[#fcfdff] px-4 text-center text-[15px] outline-none placeholder:text-[#30364a] focus:border-[#4562f0] focus:ring-2 focus:ring-[#4562f0]/20 max-[599px]:h-12 max-[599px]:text-base"
            />
          </label>

          <button
            type="button"
            className="mt-1 flex h-[48px] w-full cursor-pointer items-center justify-center rounded-[11px] bg-[#4562f0] px-5 text-[16px] font-normal text-white transition-colors hover:bg-[#3854dc] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-white active:bg-[#3049c8]"
          >
            Войти
          </button>
        </form>
      </section>
    </main>
  );
}
