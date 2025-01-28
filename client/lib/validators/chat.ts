import { z } from "zod";

export const ChatSchema = z.object({
  message: z.string().min(1, { message: "Message is required" }),
});

export type ChatValidator = z.infer<typeof ChatSchema>;
