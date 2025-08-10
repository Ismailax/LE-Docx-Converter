const LabeledInput = (
  props: React.InputHTMLAttributes<HTMLInputElement> & {
    label?: string;
    className?: string;
  }
) => {
  const { label, className, ...rest } = props;
  return (
    <label className={`flex flex-col gap-1 ${className || ""}`}>
      {label ? <span className="font-medium">{label}</span> : null}
      <input className="border rounded px-3 py-2" {...rest} />
    </label>
  );
};

export default LabeledInput;
