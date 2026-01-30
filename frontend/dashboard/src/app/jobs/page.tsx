"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { JobListItem } from "@/types/job"
import { StatusBadge } from "@/components/StatusBadge"

export default function JobsPage() {
    const [jobs, setJobs] = useState<JobListItem[]>([])

    useEffect(() => {
        fetch("http://localhost:8080/jobs?limit=50")
            .then(res => res.json())
            .then(data => setJobs(data.jobs))
    }, [])

    return (
        <div className="p-6">
            <h1 className="text-2xl font-semibold mb-4">Jobs</h1>

            <table className="w-full border">
                <thead>
                    <tr className="border-b">
                        <th>ID</th>
                        <th>Type</th>
                        <th>Status</th>
                        <th>Retries</th>
                        <th>Worker</th>
                        <th>Created</th>
                    </tr>
                </thead>
                <tbody>
                    {jobs.map(job => (
                        <tr key={job.id} className="border-b text-sm">
                            <td className="truncate max-w-[120px]">
                                <Link href={`/jobs/${job.id}`} className="text-blue-600">
                                    {job.id}
                                </Link>
                            </td>
                            <td>{job.type}</td>
                            <td>
                                <StatusBadge status={job.status} />
                            </td>
                            <td>{job.retry_count}</td>
                            <td>{job.worker_id ?? "-"}</td>
                            <td>{new Date(job.created_at).toLocaleString()}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}
