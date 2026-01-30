export function StatusBadge({ status }: { status: string }) {
  const color =
    status === "SUCCESS" ? "green" :
    status === "FAILED" || status === "DEAD" ? "red" :
    status === "RUNNING" ? "blue" :
    "yellow"

  return (
    <span className={`px-2 py-1 text-xs rounded bg-${color}-100 text-${color}-800`}>
      {status}
    </span>
  )
}
