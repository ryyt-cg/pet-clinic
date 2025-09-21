import type {Context} from "@/gateway/types.ts";

export interface Pet {
  id: number;
  name: string;
  birthdate: string;
  species: string;
  ownerId: number;
}

export interface Pets {
  context: Context;
  pets: Pet[];
}

export interface NewPet {
  name: string;
  birthdate: string;
  species: string;
  ownerId: number;
}

export interface UpdatePet extends NewPet {
  id: number;
}