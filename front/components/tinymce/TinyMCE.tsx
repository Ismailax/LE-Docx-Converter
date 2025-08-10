"use client";
import dynamic from "next/dynamic";

// ⛔️ ปิด SSR ของ TinyMCE editor ให้ทำงานเฉพาะฝั่ง client
const Editor = dynamic(() => import("./TinyMCEClient"), {
  ssr: false,
});

export default Editor;
