import { useState } from "react";
// import * as React from "react";
// import * as ReactDOM from "react-dom";
import Component from "./components/Component";
import { IndexRouteProps } from "./generated";
import "../public/test.css";
import horseImg from "../public/horse.png";

function Home({ initialCount, msg }: IndexRouteProps) {
  const [count, setCount] = useState(initialCount);

  return (
    <div className="App">
      <img src={horseImg} height={100} width={150} />
      <Component />
      <div>
        <a href="https://reactjs.org" target="_blank"></a>
      </div>
      <h1>React + Go</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">{msg}</p>
    </div>
  );
}

export default Home;
