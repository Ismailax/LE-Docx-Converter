import api from "@/lib/axios";
import type { CourseInfo } from "@/types/course_info";

const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL;

/**
 * อัปโหลดไฟล์ .docx ไปยัง backend และรับข้อมูลหลักสูตร
 * @param file ไฟล์ docx ที่เลือก
 * @returns ข้อมูลหลักสูตร (CourseInfo)
 * @throws error ถ้าอัปโหลดไม่สำเร็จ
 */
export async function uploadCourseDoc(file: File): Promise<CourseInfo> {
  const formData = new FormData();
  formData.append("file", file);

  try {
    const res = await api.post("/convert", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return res.data;
  } catch (error) {
    throw new Error("Failed to upload document");
  }
}
