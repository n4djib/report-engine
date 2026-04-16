import { env } from "#/env";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({ component: App });

function App() {
  return <div>Hello world from remote port:{env.VITE_FRONT_PORT}</div>;
}
