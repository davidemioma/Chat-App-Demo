"use client";

import React from "react";
import Link from "next/link";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { useRouter } from "next/navigation";
import { useContext, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useMutation } from "@tanstack/react-query";
import { registerHandler } from "@/lib/actions/auth";
import { zodResolver } from "@hookform/resolvers/zod";
import { AuthContext } from "@/providers/auth-provider";
import { RegisterSchema, RegisterValidator } from "@/lib/validators/register";
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

export default function SignUp() {
  const router = useRouter();

  const { authenticated } = useContext(AuthContext);

  const form = useForm<RegisterValidator>({
    resolver: zodResolver(RegisterSchema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
    },
  });

  const { mutate, isPending } = useMutation({
    mutationKey: ["register"],
    mutationFn: async (values: RegisterValidator) => {
      const result = await registerHandler(values);

      return result;
    },
    onSuccess: (res) => {
      if (res.status !== 201) {
        toast.error("Something went wrong! could not register user.");
      }

      toast.success(`Registration successfull`);

      form.reset();

      router.push("/auth/sign-in");
    },
    onError: (err) => {
      toast.error("Something went wrong! " + err.message);
    },
  });

  const onSubmit = (values: RegisterValidator) => {
    mutate(values);
  };

  useEffect(() => {
    if (authenticated) {
      router.push("/");
    }
  }, [authenticated, router]);

  return (
    <div className="w-full px-5 h-screen flex flex-col items-center justify-center p-5">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Welcome!</CardTitle>

          <CardDescription>Create an account.</CardDescription>
        </CardHeader>

        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Username</FormLabel>

                    <FormControl>
                      <Input
                        type="text"
                        placeholder="john..."
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
                {isPending ? "Loading..." : "Sign Up"}
              </Button>
            </form>
          </Form>
        </CardContent>

        <CardFooter className="text-sm">
          Already have an account?
          <Link
            className="text-blue-500 font-semibold ml-1"
            href="/auth/sign-in"
          >
            Sign In
          </Link>
        </CardFooter>
      </Card>
    </div>
  );
}
