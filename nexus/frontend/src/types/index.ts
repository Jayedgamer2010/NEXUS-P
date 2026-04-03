// User type
export interface User {
  id: number;
  uuid: string;
  username: string;
  email: string;
  role: 'admin' | 'client';
  coins: number;
  root_admin: boolean;
  created_at: string;
  updated_at: string;
}

// Node type
export interface Node {
  id: number;
  uuid: string;
  name: string;
  fqdn: string;
  scheme: 'http' | 'https';
  wings_port: number;
  memory: number;
  memory_overalloc: number;
  disk: number;
  disk_overalloc: number;
  token_id: string;
  created_at: string;
  updated_at: string;
}

// Egg type
export interface Egg {
  id: number;
  uuid: string;
  name: string;
  description?: string;
  docker_image: string;
  startup_command: string;
  created_at: string;
  updated_at: string;
}

// Allocation type
export interface Allocation {
  id: number;
  node_id: number;
  ip: string;
  port: number;
  assigned: boolean;
  server_id?: number;
  created_at: string;
  updated_at: string;
}

// Server type
export interface Server {
  id: number;
  uuid: string;
  name: string;
  user_id: number;
  node_id: number;
  egg_id: number;
  allocation_id: number;
  memory: number;
  disk: number;
  cpu: number;
  status: 'installing' | 'running' | 'stopped' | 'error';
  suspended: boolean;
  created_at: string;
  updated_at: string;
  node?: Node;
  egg?: Egg;
  allocation?: Allocation;
}

// Server resources from Wings
export interface ServerResources {
  cpu: number;
  memory: number;
  memory_max: number;
  disk: number;
  disk_max: number;
  uptime: number;
}

// Ticket type
export interface Ticket {
  id: number;
  user_id: number;
  subject: string;
  status: 'open' | 'closed' | 'pending';
  priority: 'low' | 'medium' | 'high';
  created_at: string;
  updated_at: string;
}

// CoinTransaction type
export interface CoinTransaction {
  id: number;
  user_id: number;
  amount: number;
  reason: string;
  created_at: string;
}

// API response wrapper
export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data: T;
  total?: number;
  page?: number;
  limit?: number;
}

// Paginated response
export interface PaginatedResponse<T> {
  success: boolean;
  message: string;
  data: T[];
  total: number;
  page: number;
  limit: number;
}

// Power action enum
export enum PowerAction {
  START = 'start',
  STOP = 'stop',
  RESTART = 'restart',
  KILL = 'kill'
}

// Server status colors
export const STATUS_COLORS: Record<string, string> = {
  running: '#10b981',    // green
  stopped: '#ef4444',    // red
  installing: '#f59e0b', // yellow
  error: '#dc2626',      // dark red
  suspended: '#6b7280',  // gray
};
