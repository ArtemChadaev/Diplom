export type UserRole = "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist";

export interface UserDTO {
  id: number;
  email: string;
  full_name: string;
  role: UserRole;
  ns_pv_access: boolean;
  ukep_bound: boolean;
  position?: string | null;
  department?: string | null;
  avatar_url?: string | null;
}
