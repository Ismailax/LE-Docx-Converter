type Props = { file?: File | null };

const SelectedFileLabel = ({ file }: Props) => {
  if (!file) return null;
  return (
    <p className="text-sm text-slate-600">
      Selected file: <span className="font-semibold">{file.name}</span>
    </p>
  );
};

export default SelectedFileLabel;
