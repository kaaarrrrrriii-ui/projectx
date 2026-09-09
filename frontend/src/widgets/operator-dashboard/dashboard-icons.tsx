type DashboardIconProps = {
  name: "mail" | "urgent" | "clock" | "return";
  className?: string;
};

export default function DashboardIcon({ name, className = "h-9 w-9" }: DashboardIconProps) {
  const props = { viewBox: "0 0 40 40", className, "aria-hidden": true } as const;

  if (name === "urgent") {
    return <svg {...props}><path fill="currentColor" d="M21.8 2.5c1.6 7-4.7 8.7-5.5 14.1-1.2-1-1.8-2.6-1.5-5C9.9 15 7 19.3 7 24.5 7 32 12.8 37 20 37s13-5 13-12.5c0-6.5-3.7-12.2-11.2-22Zm-1.4 29.3c-2.9 0-5.2-2-5.2-5 0-2.1 1.2-3.9 3.2-5.3-.1 1 .2 2 .7 2.6.3-2.2 2.9-2.9 2.3-5.8 3 2 4.3 5 4.3 7.8 0 3.4-2.3 5.7-5.3 5.7Z" /></svg>;
  }

  if (name === "clock") {
    return <svg {...props}><circle cx="20" cy="20" r="15" fill="currentColor" /><path d="M20 10v11h8" fill="none" stroke="white" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" /></svg>;
  }

  if (name === "return") {
    return <svg {...props} fill="none"><path d="m15 12-8 8 8 8" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round" /><path d="M8 20h15a9 9 0 0 1 9 9" stroke="currentColor" strokeWidth="3" strokeLinecap="round" /></svg>;
  }

  return <svg {...props} fill="none"><path d="m6 16 14-10 14 10v17H6V16Z" fill="currentColor" /><path d="m7 17 13 10 13-10" stroke="white" strokeWidth="2.2" strokeLinejoin="round" /><path d="M8 31h24" stroke="white" strokeWidth="2" /></svg>;
}
