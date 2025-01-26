"use server";

import { AxiosError } from "axios";
import { UserProps } from "@/types";
import axiosInstance from "../axios";
import { LoginValidator, LoginSchema } from "../validators/login";
import { RegisterValidator, RegisterSchema } from "../validators/register";

export const loginHandler = async (values: LoginValidator) => {
  try {
    const validParameters = LoginSchema.safeParse(values);

    if (!validParameters.success) {
      throw new Error("Invalid Parameters");
    }

    const res = await axiosInstance.post("/auth/sign-in", values);

    const result = (await res.data) as UserProps;

    return { status: res.status, data: result };
  } catch (err) {
    console.error("loginHandler", err);

    if (err instanceof AxiosError) {
      if (err.response?.status === 401) {
        throw new Error("Password does not match!");
      } else if (err.response?.status === 404) {
        throw new Error("User not found!");
      } else {
        throw new Error(err.response?.data);
      }
    } else {
      throw new Error("Something went wrong! Internal server error.");
    }
  }
};

export const registerHandler = async (values: RegisterValidator) => {
  try {
    const validParameters = RegisterSchema.safeParse(values);

    if (!validParameters.success) {
      throw new Error("Invalid Parameters");
    }

    const res = await axiosInstance.post("/auth/sign-up", values);

    const result = await res.data;

    return { status: res.status, data: result };
  } catch (err) {
    console.error("registerHandler", err);

    if (err instanceof AxiosError) {
      if (err.response?.status === 401) {
        throw new Error("Email already exists!");
      } else {
        throw new Error(err.response?.data);
      }
    } else {
      throw new Error("Something went wrong! Internal server error.");
    }
  }
};
