import {PET_CLINIC_HOST_URL} from "@/gateway/host-url.ts";
import {useQuery, type UseQueryResult} from "@tanstack/react-query";
import type {Veterinarian, Veterinarians} from "@/gateway/vet-io.ts";

export function GetAllVets(): UseQueryResult<Veterinarians, Error> {
  return useQuery({
    queryKey: ["all-vets"],
    queryFn: fetchAllVets,
  });
}

const fetchAllVets = async (): Promise<Veterinarians> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/vets/all`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `GetAllVets fails - status: ${response.status}`;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const vets: Veterinarians = await response.json();
  console.log(vets)
  return vets;
}

export function GetVetById(id: number): UseQueryResult<Veterinarian, Error> {
  return useQuery({
    queryKey: ["vetById", id],
    queryFn: () => fetchVetById(id),
  });
}

const fetchVetById = async (id: number): Promise<Veterinarian> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/vets/${id}`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `GetVetById fails - status: ${response.status}  `;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const vet: Veterinarian = await response.json();
  console.log(vet)
  return vet;
}

export function GetVetWithSpecialtiesById(id: number): UseQueryResult<Veterinarian, Error> {
  return useQuery({
    queryKey: ["vetWithSpecialtiesById", id],
    queryFn: () => fetchVetWithSpecialtiesById(id),
  });
}

const fetchVetWithSpecialtiesById = async (id: number): Promise<Veterinarian> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/vets/${id}/specialties`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `fetchVetWithSpecialtiesById fails - status: ${response.status}  `;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const vet: Veterinarian = await response.json();
  console.log(vet)
  return vet;
}