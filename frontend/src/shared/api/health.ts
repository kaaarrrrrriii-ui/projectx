export type HealthResponse = {
  status: string;
  uptime: string;
};

export async function getHealth(): Promise<HealthResponse | null> {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

  try {
    const response = await fetch(`${apiUrl}/health`, { cache: "no-store" });

    if (!response.ok) {
      return null;
    }

    return (await response.json()) as HealthResponse;
  } catch {
    return null;
  }
}

