export class HealthStatus {
  public constructor(
    public status: string,
    public version: string
  ) {

  }
}

export function getOK(): Promise<HealthStatus> {
  return Promise.resolve(new HealthStatus("ok", process.env.npm_package_version || ""))
}

export function getPrivateOK(): Promise<HealthStatus> {
  return Promise.resolve(new HealthStatus("ok private", process.env.npm_package_version || ""))
}