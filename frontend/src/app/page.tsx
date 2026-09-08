import { getHealth } from "@/shared/api/health";

export default async function Home() {
  const health = await getHealth();

  return (
    <main>
      <div className="/">fdfgdfgdfgd</div>
    </main>
  );
}



