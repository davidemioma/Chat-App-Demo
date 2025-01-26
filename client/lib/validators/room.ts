import { z } from "zod";

export const RoomSchema = z.object({
  name: z.string().min(1, { message: "Name is required" }),
});

export type RoomValidator = z.infer<typeof RoomSchema>;
