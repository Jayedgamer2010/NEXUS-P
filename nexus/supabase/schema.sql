-- NEXUS Supabase Schema
-- Run this SQL in your Supabase SQL editor

create table nexus_users (
  id uuid primary key default gen_random_uuid(),
  email text unique not null,
  username text unique not null,
  password_hash text not null,
  pterodactyl_id integer unique,
  role text default 'client',
  coins integer default 0,
  suspended boolean default false,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

create table tickets (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references nexus_users(id) on delete cascade,
  subject text not null,
  status text default 'open',
  priority text default 'low',
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

create table ticket_messages (
  id uuid primary key default gen_random_uuid(),
  ticket_id uuid references tickets(id) on delete cascade,
  user_id uuid references nexus_users(id),
  message text not null,
  is_staff boolean default false,
  created_at timestamptz default now()
);

create table coin_transactions (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references nexus_users(id) on delete cascade,
  amount integer not null,
  reason text not null,
  created_at timestamptz default now()
);

create table resource_packages (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  description text,
  memory integer default 0,
  disk integer default 0,
  cpu integer default 0,
  coin_cost integer not null,
  created_at timestamptz default now()
);

create table user_resources (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references nexus_users(id) on delete cascade unique,
  extra_memory integer default 0,
  extra_disk integer default 0,
  extra_cpu integer default 0,
  extra_servers integer default 0,
  updated_at timestamptz default now()
);

create table announcements (
  id uuid primary key default gen_random_uuid(),
  title text not null,
  content text not null,
  type text default 'info',
  active boolean default true,
  created_at timestamptz default now()
);

create table nexus_settings (
  key text primary key,
  value text not null
);

insert into nexus_settings (key, value) values
  ('site_name', 'NEXUS'),
  ('coins_per_day', '10'),
  ('pterodactyl_url', ''),
  ('pterodactyl_api_key', ''),
  ('time_check_interval', '30'),
  ('default_node_max_slots', '4'),
  ('free_first_start_minutes', '10'),
  ('cooldown_minutes', '2');

-- ================================================================
-- Queue-Based Server Time System
-- ================================================================

-- Server time credits
create table server_time_credits (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references nexus_users(id) on delete cascade,
  server_uuid text not null,
  minutes_remaining integer default 0,
  total_minutes_purchased integer default 0,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Server sessions (tracks active running time)
create table server_sessions (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references nexus_users(id) on delete cascade,
  server_uuid text not null,
  node_id text not null,
  started_at timestamptz,
  ends_at timestamptz,
  status text default 'queued',
  queue_position integer,
  cooldown_until timestamptz,
  first_start boolean default true,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Node slot tracking
create table node_slots (
  id uuid primary key default gen_random_uuid(),
  node_id text unique not null,
  max_active integer default 4,
  current_active integer default 0,
  updated_at timestamptz default now()
);

-- Time packages
create table time_packages (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  minutes integer not null,
  coin_cost integer not null,
  created_at timestamptz default now()
);

-- First start tracking
create table server_first_starts (
  server_uuid text primary key,
  used boolean default false,
  created_at timestamptz default now()
);

-- Insert default time packages
insert into time_packages (name, minutes, coin_cost) values
  ('Starter', 5, 10),
  ('Basic', 15, 25),
  ('Standard', 20, 50),
  ('Premium', 30, 100);

-- Indexes for time system
create index idx_server_time_credits_user_id on server_time_credits(user_id);
create index idx_server_time_credits_server_uuid on server_time_credits(server_uuid);
create index idx_server_sessions_status on server_sessions(status);
create index idx_server_sessions_server_uuid on server_sessions(server_uuid);
create index idx_server_sessions_node_id on server_sessions(node_id);
create index idx_server_sessions_ends_at on server_sessions(ends_at);
create index idx_node_slots_node_id on node_slots(node_id);
create index idx_time_packages_order on time_packages(coin_cost);

-- RLS policies for time system
alter table server_time_credits enable row level security;
alter table server_sessions enable row level security;
alter table node_slots enable row level security;
alter table time_packages enable row level security;
alter table server_first_starts enable row level security;

-- Indexes for common queries
create index idx_nexus_users_email on nexus_users(email);
create index idx_nexus_users_username on nexus_users(username);
create index idx_tickets_user_id on tickets(user_id);
create index idx_tickets_status on tickets(status);
create index idx_ticket_messages_ticket_id on ticket_messages(ticket_id);
create index idx_coin_transactions_user_id on coin_transactions(user_id);
create index idx_user_resources_user_id on user_resources(user_id);
create index idx_announcements_active on announcements(active);

-- RLS policies (adjust as needed)
alter table nexus_users enable row level security;
alter table tickets enable row level security;
alter table ticket_messages enable row level security;
alter table coin_transactions enable row level security;
alter table resource_packages enable row level security;
alter table user_resources enable row level security;
alter table announcements enable row level security;
alter table nexus_settings enable row level security;
