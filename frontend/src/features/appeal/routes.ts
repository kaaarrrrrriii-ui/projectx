import type { AppealRoleId } from "./roles";

export const appealRoutes = {
  role: "/appeal",
  topics: "/appeal/topics",
  description: "/appeal/description",
  details: "/appeal/details",
  attachments: "/appeal/attachments",
  success: "/appeal/success",
} as const;

export type AppealRouteName = keyof typeof appealRoutes;

export function getAppealRoute(
  route: AppealRouteName,
  role?: AppealRoleId | string,
) {
  const pathname = appealRoutes[route];
  return role ? `${pathname}?role=${encodeURIComponent(role)}` : pathname;
}
