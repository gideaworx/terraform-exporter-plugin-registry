import { createContext } from "react";
import { PluginRegistry } from "./Models";

export const emptyRegistry = new PluginRegistry("", "", []);

export const RegistryContext = createContext<PluginRegistry>(emptyRegistry);
RegistryContext.displayName = "RegistryContext";

