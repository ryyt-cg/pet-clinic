import type {Context, Specialty} from "@/gateway/types.ts";

export interface Veterinarian {
  id: number;
  firstName: string;
  lastName: string;
  specialties?: Specialty[];
}

export interface Veterinarians {
  context: Context;
  owners: Veterinarian[];
}
