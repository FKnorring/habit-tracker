"use client"

import React, { createContext, useContext, useEffect, useState } from 'react';
import { User, validateToken, logout as authLogout } from '@/lib/auth';
import { useRouter } from 'next/navigation';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (user: User) => void;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

interface AuthProviderProps {
  children: React.ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  const isAuthenticated = !!user;

  const login = (userData: User) => {
    setUser(userData);
  };

  const logout = () => {
    authLogout(); // Remove token from localStorage
    setUser(null);
    router.push('/login');
  };

  const refreshUser = async () => {
    try {
      const result = await validateToken();
      if (result.valid && result.user) {
        setUser(result.user);
      } else {
        setUser(null);
        authLogout(); // Clean up invalid token
      }
    } catch (error) {
      console.error('Failed to validate token:', error);
      setUser(null);
      authLogout(); // Clean up on error
    }
  };

  useEffect(() => {
    const initializeAuth = async () => {
      setIsLoading(true);
      await refreshUser();
      setIsLoading(false);
    };

    initializeAuth();
  }, []);

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated,
    login,
    logout,
    refreshUser,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
} 