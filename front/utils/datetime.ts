export function toDatetimeLocal(value?: string | null) {
  if (!value) return "";
  // รับรูปแบบ "2024-08-01 08:30:00" หรือ "2024-08-01T08:30:00"
  const v = value.trim().replace(" ", "T");
  // ให้แน่ใจว่ามีวินาที
  const hasSeconds = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}$/.test(v);
  const hasMinutes = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/.test(v);
  if (hasSeconds) return v;
  if (hasMinutes) return v + ":00";
  // ถ้าเป็นแค่วัน
  const onlyDate = /^\d{4}-\d{2}-\d{2}$/.test(v);
  if (onlyDate) return v + "T00:00:00";
  return "";
}
