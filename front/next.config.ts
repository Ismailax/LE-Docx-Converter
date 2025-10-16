import type { NextConfig } from "next";

const basePath = process.env.NEXT_PUBLIC_APP_BASEPATH || "";

const nextConfig: NextConfig = {
  output: "standalone",
  basePath,
};

export default nextConfig;
