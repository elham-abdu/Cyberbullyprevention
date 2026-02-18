export interface User {
  ID: number;
  Email: string;
  Role: string;
  CreatedAt?: string;
  UpdatedAt?: string;
}

export interface Post {
  ID: number;
  UserID: number;
  Content: string;
  ToxicityScore: number;
  IsFlagged: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface LoginResponse {
  token: string;
}

export interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<boolean>;
  register: (email: string, password: string) => Promise<boolean>;
  logout: () => void;
  isAuthenticated: boolean;
  isAdmin: boolean;
}

export interface CreatePostInput {
  content: string;
}

export interface EditPostInput {
  post_id: number;
  content: string;
}

export interface DeletePostInput {
  post_id: number;
}

export interface ToxicityResponse {
  score: number;
  flagged: boolean;
}
