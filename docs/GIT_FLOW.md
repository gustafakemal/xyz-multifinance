# Git Flow Implementation Guide
# PT XYZ Multifinance

## Branch Structure

```
main (production)
  │
  ├─ develop (integration)
  │   │
  │   ├─ feature/transaction-api
  │   ├─ feature/consumer-management
  │   ├─ feature/limit-calculation
  │   └─ feature/security-enhancement
  │
  ├─ release/v1.0.0
  │
  └─ hotfix/critical-bug-fix
```

## Branch Types

### 1. Main Branch
- **Purpose**: Production-ready code only
- **Protected**: Yes (no direct commits)
- **Deployment**: Auto-deploy to production
- **Naming**: `main`

### 2. Develop Branch
- **Purpose**: Integration branch for features
- **Protected**: Yes (require PR approval)
- **Deployment**: Auto-deploy to staging
- **Naming**: `develop`

### 3. Feature Branches
- **Purpose**: New feature development
- **Branch from**: `develop`
- **Merge to**: `develop`
- **Naming**: `feature/<feature-name>`
- **Lifetime**: Temporary (deleted after merge)

### 4. Release Branches
- **Purpose**: Prepare new production release
- **Branch from**: `develop`
- **Merge to**: `main` and `develop`
- **Naming**: `release/v<version>`
- **Lifetime**: Temporary (deleted after release)

### 5. Hotfix Branches
- **Purpose**: Emergency production fixes
- **Branch from**: `main`
- **Merge to**: `main` and `develop`
- **Naming**: `hotfix/<issue-name>`
- **Lifetime**: Temporary (deleted after merge)

## Workflow Examples

### Starting a New Feature

```bash
# 1. Update develop branch
git checkout develop
git pull origin develop

# 2. Create feature branch
git checkout -b feature/transaction-concurrent-handling

# 3. Work on feature (commit frequently)
git add .
git commit -m "feat: implement optimistic locking for transactions"

# 4. Push feature branch
git push origin feature/transaction-concurrent-handling

# 5. Create Pull Request to develop
# (via GitHub/GitLab web interface)

# 6. After PR approval and merge, delete feature branch
git checkout develop
git pull origin develop
git branch -d feature/transaction-concurrent-handling
```

### Creating a Release

```bash
# 1. Create release branch from develop
git checkout develop
git pull origin develop
git checkout -b release/v1.0.0

# 2. Update version numbers, documentation
# Update README.md, CHANGELOG.md, etc.
git add .
git commit -m "chore: prepare release v1.0.0"

# 3. Push release branch
git push origin release/v1.0.0

# 4. Create PR to main
# After testing and approval:

# 5. Merge to main
git checkout main
git merge --no-ff release/v1.0.0
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin main --tags

# 6. Merge back to develop
git checkout develop
git merge --no-ff release/v1.0.0
git push origin develop

# 7. Delete release branch
git branch -d release/v1.0.0
git push origin --delete release/v1.0.0
```

### Creating a Hotfix

```bash
# 1. Create hotfix branch from main
git checkout main
git pull origin main
git checkout -b hotfix/critical-security-fix

# 2. Fix the issue
git add .
git commit -m "fix: patch SQL injection vulnerability"

# 3. Push hotfix branch
git push origin hotfix/critical-security-fix

# 4. Merge to main
git checkout main
git merge --no-ff hotfix/critical-security-fix
git tag -a v1.0.1 -m "Hotfix version 1.0.1"
git push origin main --tags

# 5. Merge to develop
git checkout develop
git merge --no-ff hotfix/critical-security-fix
git push origin develop

# 6. Delete hotfix branch
git branch -d hotfix/critical-security-fix
git push origin --delete hotfix/critical-security-fix
```

## Commit Message Convention

Following [Conventional Commits](https://www.conventionalcommits.org/):

### Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```bash
# Feature
git commit -m "feat(transaction): add concurrent transaction handling"

# Bug fix
git commit -m "fix(limit): correct calculation of available limit"

# Documentation
git commit -m "docs(readme): update API endpoint documentation"

# Refactoring
git commit -m "refactor(repository): extract common query logic"

# Test
git commit -m "test(transaction): add unit tests for validation"

# Performance
git commit -m "perf(database): optimize transaction query with index"

# Breaking change
git commit -m "feat(api): change transaction response format

BREAKING CHANGE: Transaction API now returns different JSON structure"
```

## Pull Request Guidelines

### PR Title Format
```
[Type] Brief description

Example:
[Feature] Implement concurrent transaction handling
[Fix] Resolve limit calculation bug
[Docs] Update architecture documentation
```

### PR Description Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Manual testing completed
- [ ] All tests passing

## Checklist
- [ ] Code follows project conventions
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or documented)

## Related Issues
Closes #123
```

### PR Review Checklist

Reviewers should check:
- [ ] Code quality and readability
- [ ] Test coverage
- [ ] Documentation updates
- [ ] No security vulnerabilities
- [ ] Performance implications
- [ ] Backward compatibility

## Branch Protection Rules

### Main Branch
```yaml
Required:
  - Pull request reviews: 2 approvals
  - Status checks must pass
  - Up-to-date with base branch
  - No direct commits
  - Signed commits recommended
```

### Develop Branch
```yaml
Required:
  - Pull request reviews: 1 approval
  - Status checks must pass
  - No direct commits
```

## Versioning Strategy

Following [Semantic Versioning](https://semver.org/):

```
MAJOR.MINOR.PATCH

Example: v1.2.3
- MAJOR: Breaking changes (v1 → v2)
- MINOR: New features, backward compatible (v1.1 → v1.2)
- PATCH: Bug fixes, backward compatible (v1.2.1 → v1.2.2)
```

### When to Bump Versions

**MAJOR** version:
- Breaking API changes
- Database schema breaking changes
- Major architecture changes

**MINOR** version:
- New features
- New API endpoints
- Backward compatible changes

**PATCH** version:
- Bug fixes
- Security patches
- Performance improvements
- Documentation updates

## Git Hooks (Recommended)

### Pre-commit Hook
```bash
#!/bin/sh
# .git/hooks/pre-commit

# Run tests
go test ./...
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi

# Run linter
golangci-lint run
if [ $? -ne 0 ]; then
    echo "Linting failed. Commit aborted."
    exit 1
fi
```

### Commit-msg Hook
```bash
#!/bin/sh
# .git/hooks/commit-msg

commit_msg=$(cat "$1")

# Check commit message format
if ! echo "$commit_msg" | grep -qE "^(feat|fix|docs|style|refactor|perf|test|chore)(\(.+\))?: .+"; then
    echo "Invalid commit message format."
    echo "Use: <type>(<scope>): <subject>"
    exit 1
fi
```

## CI/CD Integration

### GitHub Actions Workflow Example

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [develop, main]
  pull_request:
    branches: [develop, main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.23
      - name: Run Tests
        run: go test ./... -v
      - name: Run Linter
        run: golangci-lint run

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build Docker Image
        run: docker build -t xyz-multifinance:${{ github.sha }} .
      - name: Push to Registry
        run: docker push xyz-multifinance:${{ github.sha }}

  deploy-staging:
    needs: build
    if: github.ref == 'refs/heads/develop'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Staging
        run: |
          # Deployment commands here

  deploy-production:
    needs: build
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Production
        run: |
          # Deployment commands here
```

## Best Practices

1. **Commit Often**: Small, atomic commits
2. **Clear Messages**: Descriptive commit messages
3. **Pull Regularly**: Keep branches up-to-date
4. **Test Before PR**: All tests should pass
5. **Review Carefully**: Thorough code reviews
6. **Clean History**: Squash commits if needed
7. **Tag Releases**: Always tag production releases
8. **Document Changes**: Update CHANGELOG.md

## Troubleshooting

### Merge Conflicts
```bash
# 1. Update your branch with latest develop
git checkout feature/your-feature
git fetch origin
git merge origin/develop

# 2. Resolve conflicts manually
# 3. Mark as resolved
git add .
git commit -m "chore: resolve merge conflicts"
```

### Undo Last Commit
```bash
# Keep changes
git reset --soft HEAD~1

# Discard changes
git reset --hard HEAD~1
```

### Revert Commit
```bash
# Create a new commit that undoes changes
git revert <commit-hash>
```

---

**Document Version**: 1.0  
**Date**: February 1, 2026  
**Author**: PT XYZ Multifinance Development Team
