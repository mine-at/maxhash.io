// SPDX-FileCopyrightText: 2025 maxhash.io <dev@maxhash.io>
//
// SPDX-License-Identifier: AGPL-3.0-only

// Take the given Bitcoin mining difficulty number and format it into a
// human-readable string with appropriate units (K, M, G, T).
function formatBitcoinDifficulty(difficulty) {
  if (typeof difficulty !== "number" || isNaN(difficulty)) return "Invalid";

  const units = [
    { value: 1e15, symbol: "P" },
    { value: 1e12, symbol: "T" },
    { value: 1e9, symbol: "G" },
    { value: 1e6, symbol: "M" },
    { value: 1e3, symbol: "K" },
  ];

  for (const unit of units) {
    if (difficulty >= unit.value) {
      return (
        (difficulty / unit.value).toFixed(2).replace(/\.00$/, "") + unit.symbol
      );
    }
  }

  return difficulty.toString();
}
