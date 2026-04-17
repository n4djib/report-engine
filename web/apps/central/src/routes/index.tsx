import { createFileRoute } from "@tanstack/react-router";
import { createApiClient } from "../api/ping";
import { usePing } from "#/lib/tanstack-query/ping";
import { env } from "#/env";
import SharedComponent from "@packages/ui-components/SharedComponent";

export const Route = createFileRoute("/")({ component: App });

export const api = createApiClient(env.VITE_API_URL, {
  axiosConfig: {
    timeout: 5000,
  },
});

function App() {
  const { data } = usePing();

  console.log("ping central:", data?.message);

  return (
    <>
      <main className="">
        <div>hello world from Central frontend</div>
        <div>
          <b>Ping (API call):</b> {data?.message}
        </div>
        imported from shared : <SharedComponent />
      </main>
    </>
  );
}
