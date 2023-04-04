import { SLA, Approval } from './Asset.model';

interface APIResponse {
  asset?: SLA,
  assets?: SLA[],
  approvals?: Approval,
  success: boolean,
  message?: string,
}

export default APIResponse;
