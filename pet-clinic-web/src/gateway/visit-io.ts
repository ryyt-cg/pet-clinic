import type {Context} from "@/gateway/types.ts";

export interface Visit {
  id: number;
  date: string;
  description: string;
  petId: number;
}

export interface Visits {
  context: Context;
  visits: Visit[];
}

export interface NewVisit {
  date: string;
  description: string;
  petId: number;
}

export interface UpdateVisit extends NewVisit {
  id: number;
}