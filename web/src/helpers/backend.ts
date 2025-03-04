import axios from "axios";
const BASE_URL: string = import.meta.env.VITE_BACKEND_URL;

const backend = axios.create({
  baseURL: BASE_URL,
  headers: { "Content-Type": "application/json" },
});

export default backend;
