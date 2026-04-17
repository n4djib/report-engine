/// <reference types="vite/client" />

interface ImportMetaEnv {
  // readonly VITE_APP_NAME: string;
  // readonly VITE_FRONT_PORT: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
