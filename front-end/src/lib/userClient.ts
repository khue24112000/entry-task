import { publicRoute, privateRoute } from "./axiosClient";
import type { User } from "@/types/api";
import { uploadImage } from "@/lib/firebase";

export async function login(
  username: string,
  password: string
) {
  try {
    const res = await privateRoute.post("/login", {
      username,
      password,
    });
    if (res.statusText !== "OK") {
      throw new Error(res.statusText);
    }
    return res.data;
  } catch (error) {
    console.error("Something wrong: ", error);
  }
}

export async function logout(username: string) {
  await privateRoute.post("/logout", { username });
}

export async function register(user: User) {
  try {
    let avatarUrl;
    if (user.avatarUrl) {
      avatarUrl = await uploadImage(user.avatarUrl);
    } else {
      avatarUrl = "";
    }
    // console.log(user);
    const res = await publicRoute.post("/register", {
      ...user,
      avatarUrl,
    });
    return res.data;
  } catch {
    throw new Error("Something wrong");
  }
}

export async function getUser(username: string) {
  try {
    const res = await privateRoute.get(
      `/profile/${username}`
    );
    return res.data;
  } catch (error) {
    throw error;
  }
}

export async function updateUser(
  username: string,
  nickname: string,
  avatarUrl: string | File
) {
  try {
    let newAvatarUrl: undefined | string = "";
    if (typeof avatarUrl !== "string") {
      newAvatarUrl = await uploadImage(avatarUrl);
    }
    const info = {
      nickname,
      avatarUrl:
        typeof avatarUrl !== "string"
          ? newAvatarUrl
          : avatarUrl,
    };
    await privateRoute.put(`/profile/${username}`, info);
    return info;
  } catch (error) {
    console.error("Something wrong", error);
  }
}

export async function checkUser() {
  try {
    const res = await privateRoute.get("/me");
    return res;
  } catch (error) {
    throw error;
  }
}
