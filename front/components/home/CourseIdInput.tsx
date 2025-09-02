type Props = {
  value: string;
  onChange: (v: string) => void;
  disabled?: boolean;
};

const CourseIdInput = ({ value, onChange, disabled }: Props) => {
  return (
    <input
      type="text"
      placeholder="Enter Course ID"
      value={value}
      onChange={(e) => onChange(e.target.value)}
      disabled={disabled}
      className="w-full border rounded-md px-3 py-2 text-center disabled:opacity-60"
    />
  );
};

export default CourseIdInput;
