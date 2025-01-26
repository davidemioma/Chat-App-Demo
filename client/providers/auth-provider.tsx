"use client";

import { UserProps } from "@/types";
import { usePathname, useRouter } from "next/navigation";
import { useState, createContext, useEffect } from "react";

export const AuthContext = createContext<{
  authenticated: boolean;
  setAuthenticated: (auth: boolean) => void;
  user: UserProps | null;
  setUser: (user: UserProps) => void;
}>({
  authenticated: false,
  setAuthenticated: () => {},
  user: null,
  setUser: () => {},
});

const AuthContextProvider = ({ children }: { children: React.ReactNode }) => {
  const router = useRouter();

  const pathname = usePathname();

  const [user, setUser] = useState<UserProps | null>(null);

  const [authenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    if (pathname.includes("/auth")) return;

    const userInfo = localStorage.getItem("user_info");

    if (!userInfo && !pathname.includes("/auth")) {
      router.push("/auth/sign-in");

      setUser(null);

      setAuthenticated(false);

      return;
    }

    const user: UserProps = JSON.parse(userInfo!);

    setUser(user);

    setAuthenticated(true);
  }, [pathname, router]);

  return (
    <AuthContext.Provider
      value={{
        authenticated: authenticated,
        setAuthenticated: setAuthenticated,
        user: user,
        setUser: setUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContextProvider;
