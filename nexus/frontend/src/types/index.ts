export interface User {
  id: number
  uuid: string
  username: string
  email: string
  name_first: string
  name_last: string
  role: string
  root_admin: boolean
  coins: number
  suspended: boolean
  created_at: string
  servers?: Server[]
}

export interface Node {
  id: number
  uuid: string
  name: string
  description: string
  fqdn: string
  scheme: string
  behind_proxy: boolean
  maintenance_mode: boolean
  memory: number
  memory_overallocate: number
  disk: number
  disk_overallocate: number
  daemon_listen: number
  daemon_sftp: number
  daemon_token_id: string
  created_at: string
  allocations?: Allocation[]
}

export interface Egg {
  id: number
  uuid: string
  author: string
  name: string
  description: string
  docker_image: string
  startup: string
  config_stop: string
  created_at: string
}

export interface Allocation {
  id: number
  node_id: number
  ip: string
  ip_alias: string
  port: number
  server_id: number | null
  notes: string
}

export interface Server {
  id: number
  uuid: string
  uuid_short: string
  name: string
  description: string
  status: string
  suspended: boolean
  memory: number
  disk: number
  cpu: number
  image: string
  user_id: number
  node_id: number
  egg_id: number
  allocation_id: number
  user?: User
  node?: Node
  egg?: Egg
  allocation?: Allocation
  created_at: string
}

export interface ServerResources {
  cpu_absolute: number
  memory_bytes: number
  disk_bytes: number
  network_rx_bytes: number
  network_tx_bytes: number
  state: string
  uptime: number
}

export interface AdminStats {
  users: number
  nodes: number
  servers: number
  running_servers: number
  recent_servers: Server[]
  recent_users: User[]
}

export interface Ticket {
  id: number
  user_id: number
  subject: string
  status: string
  priority: string
  created_at: string
  user?: User
}

export interface ApiResponse<T = any> {
  success: boolean
  message: string
  data: T
  errors?: Record<string, string>
}

export interface PaginatedMeta {
  total: number
  per_page: number
  current_page: number
  last_page: number
  from: number
  to: number
}

export interface PaginatedData<T> {
  data: T[]
  meta: PaginatedMeta
}

export type PowerAction = 'start' | 'stop' | 'restart' | 'kill'
export type ServerStatus = 'installing' | 'install_failed' | 'suspended' | 'running' | 'offline'
