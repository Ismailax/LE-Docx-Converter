import katex from "katex";

// แยก $$...$$ ก่อน แล้วค่อย $...$
const renderMathToHtml = (input: string): string => {
  if (!input) return "";

  // Block: $$...$$ (รองรับหลายบรรทัด)
  const blocks = input.split(/(\$\$[\s\S]+?\$\$)/g).map((part) => {
    const m = part.match(/^\$\$([\s\S]+?)\$\$$/);
    if (!m) {
      // Inline: $...$
      return part
        .split(/(\$[^$]+\$)/g)
        .map((p) => {
          const mi = p.match(/^\$([^$]+)\$$/);
          if (!mi) return p; // text ธรรมดา
          return katex.renderToString(mi[1], {
            throwOnError: false,
            displayMode: false,
          });
        })
        .join("");
    }
    // block
    return katex.renderToString(m[1], {
      throwOnError: false,
      displayMode: true,
    });
  });

  return blocks.join("");
};

export default renderMathToHtml;
