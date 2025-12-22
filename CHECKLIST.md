# Post-Migration Checklist

Use this checklist to verify the monorepo migration was successful.

## ✅ Basic Verification

- [ ] All tests pass: `make test`
- [ ] Quick tests run: `make test.quick`
- [ ] E2E tests pass: `make test.e2e`
- [ ] Help command works: `make help`
- [ ] Backend directory exists with all files
- [ ] Root directory contains monorepo documentation

## ✅ Configuration Files

- [ ] `backend/.kyora.yaml` exists (local config, gitignored)
- [ ] `backend/.kyora.yaml.example` exists (template, tracked)
- [ ] `backend/.air.toml` exists (hot reload config)
- [ ] `.editorconfig` exists at root (code style)
- [ ] `.gitignore` updated with backend paths

## ✅ Documentation

- [ ] `/README.md` describes monorepo structure
- [ ] `/backend/README.md` describes backend
- [ ] `/STRUCTURE.md` explains monorepo guidelines
- [ ] `/MIGRATION.md` documents what changed

## ✅ Makefile Commands

Test each command to ensure it works:

- [ ] `make dev.server` - Starts development server
- [ ] `make test` - Runs all tests
- [ ] `make test.unit` - Runs unit tests
- [ ] `make test.e2e` - Runs E2E tests
- [ ] `make test.quick` - Runs tests (no verbose)
- [ ] `make test.coverage` - Generates coverage report
- [ ] `make test.coverage.html` - Generates HTML coverage
- [ ] `make test.coverage.view` - Opens coverage in browser
- [ ] `make test.e2e.coverage` - E2E with coverage
- [ ] `make clean.coverage` - Cleans coverage files
- [ ] `make clean.backend` - Cleans backend artifacts
- [ ] `make help` - Shows help message

## ✅ Git Status

Check what's staged for commit:

```bash
git status
```

Expected changes:
- **Modified**: `Makefile`, `.gitignore`, `README.md`
- **Renamed/Moved**: All backend files to `backend/` directory
- **New**: `backend/README.md`, `STRUCTURE.md`, `MIGRATION.md`, `.editorconfig`, `backend/.kyora.yaml.example`
- **Not tracked**: `backend/.kyora.yaml` (should be in .gitignore)

## ✅ Development Workflow

Verify your development workflow still works:

1. [ ] Start dev server: `make dev.server`
2. [ ] Make a code change
3. [ ] Verify hot reload works (Air should rebuild)
4. [ ] Run tests: `make test`
5. [ ] Check coverage: `make test.coverage.view`

## ✅ IDE/Editor

- [ ] VS Code workspace still loads correctly
- [ ] Go extension finds backend code
- [ ] Debugger configurations work
- [ ] Test explorer finds tests

## ⚠️ Potential Issues to Check

- [ ] CI/CD pipelines might need path updates
- [ ] Deployment scripts might reference old paths
- [ ] Docker files might need path updates
- [ ] Any custom scripts that reference project files

## 🚀 Ready for Next Steps

Once all items are checked:

1. [ ] Commit changes with clear message:
   ```bash
   git add .
   git commit -m "refactor: restructure project as monorepo

   - Move Go backend code to backend/ directory
   - Update Makefile for monorepo structure
   - Add comprehensive documentation
   - Prepare for future frontend and mobile projects

   See MIGRATION.md for full details"
   ```

2. [ ] Update team/collaborators about new structure

3. [ ] Begin planning frontend/mobile projects

## 📝 Notes

Add any observations or issues encountered:

```
[Add notes here]
```

---

**Checklist completed**: _______________
**Completed by**: _______________
**Date**: _______________
