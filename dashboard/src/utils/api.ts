import wretch, { Wretch } from "wretch";

const BASE_URL = "http://localhost:8080";

export const api: Wretch = wretch(BASE_URL)
  .options({ credentials: "include" })
  .middlewares([
    (next) => async (url, opts) => {
      try {
        return await next(url, opts);
      } catch (error) {
        console.log("Error:", error);
        if (
          error instanceof Error &&
          "status" in error &&
          (error as any).status === 401
        ) {
          try {
            await wretch(BASE_URL)
              .url("/auth/refresh")
              .options({ credentials: "include" })
              .post()
              .res();

            return next(url, opts);
          } catch (refreshError) {
            window.location.href = "/login";
            throw refreshError;
          }
        }
        throw error;
      }
    },
  ]);

// Define the response types
export interface UserInfoResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
}
