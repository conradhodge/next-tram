var fs = require("fs");

if (fs.existsSync(".git/hooks") && !process.env.CI) {
  fs.copyFileSync("scripts/githooks/pre-commit", ".git/hooks/pre-commit");
  fs.copyFileSync("scripts/githooks/post-checkout", ".git/hooks/post-checkout");
  fs.copyFileSync("scripts/githooks/post-merge", ".git/hooks/post-merge");
}
