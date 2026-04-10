import {
    useQuery,
} from "@tanstack/react-query";
import { createApiClient } from "../../api/ping";

const API_URL = import.meta.env.VITE_API_URL;
// const CREDENTIALS = import.meta.env.VITE_CREDENTIALS;

export const api = createApiClient(import.meta.env.API_URL || "http://localhost:8080", {
    axiosConfig: {
        timeout: 5000,
    },
});

export const usePing = () => {
    return useQuery({
        queryKey: ["ping"],
        queryFn: async () => await api.get("/api/ping"),
    });
};
