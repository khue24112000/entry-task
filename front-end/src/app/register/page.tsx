"use client";

import { useRouter } from "next/navigation";
import { ChangeEvent, useState } from "react";
import type { User } from "@/types/api";
import { register } from "@/lib/userClient";

const Register = () => {
  const [registerForm, setRegisterForm] = useState({
    username: "",
    nickname: "",
    password: "",
    rptPassword: "",
  });

  const [passwordType, setPasswordType] = useState("text");
  const [avatarUrl, setAvatarUrl] = useState<File | null>(
    null
  );
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const router = useRouter();
  const passwordRegex = /^\S{8,100}$/;

  const handleChangeInput = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    if (e.target.value || e.target.value === "") {
      setRegisterForm((prevForm) => {
        return {
          ...prevForm,
          [e.target.name]: e.target.value.trim(),
        };
      });
    }
    if (e.target.name === "password") {
      setPasswordType("password");
    }
    if (
      e.target.name === "rptPassword" &&
      e.target.value !== registerForm.password
    ) {
      setError("Repeat Password don't match with password");
    } else if (e.target.name === "rptPassword") {
      setError("");
    }
  };

  const handleChangeAvatar = (
    e: ChangeEvent<HTMLInputElement>
  ) => {
    const file = e.target.files?.[0];
    if (file && file.size > 2 * 1024 * 1024) {
      alert("File quá lớn! Giới hạn 2MB");
    } else if (file) {
      setAvatarUrl(file);
    }
  };

  const handleSubmitForm = async (
    e: React.FormEvent<HTMLFormElement>
  ) => {
    e.preventDefault();
    setIsLoading(true);
    const user: User = {
      username: registerForm.username,
      password: registerForm.password,
      nickname: registerForm.nickname,
      avatarUrl,
    };
    try {
      await register(user);
      confirm("Register success. Redirect to login page?")
        ? handleClickSignin()
        : resetInput();
    } catch (error) {
      alert(`Register error: ${error.response.data.error}`);
      setRegisterForm((prevForm) => {
        return {
          ...prevForm,
          password: "",
          rptPassword: "",
        };
      });
    } finally {
      setIsLoading(false);
    }
  };
  const signupList = async () => {
    for (let i = 0; i < 200; i++) {
      await register({
        username: `khue${i}`,
        password: "123",
        nickname: "",
        avatarUrl: null,
      });
    }
  };

  const resetInput = () => {
    setRegisterForm({
      username: "",
      nickname: "",
      password: "",
      rptPassword: "",
    });
    setError("");
    setPasswordType("text");
    setAvatarUrl(null);
  };

  const handleClickSignin = () => {
    router.push("/login");
  };

  return (
    <div className="w-screen h-screen flex justify-center items-center">
      <button
        className="cursor-pointer p-2 bg-black text-white m-4"
        onClick={signupList}
      >
        Signup list
      </button>
      <div className="w-80 rounded-lg p-4 shadow-lg border border-gray-200 relative ">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-3xl text-center font-bold text-green-500 cursor-default ">
            Login
          </h2>
          <a
            onClick={handleClickSignin}
            className="cursor-pointer hover:underline text-green-500 "
          >
            Sign in
          </a>
        </div>
        <form
          onSubmit={handleSubmitForm}
          className="space-y-4"
          autoComplete="off"
        >
          <div>
            <label htmlFor="nickname">Nickname</label>
            <input
              type="text"
              name="nickname"
              value={registerForm.nickname}
              onChange={handleChangeInput}
              id="nickname"
              className="w-full rounded-lg outline-green-100 border-gray-300 border p-2 focus:bg-green-50 "
            />
          </div>
          <div>
            <label
              className="after:content-['*'] after:text-red-500 after:ml-0.5"
              htmlFor="username"
            >
              Username
            </label>
            <input
              type="text"
              name="username"
              id="username"
              value={registerForm.username}
              onChange={handleChangeInput}
              className="w-full rounded-lg outline-green-100 border-gray-300 border p-2 focus:bg-green-50 "
            />
          </div>

          <div>
            <label
              className="after:content-['*'] after:text-red-500 after:ml-0.5"
              htmlFor="password"
            >
              Password
            </label>
            <input
              type={passwordType}
              autoComplete="off"
              name="password"
              value={registerForm.password}
              onChange={handleChangeInput}
              className="w-full rounded-lg outline-green-100 border-gray-300 border p-2 focus:bg-green-50 "
            />
          </div>

          <div>
            <label
              className="after:content-['*'] after:text-red-500 after:ml-0.5"
              htmlFor="password"
            >
              Repeat Password
            </label>
            <input
              type={passwordType}
              autoComplete="off"
              name="rptPassword"
              value={registerForm.rptPassword}
              onChange={handleChangeInput}
              disabled={
                registerForm.password ? false : true
              }
              className="w-full rounded-lg disabled:bg-gray-300 outline-green-100 border-gray-300 border p-2 focus:bg-green-50 "
            />
          </div>

          <div className="">
            <label className="block" htmlFor="avatar">
              Avatar
            </label>

            <input
              type="file"
              name="avatar"
              id="avatar"
              onChange={handleChangeAvatar}
              accept="image/*"
              className="w-full hidden rounded-lg outline-green-100 border-gray-300 border p-2 focus:bg-green-50 "
            />
            <label
              htmlFor="avatar"
              className="px-2 py-1 mt-2 inline-block rounded-lg cursor-pointer bg-green-100 "
            >
              Choose image
            </label>
            {avatarUrl && (
              <p className="overflow-hidden text-ellipsis">
                {avatarUrl?.name}
              </p>
            )}
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
        {isLoading && (
          <div className="absolute rounded-lg inset-0 flex text-white text-2xl font-semibold justify-center items-center bg-gray-500/50">
            Loading...
          </div>
        )}
      </div>
    </div>
  );
};

export default Register;
