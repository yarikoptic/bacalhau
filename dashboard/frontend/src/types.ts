export interface ResourceUsageConfig {
  CPU?: string;
  Memory?: string;
  Disk?: string;
  GPU: string;
}

export interface ResourceUsageData {
  CPU?: number;
  Memory?: number;
  Disk?: number;
  GPU?: number;
}

export interface ResourceUsageProfile {
  Job?: ResourceUsageData;
  SystemUsing?: ResourceUsageData;
  SystemTotal?: ResourceUsageData;
}

export interface RunCommandResult {
  stdout: string;
  stdouttruncated: boolean;
  stderr: string;
  stderrtruncated: boolean;
  exitCode: number;
  runnerError: string;
}

export interface StorageSpec {
  StorageSource: string;
  Name?: string;
  CID?: string;
  URL?: string;
  path?: string;
  Metadata?: { [key: string]: string};
}

export interface PublishedResult {
  NodeID?: string;
  ShardIndex?: number;
  Data?: StorageSpec;
}

export interface Job {
  APIVersion: string;
  ID: string;
  RequesterNodeID?: string;
  RequesterPublicKey?: string;
  ClientID?: string;
  Spec: Spec;
  Deal: Deal;
  ExecutionPlan?: JobExecutionPlan;
  CreatedAt: string;
  JobState?: JobState;
  JobEvents?: JobEvent[];
  LocalJobEvents?: JobLocalEvent[];
}

export interface JobWithInfo {
  Job?: Job;
  JobState?: JobState;
  JobEvents?: JobEvent[];
  JobLocalEvents?: JobLocalEvent[];
}

export interface JobShard {
  Job?: Job;
  Index?: number;
}

export interface JobExecutionPlan {
  ShardsTotal?: number;
}

export interface JobShardingConfig {
  GlobPattern?: string;
  BatchSize?: number;
  GlobPatternBasePath?: string;
}

export interface JobState {
  Nodes?: { [key: string]: JobNodeState};
}

export interface JobNodeState {
  Shards?: { [key: number]: JobShardState};
}

export interface JobShardState {
  NodeId: string;
  ShardIndex: number;
  State: string;
  Status: string;
  VerificationProposal?: string;
  VerificationResult?: VerificationResult;
  PublishedResults?: StorageSpec;
  RunOutput?: RunCommandResult;
}

export interface Deal {
  Concurrency?: number;
  Confidence?: number;
  MinBids?: number;
}

export interface Spec {
  Engine?: string;
  Verifier?: string;
  Publisher?: string;
  Docker?: JobSpecDocker;
  Language?: JobSpecLanguage;
  Wasm?: JobSpecWasm;
  Resources?: ResourceUsageConfig;
  inputs?: StorageSpec[];
  Contexts?: StorageSpec[];
  outputs?: StorageSpec[];
  Annotations?: string[];
  Sharding?: JobShardingConfig;
  DoNotTrack?: boolean;
}

export interface JobSpecDocker {
  Image?: string;
  Entrypoint?: string[];
  EnvironmentVariables?: string[];
  WorkingDirectory?: string;
}

export interface JobSpecLanguage {
  Language?: string;
  LanguageVersion?: string;
  DeterministicExecution?: boolean;
  JobContext?: StorageSpec;
  Command?: string;
  ProgramPath?: string;
  RequirementsPath?: string;
}

export interface JobSpecWasm {
  EntryPoint?: string;
  Parameters?: string[];
}

export interface JobLocalEvent {
  EventName?: string;
  JobID?: string;
  ShardIndex?: number;
  TargetNodeID?: string;
}

export interface JobEvent {
  APIVersion?: string;
  JobID?: string;
  ShardIndex?: number;
  ClientID?: string;
  SourceNodeID?: string;
  TargetNodeID?: string;
  EventName?: string;
  Spec?: Spec;
  JobExecutionPlan?: JobExecutionPlan;
  Deal?: Deal;
  Status?: string;
  VerificationProposal?: string;
  VerificationResult?: VerificationResult;
  PublishedResult?: StorageSpec;
  EventTime?: string;
  SenderPublicKey?: string;
  RunOutput?: RunCommandResult;
}

export interface VerificationResult {
  Complete?: boolean;
  Result?: boolean;
}

export interface JobCreatePayload {
  ClientID?: string;
  Job?: Job;
  Context?: string;
}