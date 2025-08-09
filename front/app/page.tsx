"use client";
import { useRef, useState } from "react";
import Loading from "@/components/Loading";
import Failed from "@/components/Failed";
import CourseInfo from "@/components/CourseInfo";
import type { CourseInfo as CourseInfoType } from "@/types/course_info";
import { convertCourseDoc } from "@/lib/api/convert";

export default function Home() {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [fileName, setFileName] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);
  const [course, setCourse] = useState<CourseInfoType | null>(null);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) {
      setFileName("");
      return;
    }
    setFileName(files[0].name);
    setLoading(true);
    setError(false);
    setCourse(null);
    try {
      const data = await convertCourseDoc(files[0]);
      setCourse(data);
      setLoading(false);
    } catch {
      setError(true);
      setLoading(false);
    }
  };

  const handleReset = () => {
    setFileName("");
    setLoading(false);
    setError(false);
    setCourse(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  return (
    <div className="py-10 flex justify-center items-center min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-100">
      <div className="bg-white shadow-2xl rounded-2xl p-8 min-w-[350px] max-w-3xl w-full space-y-6 border border-gray-200">
        <h1 className="text-2xl font-bold text-center text-blue-700 mb-2">
          CMU Lifelong Education
        </h1>
        <p className="text-center text-gray-600 mb-4">
          Upload a course document{" "}
          <span className="font-semibold">(.docx)</span> to extract and display
          course information.
        </p>

        {!loading && !course && !error && (
          <div className="flex flex-col items-center gap-3">
            <label
              htmlFor="file-upload"
              className="bg-blue-600 text-white px-6 py-2 rounded-lg shadow cursor-pointer hover:bg-blue-700 transition text-base"
            >
              Choose File
            </label>
            <input
              id="file-upload"
              type="file"
              accept=".docx"
              className="hidden"
              ref={fileInputRef}
              onChange={handleFileChange}
            />
            {fileName && (
              <p className="text-xs text-gray-500">
                Selected file: <span className="font-semibold">{fileName}</span>
              </p>
            )}
          </div>
        )}

        {loading && <Loading />}
        {error && <Failed />}
        {course && <CourseInfo course={course} />}

        {(loading || error || course) && (
          <button
            className="block w-full mt-4 px-4 py-2 rounded border bg-gray-200 hover:bg-gray-300 text-gray-700 transition font-medium"
            onClick={handleReset}
          >
            Reset
          </button>
        )}
      </div>
    </div>
  );
}
