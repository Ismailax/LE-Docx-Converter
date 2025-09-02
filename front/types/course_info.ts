import type { Contact } from "./contact";

export type CourseInfo = {
  course_id: string;
  title_th: string;
  title_en: string;
  organized_by: string;
  enroll_limit: number;
  target: string[];
  rationale: string[];
  objective: string[];
  content: string[];
  evaluation: string[];
  keywords: string[];
  overview: string[];
  start_enroll: string[];
  end_enroll: string[];
  payment_deadline: string[];
  fees: number[];
  university_fees: number;
  contacts: Contact[];
  categories: string[];
};
