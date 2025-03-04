import { useEffect, useState } from "react";
import { Outlet, useNavigate } from "react-router";
import { useAuth } from "../context/AuthProvider";
import useRefreshToken from "../hooks/useRefreshToken";

export default function LoadSession() {
  const [isLoading, setLoading] = useState<boolean>(true);
  const { tokens } = useAuth();
  const refreshTokens = useRefreshToken();
  const navigate = useNavigate();

  useEffect(() => {
    async function verifyRefreshToken() {
      try {
        await refreshTokens();
      } catch (err) {
        console.log(err);
        navigate("/login");
      } finally {
        setLoading(false);
      }
    }

    if (!tokens?.access) {
      navigate("/login");
      setLoading(false);
    } else {
      verifyRefreshToken();
    }
  }, []);

  return <>{isLoading ? <p> loading... </p> : <Outlet />}</>;
}
