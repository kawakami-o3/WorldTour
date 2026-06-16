---
name: next-article
description: Use when Codex needs to create the next WorldTour article in this repository, including prompts such as "次の記事を作って", "記事を書いて", "next-article", or requests to pick an unused country, write articles/YYYY-MM-DD.md, and update INDEX.md.
---

# Next Article

## Overview

Create one WorldTour article end to end: choose the next article date, pick an unused country, research current local context, write `articles/YYYY-MM-DD.md`, update `INDEX.md`, and report the result.

`CLAUDE.md` is authoritative for article format, source standards, Japanese writing rules, and project-specific constraints. This skill only maps that workflow to Codex tools.

## Required Context

Before acting:

- Read `CLAUDE.md` completely.
- Inspect the latest article files and the top of `INDEX.md`.
- Check the worktree with `git status --short` and preserve unrelated user changes.
- Work from the repository root.

Do not ask for intermediate confirmation unless blocked by missing tools, exhausted countries, or unavailable research access.

## Workflow

1. Determine the next article date from `articles/`, not from the real calendar.

   ```bash
   LATEST=$(ls articles/*.md | sort | tail -1 | xargs basename | sed 's/\.md$//')
   NEXT=$(date -j -v+1d -f "%Y-%m-%d" "$LATEST" "+%Y-%m-%d")
   echo "$NEXT"
   ```

   If `articles/$NEXT.md` already exists, append the new article section there instead of creating a second file.

2. Pick the country.

   ```bash
   test -x bin/pickrandom || make
   bin/pickrandom
   ```

   Use the country shown after `選ばれた国:`. The command updates `data/history.json`; do not undo that unless the user explicitly asks. If the command says every country has been selected, stop and report that no article was created.

3. Research with Codex-available tools.

   - Browse the web because the article depends on current local information.
   - Use `web.search_query` and `web.open` when available.
   - Search in English, local official language(s), and relevant neighboring or stakeholder languages.
   - Include regional names, border areas, local organizations, minority groups, old country or city names, and native spellings when useful.
   - Aim for the source mix required by `CLAUDE.md`: local-language regional outlets, local-language national outlets, one neighboring or stakeholder-country source, no more than two international sources, and at least one public or official source.
   - For each source, track language, country, media type, original title, literal Japanese title translation, summary, and why it is locally grounded.
   - If browsing a source is blocked by 403/429, dynamic rendering, or navigation, use the `agent-browser` skill/CLI as a serial headless fallback, reuse one tab, and close it afterward.
   - If the source mix cannot be fully met, state the gap honestly in the audit table.

4. Write the article with `apply_patch` in small sections.

   Follow `CLAUDE.md` for the exact article format. Preserve these Codex-specific rules:

   - Title format: `# $NEXT 国名 (English / 現地語)`.
   - Write in Japanese.
   - Do not use Markdown bold syntax `**`.
   - Link to Japanese Wikipedia or other Japanese references for encyclopedic basics instead of rewriting them.
   - Focus the body on local circumstances that are thinly covered in Japanese sources.
   - Write in these chunks: title plus basic information table, local current issue, why it matters locally, historical/geographic/international background, what English reporting misses, uncertainty/source disagreement plus audit table and terminology table.

5. Update `INDEX.md`.

   Insert a new top data row under the table header and separator:

   ```markdown
   | [$NEXT](articles/$NEXT.md) | 国名 (English / 現地語) | トピックの1〜2文要約 |
   ```

   Keep the table in descending date order.

6. Verify before reporting.

   Check at minimum:

   ```bash
   git diff --check
   rg -n "\\*\\*" "articles/$NEXT.md"
   sed -n '1,12p' INDEX.md
   bin/pickrandom -remaining
   ```

   `rg` finding `**` in the new article means the article must be edited before completion.

7. Report in one line.

   Include the article path, country, article date, and remaining unselected country count.

## Common Mistakes

- Using today's real date instead of the latest article date plus one day.
- Treating Claude Code tool names such as `WebFetch`, `Agent`, or `Write/Edit` as Codex tools.
- Researching only in English or relying mainly on international media.
- Forgetting that `bin/pickrandom` mutates `data/history.json`.
- Updating the article but not `INDEX.md`.
- Leaving `**` bold markers in the article.
- Running multiple `agent-browser` sessions or leaving the browser open.
