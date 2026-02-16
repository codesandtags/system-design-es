# Node.js Cheat Sheet

## Overview and Updates

### The Modern Era: "Batteries Included" (v16–v22 Current)

This is where things get impressive. Influenced by competitors like Deno and Bun, Node.js realized it needed to stop relying on third-party libraries for basic functionality.

- **ES Modules (ESM) are now First-Class**: We can finally use `import` and `export` natively instead of `require()`.
- **Top-Level Await**: You can use `await` outside of an `async` function in module files. This is huge for database connections or config loading at startup.
- **Staff Insight**: Legacy apps (CommonJS) and new apps (ESM) often conflict. Knowing how to configure `package.json` (`"type": "module"`) is standard now.
- **Native Fetch (`global.fetch`)**: Gone are the days of `npm install node-fetch` or `axios` (unless you need interceptors). Node v18+ has `fetch` built-in, compliant with the Web Standard.
- **Native Test Runner (`node:test`)**: You don't need Jest or Mocha for everything anymore. `import { test } from 'node:test';` is built-in, fast, and supports mocking/spying as of recent versions.
- **Watch Mode (`node --watch`)**: Rip `nodemon`. Node can now watch files and restart the process natively.
- **Built-in .env Support**: Node v20+ supports `node --env-file=.env`. No more `dotenv` package for basic setups.
- **Experimental TypeScript Support**: Node can now strip types (via `--experimental-strip-types`) to run `.ts` files directly. It's early days, but it shows the direction.