import { createContext, useEffect, useState } from "react";
import { Tokens, UserProfile } from "../models/user";
import React from "react";
import backend from "../helpers/backend";

type AuthContextType = {
  user: UserProfile | null;
  tokens: Tokens | null;
  setUser: (user: UserProfile) => void;
  setTokens: (tokens: Tokens) => void;
  clearAuthContext: () => void;
  isLoggedIn: () => boolean;
};

type Props = { children: React.ReactNode };

const AuthContext = createContext<AuthContextType>({} as AuthContextType);

export const AuthProvider = ({ children }: Props) => {
  const [user, setUser] = useState<UserProfile | null>(null);
  const [tokens, setTokens] = useState<Tokens | null>(null);
  const [isReady, setIsReady] = useState(false);

  // load user and tokens from local storage
  useEffect(() => {
    const user = localStorage.getItem("user");
    const tokens = localStorage.getItem("tokens");
    if (user && tokens) {
      setTokens(JSON.parse(tokens));
      setUser(JSON.parse(user));
    }
    setIsReady(true);
  }, []);

  // update local storage and load user
  useEffect(() => {
    if (!tokens) return;
    localStorage.setItem("tokens", JSON.stringify(tokens));
    backend.defaults.headers.common["Authorization"] =
      "Bearer " + tokens?.access;
    if (!user) {
      console.log("fetching user");
      backend.get<UserProfile>("/me").then((res) => {
        if (res) {
          setUser(res.data);
        }
      });
    }
  }, [tokens]);

  useEffect(() => {
    if (!user) return;
    localStorage.setItem("user", JSON.stringify(user));
  }, [user]);

  function clearAuthContext() {
    setUser(null);
    setTokens(null);
    localStorage.removeItem("tokens");
    localStorage.removeItem("user");
  }

  function isLoggedIn() {
    return !!user;
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        tokens,
        setTokens: (tokens) => setTokens(tokens),
        setUser: (user) => setUser(user),
        clearAuthContext,
        isLoggedIn,
      }}
    >
      {isReady ? children : null}
    </AuthContext.Provider>
  );
};

export const useAuth = () => React.useContext(AuthContext);
