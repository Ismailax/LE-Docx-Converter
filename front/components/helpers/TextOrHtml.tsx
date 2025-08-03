import MathText from "@/components/helpers/MathText";
import parse, { HTMLReactParserOptions } from "html-react-parser";

interface TextOrHtmlProps {
  text: string;
}

const isHtml = (text: string) => /<\/?[a-z][\s\S]*>/i.test(text.trim());

const TextOrHtml = ({ text }: TextOrHtmlProps) => {
  if (isHtml(text)) {
    const options: HTMLReactParserOptions = {
      replace(domNode) {
        // ใส่ MathText เฉพาะ text node ที่เจอ $...$ หรือ $$...$$
        // @ts-ignore
        if (domNode.type === "text" && domNode.data?.match(/\$.*?\$/)) {
          // @ts-ignore
          return <MathText text={domNode.data} />;
        }
        return undefined;
      },
    };
    return <>{parse(text, options)}</>;
  }
  return <MathText text={text} />;
};

export { isHtml, TextOrHtml };
