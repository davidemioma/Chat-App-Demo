"use server";

import { RoomProps } from "@/types";
import axiosInstance from "../axios";

export const getRooms = async () => {
  const res = await axiosInstance.get("/rooms");

  const result = (await res.data) as RoomProps[];

  return result;
};
