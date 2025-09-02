import LabeledInput from "@/components/courseInfo/LabeledInput";
import { toDatetimeLocal } from "@/utils/datetime";

type DateTimeFieldProps = {
  value?: string;
  className?: string;
  hideIfEmpty?: boolean;
};

const DateTimeField = ({
  value,
  className,
  hideIfEmpty,
}: DateTimeFieldProps) => {
  const v = toDatetimeLocal(value || "");
  if (hideIfEmpty && !v) return <>-</>;
  return (
    <LabeledInput
      type="datetime-local"
      step={1}
      defaultValue={v}
      className={`${className} w-fit`}
    />
  );
};

export default DateTimeField;
