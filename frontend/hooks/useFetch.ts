import { isBrowser } from "@/utils/storage";

export const useFetch = () => {
  let controller
  let signal

  if (isBrowser) {
    controller = new AbortController();
    signal = controller.signal;
  }

  return {
    safeFetch: fetch,
    signal,
    controller
  }
}