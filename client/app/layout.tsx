import "./globals.css";
import { Toaster } from "sonner";
import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import QueryProvider from "@/providers/query-provider";
import AuthContextProvider from "@/providers/auth-provider";
import WebSocketProvider from "@/providers/websccket-provider";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Create Next App",
  description: "Generated by create next app",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <QueryProvider>
          <Toaster richColors position="top-center" />

          <AuthContextProvider>
            <WebSocketProvider>{children}</WebSocketProvider>
          </AuthContextProvider>
        </QueryProvider>
      </body>
    </html>
  );
}
