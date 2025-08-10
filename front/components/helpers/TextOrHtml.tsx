import { containsMath, MathText } from "@/components/helpers/MathText";
import parse, { HTMLReactParserOptions } from "html-react-parser";

interface TextOrHtmlProps {
  text: string;
}

const isHtml = (text: string) => /<\/?[a-z][\s\S]*>/i.test(text.trim());

const TextOrHtml = ({ text }: TextOrHtmlProps) => {
  if (!text) return null;
  if (isHtml(text)) {
    const options: HTMLReactParserOptions = {
      replace(domNode) {
        // แปลงเฉพาะ text node ที่มี math
        if (domNode.type === "text") {
          const data: string | undefined = domNode.data;
          if (data && containsMath(data)) {
            return <MathText text={data} />;
          }
        }
        return undefined; // คง tag/attribute เดิมทั้งหมด
      },
    };
    return <>{parse(text, options)}</>;
  }
  // ไม่ใช่ HTML → ให้ MathText จัดการ ($...$, $$...$$)
  return <MathText text={text} />;
};

export { isHtml, TextOrHtml };
