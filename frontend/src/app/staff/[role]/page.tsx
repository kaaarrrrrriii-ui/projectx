import { notFound, redirect } from "next/navigation";

const roles = new Set(["operator", "expert", "admin"]);

export default async function StaffRolePage({ params }: { params: Promise<{ role: string }> }) {
  const { role } = await params;
  if (!roles.has(role)) notFound();
  redirect("/staff");
}
