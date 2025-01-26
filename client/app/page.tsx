"use client";

import { useContext } from "react";
import { AuthContext } from "@/providers/auth-provider";

export default function Home() {
  const { user } = useContext(AuthContext);

  return <div>Home {JSON.stringify(user)}</div>;
}
