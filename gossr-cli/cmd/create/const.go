package create

const PACKAGE_JSON = `{
	"dependencies": {
	  "react": "^18.2.0",
	  "react-dom": "^18.2.0"
	},
	"devDependencies": {
	  "@types/react": "^18.2.21",
	  "@types/react-dom": "^18.2.7"
	}
}`

const TSCONFIG = `{
	"compilerOptions": {
	  "target": "ES6",
	  "lib": ["DOM", "DOM.Iterable", "ESNext"],
	  "module": "ESNext",
	  "skipLibCheck": true,
	  "moduleResolution": "bundler",
	  "allowImportingTsExtensions": true,
	  "resolveJsonModule": true,
	  "isolatedModules": true,
	  "noEmit": true,
	  "jsx": "react-jsx",
	  "strict": true,
	  "noUnusedLocals": true,
	  "noUnusedParameters": true,
	  "noFallthroughCasesInSwitch": true
	},
	"include": ["src"]
}`

const TAILWIND_CONFIG = `/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{jsx,tsx,js,ts}"],
  theme: {
    extend: {},
  },
  plugins: [],
}`

const TAILWIND_CSS = `@tailwind base;
@tailwind components;
@tailwind utilities;`

const REACT_FILE = `import { useState } from "react";
import "./index.css";

function Index({ initialCount }: IndexRouteProps) {
  const [count, setCount] = useState(initialCount);
  return (
    <div className="home">
      <div className="img-container">
        <img
          src="https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg"
          alt="Go logo"
          height={70}
          width={170}
        />
        <img
          src="https://upload.wikimedia.org/wikipedia/commons/a/a7/React-icon.svg"
          alt="React logo"
          height={80}
          width={90}
        />
      </div>
      <h1>Go + React</h1>
      <a href="https://github.com/natewong1313/go-react-ssr" target="_blank">
        View project on GitHub
      </a>
    </div>
  );
}

export default Index;`

const INDEX_CSS = `body {
	background-color: #242424;
	margin: 0 !important;
	font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
}

.home {
	height: 100vh;
	display: flex;
	justify-content: center;
	align-items: center;
	flex-direction: column;
}

.img-container {
	display: flex;
	gap: 24px;
	align-items: center;
}

h1 {
	padding-top: 20px;
	color: #fff;
	font-size: 42px;
	margin-bottom: 0;
}

.counter-container {
	padding: 28px 0px;
}

a {
	color: #888;
	text-decoration: none;
}

a:hover {
	text-decoration: underline;
}
`
