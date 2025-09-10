package types

type Output struct {
	CourseID        string    `json:"course_id"`        // รหัสหลักสูตร
	TitleTH         string    `json:"title_th"`         // ชื่อหลักสูตรภาษาไทย
	TitleEN         string    `json:"title_en"`         // ชื่อหลักสูตรภาษาอังกฤษ
	OrganizedBy     string    `json:"organized_by"`     // หน่วยงานที่เปิดหลักสูตร
	EnrollLimit     int       `json:"enroll_limit"`     // จำนวนรับสมัคร
	Target          []string  `json:"target"`           // กลุ่มเป้าหมาย
	Rationale       []string  `json:"rationale"`        // หลักการและเหตุผล
	Objective       []string  `json:"objective"`        // วัตถุประสงค์
	Content         []string  `json:"content"`          // โครงสร้างหรือเนื้อหาของหลักสูตร
	Evaluation      []string  `json:"evaluation"`       // การวัดและประเมินผล
	Keywords        []string  `json:"keywords"`         // คำสำคัญสำหรับการสืบค้น
	Overview        []string  `json:"overview"`         // คำอธิบายหลักสูตรอย่างย่อ
	StartEnroll     []string  `json:"start_enroll"`     // วันเปิดรับสมัคร
	EndEnroll       []string  `json:"end_enroll"`       // วันปิดรับสมัคร
	PaymentDeadline []string  `json:"payment_deadline"` // วันสิ้นสุดชำระเงิน
	Fees            []int     `json:"fees"`             // ค่าธรรมเนียมการอบรม
	UniversityFee   int       `json:"university_fees"`  // ค่าบำรุงมหาวิทยาลัย
	Contacts        []Contact `json:"contacts"`         // ผู้ประสานงานหลักสูตร
	Categories      []string  `json:"categories"`       // หมวดหมู่การเรียนรู้
}

type Contact struct {
	Prefix     string   `json:"prefix"`
	Name       string   `json:"name"`
	Surname    string   `json:"surname"`
	Position   string   `json:"position"`
	Department string   `json:"department"`
	Address    string   `json:"address"`
	Phones     []string `json:"phones"`
	Email      string   `json:"email"`
	Websites   []string `json:"websites"`
}
