import { createFileRoute } from '@tanstack/react-router'
import { createApiClient } from "../api/ping";
import { usePing } from '#/lib/tanstack-query/ping';

export const Route = createFileRoute('/')({ component: App })

export const api = createApiClient(import.meta.env.VITE_API_URL || "http://localhost:8080", {
  axiosConfig: {
    timeout: 5000,
  },
});

function App() {
  const { data } = usePing();

  return (
    <>
      <main className="">
        <div>
          hello world from Central frontend
        </div>
        <div><b>Ping (API call):</b> {data?.message}</div>
      </main>
    </>
  )
}
