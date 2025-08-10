import { InlineMath, BlockMath } from "react-katex";

interface MathTextProps {
  text: string;
  className?: string;
}

// ตรวจว่ามีทั้ง block/inline
const containsMath = (s: string) => /\$\$[\s\S]+?\$\$|\$[^$]+\$/m.test(s);

function renderMathInText(text: string) {
  // แยก block $$...$$ ก่อน (ครอบคลุม newline)
  const blockParts = text.split(/(\$\$[\s\S]+?\$\$)/g);

  return blockParts.map((block, i) => {
    if (/^\$\$[\s\S]+\$\$$/.test(block.trim())) {
      const math = block.replace(/^\$\$|\$\$$/g, "").trim();
      return (
        <div key={`bm-${i}`} style={{ display: "block" }}>
          <BlockMath math={math} />
        </div>
      );
    }

    // จากนั้นแยก inline $...$
    const inlineParts = block.split(/(\$[^$]+\$)/g);
    return inlineParts.map((part, j) => {
      if (/^\$[^$]+\$$/.test(part.trim())) {
        const math = part.replace(/^\$|\$$/g, "").trim();
        return <InlineMath key={`im-${i}-${j}`} math={math} />;
      }
      return <span key={`tx-${i}-${j}`}>{part}</span>;
    });
  });
}

const MathText = ({ text, className }: MathTextProps) => {
  if (!text) return null;
  return (
    <span className={className}>
      {containsMath(text) ? renderMathInText(text) : text}
    </span>
  );
};

export { containsMath, MathText };
