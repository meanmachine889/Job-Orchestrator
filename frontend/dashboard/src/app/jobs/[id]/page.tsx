"use client"

import { use, useEffect, useState } from "react"
import { JobDetail } from "@/types/job"

export default function JobDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params)
  const [job, setJob] = useState<JobDetail | null>(null)

  useEffect(() => {
    fetch(`http://localhost:8080/jobs/${id}`)
      .then(res => res.json())
      .then(setJob)
  }, [id])

  if (!job) return <p>Loading...</p>

  return (
    <div className="p-6">
      <h1 className="text-xl font-semibold mb-4">Job {job.id}</h1>

      <pre className="bg-gray-100 p-4 rounded text-sm text-black">
        {JSON.stringify(job, null, 2)}
      </pre>
    </div>
  )
}
