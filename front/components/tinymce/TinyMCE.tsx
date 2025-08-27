"use client";
import dynamic from "next/dynamic";

// ⛔️ ปิด SSR ของ TinyMCE editor ให้ทำงานเฉพาะฝั่ง client
export const Editor = dynamic(() => import("./TinyMCEClient"), {
  ssr: false,
});
