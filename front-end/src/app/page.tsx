import { redirect } from "next/navigation";
const Home = () => {
  const session = null;
  console.log(session);
  if (!session) {
    redirect("/login");
  }

  return <div>Hello Profile</div>;
};

export default Home;
