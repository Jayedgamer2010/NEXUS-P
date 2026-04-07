-- ===================================================================
-- NEXUS Dashboard - Complete Supabase PostgreSQL Schema
-- Run this SQL in your Supabase SQL Editor at:
--   https://supabase.com/dashboard/project/YOUR_PROJECT/sql
--
-- This schema uses Supabase Auth (auth.users).
-- Do NOT manage passwords manually - Supabase handles that.
-- ===================================================================

-- Profile protection is handled at the app layer (PHP backend).
-- The service_role key bypasses RLS, so backend can freely manage coins/role.
-- Users self-update through the API which only allows whitelisted fields.

-- ===================================================================
-- 1. USER PROFILES (extends auth.users)
--    Every user registered via Supabase Auth gets a profile automatically
--    via a trigger on auth.users.
-- ===================================================================

create table if not exists user_profiles (
  id uuid primary key references auth.users(id) on delete cascade,
  email text unique not null,
  username text unique not null,
  first_name text default '',
  last_name text default '',
  pterodactyl_id integer,
  role text default 'client' check (role in ('client', 'admin', 'moderator')),
  coins integer default 0,
  suspended boolean default false,
  two_factor_secret text,
  two_factor_enabled boolean default false,
  avatar text,
  background text,
  company_name text,
  vat_number text,
  address1 text,
  address2 text,
  city text,
  country text,
  state text,
  postcode text,
  referral_code text unique,
  referred_by uuid references user_profiles(id),
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Auto-create profile when a user signs up via Supabase Auth
create or replace function handle_new_user()
returns trigger as $$
begin
  insert into user_profiles (id, email, username, referral_code)
  values (
    new.id,
    new.email,
    coalesce(new.raw_user_meta_data->>'username', split_part(new.email, '@', 1)),
    'ref_' || substr(md5(random()::text || clock_timestamp()::text), 1, 10)
  );
  return new;
end;
$$ language plpgsql security definer;

-- Attach trigger on auth.users (may already exist from a previous run)
drop trigger if exists on_auth_user_created on auth.users;
create trigger on_auth_user_created
  after insert on auth.users
  for each row execute function handle_new_user();

-- Enable Supabase RLS
alter table user_profiles enable row level security;

-- Users can read their own profile
create policy "Users can view own profile"
  on user_profiles for select
  using (auth.uid() = id);

-- Users can update their own profile
create policy "Users can update own profile"
  on user_profiles for update
  using (auth.uid() = id);

-- Users can insert their own profile during signup (used by auth trigger)
create policy "Users can insert own profile during signup"
  on user_profiles for insert
  with check (auth.uid() = id);
-- Admins can view all profiles
create policy "Admins can view all profiles"
  on user_profiles for select
  using (
    exists (
      select 1 from user_profiles as up
      where up.id = auth.uid() and up.role = 'admin'
    )
  );

-- Admins can insert profiles (needed for auth trigger and manual inserts)
create policy "Admins can insert profiles"
  on user_profiles for insert
  with check (
    exists (
      select 1 from user_profiles as up
      where up.id = auth.uid() and up.role = 'admin'
    )
  );

-- Admins can update all profiles
create policy "Admins can update all profiles"
  on user_profiles for update
  using (
    exists (
      select 1 from user_profiles as up
      where up.id = auth.uid() and up.role = 'admin'
    )
  );

-- Allow the auth trigger to insert profiles (bypasses RLS during signup)
create policy "Auth trigger can insert profiles"
  on user_profiles for insert
  with check (true);


-- ===================================================================
-- 2. COIN TRANSACTIONS (audit log)
-- ===================================================================

create table if not exists coin_transactions (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete cascade,
  amount integer not null,
  reason text not null,
  metadata jsonb default '{}'::jsonb,
  created_at timestamptz default now()
);

alter table coin_transactions enable row level security;

create policy "Users can view own transactions"
  on coin_transactions for select
  using (auth.uid() = user_id);

create policy "Admins can view all transactions"
  on coin_transactions for select
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Server can insert transactions"
  on coin_transactions for insert
  with check (true);

create index idx_coin_transactions_user_id on coin_transactions(user_id);
create index idx_coin_transactions_created on coin_transactions(created_at desc);


-- ===================================================================
-- 3. USER RESOURCES (per-user resource limits)
-- ===================================================================

create table if not exists user_resources (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete cascade unique,
  extra_memory integer default 0,
  extra_disk integer default 0,
  extra_cpu integer default 0,
  extra_servers integer default 0,
  extra_backups integer default 0,
  extra_databases integer default 0,
  extra_allocations integer default 0,
  updated_at timestamptz default now()
);

alter table user_resources enable row level security;

create policy "Users can view own resources"
  on user_resources for select
  using (auth.uid() = user_id);

create policy "Admins can manage resources"
  on user_resources for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

-- Admins need select to check role, so add a separate admin select policy
create policy "Admins can view all resources"
  on user_resources for select
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create index idx_user_resources_user_id on user_resources(user_id);


-- ===================================================================
-- 4. RESOURCE PACKAGES (buyable resource upgrades)
-- ===================================================================

create table if not exists resource_packages (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  description text,
  memory integer default 0,
  disk integer default 0,
  cpu integer default 0,
  servers integer default 0,
  coin_cost integer not null,
  active boolean default true,
  created_at timestamptz default now()
);

alter table resource_packages enable row level security;

create policy "Anyone can read active packages"
  on resource_packages for select
  using (active = true);

create policy "Admins can manage packages"
  on resource_packages for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 5. ANNOUNCEMENTS
-- ===================================================================

create table if not exists announcements (
  id uuid primary key default gen_random_uuid(),
  title text not null,
  content text not null,
  type text default 'info' check (type in ('info', 'warning', 'success', 'error', 'maintenance')),
  active boolean default true,
  sort_order integer default 0,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

create table if not exists announcements_assets (
  id uuid primary key default gen_random_uuid(),
  announcement_id uuid references announcements(id) on delete cascade,
  url text not null,
  type text default 'image',
  created_at timestamptz default now()
);

create table if not exists announcements_tags (
  id uuid primary key default gen_random_uuid(),
  announcement_id uuid references announcements(id) on delete cascade,
  tag text not null
);

alter table announcements enable row level security;

create policy "Anyone can view active announcements"
  on announcements for select
  using (active = true);

create policy "Admins can manage announcements"
  on announcements for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create index idx_announcements_active on announcements(active, sort_order);


-- ===================================================================
-- 6. SETTINGS (key-value config)
-- ===================================================================

create table if not exists nexus_settings (
  key text primary key,
  value text not null,
  description text,
  updated_at timestamptz default now()
);

insert into nexus_settings (key, value, description) values
  ('site_name', 'NEXUS', 'Site display name'),
  ('coins_per_day', '10', 'Daily coin reward amount'),
  ('allow_registration', 'true', 'Whether new users can register'),
  ('allow_servers', 'true', 'Whether users can create servers'),
  ('pterodactyl_url', '', 'Pterodactyl panel URL'),
  ('time_check_interval', '30', 'Time manager check interval in seconds'),
  ('default_node_max_slots', '4', 'Default max active servers per node'),
  ('free_first_start_minutes', '10', 'Free minutes for first server start'),
  ('cooldown_minutes', '2', 'Cooldown minutes after server expiry'),
  ('turnstile_key_pub', '', 'Cloudflare Turnstile public key'),
  ('turnstile_key_secret', '', 'Cloudflare Turnstile secret key'),
  ('turnstile_enabled', 'false', 'Enable/disable Cloudflare Turnstile'),
  ('discord_client_id', '', 'Discord OAuth client ID'),
  ('discord_client_secret', '', 'Discord OAuth client secret'),
  ('github_client_id', '', 'GitHub OAuth client ID'),
  ('github_client_secret', '', 'GitHub OAuth client secret')
on conflict (key) do nothing;

alter table nexus_settings enable row level security;

create policy "Anyone can read settings"
  on nexus_settings for select
  using (true);

create policy "Admins can manage settings"
  on nexus_settings for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 7. LOCATIONS (node groups for server creation)
-- ===================================================================

create table if not exists locations (
  id bigserial primary key,
  name text not null,
  description text,
  slots integer default 0,
  used_slots integer default 0,
  pterodactyl_location_id bigint,
  node_ip text,
  status text default 'online',
  vip_only boolean default false,
  image text,
  locked boolean default false,
  deleted boolean default false,
  updated_at timestamptz default now(),
  created_at timestamptz default now()
);

alter table locations enable row level security;

create policy "Anyone can read active locations"
  on locations for select
  using (deleted = false);

create policy "Admins can manage locations"
  on locations for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 8. EGG CATEGORIES (nest categories in Pterodactyl)
-- ===================================================================

create table if not exists egg_categories (
  id bigserial primary key,
  name text not null,
  description text,
  pterodactyl_nest_id bigint,
  image text,
  enabled boolean default true,
  vip_only boolean default false,
  locked boolean default false,
  deleted boolean default false,
  updated_at timestamptz default now(),
  created_at timestamptz default now()
);

alter table egg_categories enable row level security;

create policy "Anyone can read active categories"
  on egg_categories for select
  using (deleted = false and enabled = true);

create policy "Admins can manage categories"
  on egg_categories for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 9. EGGS (server types/templates)
-- ===================================================================

create table if not exists eggs (
  id bigserial primary key,
  name text not null,
  description text,
  category_id bigint references egg_categories(id) on delete cascade,
  pterodactyl_egg_id bigint,
  enabled boolean default true,
  vip_only boolean default false,
  locked boolean default false,
  deleted boolean default false,
  image text,
  updated_at timestamptz default now(),
  created_at timestamptz default now()
);

alter table eggs enable row level security;

create policy "Anyone can read active eggs"
  on eggs for select
  using (deleted = false and enabled = true);

create policy "Admins can manage eggs"
  on eggs for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 10. SERVERS (user server tracking - NOT Pterodactyl itself)
-- ===================================================================

create table if not exists servers (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete cascade,
  name text not null,
  description text,
  uuid text unique not null,
  pterodactyl_server_id bigint,
  node_id bigint references locations(id),
  egg_id bigint references eggs(id),
  category_id bigint references egg_categories(id),
  memory integer default 1024,
  cpu integer default 100,
  disk integer default 1024,
  databases integer default 1,
  backups integer default 1,
  allocations integer default 1,
  status text default 'pending',
  deleted boolean default false,
  locked boolean default false,
  updated_at timestamptz default now(),
  created_at timestamptz default now()
);

alter table servers enable row level security;

create policy "Users can view own servers"
  on servers for select
  using (auth.uid() = user_id and deleted = false);

create policy "Admins can manage servers"
  on servers for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create index idx_servers_user_id on servers(user_id);
create index idx_servers_uuid on servers(uuid);
create index idx_servers_deleted on servers(deleted);


-- ===================================================================
-- 11. SERVER CREATION QUEUE
-- ===================================================================

create table if not exists servers_queue (
  id bigserial primary key,
  user_id uuid references user_profiles(id) on delete cascade,
  name text not null,
  description text default '',
  location_id bigint,
  category_id bigint,
  egg_id bigint,
  memory integer default 1024,
  cpu integer default 100,
  disk integer default 1024,
  databases integer default 1,
  backups integer default 1,
  allocations integer default 1,
  status text default 'pending',
  deleted boolean default false,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

alter table servers_queue enable row level security;

create policy "Users can view own queue items"
  on servers_queue for select
  using (auth.uid() = user_id and deleted = false);

create policy "Admins can manage queue"
  on servers_queue for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 12. TICKETS
-- ===================================================================

create table if not exists tickets (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete cascade,
  department_id bigint,
  subject text not null,
  status text default 'open' check (status in ('open', 'answered', 'closed', 'on-hold', 'admin_reply')),
  priority text default 'low',
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

alter table tickets enable row level security;

create policy "Users can view own tickets"
  on tickets for select
  using (auth.uid() = user_id);

create policy "Admins can view all tickets"
  on tickets for select
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Users can create tickets"
  on tickets for insert
  with check (auth.uid() = user_id);

create policy "Users can update own tickets"
  on tickets for update
  using (auth.uid() = user_id);

create index idx_tickets_user_id on tickets(user_id);
create index idx_tickets_status on tickets(status);


-- ===================================================================
-- 13. TICKET MESSAGES
-- ===================================================================

create table if not exists ticket_messages (
  id uuid primary key default gen_random_uuid(),
  ticket_id uuid references tickets(id) on delete cascade,
  user_id uuid references user_profiles(id) on delete set null,
  message text not null,
  is_staff boolean default false,
  attachments jsonb default '[]'::jsonb,
  created_at timestamptz default now()
);

alter table ticket_messages enable row level security;

create policy "Users can view own ticket messages"
  on ticket_messages for select
  using (
    exists (select 1 from tickets where tickets.id = ticket_messages.ticket_id and tickets.user_id = auth.uid())
  );

create policy "Admins can view all ticket messages"
  on ticket_messages for select
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Users can create ticket messages"
  on ticket_messages for insert
  with check (auth.uid() = user_id);

create policy "Admins can create ticket messages"
  on ticket_messages for insert
  with check (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create index idx_ticket_messages_ticket_id on ticket_messages(ticket_id);


-- ===================================================================
-- 14. TICKET DEPARTMENTS
-- ===================================================================

create table if not exists departments (
  id bigserial primary key,
  name text not null,
  description text,
  hidden boolean default false,
  created_at timestamptz default now()
);

alter table departments enable row level security;

create policy "Anyone can view active departments"
  on departments for select
  using (hidden = false);

create policy "Admins can manage departments"
  on departments for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 15. EMAIL VERIFICATION
-- ===================================================================

create table if not exists email_verification (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete cascade unique,
  email text not null,
  code text,
  expires_at timestamptz,
  verified boolean default false,
  created_at timestamptz default now()
);

alter table email_verification enable row level security;

create policy "Users can view own verification"
  on email_verification for select
  using (auth.uid() = user_id);

create policy "Users can update own verification"
  on email_verification for update
  using (auth.uid() = user_id);

create policy "Admins manage email verification"
  on email_verification for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 16. REFERRAL CODES
-- ===================================================================

create table if not exists referral_codes (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete cascade,
  code text unique not null,
  usage_count integer default 0,
  reward_coins integer default 0,
  created_at timestamptz default now()
);

alter table referral_codes enable row level security;

create policy "Users can view own referral codes"
  on referral_codes for select
  using (auth.uid() = user_id);

create policy "Admins can view all referral codes"
  on referral_codes for select
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Admins manage referral codes"
  on referral_codes for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 17. REFERRAL USES (tracking)
-- ===================================================================

create table if not exists referral_uses (
  id uuid primary key default gen_random_uuid(),
  referral_code uuid references referral_codes(id) on delete cascade,
  referred_user_id uuid references user_profiles(id) on delete cascade,
  reward_given integer default 0,
  created_at timestamptz default now()
);

alter table referral_uses enable row level security;

create policy "Users can view own referral uses"
  on referral_uses for select
  using (
    exists (
      select 1 from referral_codes
      where referral_codes.id = referral_uses.referral_code
        and referral_codes.user_id = auth.uid()
    )
  );

create policy "Admins manage referral uses"
  on referral_uses for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 18. USER ACTIVITIES (audit log)
-- ===================================================================

create table if not exists user_activities (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete cascade,
  action text not null,
  details jsonb default '{}'::jsonb,
  ip_address text,
  user_agent text,
  created_at timestamptz default now()
);

alter table user_activities enable row level security;

create policy "Users can view own activities"
  on user_activities for select
  using (auth.uid() = user_id);

create policy "Admins can view all activities"
  on user_activities for select
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create index idx_user_activities_user_id on user_activities(user_id);
create index idx_user_activities_created on user_activities(created_at desc);


-- ===================================================================
-- 19. ROLES & PERMISSIONS
-- ===================================================================

create table if not exists roles (
  id bigserial primary key,
  name text unique not null,
  description text,
  is_default boolean default false,
  created_at timestamptz default now()
);

create table if not exists roles_permissions (
  id bigserial primary key,
  role_id bigint references roles(id) on delete cascade,
  permission text not null,
  unique(role_id, permission)
);

alter table roles enable row level security;
alter table roles_permissions enable row level security;

create policy "Anyone can view roles"
  on roles for select
  using (true);

create policy "Anyone can view role permissions"
  on roles_permissions for select
  using (true);

create policy "Admins manage roles"
  on roles for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Admins manage role permissions"
  on roles_permissions for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 20. IMAGE DATABASE (category/egg/location images)
-- ===================================================================

create table if not exists image_db (
  id bigserial primary key,
  name text not null,
  image text,
  target_type text not null check (target_type in ('category', 'egg', 'location')),
  target_id bigint,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

alter table image_db enable row level security;

create policy "Anyone can view images"
  on image_db for select
  using (true);

create policy "Admins manage images"
  on image_db for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 21. REDEEM CODES
-- ===================================================================

create table if not exists redeem_codes (
  id bigserial primary key,
  code text unique not null,
  coins integer default 0,
  expires_at timestamptz,
  max_uses integer default 1,
  current_uses integer default 0,
  active boolean default true,
  created_by uuid references user_profiles(id),
  created_at timestamptz default now()
);

create table if not exists redeem_codes_redeems (
  id uuid primary key default gen_random_uuid(),
  code_id bigint references redeem_codes(id) on delete cascade,
  user_id uuid references user_profiles(id) on delete cascade,
  created_at timestamptz default now()
);

alter table redeem_codes enable row level security;

create policy "Admins manage redeem codes"
  on redeem_codes for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Users can view active redeem codes"
  on redeem_codes for select
  using (active = true);

create policy "Users can view own redeem history"
  on redeem_codes_redeems for select
  using (user_id = auth.uid());

create policy "Admins manage redeem history"
  on redeem_codes_redeems for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

-- Server policy for redeem codes redemption
create policy "Users can redeem codes"
  on redeem_codes_redeems for insert
  with check (auth.uid() = user_id);


-- ===================================================================
-- 22. PAYMENT RECORDS (Stripe / PayPal)
-- ===================================================================

create table if not exists stripe_payments (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete set null,
  stripe_payment_intent text,
  amount integer default 0,
  currency text default 'usd',
  status text default 'pending',
  description text,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

create table if not exists paypal_payments (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete set null,
  paypal_order_id text,
  amount integer default 0,
  currency text default 'usd',
  status text default 'pending',
  description text,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

alter table stripe_payments enable row level security;
alter table paypal_payments enable row level security;

create policy "Users can view own payments"
  on stripe_payments for select
  using (auth.uid() = user_id);

create policy "Users can view own PayPal payments"
  on paypal_payments for select
  using (auth.uid() = user_id);

create policy "Admins can view all payments"
  on stripe_payments for select
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Admins can view all PayPal payments"
  on paypal_payments for select
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 23. IP RELATIONSHIP (anti-fraud / multi-account detection)
-- ===================================================================

create table if not exists ip_relationship (
  id bigserial primary key,
  user_id uuid references user_profiles(id) on delete cascade,
  ip_address text not null,
  is_vpn boolean default false,
  is_proxy boolean default false,
  visits integer default 1,
  last_visit timestamptz default now(),
  created_at timestamptz default now()
);

alter table ip_relationship enable row level security;

create policy "Admins can view IP relationships"
  on ip_relationship for select
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create index idx_ip_relationship_ip on ip_relationship(ip_address);
create index idx_ip_relationship_user on ip_relationship(user_id);


-- ===================================================================
-- 24. ANNOUNCEMENT TAGS + ASSETS + REDEEM REDEEMS (RLS)
-- ===================================================================
alter table announcements_tags enable row level security;
alter table announcements_assets enable row level security;

create policy "Anyone can view announcement tags"
  on announcements_tags for select
  using (true);

create policy "Admins manage announcement tags"
  on announcements_tags for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Anyone can view announcement assets"
  on announcements_assets for select
  using (true);

create policy "Admins manage announcement assets"
  on announcements_assets for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 25. LINK REWARDS (Linkvertise, Shareus, etc)
-- ===================================================================

create table if not exists linkvertise_links (
  id bigserial primary key,
  user_id uuid references user_profiles(id) on delete cascade,
  link text not null,
  reward_coins integer default 0,
  clicks integer default 0,
  active boolean default true,
  created_at timestamptz default now()
);

create table if not exists shareus_links (
  id bigserial primary key,
  user_id uuid references user_profiles(id) on delete cascade,
  link text not null,
  reward_coins integer default 0,
  clicks integer default 0,
  active boolean default true,
  created_at timestamptz default now()
);

alter table linkvertise_links enable row level security;
alter table shareus_links enable row level security;

create policy "Anyone can view active linkvertise links"
  on linkvertise_links for select
  using (active = true);

create policy "Anyone can view active shareus links"
  on shareus_links for select
  using (active = true);

create policy "Admins manage linkvertise links"
  on linkvertise_links for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Admins manage shareus links"
  on shareus_links for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 26. TIMED TASKS (cron jobs stored in DB)
-- ===================================================================

create table if not exists timed_tasks (
  id bigserial primary key,
  task_name text unique not null,
  last_run timestamptz,
  next_run timestamptz,
  interval_minutes integer default 60,
  active boolean default true,
  created_at timestamptz default now()
);

alter table timed_tasks enable row level security;

create policy "Anyone can view timed tasks"
  on timed_tasks for select
  using (true);

create policy "Admins manage timed tasks"
  on timed_tasks for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );


-- ===================================================================
-- 27. QUEUE-BASED SERVER TIME SYSTEM
-- ===================================================================

-- Server time credits
create table if not exists server_time_credits (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete cascade,
  server_uuid text not null,
  minutes_remaining integer default 0,
  total_minutes_purchased integer default 0,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Server sessions (tracks active running time)
create table if not exists server_sessions (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references user_profiles(id) on delete cascade,
  server_uuid text not null,
  node_id text not null,
  started_at timestamptz,
  ends_at timestamptz,
  status text default 'queued' check (status in ('queued', 'active', 'cooldown', 'suspended')),
  queue_position integer,
  cooldown_until timestamptz,
  first_start boolean default true,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Node slot tracking
create table if not exists node_slots (
  id uuid primary key default gen_random_uuid(),
  node_id text unique not null,
  max_active integer default 4,
  current_active integer default 0,
  updated_at timestamptz default now()
);

-- Time packages (buyable time)
create table if not exists time_packages (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  minutes integer not null,
  coin_cost integer not null,
  active boolean default true,
  created_at timestamptz default now()
);

-- First start tracking (per-server)
create table if not exists server_first_starts (
  server_uuid text primary key,
  used boolean default false,
  created_at timestamptz default now()
);

-- Enable RLS on time system tables
alter table server_time_credits enable row level security;
alter table server_sessions enable row level security;
alter table node_slots enable row level security;
alter table time_packages enable row level security;
alter table server_first_starts enable row level security;

-- RLS: Users can see their own time data
create policy "Users can view own time credits"
  on server_time_credits for select
  using (auth.uid() = user_id);

create policy "Users can view own sessions"
  on server_sessions for select
  using (auth.uid() = user_id);

create policy "Users can view their own first starts"
  on server_first_starts for select
  using (
    exists (select 1 from servers where servers.uuid = server_first_starts.server_uuid and servers.user_id = auth.uid())
  );

-- RLS: Anyone can read time packages and node slots (needed for UI)
create policy "Anyone can view active time packages"
  on time_packages for select
  using (active = true);

create policy "Anyone can view node slots"
  on node_slots for select
  using (true);

-- RLS: Admins manage time system
create policy "Admins can manage time credits"
  on server_time_credits for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Admins can manage sessions"
  on server_sessions for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Admins can manage node slots"
  on node_slots for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Admins can manage time packages"
  on time_packages for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

create policy "Admins can manage first starts"
  on server_first_starts for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

-- Server-side can insert/update any time data (for cron/background jobs)
-- Admins manage time credits
create policy "Admins manage time credits"
  on server_time_credits for all
  using (
    exists (select 1 from user_profiles where id = auth.uid() and role = 'admin')
  );

-- Server-side can insert/update any time data (for cron/background jobs)
create policy "Server can insert time credits"
  on server_time_credits for insert
  with check (true);

create policy "Server can update time credits"
  on server_time_credits for update
  using (true);

create policy "Server can manage sessions"
  on server_sessions for insert
  with check (true);

create policy "Server can update sessions"
  on server_sessions for update
  using (true);

create policy "Server can delete sessions"
  on server_sessions for delete
  using (true);

create policy "Server can manage first starts"
  on server_first_starts for insert
  with check (true);

create policy "Server can update first starts"
  on server_first_starts for update
  using (true);

create policy "Server can manage node slots"
  on node_slots for insert
  with check (true);

create policy "Server can update node slots"
  on node_slots for update
  using (true);

-- Indexes for time system
create index idx_server_time_credits_user_id on server_time_credits(user_id);
create index idx_server_time_credits_server_uuid on server_time_credits(server_uuid);
create index idx_server_sessions_status on server_sessions(status);
create index idx_server_sessions_user_id on server_sessions(user_id);
create index idx_server_sessions_server_uuid on server_sessions(server_uuid);
create index idx_server_sessions_node_id on server_sessions(node_id);
create index idx_server_sessions_ends_at on server_sessions(ends_at);
create index idx_server_sessions_cooldown_until on server_sessions(cooldown_until);
create index idx_node_slots_node_id on node_slots(node_id);
create index idx_time_packages_order on time_packages(coin_cost, name);


-- ===================================================================
-- HELPER FUNCTIONS
-- ===================================================================

-- Add coins to a user (atomic, via Supabase service_role bypass)
create or replace function add_coins(p_user_id uuid, p_amount integer, p_reason text)
returns void as $$
begin
  update user_profiles
  set coins = coins + p_amount, updated_at = now()
  where id = p_user_id;

  insert into coin_transactions (user_id, amount, reason)
  values (p_user_id, p_amount, p_reason);
end;
$$ language plpgsql security definer;

-- Deduct coins from a user (returns false if not enough)
create or replace function deduct_coins(p_user_id uuid, p_amount integer, p_reason text)
returns boolean as $$
declare
  v_balance integer;
begin
  select coins into v_balance from user_profiles where id = p_user_id for update;

  if v_balance >= p_amount then
    update user_profiles
    set coins = coins - p_amount, updated_at = now()
    where id = p_user_id;

    insert into coin_transactions (user_id, amount, reason)
    values (p_user_id, -p_amount, p_reason);

    return true;
  else
    return false;
  end if;
end;
$$ language plpgsql security definer;


-- ===================================================================
-- DEFAULT SEED DATA
-- ===================================================================

-- Default time packages
insert into time_packages (name, minutes, coin_cost, active) values
  ('Starter', 5, 10, true),
  ('Basic', 15, 25, true),
  ('Standard', 20, 50, true),
  ('Premium', 30, 100, true)
on conflict do nothing;

-- Default resource packages
insert into resource_packages (name, description, memory, disk, cpu, servers, coin_cost, active) values
  ('Memory Boost', '+1024MB RAM', 1024, 0, 0, 0, 50, true),
  ('Disk Expansion', '+2048MB Disk', 0, 2048, 0, 0, 75, true),
  ('CPU Boost', '+50% CPU', 0, 0, 50, 0, 40, true),
  ('Server Slot', '+1 Additional Server', 0, 0, 0, 1, 150, true)
on conflict do nothing;

-- Default roles
insert into roles (name, description, is_default) values
  ('Client', 'Standard user', true),
  ('Moderator', 'Content moderator', false),
  ('Admin', 'Full system access', false)
on conflict (name) do nothing;

-- After your first user signs up, make them admin with:
-- update user_profiles set role = 'admin' where id = 'YOUR_UUID';
