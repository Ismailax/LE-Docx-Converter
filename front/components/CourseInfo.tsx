"use client";

import { useMemo } from "react";
import type { CourseInfo } from "@/types/course_info";
import Editor from "@/components/tinymce/TinyMCE";
import LabeledInput from "@/components/LabeledInput";
import Section from "@/components/Section";
import joinAsHtmlParagraphs from "@/utils/html";

interface Props {
  course: CourseInfo;
}

const CourseInfo = ({ course }: Props) => {
  const html = useMemo(
    () => ({
      target: joinAsHtmlParagraphs(course.target),
      rationale: joinAsHtmlParagraphs(course.rationale),
      objective: joinAsHtmlParagraphs(course.objective),
      content: joinAsHtmlParagraphs(course.content),
      evaluation: joinAsHtmlParagraphs(course.evaluation),
      keywords: joinAsHtmlParagraphs(course.keywords),
      overview: joinAsHtmlParagraphs(course.overview),
      start_enroll: joinAsHtmlParagraphs(course.start_enroll),
      end_enroll: joinAsHtmlParagraphs(course.end_enroll),
      payment_deadline: joinAsHtmlParagraphs(course.payment_deadline),
      categories: joinAsHtmlParagraphs(course.categories),
    }),
    [course]
  );
  console.log("course:", course);

  return (
    <div className="w-full max-w-5xl mx-auto bg-white rounded-xl border border-gray-200 p-8 shadow space-y-8 mt-4">
      <h2 className="text-2xl font-semibold text-blue-700">
        Course Information
      </h2>

      <LabeledInput label="ชื่อหลักสูตร (TH)" defaultValue={course.title_th} />
      <LabeledInput
        label="ชื่อหลักสูตร (EN)"
        defaultValue={course.title_en || ""}
      />
      <LabeledInput
        label="ดำเนินการโดย"
        defaultValue={course.organized_by}
        className="md:col-span-2"
      />
      <LabeledInput
        type="number"
        label="จำนวนผู้เข้าร่วม"
        defaultValue={String(course.enroll_limit)}
      />

      <Section title="กลุ่มเป้าหมาย">
        <Editor value={html.target} />
      </Section>
      <Section title="เหตุผลในการจัดหลักสูตร">
        <Editor value={html.rationale} />
      </Section>
      <Section title="วัตถุประสงค์ของหลักสูตร">
        <Editor value={html.objective} />
      </Section>
      <Section title="โครงสร้างหรือเนื้อหาของหลักสูตร">
        <Editor value={html.content} />
      </Section>
      <Section title="การประเมินผลตลอดหลักสูตร">
        <Editor value={html.evaluation} />
      </Section>
      <Section title="คำสำคัญสำหรับการสืบค้น (keyword)">
        <Editor value={html.keywords} />
      </Section>
      <Section title="คำอธิบายหลักสูตรอย่างย่อ">
        <Editor value={html.overview} />
      </Section>
      <Section title="วันเริ่มต้นการลงทะเบียน">
        <Editor value={html.start_enroll} />
      </Section>
      <Section title="วันสิ้นสุดการลงทะเบียน">
        <Editor value={html.end_enroll} />
      </Section>
      <Section title="วันสิ้นสุดการชำระเงิน">
        <Editor value={html.payment_deadline} />
      </Section>

      <Section title="ค่าธรรมเนียมในการอบรม (บาท)">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
          {(course.fees ?? []).length ? (
            <LabeledInput defaultValue={(course.fees ?? []).join(", ")} />
          ) : (
            <LabeledInput placeholder="0" />
          )}
        </div>
      </Section>

      <LabeledInput
        type="number"
        label="ค่าบำรุงมหาวิทยาลัย (บาท)"
        defaultValue={String(course.university_fees)}
      />

      <div className="space-y-3">
        <div className="font-medium">ข้อมูลในการติดต่อสอบถาม</div>
        {(course.contacts ?? []).length ? (
          course.contacts!.map((c, i) => (
            <div
              key={i}
              className="grid grid-cols-1 md:grid-cols-2 gap-3 border rounded p-3"
            >
              <LabeledInput label="คำนำหน้า" defaultValue={c.prefix} />
              <LabeledInput label="ชื่อ" defaultValue={c.name} />
              <LabeledInput label="นามสกุล" defaultValue={c.surname} />
              <LabeledInput label="ตำแหน่ง" defaultValue={c.position} />
              <LabeledInput label="หน่วยงาน" defaultValue={c.department} />
              <LabeledInput label="อีเมล" defaultValue={c.email} />
              <LabeledInput
                label="ที่อยู่"
                defaultValue={c.address}
                className="md:col-span-2"
              />
              <LabeledInput
                label="โทรศัพท์"
                defaultValue={(c.phones || []).join(", ")}
              />
              <LabeledInput
                label="เว็บไซต์"
                defaultValue={(c.websites || []).join(", ")}
              />
            </div>
          ))
        ) : (
          <div className="text-gray-500">-</div>
        )}
      </div>

      <Section title="หมวดหมู่การเรียนรู้">
        <Editor value={html.categories} />
      </Section>
    </div>
  );
};

export default CourseInfo;
