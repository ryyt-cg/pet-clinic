import {PET_CLINIC_HOST_URL} from "@/gateway/host-url.ts";
import {useQuery, type UseQueryResult} from "@tanstack/react-query";
import type {Owner, Owners} from "@/gateway/owner-io.ts";

export function GetAllOwners(): UseQueryResult<Owners, Error> {
  return useQuery({
    queryKey: ["all-owners"],
    queryFn: fetchAllPets,
  });
}

const fetchAllPets = async (): Promise<Owners> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/owners/all`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `GetAllOwners fails - status: ${response.status}`;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const owners: Owners = await response.json();
  console.log(owners)
  return owners;
}

export function GetOwnerById(id: number): UseQueryResult<Owner, Error> {
  return useQuery({
    queryKey: ["ownerById", id],
    queryFn: () => fetchOwnerById(id),
  });
}

const fetchOwnerById = async (id: number): Promise<Owner> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/owners/${id}`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `GetOwnerById fails - status: ${response.status}  `;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const owner: Owner = await response.json();
  console.log(owner)
  return owner;
}

export function GetOwnerWithPetsById(id: number): UseQueryResult<Owner, Error> {
  return useQuery({
    queryKey: ["ownerWithPetsById", id],
    queryFn: () => fetchOwnerWithPetsById(id),
  });
}

const fetchOwnerWithPetsById = async (id: number): Promise<Owner> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/owners/${id}/pets`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `GetOwnerWithPetsById fails - status: ${response.status}  `;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const owner: Owner = await response.json();
  console.log(owner)
  return owner;
}