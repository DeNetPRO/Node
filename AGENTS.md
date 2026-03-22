# AGENTS.md - DeNet Node Documentation Repository

## Repository Overview
This repository contains documentation, guides, and assets for the DeNet Node project - a decentralized data storage protocol. The content is primarily markdown documentation with supporting images and static web assets.

## File Structure
- `/readme.md` - Main project README
- `/guides/` - Documentation guides (installation, configuration, monitoring, etc.)
- `/assets/` - Images and static assets used in documentation
- `/scripts/` - Shell scripts for installer/manager
- `/guides/css/` - Stylesheets for web guides
- `/guides/js/` - JavaScript for web guides

## Content Guidelines

### Markdown Files
- Use standard GitHub Flavored Markdown
- Keep line wrapping at ~100 characters for readability
- Use relative paths for internal links: `./guides/monitoring.md`
- Use absolute paths for external resources
- Images stored in `/assets/` referenced as `assets/image.png` or `../assets/image.png`

### Naming Conventions
- Files: lowercase with hyphens (e.g., `install-denode-windows.md`)
- Directories: lowercase with hyphens
- Headings: Sentence case (capitalize first letter and proper nouns)

### Link Structure
- Internal links: `./filename.md` or `../filename.md`
- Anchor links: Use kebab-case for heading anchors
- External links: Include protocol (https://)

### Image Assets
- Store in `/assets/` directory
- Use descriptive names: `install-step-1.png`, `config-screen.png`
- WebP format preferred in `/assets/webp/` for modern browsers

### Scripts (Shell)
- Use bash shebang: `#!/bin/bash`
- Follow POSIX shell conventions where possible
- Include error handling with `set -e`
- Comment sections clearly

## Writing Guidelines

### Tone and Style
- Clear, instructional language
- Use imperative mood for steps ("Click the button", not "You should click")
- Include warnings and tips using blockquotes
- Provide platform-specific instructions where needed

### Documentation Structure
1. Start with overview/introduction
2. Include prerequisites/requirements
3. Step-by-step instructions with numbered lists
4. Include expected outputs/screenshots
5. Add troubleshooting sections for common issues

### Code/Command Blocks
```bash
# Use language-specific syntax highlighting
command --flag "argument"
```

### Tables
- Use for version matrices, platform support, download links
- Include headers and alignment
- Keep cell content concise

## Version Management
- Document version numbers in release notes
- Keep download links updated with latest releases
- Note platform-specific version requirements

## Common Tasks

### Adding a New Guide
1. Create markdown file in `/guides/`
2. Use descriptive hyphenated filename
3. Add link to main readme.md Table of Contents
4. Include necessary screenshots in `/assets/`
5. Update any cross-references

### Updating Download Links
1. Update version numbers in all affected files
2. Verify all platform links are current
3. Update architecture-specific instructions if needed

### Image Management
1. Add new images to `/assets/`
2. Create WebP versions in `/assets/webp/` when possible
3. Compress images before adding (optimize for web)
4. Update references in documentation

## Quality Checks
- Verify all internal links resolve correctly
- Check image paths and loading
- Ensure platform-specific sections are accurate
- Validate code blocks and commands
- Review for consistent terminology

## No Build Process Required
This is a documentation repository - no compilation, linting, or testing is required. Changes are deployed directly from the main branch.

## Development Commands

### Previewing Changes
```bash
# Serve locally with any static server
python3 -m http.server 8000
# or
npx serve .
```

### Link Checking
```bash
# Manual verification - check all internal links resolve
# Use browser dev tools or markdown link checker extensions
```

## Git Workflow
```bash
# Standard workflow for documentation changes
git checkout -b docs/update-guide-name
# Make changes
git add .
git commit -m "docs: update installation guide with new screenshots"
git push origin docs/update-guide-name
# Create PR for review
```

## Commit Message Conventions
- Prefix with `docs:` for documentation changes
- Keep messages concise and descriptive
- Example: `docs: add Windows installation guide`

## Agent Guidelines

### When Adding Documentation
1. Check existing guides in `/guides/` for similar content
2. Follow established structure and formatting
3. Add screenshots to `/assets/` and create WebP versions
4. Update Table of Contents in readme.md if adding new pages
5. Verify all links work before committing

### When Updating Existing Content
1. Maintain consistent tone and style
2. Update screenshots if UI/commands change
3. Check cross-references are still valid
4. Keep version numbers current

### Image Best Practices
- Use WebP format for smaller file sizes
- Compress images before adding
- Name files descriptively: `feature-name-step.png`
- Include alt text in markdown: `![description](path)`

### Platform-Specific Content
- Use tabs or clear sections for different OS instructions
- Label commands with language: ```bash, ```powershell
- Test commands on each platform before documenting

## Troubleshooting Common Issues

### Broken Links
- Run link checker or manually verify after edits
- Use relative paths for internal links
- Check anchor link casing (GitHub uses lowercase)

### Image Issues
- Verify file exists at specified path
- Check file permissions
- Ensure WebP fallbacks exist for critical images

### Formatting Problems
- Keep lines at ~100 chars for readability
- Use consistent heading hierarchy (## then ###)
- Maintain consistent list formatting
