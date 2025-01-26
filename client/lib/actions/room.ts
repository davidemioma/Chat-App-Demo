"use server";

import { AxiosError } from "axios";
import { RoomProps } from "@/types";
import axiosInstance from "../axios";
import { RoomSchema, RoomValidator } from "../validators/room";

export const createRoomHandler = async (values: RoomValidator) => {
  try {
    const validParameters = RoomSchema.safeParse(values);

    if (!validParameters.success) {
      throw new Error("Invalid Parameters");
    }

    const res = await axiosInstance.post("/rooms/create", values);

    const result = (await res.data) as RoomProps;

    return { status: res.status, data: result };
  } catch (err) {
    console.error("createRoomHandler", err);

    if (err instanceof AxiosError) {
      throw new Error(`Something went wrong! ${err.response?.data}`);
    } else {
      throw new Error("Something went wrong! Internal server error.");
    }
  }
};
