import { useState } from "react";
// import * as React from "react";
// import * as ReactDOM from "react-dom";
import Component from "./Component";
type Props = {
  initialCount: number;
};
function Home({ initialCount }: Props) {
  const [count, setCount] = useState(initialCount);

  return (
    <div className="App">
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
      <p className="read-the-docs">
        Click on the Rspack and React logos to learn more
      </p>
    </div>
  );
}

export default Home;
