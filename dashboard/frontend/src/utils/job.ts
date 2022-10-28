import {
  Job,
  JobShardState,
} from '../types'

export const getShortId = (id: string, length = 8) => {
  return id.slice(0, length)
}

export const getJobShardState = (job: Job): JobShardState | undefined => {
  const nodeStates = job.JobState?.Nodes || {}
  const nodeId = Object.keys(nodeStates)[0]
  if(!nodeId) return
  const nodeState = nodeStates[nodeId]
  const shardState = (nodeState?.Shards || {})['0']
  if(!shardState) return
  return shardState
}

export const getShardStateTitle = (shardState: JobShardState | undefined): string => {
  return shardState ?
    shardState.State :
    'Unknown'
}

export const getJobStateTitle = (job: Job) => getShardStateTitle(getJobShardState(job))