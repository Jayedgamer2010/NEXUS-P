// ─── Server status colors (Material / MUI compatible palette names) ────────────

export const STATUS_COLORS: Record<string, string> = {
  running: "green",
  offline: "red",
  installing: "yellow",
  install_failed: "red",
  suspended: "gray",
  starting: "blue",
  stopping: "orange",
};

// ─── Available power actions for a server ──────────────────────────────────────

export const POWER_ACTIONS = ["start", "stop", "restart", "kill"] as const;

// ─── Role display colors ──────────────────────────────────────────────────────

export const ROLE_COLORS: Record<string, string> = {
  admin: "purple",
  client: "blue",
};
