# File Map

This file maps the relevant OMugs `Source` files to the file organization used by this project.

For the purposes of this project:
- Ignore `C:\OMugs\Source\Osi`
- Ignore `C:\OMugs\Source\WinApp`
- Treat `C:\OMugs\Source\Tools\Validate.cpp` / `Validate.h` as the counterpart for this project's `Source\Server\Validate.go`

## OMugs to HolyQuestGo

| OMugs source | HolyQuestGo source |
| --- | --- |
| `C:\OMugs\Source\Server\BigDog.cpp` / `BigDog.h` | `C:\Projects\HolyQuestGo\Source\Server\BigDog.go` |
| `C:\OMugs\Source\Server\Calendar.cpp` / `Calendar.h` | `C:\Projects\HolyQuestGo\Source\Server\Calendar.go` |
| `C:\OMugs\Source\Server\Color.h` | `C:\Projects\HolyQuestGo\Source\Server\Color.go` |
| `C:\OMugs\Source\Server\Communication.cpp` / `Communication.h` | `C:\Projects\HolyQuestGo\Source\Server\Communication.go` |
| `C:\OMugs\Source\Server\Config.h` | `C:\Projects\HolyQuestGo\Source\Server\Config.go` |
| `C:\OMugs\Source\Server\Descriptor.cpp` / `Descriptor.h` | `C:\Projects\HolyQuestGo\Source\Server\Descriptor.go` |
| `C:\OMugs\Source\Server\Dnode.cpp` / `Dnode.h` | `C:\Projects\HolyQuestGo\Source\Server\Dnode.go` |
| `C:\OMugs\Source\Server\Help.cpp` / `Help.h` | `C:\Projects\HolyQuestGo\Source\Server\Help.go` |
| `C:\OMugs\Source\Server\Log.cpp` / `Log.h` | `C:\Projects\HolyQuestGo\Source\Server\Log.go` |
| `C:\OMugs\Source\Server\Mobile.cpp` / `Mobile.h` | `C:\Projects\HolyQuestGo\Source\Server\Mobile.go` |
| `C:\OMugs\Source\Server\Object.cpp` / `Object.h` | `C:\Projects\HolyQuestGo\Source\Server\Object.go` |
| `C:\OMugs\Source\Server\Player.cpp` / `Player.h` | `C:\Projects\HolyQuestGo\Source\Server\Player.go` |
| `C:\OMugs\Source\Server\Room.cpp` / `Room.h` | `C:\Projects\HolyQuestGo\Source\Server\Room.go` |
| `C:\OMugs\Source\Server\Shop.cpp` / `Shop.h` | `C:\Projects\HolyQuestGo\Source\Server\Shop.go` |
| `C:\OMugs\Source\Server\Social.cpp` / `Social.h` | `C:\Projects\HolyQuestGo\Source\Server\Social.go` |
| `C:\OMugs\Source\Server\Utility.cpp` / `Utility.h` | `C:\Projects\HolyQuestGo\Source\Server\Utility.go` |
| `C:\OMugs\Source\Tools\Validate.cpp` / `Validate.h` | `C:\Projects\HolyQuestGo\Source\Server\Validate.go` |
| `C:\OMugs\Source\Server\Violence.cpp` / `Violence.h` | `C:\Projects\HolyQuestGo\Source\Server\Violence.go` |
| `C:\OMugs\Source\Server\World.cpp` / `World.h` | `C:\Projects\HolyQuestGo\Source\Server\World.go` |

## Notes

- Most OMugs modules that were split across `.cpp` and `.h` are represented here as a single `.go` file.
- The main intentional location mismatch is `Validate`, which moved from `OMugs\Source\Tools` into this project's `Source\Server`.
