"use client";

import React, { useContext, useEffect } from "react";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { zodResolver } from "@hookform/resolvers/zod";
import { AuthContext } from "@/providers/auth-provider";
import { WebsocketContext } from "@/providers/websccket-provider";
import { ChatValidator, ChatSchema } from "@/lib/validators/chat";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";

const ChatForm = () => {
  const router = useRouter();

  const { user } = useContext(AuthContext);

  const { conn } = useContext(WebsocketContext);

  const form = useForm<ChatValidator>({
    resolver: zodResolver(ChatSchema),
    defaultValues: {
      message: "",
    },
  });

  const onSubmit = async (values: ChatValidator) => {
    if (!conn) {
      router.push("/");

      return;
    }

    await conn.send(values.message);

    form.reset();
  };

  useEffect(() => {
    if (!user) {
      return router.push("/auth/sign-in");
    }
  }, [user, router]);

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="w-full flex items-center gap-4"
      >
        <div className="w-full flex-1">
          <FormField
            control={form.control}
            name="message"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input placeholder="Write something..." {...field} />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <Button type="submit">Send</Button>
      </form>
    </Form>
  );
};

export default ChatForm;
