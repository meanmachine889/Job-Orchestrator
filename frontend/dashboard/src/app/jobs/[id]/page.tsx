"use client"

import { use, useEffect, useState } from "react"
import { JobDetail } from "@/types/job"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ArrowLeft } from "lucide-react"
import Link from "next/link"

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

export default function JobDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)
  const [job, setJob] = useState<JobDetail | null>(null)

  useEffect(() => {
    fetch(`http://localhost:8080/jobs/${id}`)
      .then(res => res.json())
      .then(setJob)
  }, [id])

  if (!job) {
    return (
      <div className="p-6">
        <p className="text-sm text-muted-foreground">Loading...</p>
      </div>
    )
  }

  return (
    <div className="p-6">
      <Link
        href="/jobs"
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-4"
      >
        <ArrowLeft className="size-4" />
        Back to jobs
      </Link>

      <div className="flex items-center gap-3 mb-6">
        <h1 className="text-lg font-normal text-muted-foreground">Job Details</h1>
        <Badge variant={getStatusVariant(job.status)} className="font-normal">
          {job.status}
        </Badge>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-normal text-muted-foreground">General</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">ID</span>
              <span className="font-mono text-xs">{job.id}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Type</span>
              <span>{job.type}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Worker</span>
              <span className="font-mono text-xs">{job.worker_id ?? "—"}</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-normal text-muted-foreground">Execution</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Retries</span>
              <span>{job.retry_count} / {job.max_retries}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Timeout</span>
              <span>{job.timeout_seconds}s</span>
            </div>
            {job.error && (
              <div className="flex justify-between">
                <span className="text-muted-foreground">Error</span>
                <span className="text-destructive text-xs">{job.error}</span>
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-normal text-muted-foreground">Timestamps</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Created</span>
              <span>{new Date(job.created_at).toLocaleString()}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Updated</span>
              <span>{new Date(job.updated_at).toLocaleString()}</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-normal text-muted-foreground">Payload</CardTitle>
          </CardHeader>
          <CardContent>
            <pre className="text-xs font-mono bg-muted p-3 rounded overflow-auto max-h-40">
              {JSON.stringify(job.payload, null, 2)}
            </pre>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
