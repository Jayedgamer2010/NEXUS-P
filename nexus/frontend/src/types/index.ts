// ─── Status colors ─────────────────────────────────────────────────────────────

export const STATUS_COLORS: Record<string, string> = {
  running: '#10b981',
  offline: '#6b7280',
  installing: '#f59e0b',
  install_failed: '#ef4444',
  suspended: '#8b5cf6',
  starting: '#3b82f6',
  stopping: '#f97316',
} as const;

// ─── Power actions ──────────────────────────────────────────────────────────────

export const POWER_ACTIONS = {
  START: 'start',
  STOP: 'stop',
  RESTART: 'restart',
  KILL: 'kill',
} as const;

// ─── Pagination ────────────────────────────────────────────────────────────────

export interface PaginationMeta {
  total: number;
  per_page: number;
  current_page: number;
  last_page: number;
}

// ─── Generic API wrappers ──────────────────────────────────────────────────────

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface PaginatedData<T> {
  data: T[];
  meta: PaginationMeta;
}

// ─── Roles ─────────────────────────────────────────────────────────────────────

export type Role = "admin" | "client";

// ─── User ──────────────────────────────────────────────────────────────────────

export interface User {
  id: number;
  uuid: string;
  username: string;
  email: string;
  role: Role;
  root_admin: boolean;
  coins: number;
  name_first: string;
  name_last: string;
  language: string;
  created_at: string;
}

export interface UserBrief {
  id: number;
  username: string;
  email: string;
}

// ─── Server statuses ──────────────────────────────────────────────────────────

export type ServerStatus =
  | "running"
  | "offline"
  | "installing"
  | "install_failed"
  | "suspended"
  | "starting"
  | "stopping";

// ─── Power actions ────────────────────────────────────────────────────────────

export type PowerAction = "start" | "stop" | "restart" | "kill";

// ─── Allocation ────────────────────────────────────────────────────────────────

export interface Allocation {
  id: number;
  node_id: number;
  ip: string;
  ip_alias: string | null;
  port: number;
  server_id: number | null;
  assigned: boolean;
  server_name?: string;
  notes: string;
  created_at: string;
}

export interface AllocationBrief {
  id: number;
  ip: string;
  port: number;
}

// ─── Server ────────────────────────────────────────────────────────────────────

export interface Server {
  id: number;
  uuid: string;
  uuid_short: string;
  name: string;
  description: string;
  user_id?: number;
  status: ServerStatus;
  suspended: boolean;
  memory: number;
  disk: number;
  cpu: number;
  node?: NodeBrief;
  egg?: EggBrief;
  user?: UserBrief;
  allocation?: AllocationBrief;
  created_at: string;
}

// ─── Server resources (live stats from daemon) ────────────────────────────────

export interface ServerResources {
  cpu_absolute: number;
  memory_bytes: number;
  memory_limit_bytes: number;
  disk_bytes: number;
  disk_limit_bytes: number;
  network_rx_bytes: number;
  network_tx_bytes: number;
  state: ServerStatus;
  uptime: number;
}

// ─── Node ──────────────────────────────────────────────────────────────────────

export interface Node {
  id: number;
  uuid: string;
  public: boolean;
  name: string;
  description: string;
  location_id: number;
  fqdn: string;
  scheme: string;
  behind_proxy: boolean;
  maintenance_mode: boolean;
  memory: number;
  disk: number;
  daemon_listen: number;
  daemon_sftp: number;
  used_memory: number;
  used_disk: number;
  server_count: number;
  created_at: string;
}

export interface NodeBrief {
  id: number;
  uuid: string;
  name: string;
}

// ─── Egg ───────────────────────────────────────────────────────────────────────

export interface Egg {
  id: number;
  uuid: string;
  author: string;
  name: string;
  description: string;
  docker_image: string;
  docker_images: Record<string, string>;
  startup: string;
  created_at: string;
}

export interface EggBrief {
  id: number;
  uuid: string;
  name: string;
}

// ─── Admin Stats ──────────────────────────────────────────────────────────────

export interface NodeStatusEntry {
  node_id: number;
  name: string;
  status: "online" | "offline";
}

export interface AdminStats {
  users: number;
  nodes: number;
  servers: number;
  running_servers: number;
  node_statuses: NodeStatusEntry[];
}
