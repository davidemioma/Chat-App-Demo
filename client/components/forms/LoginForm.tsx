"use client";

import React from "react";
import { toast } from "sonner";
import { BackBtn } from "../BackBtn";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { loginHandler } from "@/lib/actions/auth";
import { useMutation } from "@tanstack/react-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginSchema, LoginValidator } from "@/lib/validators/login";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const LoginForm = () => {
  const router = useRouter();

  const form = useForm<LoginValidator>({
    resolver: zodResolver(LoginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const { mutate, isPending } = useMutation({
    mutationKey: ["login"],
    mutationFn: async (values: LoginValidator) => {
      const result = await loginHandler(values);

      return result;
    },
    onSuccess: (res) => {
      if (res.status !== 200) {
        toast.error("Something went wrong! could not login.");
      }

      toast.success(`Login successfull. Hello ${res.data?.username}`);

      form.reset();

      router.push("/");
    },
    onError: (err) => {
      toast.error("Something went wrong! " + err.message);
    },
  });

  const onSubmit = (values: LoginValidator) => {
    mutate(values);
  };

  return (
    <div className="w-full px-5 h-screen flex flex-col items-center justify-center p-5">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Welcome Back!</CardTitle>

          <CardDescription>Sign In to your account.</CardDescription>
        </CardHeader>

        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>

                    <FormControl>
                      <Input
                        type="email"
                        placeholder="test@mail.com"
                        disabled={isPending}
                        {...field}
                      />
                    </FormControl>

                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>

                    <FormControl>
                      <Input
                        type="password"
                        placeholder="*******"
                        disabled={isPending}
                        {...field}
                      />
                    </FormControl>

                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" disabled={isPending}>
                {isPending ? "Loading..." : "Sign In"}
              </Button>
            </form>
          </Form>
        </CardContent>

        <CardFooter>
          <BackBtn label="Don't have an account?" href="/auth/sign-up" />
        </CardFooter>
      </Card>
    </div>
  );
};

export default LoginForm;
