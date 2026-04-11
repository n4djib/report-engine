import { makeApi, Zodios, type ZodiosOptions } from "@zodios/core";
import { z } from "zod";

const endpoints = makeApi([
  {
    method: "get",
    path: "/api/ping",
    alias: "ping_pong",
    requestFormat: "json",
    response: z.object({ message: z.string() }).strict(),
  },
]);

export const PingApi = new Zodios(endpoints);

export function createApiClient(baseUrl: string, options?: ZodiosOptions) {
  return new Zodios(baseUrl, endpoints, options);
}
