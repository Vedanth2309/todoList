import axios from "axios";
const api = axios.create({ baseURL: "/api" });
api.interceptors.request.use((c) => {
  const t = localStorage.getItem("token");
  if (t) c.headers.Authorization = "Bearer " + t;
  return c;
});
api.interceptors.response.use(
  (r) => r.data.data,
  (e) => {
    if (e.response?.status === 401 && !e.config.url.includes("/auth/")) {
      localStorage.removeItem("token");
      location.href = "/login";
    }
    return Promise.reject(
      new Error(
        e.response?.data?.message || "Network error — is the API running?",
      ),
    );
  },
);
export default api;
