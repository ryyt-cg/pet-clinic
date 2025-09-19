import {PET_CLINIC_HOST_URL} from "@/gateway/host-url.ts";
import {useQuery} from "@tanstack/react-query";

export function GetAllPets() {
  return useQuery({
    queryKey: ["all-pets"],
    queryFn: fetchAllPets,
  });
}

const fetchAllPets = async () => {
  const response = await fetch(`${PET_CLINIC_HOST_URL}/v1/pets/all`, {
    headers: {
      "Access-Control-Allow-Origin": "*",
    },
  });
  if (!response.ok) {
    throw new Error("Network response was not ok");
  }
  console.log(response.json());
  return response.json();
}

