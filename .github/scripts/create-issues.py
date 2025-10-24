#!/usr/bin/env python3
"""
Script to parse CODE_REVIEW_ISSUES.md and create GitHub issues.
"""
import os
import re
import sys
from github import Github


def parse_issues_file(filepath):
    """Parse the CODE_REVIEW_ISSUES.md file and extract structured issue data."""
    with open(filepath, 'r') as f:
        content = f.read()

    # Parse issues from the markdown file
    # Split by "### Issue N:" pattern
    issue_pattern = r'### Issue (\d+): ([^\n]+)\n\*\*Title\*\*: ([^\n]+)\n\n\*\*Description\*\*:\n(.*?)\n\n\*\*Goals\*\*:\n(.*?)\n\n\*\*Files to modify\*\*:\n(.*?)\n\n\*\*Acceptance Criteria\*\*:\n(.*?)\n\n\*\*Labels\*\*: ([^\n]+)\n\*\*Estimated effort\*\*: ([^\n]+)'
    
    issues = []
    for match in re.finditer(issue_pattern, content, re.DOTALL):
        issue_num = int(match.group(1))
        issue_short_title = match.group(2).strip()
        issue_title = match.group(3).strip()
        description = match.group(4).strip()
        goals = match.group(5).strip()
        files = match.group(6).strip()
        acceptance = match.group(7).strip()
        labels_str = match.group(8).strip()
        effort = match.group(9).strip()
        
        issues.append({
            'number': issue_num,
            'short_title': issue_short_title,
            'title': issue_title,
            'description': description,
            'goals': goals,
            'files': files,
            'acceptance': acceptance,
            'labels': [l.strip() for l in labels_str.split(',')],
            'effort': effort
        })
    
    return issues


def create_issue_body(issue):
    """Create the GitHub issue body from parsed issue data."""
    body = f"""## Description

{issue['description']}

## Goals

{issue['goals']}

## Files to modify

{issue['files']}

## Acceptance Criteria

{issue['acceptance']}

---
**Estimated Effort**: {issue['effort']}
**Related to**: Code Review 2025-10-24
"""
    return body


def main():
    # Configuration from environment
    dry_run = os.getenv('DRY_RUN', 'true').lower() == 'true'
    start_issue = int(os.getenv('START_ISSUE', '1'))
    end_issue = int(os.getenv('END_ISSUE', '20'))
    
    print(f"Configuration:")
    print(f"  Dry run: {dry_run}")
    print(f"  Issue range: {start_issue} to {end_issue}")
    print()

    # Parse issues file
    issues_file = 'code-reviews/2025-10-24-comprehensive-review/CODE_REVIEW_ISSUES.md'
    issues = parse_issues_file(issues_file)
    
    print(f"Parsed {len(issues)} issues from CODE_REVIEW_ISSUES.md")
    print()

    # Filter issues based on range
    filtered_issues = [issue for issue in issues if start_issue <= issue['number'] <= end_issue]
    print(f"Creating {len(filtered_issues)} issues (range: {start_issue}-{end_issue})")
    print()

    if not dry_run:
        # Initialize GitHub client
        github_token = os.getenv('GITHUB_TOKEN')
        if not github_token:
            print("ERROR: GITHUB_TOKEN environment variable not set")
            sys.exit(1)
        
        g = Github(github_token)
        
        github_repo = os.getenv('GITHUB_REPOSITORY')
        if not github_repo:
            print("ERROR: GITHUB_REPOSITORY environment variable not set")
            sys.exit(1)
        
        repo = g.get_repo(github_repo)
        
        # Get existing labels
        existing_labels = {label.name for label in repo.get_labels()}

    # Process each issue
    created_count = 0
    for issue in filtered_issues:
        body = create_issue_body(issue)
        
        print(f"{'[DRY RUN] ' if dry_run else ''}Issue #{issue['number']}: {issue['title']}")
        print(f"  Labels: {', '.join(issue['labels'])}")
        print(f"  Effort: {issue['effort']}")
        
        if dry_run:
            print(f"  Would create issue with title: {issue['title']}")
            print()
        else:
            # Create labels if they don't exist
            for label in issue['labels']:
                if label not in existing_labels:
                    print(f"  Creating label: {label}")
                    # Use default color for new labels
                    repo.create_label(name=label, color="0366d6")
                    existing_labels.add(label)
            
            # Create the issue
            try:
                created_issue = repo.create_issue(
                    title=issue['title'],
                    body=body,
                    labels=issue['labels']
                )
                print(f"  ✓ Created issue #{created_issue.number}: {created_issue.html_url}")
                print()
                created_count += 1
            except Exception as e:
                print(f"  ✗ Error creating issue: {e}")
                print()

    print("Done!")
    if dry_run:
        print("\nThis was a DRY RUN. No issues were created.")
        print("To create issues, run this workflow again with dry_run=false")
    else:
        print(f"\n✓ Successfully created {created_count} of {len(filtered_issues)} issues")


if __name__ == "__main__":
    main()
