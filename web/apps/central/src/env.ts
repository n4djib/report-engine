import { z } from "zod";

const envSchema = z.object({
  VITE_API_URL: z.string().url(),
  VITE_APP_NAME: z.string().min(1),
  VITE_FRONT_PORT: z.number(),
});

const _env = {
  VITE_API_URL: import.meta.env.VITE_API_URL,
  VITE_APP_NAME: import.meta.env.VITE_APP_NAME,
  VITE_FRONT_PORT: Number(import.meta.env.VITE_FRONT_PORT), // docker sets environments as a string
};

// Validate at runtime
let env: z.infer<typeof envSchema>;

try {
  env = envSchema.parse(_env);
} catch (err) {
  console.error("❌ Invalid environment variables:", err);

  // Show error in UI
  document.body.innerHTML = `
    <div style="font-family: sans-serif; padding: 20px;">
        <h1 style="color: red;">Environment Error</h1>
        <pre style="
        background: #111;
        color: #0f0;
        padding: 15px;
        border-radius: 8px;
        overflow: auto;
        ">${JSON.stringify(err, null, 2)}</pre>
    </div>
    `;

  throw err; // still crash (good for dev)
}

export { env };
