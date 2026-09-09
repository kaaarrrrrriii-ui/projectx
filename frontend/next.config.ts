import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  async redirects() {
    return [
      {
        source: "/appeal/appeal2",
        destination: "/appeal/topics",
        permanent: true,
      },
      {
        source: "/appeal/appeal2/appeal3",
        destination: "/appeal/description",
        permanent: true,
      },
      {
        source: "/appeal/appeal2/appeal3/appeal4",
        destination: "/appeal/details",
        permanent: true,
      },
      {
        source: "/appeal/appeal2/appeal3/appeal4/mediaAdd",
        destination: "/appeal/attachments",
        permanent: true,
      },
      {
        source: "/appeal/appeal2/appeal3/appeal4/mediaAdd/success",
        destination: "/appeal/success",
        permanent: true,
      },
    ];
  },
};

export default nextConfig;

