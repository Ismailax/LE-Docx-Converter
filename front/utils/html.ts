const isBlockHtml = (t: string) =>
  /^\s*<\s*(p|table|thead|tbody|tr|th|td|ul|ol|li|blockquote|div|section|article|h[1-6])\b/i.test(
    (t || "").trim()
  );

const joinAsHtmlParagraphs = (items?: string[]) => {
  if (!items || items.length === 0) return "";

  return items
    .map((s) => (s ?? "").trim())
    .filter((s) => s.length > 0) // ตัดช่องว่าง/บรรทัดว่างทิ้ง
    .map((s) => {
      if (isBlockHtml(s)) return s; // เป็น block HTML อยู่แล้ว → ปล่อยผ่าน
      // มี inline tag เล็กๆ ได้ แต่ไม่ใช่ block → ครอบเป็น <p>
      return `<p>${s.replace(/\n/g, "<br/>")}</p>`;
    })
    .join("");
};

export default joinAsHtmlParagraphs;
