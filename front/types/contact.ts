export interface Contact {
  prefix: string;
  name: string;
  surname: string;
  position: string;
  department: string;
  address: string;
  phones: string[] | null;
  email: string;
  websites: string[] | null;
}