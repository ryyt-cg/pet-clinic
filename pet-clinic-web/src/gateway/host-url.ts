import axios from "axios";

export const PET_CLINIC_HOST_URL = 'http://local.capgroup.com:8092/api/pet-clinic';

export const axiosClient = axios.create({
  baseURL: PET_CLINIC_HOST_URL,
  headers: {
    "Access-Control-Allow-Origin": "*",
    "Content-Type": "application/json",
  }
});