type StepItemProps = {
  step: number | string;
  children: React.ReactNode;
};

const StepItem = ({ step, children }: StepItemProps) => (
  <div className="px-6 py-5 flex items-start gap-4">
    <span className="flex-none inline-flex items-center justify-center size-7 rounded-full bg-purple-100 text-purple-700 text-sm font-semibold">
      {step}
    </span>
    <div className="flex flex-col gap-2">{children}</div>
  </div>
);

export default StepItem;
