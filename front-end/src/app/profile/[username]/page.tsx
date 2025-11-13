"use client";

import {
  getUser,
  logout,
  updateUser,
} from "@/lib/userClient";
import Image from "next/image";
import { useParams, useRouter } from "next/navigation";
import React, {
  useEffect,
  useState,
  ChangeEvent,
  FormEvent,
} from "react";

interface User {
  username: string;
  nickname: any;
  avatarUrl: any;
}

const Profile = () => {
  const router = useRouter();
  const { username } = useParams<{ username: string }>();
  const [profile, setProfile] = useState<User>({
    username: "",
    nickname: "",
    avatarUrl: "",
  });
  const [isModal, setIsModal] = useState(false);
  const [newAvatarUrl, setNewAvatarUrl] =
    useState<File | null>(null);
  const [newNickName, setNewNickname] = useState<
    string | undefined
  >("");
  const [isLoading, setIsLoading] = useState(false);

  const handleChangeProfile = () => {
    setIsModal(true);
    setNewNickname(profile.nickname);
    setNewAvatarUrl(null);
  };

  const handleCloseModal = () => {
    setIsModal(false);
  };

  const handleChangeNickname = (
    e: ChangeEvent<HTMLInputElement>
  ) => {
    if (e.target.value || e.target.value === "") {
      setNewNickname(e.target.value.trim());
    }
  };

  const handleChangeAvatar = (
    e: ChangeEvent<HTMLInputElement>
  ) => {
    const file = e.target.files?.[0];
    if (file && file.size > 2 * 1024 * 1024) {
      alert("File quá lớn! Giới hạn 2MB");
    } else if (file) {
      setNewAvatarUrl(file);
    }
  };

  const handleUpdate = async (
    e: FormEvent<HTMLFormElement>
  ) => {
    e.preventDefault();
    setIsLoading(true);
    try {
      const nickname = newNickName
        ? newNickName
        : profile.nickname;
      const avatarUrl = newAvatarUrl
        ? newAvatarUrl
        : profile.avatarUrl;
      const res = await updateUser(
        username,
        nickname,
        avatarUrl
      );
      setProfile({
        username,
        nickname: res?.nickname,
        avatarUrl: res!.avatarUrl,
      });
      handleCloseModal();
    } catch (error) {
      alert(`Update error: ${error.response?.data.error}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      if (confirm("Logout confirm:")) {
        await logout(username);
        router.push("/login");
      }
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await getUser(username);
        if (userData) {
          setProfile({
            username: userData.username,
            nickname: userData.nickname,
            avatarUrl: userData.avatar_url,
          });
        }
      } catch (error) {
        if (error.status === 404 || error.status === 401) {
          router.push("/login");
        }
      }
    };
    fetchUser();
  }, [username]);

  return (
    <div className="min-h-screen flex justify-center items-center relative">
      <div className="w-80 flex flex-col items-center">
        <Image
          src={
            profile.avatarUrl
              ? profile.avatarUrl
              : "https://www.svgrepo.com/show/452030/avatar-default.svg"
          }
          alt="avatar"
          width={200}
          height={200}
          loading="eager"
          className="w-52 h-52 rounded-full object-cover bg-gray-200 mb-5"
        />

        <p className="text-2xl font-semibold mb-10">
          {profile.nickname}
        </p>
        <button
          onClick={handleChangeProfile}
          className="w-full p-1 text-center border border-gray-300 rounded-lg flex-1 cursor-pointer hover:bg-gray-200 mb-5"
        >
          Change Profile
        </button>
        <button
          onClick={handleLogout}
          className="w-full bg-green-500 text-white font-semibold px-4 py-2 rounded-lg cursor-pointer hover:bg-green-600"
        >
          Log out
        </button>
      </div>
      {isModal && (
        <div className="absolute inset-0 bg-gray-700/80 flex items-center justify-center">
          <div className="max-w-full w-64 sm:w-72 p-4 relative bg-white rounded-lg flex flex-col items-center shadow-lg">
            <h2 className="text-xl sm:text-2xl text-green-500 font-semibold mb-5">
              Update Profile
            </h2>
            <form
              onSubmit={handleUpdate}
              className="space-y-4"
            >
              <div>
                <label htmlFor="nickname">Nickname</label>
                <input
                  type="text"
                  name="nickname"
                  value={newNickName}
                  onChange={handleChangeNickname}
                  id="nickname"
                  className="w-full rounded-lg outline-green-100 border-gray-300 border p-2 focus:bg-green-50 "
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
                  className="px-2 py-1 mt-2 inline-block rounded-lg cursor-pointer bg-green-100 hover:bg-green-200"
                >
                  Choose image
                </label>
                {newAvatarUrl && (
                  <p className=" max-w-64 mt-2 overflow-hidden text-ellipsis">
                    {newAvatarUrl?.name}
                  </p>
                )}
              </div>
              <div>
                <button
                  type="submit"
                  className="w-full bg-green-500 rounded-lg p-2 text-white font-semibold cursor-pointer hover:bg-green-600"
                >
                  Update
                </button>
              </div>
            </form>
            {isLoading && (
              <div className="absolute rounded-lg inset-0 flex text-white text-2xl font-semibold justify-center items-center bg-slate-400/50">
                Loading...
              </div>
            )}
            <span
              className="absolute right-0 top-0 px-2 py-px rounded-full text-xl cursor-pointer hover:bg-gray-200"
              onClick={handleCloseModal}
            >
              &times;
            </span>
          </div>
        </div>
      )}
    </div>
  );
};

export default Profile;
