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

type AppealRouteParams = {
  role?: AppealRoleId | string;
  topic?: string;
};

export function getAppealRoute(
  route: AppealRouteName,
  params: AppealRouteParams = {},
) {
  const searchParams = new URLSearchParams();

  if (params.role) searchParams.set("role", params.role);
  if (params.topic) searchParams.set("topic", params.topic);

  const query = searchParams.toString();
  return query ? `${appealRoutes[route]}?${query}` : appealRoutes[route];
}
