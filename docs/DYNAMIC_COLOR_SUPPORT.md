# Dynamic Color Support Implementation Summary

## Overview
Implemented automatic terminal color capability detection with graceful fallback support. The system now seamlessly adapts colors from 24-bit TrueColor (16.7 million colors) down to 256-color or 16-color ANSI palettes based on the terminal's capabilities.

## Changes Made

### 1. Terminal Capability Detection (`terminal/terminal.go`)
- Added `ColorCapability` type enum with three levels:
  - `ColorTrueColor` (0): Full 24-bit RGB support
  - `Color256` (1): 256-color palette support
  - `Color16` (2): Basic 16 ANSI colors
- Added `DetectColorCapability()` function that:
  - Checks `COLORTERM` environment variable for "truecolor" or "24bit"
  - Falls back to checking `TERM` variable for "256color" suffix
  - Falls back to checking `TERM` for "color" suffix
  - Defaults to `Color256` if no environment variables are set

### 2. Color Conversion Utilities (`style/color/conversion.go`) - NEW FILE
- **ANSI 16-color palette**: Defined standard ANSI colors (0-15)
- **xterm 256-color palette**: Generated programmatically:
  - First 16 colors: ANSI colors
  - Colors 16-231: 6×6×6 RGB color cube
  - Colors 232-255: 24 grayscale colors
- **Color distance algorithm**: Euclidean distance in RGB space
- **Conversion functions**:
  - `ToNearest256(c Color) uint8`: Maps RGB to nearest 256-color index
  - `ToNearest16(c Color) uint8`: Maps RGB to nearest ANSI 16-color

### 3. ANSI Color Output Functions (`ansi/ansi.go`)
Added capability-aware color output functions:
- `ForegroundColor256(c Color)`: 256-color foreground (format: `38;5;n`)
- `BackgroundColor256(c Color)`: 256-color background (format: `48;5;n`)
- `ForegroundColor16(c Color)`: 16-color ANSI foreground (format: `3n` or `9n`)
- `BackgroundColor16(c Color)`: 16-color ANSI background (format: `4n` or `10n`)
- `ForegroundColorAuto(c Color, capability)`: Auto-selects based on terminal capability
- `BackgroundColorAuto(c Color, capability)`: Auto-selects based on terminal capability

Existing TrueColor functions remain unchanged for backward compatibility:
- `ForegroundColor(c Color)`: 24-bit RGB (format: `38;2;R;G;B`)
- `BackgroundColor(c Color)`: 24-bit RGB (format: `48;2;R;G;B`)

### 4. Renderer Integration (`renderer.go`)
- Added `colorCapability` field to `Renderer` struct
- Updated `NewRenderer()` to auto-detect capability on initialization via `terminal.DetectColorCapability()`
- Modified `outputCell()` to use capability-aware functions:
  - Calls `ansi.ForegroundColorAuto()` instead of `ansi.ForegroundColor()`
  - Calls `ansi.BackgroundColorAuto()` instead of `ansi.BackgroundColor()`
  - Passes detected `colorCapability` to both functions

## How It Works

1. **At Startup**: `NewRenderer()` calls `terminal.DetectColorCapability()` to detect what the terminal supports
2. **During Rendering**: Each cell's colors are converted using the capability-aware functions:
   - TrueColor terminals: RGB colors passed through directly (highest quality)
   - 256-color terminals: RGB colors converted to nearest palette entry
   - 16-color terminals: RGB colors converted to nearest ANSI color
3. **Automatic Fallback**: No configuration needed - detection happens automatically

## Color Conversion Quality

The implementation uses Euclidean distance in RGB space for color matching:
- **Fast**: O(n) where n is palette size (16 or 256)
- **Reliable**: Mathematically sound color distance metric
- **Perceptually reasonable**: Works well for most use cases

For future enhancement, more sophisticated algorithms like CIEDE2000 could be implemented for perceptually-accurate color matching.

## Environment Variable Detection

The detection strategy is robust and follows industry standards:

1. **COLORTERM** (highest priority): Set to "truecolor" or "24bit"
   - Modern terminals like iTerm2, GNOME Terminal, Windows Terminal set this
2. **TERM** variable suffix: Checked for "256color"
   - xterm-256color, screen-256color, etc.
3. **TERM** variable contains "color"
   - Basic color support indicator
4. **Default**: Falls back to 256-color support

## Backward Compatibility

All changes are backward compatible:
- Existing color API remains unchanged
- TrueColor functions still available for direct use
- New capability-aware functions coexist with old ones
- Renderer automatically uses new functions without user changes needed

## Testing

- All code compiles successfully with `go build ./...`
- Example projects build without errors
- Ready for integration testing with actual terminals

