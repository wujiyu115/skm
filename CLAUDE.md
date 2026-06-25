# SKM Development Notes

## TypeScript Checking

Always use `cd web && npx tsc -b` (or `npm run build`) for type checking — NOT `npx tsc --noEmit`.

CI runs `npm run build` which is `tsc -b && vite build`. The `-b` flag uses project references (tsconfig.json → tsconfig.app.json) which enables `noUnusedLocals` and `noUnusedParameters`. Plain `tsc --noEmit` skips these checks, causing CI failures after push.

## Build Commands

```bash
# Frontend type check (matches CI)
cd web && npx tsc -b

# Full frontend build
cd web && npm run build

# Go build
go build ./...

# Go tests
go test ./...
```
