# Code Review Issues Creation Workflow

## Overview

This GitHub Action workflow automates the creation of GitHub issues from the code review document `CODE_REVIEW_ISSUES.md`.

## Usage

### Manual Trigger

The workflow uses `workflow_dispatch` trigger, which means it must be triggered manually from the GitHub Actions UI.

### Steps to Run

1. Go to the repository on GitHub
2. Click on "Actions" tab
3. Select "Create Code Review Issues" workflow
4. Click "Run workflow" button
5. Configure the run parameters:
   - **dry_run**: Choose `true` to preview what will be created without actually creating issues (recommended first run)
   - **start_issue**: Start from issue number (1-20)
   - **end_issue**: End at issue number (1-20)
6. Click "Run workflow" to start

### Parameters

- **dry_run** (default: `true`)
  - `true`: Show what would be created without actually creating issues
  - `false`: Create the issues in GitHub
  - Recommended to run with `true` first to verify everything looks correct

- **start_issue** (default: `1`)
  - Starting issue number from CODE_REVIEW_ISSUES.md
  - Valid range: 1-20

- **end_issue** (default: `20`)
  - Ending issue number from CODE_REVIEW_ISSUES.md
  - Valid range: 1-20

### Examples

**Example 1: Preview all issues (dry run)**
- dry_run: `true`
- start_issue: `1`
- end_issue: `20`

This will parse all 20 issues and show what would be created.

**Example 2: Create only high-priority issues (issues 1-4)**
- dry_run: `false`
- start_issue: `1`
- end_issue: `4`

This will create issues 1-4 from the code review.

**Example 3: Create all issues**
- dry_run: `false`
- start_issue: `1`
- end_issue: `20`

This will create all 20 issues from the code review.

## What the Workflow Does

1. **Parses** the `CODE_REVIEW_ISSUES.md` file
2. **Extracts** structured information:
   - Title
   - Description
   - Goals
   - Files to modify
   - Acceptance criteria
   - Labels
   - Estimated effort
3. **Creates** GitHub issues with:
   - Proper title and body
   - Appropriate labels (creates labels if they don't exist)
   - Reference to the code review

## Issue Format

Each created issue will have:

```markdown
## Description
[Issue description]

## Goals
[Issue goals]

## Files to modify
[Files to modify]

## Acceptance Criteria
[Acceptance criteria checklist]

---
**Estimated Effort**: [Small/Medium/Large]
**Related to**: Code Review 2025-10-24
```

## Permissions

The workflow requires:
- `issues: write` - To create issues
- `contents: read` - To read the repository files

## Troubleshooting

### Labels are not created
The workflow automatically creates labels if they don't exist. If you see an error about labels, ensure the workflow has `issues: write` permission.

### Parsing errors
If issues are not parsed correctly, verify that `CODE_REVIEW_ISSUES.md` follows the expected format with proper markdown headers and fields.

## Notes

- The workflow uses PyGithub library to interact with the GitHub API
- Issues are created using the GitHub token provided by the workflow
- Labels are created automatically with a default blue color (0366d6)
- The workflow processes issues sequentially to avoid rate limiting

## Workflow Location

The workflow is located at: `.github/workflows/create-code-review-issues.yml`

## Related Files

- `code-reviews/2025-10-24-comprehensive-review/CODE_REVIEW_ISSUES.md` - Source document with all issues
- `code-reviews/2025-10-24-comprehensive-review/NEXT_STEPS.md` - Guide for creating issues
