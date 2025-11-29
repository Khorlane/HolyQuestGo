package server
//***********************************************************
// ANSI color
// Command format: Esc[x;ym
//
// In C, Escape (Esc) is coded as \x1B
// so the command becomes \x1B[x;ym
// The m at the end is a constant, leaving just the x and y
// values to be specified. There are many possible
// combinations, but 'x' = 1 and 'y' = 31–36 is common.
//
// These codes change the color and the color remains
// the same until it is changed again. So be sure to
// 'reset' the color to 'normal' as appropriate.
//
// ANSI color codes:
//    Control
//  0 - normal
//  1 - bold (bright)
//  4 - underline
//  5 - blink
//  7 - reverse video
//  8 - invisible
//    Foreground colors
// 30 - black
// 31 - red
// 32 - green
// 33 - yellow
// 34 - blue
// 35 - magenta
// 36 - cyan
// 37 - white
//    Background colors
// 40 - black
// 41 - red
// 42 - green
// 43 - yellow
// 44 - blue
// 45 - magenta
// 46 - cyan
// 47 - white
//***********************************************************

// Define color codes
const NORMAL    = "\x1B[0;m"   // &N
const BBLACK    = "\x1B[1;30m" // &K
const BRED      = "\x1B[1;31m" // &R
const BGREEN    = "\x1B[1;32m" // &G
const BYELLOW   = "\x1B[1;33m" // &Y
const BBLUE     = "\x1B[1;34m" // &B
const BMAGENTA  = "\x1B[1;35m" // &M
const BCYAN     = "\x1B[1;36m" // &C
const BWHITE    = "\x1B[1;37m" // &W

// Color constants
const Normal        = NORMAL
const BrightBlack   = BBLACK
const BrightRed     = BRED
const BrightGreen   = BGREEN
const BrightYellow  = BYELLOW
const BrightBlue    = BBLUE
const BrightMagenta = BMAGENTA
const BrightCyan    = BCYAN
const BrightWhite   = BWHITE
