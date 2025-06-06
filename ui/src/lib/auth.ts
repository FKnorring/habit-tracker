const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface User {
  id: string;
  email: string;
  username: string;
  created_at: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterCredentials {
  email: string;
  password: string;
  username: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

// Helper function to get auth headers
function getAuthHeaders(): Record<string, string> {
  const token = localStorage.getItem('auth_token');
  return {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` })
  };
}

// Register a new user
export async function register(credentials: RegisterCredentials): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Registration failed');
  }

  const data = await response.json();
  
  // Store token in localStorage
  if (data.token) {
    localStorage.setItem('auth_token', data.token);
  }
  
  return data;
}

// Login user
export async function login(credentials: LoginCredentials): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Login failed');
  }

  const data = await response.json();
  
  // Store token in localStorage
  if (data.token) {
    localStorage.setItem('auth_token', data.token);
  }
  
  return data;
}

// Get user profile
export async function getProfile(): Promise<User> {
  const response = await fetch(`${API_BASE_URL}/auth/profile`, {
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    throw new Error('Failed to get profile');
  }

  return response.json();
}

// Validate token
export async function validateToken(): Promise<{ valid: boolean; user?: User }> {
  const token = localStorage.getItem('auth_token');
  if (!token) {
    return { valid: false };
  }

  try {
    const response = await fetch(`${API_BASE_URL}/auth/validate`, {
      headers: getAuthHeaders(),
    });

    if (!response.ok) {
      return { valid: false };
    }

    const data = await response.json();
    return { valid: data.valid, user: data.user };
  } catch (error: unknown) {
    console.error('Error validating token:', error);
    return { valid: false };
  }
}

// Logout user
export function logout(): void {
  localStorage.removeItem('auth_token');
}

// Check if user is authenticated
export function isAuthenticated(): boolean {
  return !!localStorage.getItem('auth_token');
}

// Get stored auth token
export function getAuthToken(): string | null {
  return localStorage.getItem('auth_token');
} 