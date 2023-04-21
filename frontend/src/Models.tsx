export enum TargetArch {
  OSXAMD64 = "darwin/amd64",
  OSXARM64 = "darwin/arm64",
  LinuxAMD64 = "linux/amd64",
  LinuxARM64 = "linux/arm64",
  WindowsAMD64 = "windows/amd64",
  WindowsARM64 = "windows/arm64",
  MultiArch = "multi-arch",
}

export enum PluginType {
  Native = "native",
  Python = "python",
  NodeJS = "nodejs",
  Java = "java",
}

export interface RegistryInfo {
  name?: string;
  baseURL?: string;
}

export interface ExtraCommandLineArgs {
  args?: string[];
}

export interface NeedsChecksum {
  sha256sum: string;
}

export interface PluginExecutable {
  locator: string; // binary/jar name/url or python/nodejs module name
  type: PluginType;
  info?: ExtraCommandLineArgs | NeedsChecksum;
}

export type DownloadInfo = Partial<Record<TargetArch, PluginExecutable>>;

export interface PluginVersion {
  version: string;
  download: DownloadInfo;
}

export interface PluginAuthor {
  name: string;
  email: string;
  company?: string;
}

export class Plugin {
  constructor(
    public readonly name: string,
    public readonly lastUpdated: Date,
    public readonly description: string,
    public readonly authors: PluginAuthor[],
    public readonly homepage: string,
    public readonly versions: PluginVersion[]
  ) {}
}

export class PluginRegistry implements RegistryInfo {
  constructor(
    public readonly name: string,
    public readonly baseURL: string,
    public readonly plugins: Plugin[]
  ) {}
}
