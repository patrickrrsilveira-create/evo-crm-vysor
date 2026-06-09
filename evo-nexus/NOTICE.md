# Third-Party Notices

EvoNexus includes open source software developed by third parties.
This file lists attributions for all third-party code included in this workspace.

---

## Engineering Layer — derived from oh-my-claudecode (OMC)

**19 of the 21 engineering agents** and **all `dev-*` skills** in this workspace are derived from
[**oh-my-claudecode**](https://github.com/yeachan-heo/oh-my-claudecode) (OMC),
an open source multi-agent orchestration framework for Claude Code by **Yeachan Heo**.

The 2 remaining engineering agents (`helm-conductor` and `mirror-retro`) plus the canonical
6-phase workflow (`.claude/rules/dev-phases.md`) are EvoNexus-native additions introduced in
v0.13.0 — they are not derived from OMC and are original EvoNexus work.

- **Source:** https://github.com/yeachan-heo/oh-my-claudecode
- **Version imported:** v4.11.4
- **License:** MIT
- **Copyright:** Copyright (c) 2025 Yeachan Heo
- **Date imported:** 2026-04-10

### Agents derived from OMC

| EvoNexus name | Original OMC name |
|---|---|
| apex-architect | architect |
| echo-analyst | analyst |
| compass-planner | planner |
| raven-critic | critic |
| bolt-executor | executor |
| hawk-debugger | debugger |
| lens-reviewer | code-reviewer |
| zen-simplifier | code-simplifier |
| vault-security | security-reviewer |
| grid-tester | test-engineer |
| probe-qa | qa-tester |
| oath-verifier | verifier |
| trail-tracer | tracer |
| scout-explorer | explore |
| flow-git | git-master |
| scroll-docs | document-specialist |
| quill-writer | writer |
| canvas-designer | designer |
| prism-scientist | scientist |

### Skills derived from OMC (`dev-*` namespace)

`dev-autopilot`, `dev-plan`, `dev-ralplan`, `dev-deep-dive`, `dev-deep-interview`,
`dev-external-context`, `dev-trace`, `dev-verify`, `dev-ultraqa`, `dev-visual-verdict`,
`dev-ai-slop-cleaner`, `dev-sciomc`, `dev-team`, `dev-ccg`, `dev-ralph`,
`dev-mcp-setup`, `dev-deepinit`, `dev-project-session-manager`,
`dev-configure-notifications`, `dev-release`, `dev-remember`, `dev-ask`,
`dev-learner`, `dev-skillify`, `dev-cancel`

### Modifications made to OMC content

- Renamed agents from generic names (e.g., `architect`) to themed names
  (e.g., `apex-architect`) to fit EvoNexus naming convention.
- Added `dev-` prefix to all imported skills for namespace consistency with
  existing EvoNexus prefixes (`fin-`, `social-`, `legal-`, etc.).
- Adapted memory structure to match EvoNexus per-agent memory pattern
  (`.claude/agent-memory/<agent-name>/`).
- Removed runtime dependencies on OMC TypeScript `src/` — only the markdown
  definitions (`agents/*.md`, `skills/*/SKILL.md`) were imported.
- Discarded OMC skills that overlapped with EvoNexus builtins or were
  OMC meta-skills (`omc-setup`, `omc-doctor`, `omc-reference`, `hud`, `debug`,
  `skill`, `wiki`, `writer-memory`, `ultrawork`, `omc-teams`, `learn-about-omc`).

### MIT License

```
MIT License

Copyright (c) 2025 Yeachan Heo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
