const LabeledInput = (
  props: React.InputHTMLAttributes<HTMLInputElement> & {
    label?: string;
    className?: string;
    inputClassName?: string;
  }
) => {
  const { label, className, inputClassName, ...rest } = props;
  return (
    <label className={`flex flex-col gap-2 ${className || ""}`}>
      {label ? <span className="font-medium">{label}</span> : null}
      <input
        className={`border rounded px-3 py-2 ${inputClassName || ""}`}
        {...rest}
      />
    </label>
  );
};

export default LabeledInput;
