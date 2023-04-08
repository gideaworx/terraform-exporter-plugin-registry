import {
  Box, FormControl,
  Grid,
  InputLabel,
  Link,
  MenuItem,
  Paper,
  Select,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  Typography
} from "@mui/material";
import cloneDeep from "lodash.clonedeep";
import { Dispatch, SetStateAction, useEffect, useState } from "react";
import { compareLoose } from "semver";
import { RegistryContext } from "./Context";
import { Plugin, PluginAuthor, PluginVersion } from "./Models";

const pluginSortVals = {
  method: {
    name: "Name",
    age: "Last Updated",
  },
  direction: {
    asc: "Ascending (A-Z, oldest to newest)",
    desc: "Descending (Z-A, newest to oldest)",
  },
};

interface pluginHeaderProps {
  setSort: Dispatch<SetStateAction<"name" | "age">>;
  setDirection: Dispatch<SetStateAction<"asc" | "desc">>;
}

function PluginHeader(props: pluginHeaderProps) {
  const { setSort, setDirection } = props;

  return (
    <Grid
      container
      direction="row"
      justifyContent="flex-end"
      alignItems="center"
      sx={{
        pt: "12px",
      }}
    >
      <Grid item xs={1}>
        Sort:
      </Grid>
      <Grid item xs={2}>
        <FormControl fullWidth>
          <InputLabel id="sort-method">Sort By</InputLabel>
          <Select
            labelId="sort-method"
            label="Sort By"
            onChange={(event) => {
              setSort(event.target.value as "name" | "age");
            }}
          >
            {Object.entries(pluginSortVals.method).map(([value, name]) => (
              <MenuItem value={value}>{name}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>
      <Grid item xs={2}>
        <FormControl fullWidth>
          <InputLabel id="sort-direction">Sort Direction</InputLabel>
          <Select
            labelId="sort-direction"
            label="Sort Direction"
            onChange={(event) => {
              setDirection(event.target.value as "asc" | "desc");
            }}
          >
            {Object.entries(pluginSortVals.direction).map(([value, name]) => (
              <MenuItem value={value}>{name}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Grid>
    </Grid>
  );
}

export function PluginContainer(props: { plugins: Plugin[] }) {
  const [sort, setSort] = useState<"name" | "age">("age");
  const [direction, setDirection] = useState<"asc" | "desc">("desc");

  const { plugins } = props;

  const sortByName = (a: Plugin, b: Plugin) => {
    return a.name.localeCompare(b.name, undefined, { sensitivity: "base" });
  };

  const sortByAge = (a: Plugin, b: Plugin) => {
    return a.lastUpdated.getTime() - b.lastUpdated.getTime();
  };

  useEffect(() => {
    plugins.sort(sort === "name" ? sortByName : sortByAge);
    if (direction === "desc") {
      plugins.reverse();
    }
  }, [sort, direction, plugins]);

  return (
    <Stack
      direction="column"
      justifyContent="flex-start"
      alignItems="flex-start"
      spacing={2}
    >
      <PluginHeader setDirection={setDirection} setSort={setSort} />
      {plugins.map((plugin) => (
        <PluginPane plugin={plugin} />
      ))}
    </Stack>
  );
}

function PluginAuthorDisplay(props: { author: PluginAuthor }) {
  const { author } = props;

  return author.email.trim() !== "" ? (
    <Link target="_blank" rel="noreferrer" href={`mailto:${author.email}`}>
      {author.name}
    </Link>
  ) : (
    <span>{author.name}</span>
  );
}

export function PluginPane(props: { plugin: Plugin }) {
  const { plugin } = props;

  const versions = cloneDeep(plugin.versions);
  versions.sort((v1: PluginVersion, v2: PluginVersion) => {
    return compareLoose(v2.version, v1.version);
  });

  if (versions.length === 0) {
    return <></>;
  }

  const latestVersion = versions[0];

  return (
    <RegistryContext.Consumer>
      {(registry) => (
        <div style={{padding: "4px 24px", width: "100%", boxSizing: "border-box"}}>
        <Paper sx={{
          width: "100%",
          p: "4px 24px",
        }}>
          <Stack
            direction="row"
            justifyContent="space-between"
            alignItems="flex-start"
            spacing={1}
          >
            <Typography
              variant="h5"
              style={{ display: "inline" }}
              component="div"
            >
              {plugin.name}{" — "}
              <Typography variant="h6" style={{ display: "inline" }}>
                v{latestVersion.version}
              </Typography>
            </Typography>

            <Typography variant="h6" style={{ display: "inline" }}>
              {plugin.lastUpdated.toISOString().substring(0, 10)}
            </Typography>
          </Stack>
          <Stack
            direction="column"
            justifyContent="flex-start"
            alignItems="stretch"
            spacing={2}
          >
            <Box component="h4">{plugin.description}</Box>
            <Box sx={{ width: "100%" }}>
              <TableContainer>
                <Table>
                  <TableBody>
                    <TableRow>
                      <TableCell align="right">Install</TableCell>
                      <TableCell align="left">
                        <code>
                          terraform-exporter install -r "{registry.name}"{" "}
                          {plugin.name}
                        </code>
                      </TableCell>
                    </TableRow>
                    {plugin.authors.some((a) => a.company) && (
                      <TableRow>
                        <TableCell align="right">Company</TableCell>
                        <TableCell align="left">
                          {plugin.authors
                            .filter((a) => a.company)
                            .map((a) => a.company)
                            .join(", ")}
                        </TableCell>
                      </TableRow>
                    )}
                    {plugin.homepage !== "" && (
                      <TableRow>
                        <TableCell align="right">Homepage</TableCell>
                        <TableCell align="left">
                          <Link
                            target="_blank"
                            rel="noreferrer"
                            href={plugin.homepage}
                          >
                            {plugin.homepage}
                          </Link>
                        </TableCell>
                      </TableRow>
                    )}
                    {plugin.authors.length > 0 && (
                      <TableRow>
                        <TableCell align="right">Authors</TableCell>
                        <TableCell align="left">
                          {plugin.authors
                            .map((author) => (
                              <PluginAuthorDisplay author={author} />
                            ))
                            .reduce((prev, next) => (
                              <>
                                {prev}
                                <span>,&nbsp;</span>
                                {next}
                              </>
                            ))}
                        </TableCell>
                      </TableRow>
                    )}
                    <TableRow>
                      <TableCell align="right">Platforms</TableCell>
                      <TableCell align="left">
                        {Object.keys(latestVersion.download).join(", ")}
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </TableContainer>
            </Box>
          </Stack>
        </Paper>
        </div>
      )}
    </RegistryContext.Consumer>
  );
}
