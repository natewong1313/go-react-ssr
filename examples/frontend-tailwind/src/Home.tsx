import { useState } from "react";
import { IndexRouteProps } from "./generated";
import GoLogo from "../public/go.png";
import ReactLogo from "../public/react.png";
import Counter from "./components/Counter";

function Home({ initialCount }: IndexRouteProps) {
  const [count, setCount] = useState(initialCount);

  return (
    <div className="h-screen items-center flex-col flex justify-center bg-zinc-800">
      <div className="flex gap-6 items-center">
        <img src={GoLogo} alt="Go logo" className="h-[70px] w-[170px]"/>
        <img src={ReactLogo} alt="React logo"  className="h-[80px] w-[90px]"/>
      </div>
      <h1 className="pt-6 mt-6 text-white text-5xl font-semibold">Go + React</h1>
      <div className="py-7">
        <Counter count={count} increment={() => setCount(count + 1)} />
      </div>
      <a className="text-zinc-500 hover:underline" href="https://github.com/natewong1313/go-react-ssr" target="_blank">
        View project on GitHub
      </a>
    </div>
  );
}

export default Home;
