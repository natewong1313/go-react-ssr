import { useState } from "react";
import { IndexRouteProps } from "./generated";
import GoLogo from "../public/go.png";
import ReactLogo from "../public/react.png";
import "./Home.css";
import Counter from "./components/Counter";

function Home({ initialCount }: IndexRouteProps) {
  const [count, setCount] = useState(initialCount);

  return (
    <div className="home">
      <div className="img-container">
        <img src={GoLogo} alt="Go logo" height={70} width={170} />
        <img src={ReactLogo} alt="React logo" height={80} width={90} />
      </div>
      <h1>Go + React</h1>
      <div className="counter-container">
        <Counter count={count} increment={() => setCount(count + 1)} />
      </div>
      <a href="https://github.com/natewong1313/go-react-ssr" target="_blank">
        View project on GitHub
      </a>
    </div>
  );
}

export default Home;
