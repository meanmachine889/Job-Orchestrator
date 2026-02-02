"use client"

import { useEffect, useState, useCallback } from "react"
import { Worker } from "@/types/worker"
import { Badge } from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

export default function WorkersPage() {
  const [workers, setWorkers] = useState<Worker[]>([])

  const fetchWorkers = useCallback(async () => {
    try {
      const res = await fetch("http://localhost:8080/workers")
      const data = await res.json()
      // Filter to show only online workers
      const onlineWorkers = (data.workers ?? []).filter(
        (w: Worker) => w.status === "ONLINE"
      )
      setWorkers(onlineWorkers)
    } catch (err) {
      console.error("Failed to fetch workers:", err)
    }
  }, [])

  useEffect(() => {
    fetchWorkers()
    const interval = setInterval(fetchWorkers, 3000)
    return () => clearInterval(interval)
  }, [fetchWorkers])

  const onlineCount = workers.length

  return (
    <div className="p-6">
      <div className="flex items-center gap-3 mb-6">
        <h1 className="text-lg font-normal text-muted-foreground">Workers</h1>
        <span className="text-xs text-muted-foreground">
          {onlineCount} online
        </span>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="font-normal">ID</TableHead>
            <TableHead className="font-normal">Hostname</TableHead>
            <TableHead className="font-normal">Status</TableHead>
            <TableHead className="font-normal">Last Heartbeat</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {workers.length === 0 ? (
            <TableRow>
              <TableCell colSpan={4} className="text-center text-sm text-muted-foreground py-8">
                No online workers
              </TableCell>
            </TableRow>
          ) : (
            workers.map(w => (
              <TableRow key={w.id}>
                <TableCell className="font-mono text-xs truncate max-w-[140px]">
                  {w.id}
                </TableCell>
                <TableCell className="text-sm">{w.hostname}</TableCell>
                <TableCell>
                  <Badge variant="default" className="font-normal">
                    {w.status}
                  </Badge>
                </TableCell>
                <TableCell className="text-sm text-muted-foreground">
                  {new Date(w.last_heartbeat).toLocaleString()}
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  )
}
