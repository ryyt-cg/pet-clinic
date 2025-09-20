import type {Pet} from "@/gateway/pet-io.ts";
import type {Context} from "@/gateway/types.ts";

export interface Owner {
  id: number;
  firstName: string;
  lastName: string;
  address: string;
  city: string;
  telephone: string;
  pets?: Pet[];
}

export interface Owners {
  context: Context;
  owners: Owner[];
}

export interface NewOwner {
  firstName: string;
  lastName: string;
  address: string;
  city: string;
  telephone: string;
}

export interface UpdateOwner extends NewOwner {
  id: number;
}