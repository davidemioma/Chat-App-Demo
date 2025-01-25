import { redirect } from "next/navigation";
import { getCurrentUser } from "@/lib/data/auth";

export default async function Home() {
  const currentUser = await getCurrentUser();

  if (!currentUser || "error" in currentUser) {
    return redirect("/auth/sign-in");
  }

  return <div>Home {JSON.stringify(currentUser)}</div>;
}
