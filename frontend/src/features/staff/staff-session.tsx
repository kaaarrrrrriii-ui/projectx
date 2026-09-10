"use client";

import { useEffect, useState } from "react";
import { clearStaffSession, getStaffSession, type StaffRole } from "@/shared/api/staff-api";

export function StaffIdentity({ role, label, fallback }: { role: StaffRole; label: string; fallback: string }) {
  const [session, setSession] = useState<ReturnType<typeof getStaffSession>>(null);
  const name = session?.user.full_name || fallback;

  useEffect(() => {
    queueMicrotask(() => {
      const current = getStaffSession();
      if (!current || current.user.role !== role) {
        // eslint-disable-next-line @next/next/no-location-assign-relative-destination
        window.location.assign("/staff");
        return;
      }
      setSession(current);
    });
  }, [role]);

  return (
    <div className="mb-7 px-2 max-[799px]:mb-3">
      <p className="font-semibold text-[#4562f0]">{name}</p>
      <p className="mt-0.5 text-xs text-[#646d86]">{label}</p>
    </div>
  );
}

export function StaffLogout({ className = "" }: { className?: string }) {
  return (
    <button
      type="button"
      onClick={() => {
        clearStaffSession();
        // eslint-disable-next-line @next/next/no-location-assign-relative-destination
        window.location.assign("/");
      }}
      className={`mt-auto flex min-h-10 cursor-pointer items-center justify-center rounded-xl border border-[#4562f0] px-4 text-sm font-medium text-[#4562f0] transition-colors hover:bg-[#4562f0] hover:text-white ${className}`}
    >
      Выйти
    </button>
  );
}
