import axios from "axios";
import type { CourseInfo } from "@/types/course_info";

/**
 * อัปโหลดไฟล์ .docx ไปยัง backend และรับข้อมูลหลักสูตร
 * @param file ไฟล์ docx ที่เลือก
 * @returns ข้อมูลหลักสูตร (CourseInfo)
 * @throws error ถ้าอัปโหลดไม่สำเร็จ
 */
export async function uploadCourseDoc(file: File): Promise<CourseInfo> {
  const formData = new FormData();
  formData.append("file", file);

  const res = await axios.post(
    process.env.NEXT_PUBLIC_BACKEND_URL + "/convert",
    formData,
    {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      timeout: 60000,
    }
  );
  return res.data;
}
