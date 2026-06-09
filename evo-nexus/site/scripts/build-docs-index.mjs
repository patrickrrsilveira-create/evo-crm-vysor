/**
 * Scans public/docs/ and generates public/docs-index.json
 * Used by the Docs page to render the documentation tree.
 */
import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const DOCS_DIR = path.resolve(__dirname, "../public/docs");
const OUTPUT = path.resolve(__dirname, "../public/docs-index.json");

const SECTION_ORDER = [
  "getting-started",
  "guides",
  "dashboard",
  "agents",
  "skills",
  "routines",
  "integrations",
  "real-world",
  "reference",
];

const SECTION_TITLES = {
  "getting-started": "Getting Started",
  guides: "Guides",
  dashboard: "Dashboard",
  agents: "Agents",
  skills: "Skills",
  routines: "Routines",
  integrations: "Integrations",
  "real-world": "Real World",
  reference: "Reference",
};

const TOP_LEVEL_ORDER = ["introduction.md", "getting-started.md", "architecture.md"];

function titleFromMd(filePath) {
  try {
    const text = fs.readFileSync(filePath, "utf-8");
    for (const line of text.split("\n")) {
      const m = line.match(/^#\s+(.+)/);
      if (m) return m[1].trim();
    }
  } catch {}
  return path.basename(filePath, ".md").replace(/[-_]/g, " ").replace(/\b\w/g, (c) => c.toUpperCase());
}

function buildTree() {
  if (!fs.existsSync(DOCS_DIR)) {
    console.warn("No docs/ directory found at", DOCS_DIR);
    fs.writeFileSync(OUTPUT, JSON.stringify({ sections: [] }));
    return;
  }

  const sections = [];

  // Top-level .md files -> "Getting Started"
  const topFiles = fs
    .readdirSync(DOCS_DIR)
    .filter((f) => f.endsWith(".md") && fs.statSync(path.join(DOCS_DIR, f)).isFile())
    .sort((a, b) => {
      const ia = TOP_LEVEL_ORDER.indexOf(a);
      const ib = TOP_LEVEL_ORDER.indexOf(b);
      return (ia === -1 ? 999 : ia) - (ib === -1 ? 999 : ib) || a.localeCompare(b);
    });

  if (topFiles.length) {
    sections.push({
      title: "Getting Started",
      slug: "getting-started",
      children: topFiles.map((f) => ({
        title: titleFromMd(path.join(DOCS_DIR, f)),
        slug: f.replace(".md", ""),
        path: f,
      })),
    });
  }

  // Subdirectories
  const subdirs = fs
    .readdirSync(DOCS_DIR)
    .filter((d) => d !== "imgs" && fs.statSync(path.join(DOCS_DIR, d)).isDirectory())
    .sort((a, b) => {
      const ia = SECTION_ORDER.indexOf(a);
      const ib = SECTION_ORDER.indexOf(b);
      return (ia === -1 ? 999 : ia) - (ib === -1 ? 999 : ib) || a.localeCompare(b);
    });

  for (const dir of subdirs) {
    const dirPath = path.join(DOCS_DIR, dir);
    const mdFiles = fs
      .readdirSync(dirPath)
      .filter((f) => f.endsWith(".md"))
      .sort();

    if (!mdFiles.length) continue;

    sections.push({
      title: SECTION_TITLES[dir] || dir.replace(/[-_]/g, " ").replace(/\b\w/g, (c) => c.toUpperCase()),
      slug: dir,
      children: mdFiles.map((f) => ({
        title: titleFromMd(path.join(dirPath, f)),
        slug: `${dir}/${f.replace(".md", "")}`,
        path: `${dir}/${f}`,
      })),
    });
  }

  fs.writeFileSync(OUTPUT, JSON.stringify({ sections }, null, 2));
  const total = sections.reduce((n, s) => n + s.children.length, 0);
  console.log(`docs-index.json: ${sections.length} sections, ${total} pages`);
}

buildTree();
