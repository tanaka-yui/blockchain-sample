import * as path from "path";
import { fileURLToPath } from "url";

import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv, UserConfig } from "vite";

const genCfg = (mode: string): UserConfig => {
  const env = loadEnv(mode, process.cwd());

  const cfg: UserConfig = {
    envDir: path.join(__dirname, "./"),
    root: "./",
    plugins: [react()],
    resolve: {
      alias: {
        "~": fileURLToPath(new URL("./src", import.meta.url)),
      },
    },
    build: {
      sourcemap: env.VITE_SOURCE_MAP === "true",
      chunkSizeWarningLimit: 1000,
    },
  };

  cfg.server = {
    host: "0.0.0.0",
    port: Number(env.VITE_PORT || 3000),
    open: true,
    proxy: {
      "^/api": {
        target: `http://${env.VITE_TARGET}`,
        changeOrigin: true,
        secure: true,
        xfwd: true,
        timeout: 5000,
        proxyTimeout: 5000,
      },
    },
  };

  return cfg;
};

export default defineConfig(({ mode }) => {
  return genCfg(mode);
});
