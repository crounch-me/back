
export class HealthStatus {
  public constructor(
    public status: string
  ) {

  }
}

export function getOK(): Promise<HealthStatus> {
  return Promise.resolve(new HealthStatus("ok"))
}

export function getPrivateOK(): Promise<HealthStatus> {
  return Promise.resolve(new HealthStatus("ok private"))
}