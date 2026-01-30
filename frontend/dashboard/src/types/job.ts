export type JobListItem = {
  id: string
  type: string
  status: string
  retry_count: number
  worker_id: string | null
  created_at: string
}

export type JobDetail = {
  id: string
  type: string
  status: string
  payload: any
  retry_count: number
  max_retries: number
  timeout_seconds: number
  worker_id: string | null
  error?: string
  created_at: string
  updated_at: string
}
