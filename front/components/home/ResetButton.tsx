type Props = {
  show: boolean;
  onReset: () => void;
};

const ResetButton = ({ show, onReset }: Props) => {
  if (!show) return null;
  return (
    <button
      className="block w-full max-w-5xl mx-auto mt-4 px-4 py-2 rounded-lg bg-gray-200 hover:bg-gray-300 text-gray-700 transition font-medium"
      onClick={onReset}
    >
      Reset
    </button>
  );
};

export default ResetButton;
