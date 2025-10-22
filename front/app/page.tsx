"use client";
import { useRef, useState } from "react";
import Image from "next/image";
import Loading from "@/components/Loading";
import Failed from "@/components/Failed";
import CourseInformation from "@/components/CourseInformation";
import type { CourseInfo as CourseInfoType } from "@/types/course_info";
import { convertCourseDoc } from "@/lib/api/convert";
import Card from "@/components/home/Card";
import CourseIdInput from "@/components/home/CourseIdInput";
import FilePicker from "@/components/home/FilePicker";
import SelectedFileLabel from "@/components/home/SelectedFileLabel";
import StepItem from "@/components/home/StepItem";
import ConvertButton from "@/components/home/ConvertButton";
import ResetButton from "@/components/home/ResetButton";

const base = process.env.NEXT_PUBLIC_APP_BASEPATH || "";

const Home = () => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [courseId, setCourseId] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);
  const [course, setCourse] = useState<CourseInfoType | null>(null);

  const handleConvert = async () => {
    if (!courseId || !file) return;
    setLoading(true);
    setError(false);
    setCourse(null);
    try {
      const data = await convertCourseDoc(courseId, file);
      setCourse(data);
    } catch {
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    setCourseId("");
    setFile(null);
    setLoading(false);
    setError(false);
    setCourse(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const showPickers = !loading && !course && !error;

  return (
    <div className="py-10 flex justify-center items-center min-h-screen bg-[#F8F4FF] from-blue-50 via-white to-purple-100">
      <div className="w-full h-full space-y-6">
        <div className="flex flex-col items-center space-y-4">
          <Image
            src={`${base}/logo.png`}
            alt="Logo"
            width={240}
            height={240}
            priority
          />
          <h1 className="text-2xl font-bold text-purple-600">
            Course Document Converter
          </h1>
          <p className="text-slate-600">
            This is a demonstration frontend for testing the document conversion
            API, <br /> not the official Lifelong Education admin interface.
          </p>
        </div>

        {showPickers && (
          <div className="space-y-6">
            <Card className="max-w-xl mx-auto">
              <StepItem step={1}>
                <label className="block text-base font-medium text-slate-700">
                  Course ID (กรอกเลขอะไรก็ได้)
                </label>
                <CourseIdInput value={courseId} onChange={setCourseId} />
              </StepItem>
            </Card>

            <Card className="max-w-xl mx-auto">
              <StepItem step={2}>
                <label className="block text-base font-medium text-slate-700">
                  <p>เลือกไฟล์หลักสูตร</p>
                  <p className="text-sm text-slate-500">
                    (รองรับเฉพาะไฟล์ .docx ขนาดไม่เกิน 10 MB)
                  </p>
                </label>
                <FilePicker onPick={setFile} />
                <SelectedFileLabel file={file ?? undefined} />
              </StepItem>
            </Card>

            <div className="max-w-xl mx-auto text-center">
              <ConvertButton
                courseId={courseId}
                file={file}
                loading={loading}
                onConvert={handleConvert}
              />
              {!courseId || !file ? (
                <p className="text-xs text-amber-600 mt-2">
                  กรุณากรอก Course ID และเลือกไฟล์ก่อนกด Convert
                </p>
              ) : null}
            </div>
          </div>
        )}

        {loading && <Loading />}
        {error && <Failed />}
        {course && <CourseInformation course={course} />}

        <ResetButton show={Boolean(error || course)} onReset={handleReset} />
      </div>
    </div>
  );
};

export default Home;
