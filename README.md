# HolyQuestGo

HolyQuest — powered by the Online Multi-User Game Server (OMugs) — is a MUD server originally written from scratch in C++. Development of OMugs began in 2002, and the first live release was launched on **March 7, 2003**. By the end of that year, the codebase was feature-complete and supported **61 player commands**.

This repository continues the work of the OMugs modernization project:  
https://github.com/Khorlane/OMugs

That project removed MFC (Microsoft Foundation Classes) from the original C++ code, but the codebase still relied on `Winsock2.h` for networking. Rather than continue maintaining platform-specific C++ code, this project explores a full conversion of OMugs into **Go**, with the goals of improving portability, reducing complexity, and preserving original game behavior.

---

## Conversion Plan

The conversion approach is intentionally simple and incremental:

- Each C++ module (`*.h`, `*.cpp`) is translated into a functionally equivalent Go module (`*.go`).
- The structure and flow of the original OMugs code is preserved where practical.
- All functionality is implemented using **free functions** instead of methods, following the design of the original C++ code.
- Visual Studio Code is used for side-by-side translation, comparison, and debugging.

The end goal is a fully working Go implementation of the OMugs engine, retaining the gameplay feel and command system of the original OMugs server.

---

## Build Instructions

From the project root:

```sh
del HolyQuestGo.exe
go build -o HolyQuestGo.exe
./HolyQuestGo