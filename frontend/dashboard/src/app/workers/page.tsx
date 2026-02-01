"use client"

import { useEffect, useState } from "react"
import { Worker } from "@/types/worker"

export default function WorkersPage() {
  const [workers, setWorkers] = useState<Worker[]>([])

  const fetchWorkers = async () => {
    const res = await fetch("http://localhost:8080/workers")
    const data = await res.json()
    setWorkers(data.workers)
  }

  useEffect(() => {
    fetchWorkers()
    const interval = setInterval(fetchWorkers, 5000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Workers</h1>

      <table className="w-full border">
        <thead>
          <tr className="border-b">
            <th>ID</th>
            <th>Hostname</th>
            <th>Status</th>
            <th>Last Heartbeat</th>
          </tr>
        </thead>
        <tbody>
          {workers.map(w => (
            <tr key={w.id} className="border-b text-sm">
              <td className="truncate max-w-[140px]">{w.id}</td>
              <td>{w.hostname}</td>
              <td>
                <span
                  className={`px-2 py-1 rounded text-xs ${
                    w.status === "ONLINE"
                      ? "bg-green-100 text-green-800"
                      : "bg-red-100 text-red-800"
                  }`}
                >
                  {w.status}
                </span>
              </td>
              <td>{new Date(w.last_heartbeat).toLocaleString()}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
