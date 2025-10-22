# PowerShell Development Guidelines

### Script Output
- **Never use emojis in PowerShell scripts** - Emojis can cause encoding issues and display problems across different PowerShell versions and terminals
- Use plain text for all output messages
- Use `Write-Host` for user-facing messages
- Use standard status messages like "Success", "Error", "Warning" instead of emoji indicators
