// SPDX-FileCopyrightText: 2025 maxhash.io <dev@maxhash.io>
//
// SPDX-License-Identifier: AGPL-3.0-only

module.exports = {
  overrides: [
    {
      files: ["*.html"],
      options: {
        tabWidth: 2,
        useTabs: false,
        printWidth: 120,
        singleQuote: false,
        bracketSpacing: true,
        htmlWhitespaceSensitivity: "css",
      },
    },
  ],
};
