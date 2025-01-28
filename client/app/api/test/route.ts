import { NextResponse } from "next/server";

export async function GET(request: Request) {
  const response = NextResponse.json({ message: "Cookie set!" });

  response.cookies.set("myCookie", "cookieValue", { path: "/" });

  return response;
}
