"use client";

import React, { useContext, useEffect, useState } from "react";
import { cn } from "@/lib/utils";
import { useRouter } from "next/navigation";
import { getRoomUsers } from "@/lib/data/rooms";
import { MessageProps, RoomUserProps } from "@/types";
import { AuthContext } from "@/providers/auth-provider";
import { WebsocketContext } from "@/providers/websccket-provider";

type Props = {
  roomId: string;
};

const ChatBody = ({ roomId }: Props) => {
  const router = useRouter();

  const { user } = useContext(AuthContext);

  const { conn } = useContext(WebsocketContext);

  const [users, setUsers] = useState<RoomUserProps[]>([]);

  const [messages, setMessages] = useState<MessageProps[]>([]);

  useEffect(() => {
    if (!conn) {
      router.push("/");

      return;
    }

    const getUsers = async () => {
      try {
        const res = await getRoomUsers(roomId);

        setUsers(res);
      } catch (err) {
        console.error(`Get Users Err: ${err}`);

        setUsers([]);
      }
    };

    getUsers();
  }, [conn, roomId, router]);

  useEffect(() => {
    if (!user) {
      router.push("/auth/sign-in");

      return;
    }

    if (!conn) {
      router.push("/");

      return;
    }

    conn.onmessage = (message: MessageEvent<string>) => {
      const m: MessageProps = JSON.parse(message.data);

      console.log("Message Recieved: ", m);

      if (m.content === "A new user has joined the room") {
        setUsers((prev) => [
          {
            id: m.clientId,
            username: m.username,
          },
          ...prev,
        ]);
      }

      if (m.content == "user left the chat") {
        const updatedUsers = users.filter((user) => user.id !== m.clientId);

        setUsers(updatedUsers);

        setMessages((prev) => [...prev, m]);

        return;
      }

      // Handle regular messages
      setMessages((prev) => [
        ...prev,
        { ...m, type: m.clientId === user.id ? "sent" : "recieved" },
      ]);
    };

    conn.onopen = () => {
      console.log("WebSocket connection established.");
    };

    conn.onclose = () => {
      console.log("Connection closed");
    };

    conn.onerror = (err) => {
      console.log(`Connection err: ${err}`);
    };
  }, [conn, router, user, users]);

  // if (messages.length === 0) {
  //   return (
  //     <div className="w-full h-full flex items-center justify-center text-muted-foreground">
  //       No messages available!
  //     </div>
  //   );
  // }

  return (
    <>
      <div>Users: {JSON.stringify(users)}</div>

      <div>Messages: {JSON.stringify(messages)}</div>

      {messages.map((msg, i) => (
        <div
          key={i}
          className={cn(
            "mt-3",
            msg.type === "recieved" && "w-full flex flex-col items-end"
          )}
        >
          <p className="text-sm text-muted-foreground">{msg.username}</p>

          <div
            className={cn(
              "w-fit mt-1 py-1 px-4 rounded-md",
              msg.type === "recieved"
                ? "bg-gray-300 text-dark-secondary"
                : "bg-blue-500 text-white"
            )}
          >
            {msg.content}
          </div>
        </div>
      ))}
    </>
  );
};

export default ChatBody;
