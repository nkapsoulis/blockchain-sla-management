interface SLA {
  id: string,
  name: string,
  state: string,
  assessment: Assessment,
  details: Detail,
}

interface Detail {
  id: string,
  type: string,
  name: string,
  provider: Entity,
  client: Entity,
  creation: string,
  guarantees: Guarantee[],
  service: string,
}

interface Assessment {
  firstExecution: string,
  lastExecution: string,
}

interface Entity {
  id: string,
  name: string,
}

interface Importance {
  name: string,
  constraint: string,
}

interface Guarantee {
  name: string,
  constraint: string,
  importance: Importance[],
}

interface Violation {
  id: string,
  slaid: string,
  guaranteeID: string,
  datetime: string,
  constraint: string,
  values: Value[],
  importanceName: string,
  importance: number,
  appID: string,
}

interface Value {
  key: string,
  value: number,
  datetime: string,
}

interface Approval {
  providerApproved: boolean,
  consumerApproved: boolean,
}

interface User {
  name: string,
  pubKey: string,
  balance: string,
  providerOf: string,
  clientOf: string,
}

export type {
  SLA,
  Violation,
  Approval,
  User,
};
