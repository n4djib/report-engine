import { useQuery } from "@tanstack/react-query";
import { createApiClient } from "../../api/ping";
// import { env } from "#/env";

export const api = createApiClient(
  // env.VITE_API_URL || "http://localhost:8080",
  "/", // the api url is set in vite proxy
  {
    axiosConfig: {
      timeout: 5000,
    },
  },
);

export const usePing = () => {
  return useQuery({
    queryKey: ["ping"],
    queryFn: async () => await api.get("/api/ping"),
  });
};
