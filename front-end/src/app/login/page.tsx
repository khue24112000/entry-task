"use client";

import { checkUser, login } from "@/lib/userClient";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const Login = () => {
  const [passwordType, setPasswordType] = useState("text");
  const [error, setError] = useState("");
  const router = useRouter();
  const handleSubmitForm = async (
    e: React.FormEvent<HTMLFormElement>
  ) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const username = formData
      .get("username")
      ?.toString()
      .trim();
    const password = formData
      .get("password")
      ?.toString()
      .trim();

    if (username && password) {
      try {
        const res = await login(username, password);
        if (res.message) {
          alert("login success");
          router.push(`/profile/${username}`);
        }
      } catch {
        alert("Login failed");
      }
    } else {
      setError("Missing username or password");
    }
  };

  const handleFocusPwd = () => {
    if (passwordType === "text") {
      setPasswordType("password");
    }
  };

  const handleClickSignup = () => {
    router.push("/register");
  };

  useEffect(() => {
    const fetch = async () => {
      try {
        const res = await checkUser();
        if (res.status === 200) {
          router.push(`/profile/${res.data.username}`);
        }
      } catch (error) {
        console.log("User not logged in", error);
      }
    };
    fetch();
  }, [router]);
  return (
    <div className="w-screen h-screen flex justify-center items-center">
      <div className="w-80 rounded-lg p-4 shadow-lg border border-gray-200 ">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-3xl text-center font-bold text-green-500 cursor-default ">
            Login
          </h2>
          <a
            onClick={handleClickSignup}
            className="cursor-pointer hover:underline text-green-500 "
          >
            Sign up
          </a>
        </div>
        <form
          onSubmit={handleSubmitForm}
          className="space-y-4"
        >
          <div>
            <label htmlFor="username">Username</label>
            <input
              type="text"
              name="username"
              className="w-full rounded-lg outline-green-100 border-gray-300 border p-2 focus:bg-green-50 "
            />
          </div>

          <div>
            <label htmlFor="password">Password</label>
            <input
              type={passwordType}
              name="password"
              autoComplete="off"
              onFocus={handleFocusPwd}
              className="w-full rounded-lg outline-green-100 border-gray-300 border p-2 focus:bg-green-50 "
            />
          </div>

          {error && (
            <p className="text-sm font-semibold text-red-500 ml-2 -mt-2">
              {error}
            </p>
          )}

          <button
            type="submit"
            className="bg-green-500 text-white px-4 py-2 mt-4 rounded-lg w-full font-semibold cursor-pointer"
          >
            Sign In
          </button>
        </form>
      </div>
    </div>
  );
};

export default Login;
