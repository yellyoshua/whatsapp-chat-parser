import process from "process"

export const IS_DEVELOPMENT: boolean = String(process.env.NODE_ENV).includes("development");

export const API_URL: string = process.env.SERVER || "http://localhost:3000"