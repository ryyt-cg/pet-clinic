import {PET_CLINIC_HOST_URL} from "@/gateway/host-url.ts";
import {useQuery, type UseQueryResult} from "@tanstack/react-query";
import type {Visit, Visits} from "@/gateway/visit-io.ts";

export function GetAllVisits(): UseQueryResult<Visits, Error> {
  return useQuery({
    queryKey: ["all-visits"],
    queryFn: fetchAllVisits,
  });
}

const fetchAllVisits = async (): Promise<Visits> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/visits/all`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `GetAllVisits fails - status: ${response.status}`;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const visits: Visits = await response.json();
  console.log(visits)
  return visits;
}

export function GetVisitById(id: number): UseQueryResult<Visit, Error> {
  return useQuery({
    queryKey: ["visitById", id],
    queryFn: () => fetchVisitById(id),
  });
}

const fetchVisitById = async (id: number): Promise<Visit> => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/visits/${id}`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    const errorMessage: string  = `GetVisitById fails - status: ${response.status}  `;
    console.error(errorMessage);
    throw new Error(errorMessage);
  }

  const pet: Visit = await response.json();
  console.log(pet)
  return pet;
}