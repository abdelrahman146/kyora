import ky from "ky";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/v1";

export const apiClient = ky.create({
  prefixUrl: API_BASE_URL,
  timeout: 30000,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = sessionStorage.getItem("access_token");
        if (token) {
          request.headers.set("Authorization", `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      async (_request, _options, response) => {
        if (response.status === 401) {
          window.location.href = "/login";
        }
        return response;
      },
    ],
  },
});

export default apiClient;
