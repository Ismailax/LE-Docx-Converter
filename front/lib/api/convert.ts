import api from "@/lib/axios";
import type { CourseInfo } from "@/types/course_info";

/**
 * อัปโหลดไฟล์ .docx ไปยัง backend และรับข้อมูลหลักสูตร
 * @param courseId รหัสหลักสูตร
 * @param file ไฟล์ docx ที่เลือก
 * @returns ข้อมูลหลักสูตร (CourseInfo)
 * @throws error ถ้าอัปโหลดไม่สำเร็จ
 */
export async function convertCourseDoc(
  courseId: string,
  file: File
): Promise<CourseInfo> {
  const formData = new FormData();
  formData.append("file", file);

  const res = await api.post(`/convert/${courseId}`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  return res.data;
}
