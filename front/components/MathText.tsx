import { InlineMath, BlockMath } from "react-katex";

/**
 * แยก text เป็น BlockMath, InlineMath, และ text ธรรมดา
 */
function renderMathInText(text: string) {
  // ตัด display math ก่อน ($$...$$)
  const blockParts = text.split(/(\$\$[\s\S]+?\$\$)/g);
  return blockParts.map((block, i) => {
    if (/^\$\$[\s\S]+\$\$$/.test(block.trim())) {
      // Display math (block)
      return (
        <BlockMath key={i} math={block.replace(/^\$\$|\$\$$/g, "").trim()} />
      );
    } else {
      // ตัด inline math ($...$)
      const inlineParts = block.split(/(\$[^$]+\$)/g);
      return inlineParts.map((part, j) =>
        /^\$[^$]+\$$/.test(part.trim()) ? (
          <InlineMath key={j} math={part.replace(/^\$|\$$/g, "").trim()} />
        ) : (
          <span key={j}>{part}</span>
        )
      );
    }
  });
}

interface MathTextProps {
  text: string;
  className?: string;
}
const MathText: React.FC<MathTextProps> = ({ text, className }) => {
  return <span className={className}>{renderMathInText(text)}</span>;
};
export default MathText;
