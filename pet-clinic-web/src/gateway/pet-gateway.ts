import {PET_CLINIC_HOST_URL} from "@/gateway/host-url.ts";
import {useQuery, type UseQueryResult} from "@tanstack/react-query";
import type {Pet, Pets} from "@/gateway/pet-io.ts";

export function GetAllPets(): UseQueryResult<Pets, Error> {
  return useQuery({
    queryKey: ["all-pets"],
    queryFn: fetchAllPets,
  });
}

const fetchAllPets = async (): Promise<Pets> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/pets/all`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `GetAllPets fails - status: ${response.status}`;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const pets: Pets = await response.json();
  console.log(pets)
  return pets;
}

export function GetPetById(id: number): UseQueryResult<Pet, Error> {
  return useQuery({
    queryKey: ["petById", id],
    queryFn: () => fetchPetById(id),
  });
}

const fetchPetById = async (id: number): Promise<Pet> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/pets/${id}`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `GetPetById fails - status: ${response.status}  `;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const pet: Pet = await response.json();
  console.log(pet)
  return pet;
}

export function GetPetWithVisitsById(id: number): UseQueryResult<Pet, Error> {
  return useQuery({
    queryKey: ["petWithVisitsById", id],
    queryFn: () => fetchPetWithVisitsById(id),
  });
}

const fetchPetWithVisitsById = async (id: number): Promise<Pet> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/pets/${id}/visits`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `GetPetWithVisitsById fails - status: ${response.status}  `;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const pet: Pet = await response.json();
  console.log(pet)
  return pet;
}