/**
 * Format a byte count into a human-readable string.
 *   512        -> "512 B"
 *   524288     -> "512 KB"
 *   536870912  -> "512 MB"
 *   1073741824 -> "1 GB"
 */
export function formatBytes(bytes: number, decimals: number = 2): string {
  if (bytes === 0) return "0 B";

  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const clampedIndex = Math.min(i, sizes.length - 1);
  const value = bytes / Math.pow(k, clampedIndex);
  return `${value.toFixed(decimals)} ${sizes[clampedIndex]}`;
}

/**
 * Format a CPU usage value (0-100+) as a percentage string.
 *   45.2  -> "45.2%"
 *   100   -> "100.0%"
 */
export function formatCPU(cpu: number): string {
  return `${cpu.toFixed(1)}%`;
}

/**
 * Format an ISO-8601 date string into a readable date-time.
 *   "2026-04-03T16:32:00Z" -> "Apr 3, 2026 at 4:32 PM"
 */
export function formatDate(dateStr: string): string {
  const date = new Date(dateStr);
  const months = [
    "Jan", "Feb", "Mar", "Apr", "May", "Jun",
    "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
  ];
  const month = months[date.getMonth()];
  const day = date.getDate();
  const year = date.getFullYear();

  let hours = date.getHours();
  const minutes = date.getMinutes();
  const ampm = hours >= 12 ? "PM" : "AM";
  hours = hours % 12 || 12;
  const minutesStr = minutes.toString().padStart(2, "0");

  return `${month} ${day}, ${year} at ${hours}:${minutesStr} ${ampm}`;
}

/**
 * Format a duration in seconds into a compact human-readable string.
 *   180600  -> "2d 1h 10m"
 *   3661    -> "1h 1m 1s"
 *   45      -> "45s"
 */
export function formatUptime(seconds: number): string {
  if (seconds < 0) return "0s";

  const d = Math.floor(seconds / 86400);
  const h = Math.floor((seconds % 86400) / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);

  const parts: string[] = [];
  if (d > 0) parts.push(`${d}d`);
  if (h > 0) parts.push(`${h}h`);
  if (m > 0) parts.push(`${m}m`);
  if (s > 0 || parts.length === 0) parts.push(`${s}s`);

  return parts.join(" ");
}
