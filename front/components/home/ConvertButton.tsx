"use client";
import { useMemo } from "react";

type Props = {
  courseId: string;
  file: File | null;
  loading: boolean;
  onConvert: () => void;
};

const ConvertButton = ({ courseId, file, loading, onConvert }: Props) => {
  const canConvert = useMemo(
    () => Boolean(courseId && file && !loading),
    [courseId, file, loading]
  );

  return (
    <button
      type="button"
      onClick={canConvert ? onConvert : undefined}
      disabled={!canConvert}
      aria-disabled={!canConvert}
      title={
        !courseId ? "กรอก Course ID ก่อน" : !file ? "เลือกไฟล์ .docx ก่อน" : ""
      }
      className={`w-xl mt-2 px-6 py-2 rounded-md transition 
        ${
          canConvert
            ? "bg-purple-600 text-white hover:bg-purple-700 cursor-pointer"
            : "bg-purple-300 text-white cursor-not-allowed opacity-60"
        }`}
    >
      Convert
    </button>
  );
};

export default ConvertButton;
