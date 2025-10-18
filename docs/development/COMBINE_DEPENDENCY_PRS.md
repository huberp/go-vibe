# Combine Dependency PRs Workflow

This document explains how to use the automated dependency PR combination workflows.

## Overview

The go-vibe project includes two GitHub Actions workflows to manage dependency updates:

1. **Combine Dependency PRs** (`combine-dependency-prs.yml`) - Combines multiple dependency PRs into one
2. **Cleanup Combined PRs** (`cleanup-combined-prs.yml`) - Automatically closes source PRs after merge

## How It Works

### Step 1: Automatic Dependency PRs

Dependabot creates individual PRs for Go dependency updates:
- Each PR is labeled with `dependencies` and `go`
- Each PR title starts with `chore(deps):`
- Build and Test workflows run automatically

### Step 2: Manual Combination (When Ready)

When you're ready to combine multiple dependency PRs:

1. Navigate to **Actions** tab in GitHub
2. Select **"Combine Dependency PRs"** workflow
3. Click **"Run workflow"** button
4. Configure options:
   - **Branch name**: Name for the combined branch (default: `combined-dependency-updates`)
   - **Delete stale branch**: Whether to delete existing combined branch (default: `true`)
5. Click **"Run workflow"**

### Step 3: Automatic Filtering

The workflow automatically:
- ✅ Finds all open PRs with labels `dependencies` AND `go`
- ✅ Verifies PR titles start with `chore(deps):` or `chore:`
- ✅ Checks that Build and Test workflows have passed (status: SUCCESS)
- ✅ Combines only PRs that meet ALL criteria

### Step 4: Combined PR Creation

The workflow:
- Creates a new branch with all changes
- Opens a combined PR titled `chore(deps): Combined dependency updates`
- Lists all source PRs in the description
- Adds labels `dependencies` and `go`
- Comments on source PRs noting they've been combined

### Step 5: Review and Merge

1. Review the combined PR
2. Wait for Build and Test workflows to pass
3. Merge the combined PR

### Step 6: Automatic Cleanup

After merging the combined PR:
- Source PRs are automatically closed
- Source branches are automatically deleted
- Combined branch is automatically deleted
- Comments are added to source PRs confirming the merge

## Requirements

The workflow only combines PRs that meet ALL of these criteria:

- ✅ PR has label: `dependencies`
- ✅ PR has label: `go`
- ✅ PR title starts with: `chore(deps):` or `chore:`
- ✅ Build workflow status: `SUCCESS`
- ✅ Test workflow status: `SUCCESS`

## Configuration

### Repository Settings (Required)

**Important:** Before using this workflow, you must enable GitHub Actions to create pull requests:

1. Go to your repository **Settings** → **Actions** → **General**
2. Scroll down to **Workflow permissions**
3. Select **"Read and write permissions"**
4. Check **"Allow GitHub Actions to create and approve pull requests"**
5. Click **Save**

Without this setting enabled, the workflow will fail with: `"GitHub Actions is not permitted to create or approve pull requests"`

### Combine Dependency PRs Workflow

You can customize the workflow inputs when running manually:

| Input | Description | Default | Required |
|-------|-------------|---------|----------|
| `combineBranchName` | Name of the branch to combine PRs into | `combined-dependency-updates` | Yes |
| `deleteStaleBranch` | Delete existing combined branch if it exists | `true` | Yes |

### Cleanup Workflow

The cleanup workflow runs automatically when:
- A PR is closed via merge
- The PR title starts with `chore(deps): Combined dependency updates`
- No configuration needed

## Example Scenarios

### Scenario 1: Multiple Passing PRs

**Situation:**
- PR #10: `chore(deps): bump github.com/gin-gonic/gin` - Status: ✅ SUCCESS
- PR #11: `chore(deps): bump go.uber.org/zap` - Status: ✅ SUCCESS  
- PR #12: `chore(deps): bump gorm.io/gorm` - Status: ✅ SUCCESS

**Result:**
- All 3 PRs are combined into one combined PR
- Combined PR includes all dependency updates
- After merge, PRs #10, #11, #12 are automatically closed

### Scenario 2: Mixed Status PRs

**Situation:**
- PR #10: `chore(deps): bump github.com/gin-gonic/gin` - Status: ✅ SUCCESS
- PR #11: `chore(deps): bump go.uber.org/zap` - Status: ❌ FAILURE
- PR #12: `chore(deps): bump gorm.io/gorm` - Status: ✅ SUCCESS

**Result:**
- Only PRs #10 and #12 are combined
- PR #11 is excluded (failed tests)
- PR #11 remains open for investigation

### Scenario 3: Wrong Labels

**Situation:**
- PR #10: `chore(deps): bump github.com/gin-gonic/gin` - Labels: `dependencies`, `go` - ✅
- PR #11: `chore(deps): bump actions/checkout` - Labels: `dependencies`, `github-actions` - ❌

**Result:**
- Only PR #10 is combined
- PR #11 is excluded (missing `go` label, has `github-actions` instead)

### Scenario 4: Merge Conflicts

**Situation:**
- PR #10 and #11 both modify the same line in `go.sum`

**Result:**
- One PR is merged successfully
- Conflicting PR is listed in the combined PR description
- Conflicting PR remains open for manual resolution

## Troubleshooting

### GitHub Actions Not Permitted to Create PRs

**Error:** `"GitHub Actions is not permitted to create or approve pull requests"`

**Cause:** Repository settings prevent GitHub Actions from creating pull requests.

**Solution:**
1. Go to repository **Settings** → **Actions** → **General**
2. Scroll to **Workflow permissions**
3. Enable **"Read and write permissions"**
4. Check **"Allow GitHub Actions to create and approve pull requests"**
5. Click **Save**
6. Re-run the failed workflow

This is a **required one-time setup** for the workflow to function.

### No PRs Found

**Error:** "No PRs matched criteria"

**Solutions:**
- Verify PRs have both `dependencies` AND `go` labels
- Check PR titles start with `chore(deps):` or `chore:`
- Ensure Build and Test workflows have passed
- Check that PRs are still open

### Branch Already Exists

**Error:** "Branch combined-dependency-updates already exists"

**Solutions:**
- Set `deleteStaleBranch` to `true` when running workflow
- Or manually delete the branch and re-run

### Merge Conflicts

**Error:** Some PRs listed as "excluded due to merge conflicts"

**Solutions:**
- Review conflicting PRs individually
- Update dependencies manually to resolve conflicts
- The combined PR includes all non-conflicting PRs

### Only One PR Found

**Error:** "Only one PR matched criteria - no combining needed"

**Solutions:**
- This is expected when only one dependency PR is ready
- No action needed - just merge the single PR normally

## Manual Workflow Steps

If you need to run this manually:

1. **Go to Actions tab**
   ```
   https://github.com/huberp/go-vibe/actions/workflows/combine-dependency-prs.yml
   ```

2. **Click "Run workflow"**

3. **Configure inputs** (or use defaults)

4. **Wait for workflow to complete** (usually < 1 minute)

5. **Review the combined PR** (link will be in workflow logs)

6. **Merge the combined PR** when ready

7. **Source PRs automatically close** after merge

## Benefits

- ✅ **Reduced PR clutter** - Combine multiple dependency updates into one
- ✅ **Safety first** - Only combines PRs that pass all tests
- ✅ **Automatic cleanup** - No manual PR closing needed
- ✅ **Conflict detection** - Clearly identifies PRs with merge conflicts
- ✅ **Full audit trail** - Comments on all PRs for transparency

## Limitations

- Only combines PRs with labels `dependencies` AND `go`
- Requires PRs to have passing Build and Test workflows
- Cannot auto-resolve merge conflicts (PRs with conflicts are excluded)
- Manual trigger required (workflow doesn't run on schedule)

## Future Enhancements

Potential improvements for the future:

- Schedule automatic runs (e.g., weekly)
- Support for other dependency types (Docker, GitHub Actions)
- Auto-merge combined PR if all checks pass
- Slack/Teams notifications
- Custom label combinations
