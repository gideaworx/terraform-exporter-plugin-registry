import {
  createTheme, CssBaseline, Theme,
  ThemeProvider
} from "@mui/material";
import { useEffect, useState } from "react";
import { emptyRegistry, RegistryContext } from "./Context";
import { PluginRegistry } from "./Models";
import { PageHeader } from "./PageComponents";
import { parseRegistry } from "./parser";
import { PluginContainer } from "./PluginComponents";

const lightTheme = createTheme({
  palette: {
    mode: "light",
  },
});

const darkTheme = createTheme({
  palette: {
    mode: "dark",
  },
});

export function Registry() {
  const [registry, setRegistry] = useState<PluginRegistry>(emptyRegistry);
  const [theme, setTheme] = useState<Theme>(darkTheme);

  useEffect(() => {
    parseRegistry().then(setRegistry);
  }, []);

  const registryName = registry.name;
  const headerProps = { registryName, theme, lightTheme, darkTheme, setTheme };

  return (
    <RegistryContext.Provider value={registry}>
      <ThemeProvider theme={theme}>
        <CssBaseline enableColorScheme />
        <PageHeader {...headerProps} />
        <PluginContainer plugins={registry.plugins} />
      </ThemeProvider>
    </RegistryContext.Provider>
  );
}
