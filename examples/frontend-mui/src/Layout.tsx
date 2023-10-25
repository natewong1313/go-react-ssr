import createCache from "@emotion/cache";
import { CacheProvider } from "@emotion/react";

const cache = createCache({ key: "css" })

export default function Layout({ children }: { children: React.ReactNode }) {
  console.log("Hello from Layout");
  return <CacheProvider value={cache}>{children}</CacheProvider>;
}
