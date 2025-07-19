interface TextOrHtmlProps {
  text: string;
}

const TextOrHtml = ({ text }: TextOrHtmlProps) => {
  // เช็คถ้าเป็น HTML tag จะ render แบบ html, ถ้าไม่ใช่ให้แสดงเป็น text
  if (/</.test(text) && /<\/?[a-z][\s\S]*>/i.test(text)) {
    return <span dangerouslySetInnerHTML={{ __html: text }} />;
  }
  return <span>{text}</span>;
};

export default TextOrHtml;
