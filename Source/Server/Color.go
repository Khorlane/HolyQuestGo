//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Color.go                                              *
// Usage:     Define ascii color codes                              *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

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
const NORMAL        string = "\x1B[0;m"   // &N
const BBLACK        string = "\x1B[1;30m" // &K
const BRED          string = "\x1B[1;31m" // &R
const BGREEN        string = "\x1B[1;32m" // &G
const BYELLOW       string = "\x1B[1;33m" // &Y
const BBLUE         string = "\x1B[1;34m" // &B
const BMAGENTA      string = "\x1B[1;35m" // &M
const BCYAN         string = "\x1B[1;36m" // &C
const BWHITE        string = "\x1B[1;37m" // &W

// Color constants
const Normal        string = NORMAL
const BrightBlack   string = BBLACK
const BrightRed     string = BRED
const BrightGreen   string = BGREEN
const BrightYellow  string = BYELLOW
const BrightBlue    string = BBLUE
const BrightMagenta string = BMAGENTA
const BrightCyan    string = BCYAN
const BrightWhite   string = BWHITE