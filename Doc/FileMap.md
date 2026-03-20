# File Map

This file maps all OMugs project folders/files to all HolyQuestGo folders/files
Project paths:
- OMugs:       `C:\OMugs\`
- HolyQuestGo: `C:\Projects\HolyQuestGo\`

## Top Level OMugs to HolyQuestGo

| OMugs | HolyQuestGo | Note |
| ---   | ---         | ---  |
| .git                      | .git                       | git repo |
| .gitignore                | .gitignore                 | git ignore |
| .vs                       | .vscode                    | editor / workspace settings |
| Debug                     | n/a                        | Visual Studio build output |
| Doc                       | Doc                        | Documentation |
| HomeDir.txt               | HomeDir.txt                | Path to server home |
| Library                   | Library                    | Game content library |
| LICENSE                   | LICENSE                    | Project license |
| OMugs.slnx                | n/a                        | Visual Studio solution |
| OMugs.vcxproj             | n/a                        | Visual Studio C++ project |
| OMugs.vcxproj.filters     | n/a                        | Visual Studio project filters |
| OMugs.vcxproj.user        | n/a                        | Visual Studio user settings |
| README.md                 | README.md                  | Project readme |
| ReadMe.txt                | n/a                        | Legacy OMugs readme |
| Res                       | n/a                        | Windows resources |
| Running                   | Running                    | Server execution state |
| Source                    | Source                     | Server codebase |
| Tmp                       | n/a                        | Temporary / scratch files |
| Utility                   | Utility                    | Utility scripts |
| WebSite                   | WebSite                    | HolyQuest website files |
| zNotes.txt                | n/a                        | Developer notes |
| n/a                       | Build.ps1                  | PowerShell build script |
| n/a                       | chat_context.txt           | Chat session continuity notes |
| n/a                       | go.mod                     | Go module definition |
| n/a                       | go.sum                     | Go module checksums |
| n/a                       | HolyQuestGo.bin            | Built binary artifact |
| n/a                       | HolyQuestGo.code-workspace | VS Code workspace file |
| n/a                       | HolyQuestGo.exe            | Built Windows executable |
| n/a                       | main.go                    | Go program entry point |

## Source folder
| OMugs  | HolyQuestGo | Note |
| ---    | ---         | ---  |
| Osi    | n/a         | OMugs Script Interpreter |
| Server | Server      | Core server code |
| Tools  | Server      | Partial - Validate only |
| WinApp | n/a         | Window application code |

## Server folder
| OMugs | HolyQuestGo | Note |
| ---   | ---         | ---  |
| BigDog.cpp / BigDog.h               | BigDog.go        | Main server loop |
| Calendar.cpp / Calendar.h           | Calendar.go      | Game calendar |
| Color.h                             | Color.go         | Color codes |
| Communication.cpp / Communication.h | Communication.go | Player I/O |
| Config.h                            | Config.go        | Constants and globals |
| Descriptor.cpp / Descriptor.h       | Descriptor.go    | Login descriptors |
| Dnode.cpp / Dnode.h                 | Dnode.go         | Connection nodes |
| Help.cpp / Help.h                   | Help.go          | Help text |
| Log.cpp / Log.h                     | Log.go           | Logging |
| Mobile.cpp / Mobile.h               | Mobile.go        | Mobiles |
| Object.cpp / Object.h               | Object.go        | Objects |
| Player.cpp / Player.h               | Player.go        | Players |
| Room.cpp / Room.h                   | Room.go          | Rooms |
| Shop.cpp / Shop.h                   | Shop.go          | Shops |
| Social.cpp / Social.h               | Social.go        | Social commands |
| Utility.cpp / Utility.h             | Utility.go       | Shared helpers |
| Violence.cpp / Violence.h           | Violence.go      | Combat |
| World.cpp / World.h                 | World.go         | World state |
| n/a                                 | Validate.go      | Validation - from Tools |


## Tools folder
| OMugs | HolyQuestGo | Note |
| ---   | ---         | ---  |
| GenerateRooms.cpp / GenerateRooms.h | n/a         | zMapper room translation |
| LineCount.cpp / LineCount.h         | n/a         | Count source lines |
| WhoIsOnline.cpp / WhoIsOnline.h     | n/a         | Create statswho.xml |
| Validate.cpp / Validate.h           | Validate.go | Validation - moved to Server |
