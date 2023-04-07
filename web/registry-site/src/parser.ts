import * as yaml from "js-yaml";
import { Plugin, PluginRegistry } from "./Models";

export async function parseRegistry(): Promise<PluginRegistry> {
  const response = await fetch("./index.yaml");
  const text = await response.text();

  const rawRegistry = yaml.load(text, { schema: yaml.DEFAULT_SCHEMA }) as any;
  const plugins = rawRegistry.plugins.map((plugin: any) => {
    return new Plugin(plugin.name, new Date(plugin.updated), plugin.description, plugin.authors, plugin.homepage, plugin.versions);
  })

  return new PluginRegistry(rawRegistry.name, rawRegistry.baseURL, plugins);
}
