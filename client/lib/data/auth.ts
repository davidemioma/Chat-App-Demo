"use server";

import axios from "axios";
import { cache } from "react";
import { UserProps } from "@/types";
import axiosInstance from "../axios";

export const getCurrentUser = cache(async () => {
  try {
    const res = await axiosInstance.get(
      `${process.env.NEXT_PUBLIC_BASE_API_URL}/auth/user`
    );

    const result = (await res.data) as UserProps;

    return result;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      console.error(`getCurrentUser error: ${error.message}`);

      if (error.response?.status === 404) {
        return { error: "Token not found. Please check your credentials." };
      }

      if (error.response?.status === 401) {
        return { error: "Unauthorized. Please check your credentials." };
      }
    }

    return { error: "An unexpected error occurred. Please try again later." };
  }
});
