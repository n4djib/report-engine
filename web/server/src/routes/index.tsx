import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useState } from 'react';
import { createApiClient } from "../api/ping";

export const Route = createFileRoute('/')({ component: App })

// TODO make this configurable via env vars
export const api = createApiClient(import.meta.env.VITE_API_URL || "http://localhost:8080", {
  axiosConfig: {
    timeout: 5000,
  },
});

function App() {
  const [message, setMessage] = useState<any>("Loading...");

  useEffect(() => {
    async function load() {
      try {
        // This call is type-safe and runtime-validated!
        const response = await api.get("/api/ping");
        console.log("response: ", response)
        setMessage(response.message);

        // FIXME: this is just to test the type-safety and runtime validation, remove it later
        const name = response.name 
        console.log("--name:", name)

      } catch (err) {
        console.error("Validation failed or Network error:", err);
      }
    }
    load();
  }, []);

  return (
    <>
      <main className="">
      <div>
        hello world from Central frontend
      </div>
      <div><b>Ping (API call):</b> {message}</div>
      </main>
    </>

  )
}
