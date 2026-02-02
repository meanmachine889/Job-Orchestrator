"use client"

import { useEffect, useState, useCallback } from "react"
import { useRouter } from "next/navigation"
import { JobListItem } from "@/types/job"
import { Badge } from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { CreateJobDialog } from "@/components/create-job-dialog"

function getStatusVariant(status: string): "default" | "secondary" | "destructive" | "outline" {
  switch (status) {
    case "SUCCESS":
      return "default"
    case "FAILED":
    case "DEAD":
      return "destructive"
    case "RUNNING":
      return "secondary"
    default:
      return "outline"
  }
}

export default function JobsPage() {
  const [jobs, setJobs] = useState<JobListItem[]>([])
  const router = useRouter()

  const fetchJobs = useCallback(async () => {
    try {
      const res = await fetch("http://localhost:8080/jobs?limit=50")
      const data = await res.json()
      setJobs(data.jobs ?? [])
    } catch (err) {
      console.error("Failed to fetch jobs:", err)
    }
  }, [])

  useEffect(() => {
    fetchJobs()
    const interval = setInterval(fetchJobs, 2000)
    return () => clearInterval(interval)
  }, [fetchJobs])

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-lg font-normal text-muted-foreground">Jobs</h1>
        <CreateJobDialog onJobCreated={fetchJobs} />
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="font-normal">ID</TableHead>
            <TableHead className="font-normal">Type</TableHead>
            <TableHead className="font-normal">Status</TableHead>
            <TableHead className="font-normal">Retries</TableHead>
            <TableHead className="font-normal">Worker</TableHead>
            <TableHead className="font-normal">Created</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {jobs.map(job => (
            <TableRow
              key={job.id}
              className="cursor-pointer hover:bg-muted/50"
              onClick={() => router.push(`/jobs/${job.id}`)}
            >
              <TableCell className="font-mono text-xs truncate max-w-[120px]">
                {job.id}
              </TableCell>
              <TableCell className="text-sm">{job.type}</TableCell>
              <TableCell>
                <Badge variant={getStatusVariant(job.status)} className="font-normal">
                  {job.status}
                </Badge>
              </TableCell>
              <TableCell className="text-sm">{job.retry_count}</TableCell>
              <TableCell className="text-sm text-muted-foreground">
                {job.worker_id ?? "—"}
              </TableCell>
              <TableCell className="text-sm text-muted-foreground">
                {new Date(job.created_at).toLocaleString()}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
