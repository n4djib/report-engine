import { defineConfig, loadEnv } from "vite";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";

import viteReact from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default ({ mode }: { mode: string }) => {
  const env = loadEnv(mode, process.cwd(), "");
  // console.log("\n ----------", env.VITE_API_URL, "\n");
  // console.log(" ----------", env.VITE_FRONT_PORT, "\n");

  return defineConfig({
    resolve: { tsconfigPaths: true },
    plugins: [
      devtools(),
      tailwindcss(),
      tanstackRouter({ target: "react", autoCodeSplitting: true }),
      viteReact(),
    ],
    server: {
      port: Number(env.VITE_FRONT_PORT),
      proxy: {
        // "/api": {
        //   target: env.VITE_API_URL,
        //   changeOrigin: true,
        //   secure: false,
        // },
        "/api": env.VITE_API_URL,
      },
    },
  });
};

// const config = defineConfig({
//   resolve: { tsconfigPaths: true },
//   plugins: [
//     devtools(),
//     tailwindcss(),
//     tanstackRouter({ target: "react", autoCodeSplitting: true }),
//     viteReact(),
//   ],
//   server: {
//     // proxy: {
//     //   "/api": {
//     //     target: "http://localhost:8080",
//     //     changeOrigin: true,
//     //   },
//     // },
//     proxy: {
//       "/api": import.meta.env.VITE_API_URL,
//     },
//   },
// });

// export default config;
