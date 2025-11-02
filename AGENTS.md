# Repository Guidelines

## Project Structure & Module Organization
The repo is documentation-first: `docs/` holds requirements, technical design, and the interactive `prototype.html`. Keep these files canonical and funnel implementation notes back into the design docs. When the React + TypeScript app lands, place runtime code in `src/`, static assets in `public/`, and feature tests in `src/<feature>/__tests__`, mirroring the hierarchy described in `docs/技术设计文档.md`.

## Build, Test, and Development Commands
Bootstrap the app with Vite using `npm create vite@latest . -- --template react-ts`. Install dependencies with `npm install`, then start the dev server via `npm run dev` (hot reload on `http://localhost:5173`). Add `npm run build` for production bundles and `npm run preview` to verify the output; document any new scripts in `package.json`.

## Coding Style & Naming Conventions
The technical design locks in ESLint + Prettier; enforce the shared config with a pre-commit hook. Use 2-space indentation, TypeScript strict mode, and functional React components. Name components and files with `PascalCase`, hooks with `useCamelCase`, Zustand stores with `<name>Store.ts`, and utility modules with `camelCase`. Keep state slices thin and co-locate Zustand selectors near their consumers to match the design.

## Testing Guidelines
Adopt Vitest and React Testing Library for unit and integration coverage, reserving Playwright (or similar) for drag-and-drop and keyboard shortcuts. Name files `*.test.ts` for logic and `*.test.tsx` for component suites. Aim for 80% line coverage, prioritizing flows listed in `docs/需求说明.md`. Run `npm run test` locally before pushing and grow coverage expectations as the suite matures.

## Commit & Pull Request Guidelines
Follow Conventional Commits (`feat:`, `fix:`, `chore:`, etc.) so automated changelogs stay predictable. Keep commits focused—separate documentation updates from code changes. Pull requests should summarize scope, link the relevant requirement or design section, and include screenshots or terminal output for user-facing changes. Request review from a product stakeholder for requirement alignment and an engineer for implementation details when behavior shifts.

## Documentation & Knowledge Sharing
This repository treats documentation as the source of truth. Update `docs/技术设计文档.md` whenever architecture decisions deviate and record rationale inline. Capture open questions or follow-up tasks in a shared tracker so the next agent can pick up context quickly.
