import type { CourseInfo } from "@/types/course_info";
import TextOrHtml from "@/components/TextOrHtml";

interface Props {
  course: CourseInfo;
}

const isHtml = (str: string) => /<\/?[a-z][\s\S]*>/i.test(str.trim());

const CourseInfo = ({ course }: Props) => {
  return (
    <div className="w-full max-w-5xl mx-auto bg-white rounded-xl border border-gray-200 p-8 shadow space-y-6 mt-4 max-h-[70vh] overflow-y-auto">
      <h2 className="text-2xl font-semibold text-blue-700 mb-3">
        Course Information
      </h2>

      {/* Block fields, one per line */}
      <div className="space-y-6 text-sm">
        <div>
          <b>ชื่อหลักสูตร (TH):</b> {course.title_th}
        </div>
        <div>
          <b>ชื่อหลักสูตร (EN):</b> {course.title_en || "-"}
        </div>
        <div>
          <b>ดำเนินการโดย:</b> {course.organized_by}
        </div>
        <div>
          <b>จำนวนผู้เข้าร่วม:</b> {course.enroll_limit}
        </div>
        <div>
          <b>กลุ่มเป้าหมาย:</b> {course.target?.join(", ")}
        </div>
        <div>
          <b>เหตุผลในการจัดหลักสูตร:</b>
          <div className="space-y-4 mt-2">
            {course.rationale?.map((item, idx) => (
              <div key={idx} className="pl-4">
                {item}
              </div>
            ))}
          </div>
        </div>
        <div>
          <b>วัตถุประสงค์ของหลักสูตร:</b>
          <div className="space-y-4 mt-2">
            {course.objective?.map((item, idx) => (
              <div key={idx} className="pl-4">
                {item}
              </div>
            ))}
          </div>
        </div>
        {/* Content */}
        <div>
          <b>โครงสร้างหรือเนื้อหาของหลักสูตร:</b>
          <div className="space-y-4 mt-2">
            {course.content?.map((item, idx) =>
              isHtml(item) ? (
                <div
                  key={idx}
                  className="border border-gray-200 rounded px-3 py-2 bg-gray-50"
                >
                  <TextOrHtml text={item} />
                </div>
              ) : (
                <p key={idx} className="pl-2">
                  {item}
                </p>
              )
            )}
          </div>
        </div>

        {/* Evaluation */}
        <div>
          <b>การประเมินผลตลอดหลักสูตร:</b>
          <div className="space-y-4 mt-2">
            {course.evaluation?.map((item, idx) =>
              isHtml(item) ? (
                <div
                  key={idx}
                  className="border border-gray-200 rounded px-3 py-2 bg-gray-50"
                >
                  <TextOrHtml text={item} />
                </div>
              ) : (
                <p key={idx} className="pl-2">
                  {item}
                </p>
              )
            )}
          </div>
        </div>
        <div>
          <b>คำสำคัญสำหรับการสืบค้น:</b> {course.keywords?.join(", ")}
        </div>
        <div>
          <b>คำอธิบายหลักสูตรอย่างย่อ:</b>
          <div className="space-y-4 mt-2">
            {course.overview?.map((item, idx) => (
              <div key={idx} className="pl-4">
                {item}
              </div>
            ))}
          </div>
        </div>
        <div>
          <b>วันเริ่มต้นการลงทะเบียน:</b> {course.start_enroll?.join(", ")}
        </div>
        <div>
          <b>วันสิ้นสุดการลงทะเบียน:</b> {course.end_enroll?.join(", ")}
        </div>
        <div>
          <b>วันสิ้นสุดการชำระเงิน:</b> {course.payment_deadline?.join(", ")}
        </div>
        <div>
          <b>ค่าธรรมเนียมในการอบรม:</b> {course.fees?.join(", ")} บาท
        </div>
        <div>
          <b>ค่าบำรุงมหาวิทยาลัย:</b> {course.university_fees} บาท
        </div>
      </div>

      <div>
        <b>ข้อมูลในการติดต่อสอบถาม:</b>
        <div className="mt-2 space-y-2">
          {course.contacts?.length ? (
            course.contacts.map((contact, idx) => (
              <div
                key={idx}
                className="bg-gray-50 border border-gray-200 rounded p-3 shadow-sm text-sm"
              >
                {!contact.department && (
                  <div>
                    <b>ชื่อ-สกุล:</b> {contact.prefix}
                    {contact.name} {contact.surname}
                  </div>
                )}
                {contact.position && (
                  <div>
                    <b>ตำแหน่ง:</b> {contact.position}
                  </div>
                )}
                {contact.department && (
                  <div>
                    <b>หน่วยงาน:</b> {contact.department}
                  </div>
                )}
                {contact.email && (
                  <div>
                    <b>อีเมล:</b> {contact.email}
                  </div>
                )}
                {contact.phones && contact.phones.length > 0 && (
                  <div>
                    <b>โทรศัพท์:</b> {contact.phones.join(", ")}
                  </div>
                )}
                {contact.websites && contact.websites.length > 0 && (
                  <div>
                    <b>เว็บไซต์:</b> {contact.websites.join(", ")}
                  </div>
                )}
                {contact.address && (
                  <div>
                    <b>ที่อยู่:</b> {contact.address}
                  </div>
                )}
              </div>
            ))
          ) : (
            <div className="text-gray-400">-</div>
          )}
        </div>
      </div>

      <div>
        <b>หมวดหมู่การเรียนรู้:</b> {course.categories?.join(", ")}
      </div>
    </div>
  );
};

export default CourseInfo;
