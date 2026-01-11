#!/usr/bin/env node

import fs from "node:fs/promises";
import path from "node:path";

function parseArgs(argv) {
  const args = {
    root: process.cwd(),
    quiet: false,
  };

  for (let i = 0; i < argv.length; i += 1) {
    const v = argv[i];
    if (v === "--root") {
      const next = argv[i + 1];
      if (!next) throw new Error("Missing value for --root");
      args.root = next;
      i += 1;
      continue;
    }
    if (v === "--quiet") {
      args.quiet = true;
    }
    if (v === "--help" || v === "-h") {
      args.help = true;
    }
  }

  return args;
}

async function fileExists(p) {
  try {
    await fs.access(p);
    return true;
  } catch {
    return false;
  }
}

async function listJsonBasenames(dir) {
  const entries = await fs.readdir(dir, { withFileTypes: true });
  return entries
    .filter((e) => e.isFile() && e.name.endsWith(".json"))
    .map((e) => e.name)
    .sort();
}

function collectKeyPaths(value, prefix, out) {
  if (value === null || value === undefined) {
    out.add(prefix);
    return;
  }

  if (Array.isArray(value)) {
    // Treat arrays as a leaf (structure parity is what matters here).
    out.add(prefix);
    return;
  }

  if (typeof value !== "object") {
    out.add(prefix);
    return;
  }

  const keys = Object.keys(value);
  if (keys.length === 0) {
    out.add(prefix);
    return;
  }

  for (const key of keys) {
    const nextPrefix = prefix ? `${prefix}.${key}` : key;
    collectKeyPaths(value[key], nextPrefix, out);
  }
}

function diffSets(a, b) {
  const missingInB = [];
  const missingInA = [];
  for (const k of a) {
    if (!b.has(k)) missingInB.push(k);
  }
  for (const k of b) {
    if (!a.has(k)) missingInA.push(k);
  }
  missingInA.sort();
  missingInB.sort();
  return { missingInA, missingInB };
}

async function readJson(filePath) {
  const raw = await fs.readFile(filePath, "utf8");
  try {
    return JSON.parse(raw);
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err);
    throw new Error(`Invalid JSON: ${filePath}: ${msg}`);
  }
}

async function extractInitImportedNamespaces(initTsPath) {
  const src = await fs.readFile(initTsPath, "utf8");

  const ar = new Set();
  const en = new Set();

  // Example: import arCommon from './ar/common.json'
  const importRe = /import\s+\w+\s+from\s+'\.\/(ar|en)\/([^']+)\.json'/g;
  let m;
  while ((m = importRe.exec(src)) !== null) {
    const locale = m[1];
    const ns = m[2];
    if (locale === "ar") ar.add(`${ns}.json`);
    if (locale === "en") en.add(`${ns}.json`);
  }

  return { ar, en };
}

async function main() {
  const args = parseArgs(process.argv.slice(2));

  if (args.help) {
    process.stdout.write(
      [
        "check-portal-i18n-parity",
        "",
        "Usage:",
        "  node .github/skills/i18n-parity-portal-web/scripts/check-portal-i18n-parity.mjs",
        "  node .github/skills/i18n-parity-portal-web/scripts/check-portal-i18n-parity.mjs --root /path/to/kyora",
        "",
        "Options:",
        "  --root  Repository root (defaults to cwd)",
        "  --quiet Only print errors",
        "",
      ].join("\n")
    );
    return;
  }

  const repoRoot = path.resolve(args.root);
  const portalI18nDir = path.join(repoRoot, "portal-web", "src", "i18n");
  const arDir = path.join(portalI18nDir, "ar");
  const enDir = path.join(portalI18nDir, "en");
  const initTs = path.join(portalI18nDir, "init.ts");

  if (!(await fileExists(arDir)) || !(await fileExists(enDir))) {
    throw new Error(
      `Could not find portal-web i18n dirs under: ${portalI18nDir}`
    );
  }
  if (!(await fileExists(initTs))) {
    throw new Error(`Could not find: ${initTs}`);
  }

  const arFiles = await listJsonBasenames(arDir);
  const enFiles = await listJsonBasenames(enDir);

  const arSet = new Set(arFiles);
  const enSet = new Set(enFiles);

  let hasErrors = false;

  const nsDiff = diffSets(arSet, enSet);
  if (nsDiff.missingInB.length > 0 || nsDiff.missingInA.length > 0) {
    hasErrors = true;
    if (!args.quiet) {
      console.log("Namespace mismatch (JSON filenames):");
    }
    if (nsDiff.missingInB.length > 0) {
      console.error(`- Missing in en/: ${nsDiff.missingInB.join(", ")}`);
    }
    if (nsDiff.missingInA.length > 0) {
      console.error(`- Missing in ar/: ${nsDiff.missingInA.join(", ")}`);
    }
  }

  // Check init.ts imports match disk
  const initImports = await extractInitImportedNamespaces(initTs);
  const initArDiff = diffSets(arSet, initImports.ar);
  const initEnDiff = diffSets(enSet, initImports.en);

  if (initArDiff.missingInB.length > 0 || initEnDiff.missingInB.length > 0) {
    hasErrors = true;
    if (!args.quiet) {
      console.log("init.ts imports do not cover all namespaces on disk:");
    }
    if (initArDiff.missingInB.length > 0) {
      console.error(
        `- Missing ar imports in init.ts: ${initArDiff.missingInB.join(", ")}`
      );
    }
    if (initEnDiff.missingInB.length > 0) {
      console.error(
        `- Missing en imports in init.ts: ${initEnDiff.missingInB.join(", ")}`
      );
    }
  }

  const commonFiles = arFiles.filter((f) => enSet.has(f));
  for (const file of commonFiles) {
    const arJson = await readJson(path.join(arDir, file));
    const enJson = await readJson(path.join(enDir, file));

    const arKeys = new Set();
    const enKeys = new Set();

    collectKeyPaths(arJson, "", arKeys);
    collectKeyPaths(enJson, "", enKeys);

    const { missingInA, missingInB } = diffSets(arKeys, enKeys);
    if (missingInA.length === 0 && missingInB.length === 0) continue;

    hasErrors = true;
    if (!args.quiet) {
      console.log(`Key mismatch in ${file}:`);
    }
    if (missingInB.length > 0) {
      console.error(`- Missing in en: ${missingInB.join(", ")}`);
    }
    if (missingInA.length > 0) {
      console.error(`- Missing in ar: ${missingInA.join(", ")}`);
    }
  }

  if (hasErrors) {
    process.exitCode = 1;
    return;
  }

  if (!args.quiet) {
    console.log("OK: portal-web i18n parity check passed.");
  }
}

main().catch((err) => {
  const msg = err instanceof Error ? err.message : String(err);
  console.error(msg);
  process.exitCode = 1;
});
