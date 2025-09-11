"use client";
import { ChangeEvent, useRef } from "react";

type Props = {
  onPick: (file: File | null) => void;
  disabled?: boolean;
};

const FilePicker = ({ onPick, disabled }: Props) => {
  const ref = useRef<HTMLInputElement>(null);

  const onChange = (e: ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files && e.target.files[0] ? e.target.files[0] : null;
    onPick(f);
  };

  return (
    <>
      <label
        htmlFor="file-upload"
        className={`w-fit text-white px-6 py-2 rounded-md shadow cursor-pointer transition text-base text-center
          ${
            disabled
              ? "bg-purple-300 cursor-not-allowed"
              : "bg-purple-600 hover:bg-purple-700"
          }`}
      >
        Choose File
      </label>
      <input
        id="file-upload"
        type="file"
        accept=".docx"
        className="hidden"
        ref={ref}
        onChange={onChange}
        disabled={disabled}
      />
    </>
  );
};

export default FilePicker;
