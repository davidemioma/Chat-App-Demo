"use server";

import axiosInstance from "../axios";
import { RoomProps, RoomUserProps } from "@/types";

export const getRooms = async () => {
  const res = await axiosInstance.get("/rooms");

  const result = (await res.data) as RoomProps[];

  return result;
};

export const getRoomUsers = async (roomId: string) => {
  const res = await axiosInstance.get(`/rooms/${roomId}/clients`);

  const result = (await res.data) as RoomUserProps[];

  return result;
};
