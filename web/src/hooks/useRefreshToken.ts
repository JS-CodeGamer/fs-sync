import { useNavigate } from "react-router";
import { useAuth } from "../context/AuthProvider";
import backend from "../helpers/backend";
import { Tokens } from "../models/user";

export default function useRefreshToken() {
  const { tokens, setTokens } = useAuth();
  return async () => {
    return backend
      .get<Tokens>("/refresh-token", {
        headers: {
          Authorization: "Bearer " + tokens?.refresh,
        },
      })
      .then((res) => {
        setTokens(res?.data);
      })
      .catch(() => {
        const navigate = useNavigate();
        navigate(
          "/login?return=" + encodeURIComponent(window.location.pathname),
        );
      });
  };
}
