"use client";

import { useContext } from "react";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/router";
import { getRooms } from "@/lib/data/rooms";
import { Button } from "@/components/ui/button";
import { useQuery } from "@tanstack/react-query";
import CreateRoom from "@/components/forms/CreateRoom";
import { AuthContext } from "@/providers/auth-provider";
import { Card, CardContent } from "@/components/ui/card";

export default function Home() {
  const router = useRouter();

  const { user } = useContext(AuthContext);

  const {
    data: rooms,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["get-rooms"],
    queryFn: async () => {
      const rooms = await getRooms();

      return rooms;
    },
  });

  return (
    <div className="w-full max-w-5xl mx-auto min-h-screen flex flex-col gap-10 py-7 px-5">
      <h1 className="text-2xl font-bold">Welcome {user?.username}</h1>

      <div className="w-full flex items-center justify-center">
        <CreateRoom />
      </div>

      <div className="space-y-5">
        <h1 className="text-xl font-bold">Available Rooms</h1>

        {isLoading && (
          <div className="w-full py-10 flex items-center justify-center">
            <Loader2 className="w-5 h-5 animate-spin" />
          </div>
        )}

        {isError && (
          <div className="w-full py-10 flex items-center justify-center">
            <span className="text-muted-foreground">No rooms Available!</span>
          </div>
        )}

        {!isLoading && !isError && !rooms && (
          <div className="w-full py-10 flex items-center justify-center">
            <span className="text-muted-foreground">No rooms Available!</span>
          </div>
        )}

        {!isLoading && !isError && rooms && rooms.length > 0 && (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5">
            {rooms.map((room) => (
              <Card key={room.id}>
                <CardContent className="flex items-center justify-between gap-2 py-4">
                  <div>
                    <h2 className="font-bold">Room</h2>

                    <p className="text-sm text-muted-foreground">{room.name}</p>
                  </div>

                  <Button>Join</Button>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
