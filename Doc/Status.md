# OMugs – Online Multi-User Game Server
# Project Status

## Table of Contents

- [License](#license)
- [Introduction](#introduction)
- [Server Components Completed](#server-components-completed)
- [OMugs Source Code and Headers](#omugs-source-code-and-headers)
	- [Server cpp files](#server-cpp-files)
	- [Server header files](#server-header-files)
	- [Osi cpp files](#osi-cpp-files)
	- [Osi header files](#osi-header-files)
	- [Tools cpp files](#tools-cpp-files)
	- [Tools header files](#tools-header-files)
	- [WinApp cpp files](#winapp-cpp-files)
	- [WinApp header files](#winapp-header-files)
- [Summary](#summary)

# License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

For more information, please refer to <http://unlicense.org/>

Author  
Stephen L Bryant

Revision 1.0 April    22, 2001
Revision 1.1 December  4, 2002
Revision 1.2 December 22, 2025

Revision tracking after December 22, 2025 is maintained via git.

# Introduction

This document's purpose is to document the progress of the OMugs
project. Most of the components listed are also game commands,
exceptions are components like 'logon players', which refer to the code
that handles a player's from initial connection state to the 'ready to
play' state.

# Server Components Completed

<table style="width:31%;">
<colgroup>
<col style="width: 17%" />
<col style="width: 14%" />
</colgroup>
<thead>
<tr>
<th><h4 id="component">Component</h4></th>
<th><h4 id="date">Date</h4></th>
</tr>
</thead>
<tbody>
<tr>
<td>Advance</td>
<td>02/25/2003</td>
</tr>
<tr>
<td>Afk</td>
<td>12/19/2002</td>
</tr>
<tr>
<td>Armor Class</td>
<td>04/10/2003</td>
</tr>
<tr>
<td>Assist</td>
<td>11/27/2002</td>
</tr>
<tr>
<td>Buy</td>
<td>03/11/2003</td>
</tr>
<tr>
<td>Calendar</td>
<td>01/08/2004</td>
</tr>
<tr>
<td>Chat</td>
<td>07/11/2003</td>
</tr>
<tr>
<td>Color</td>
<td>05/01/2002</td>
</tr>
<tr>
<td>Consider</td>
<td>10/09/2003</td>
</tr>
<tr>
<td>Delete</td>
<td>03/14/2003</td>
</tr>
<tr>
<td>Destroy</td>
<td>07/03/2002</td>
</tr>
<tr>
<td>Drink</td>
<td>07/09/2003</td>
</tr>
<tr>
<td>Drop</td>
<td>07/05/2002</td>
</tr>
<tr>
<td>Eat</td>
<td>07/09/2003</td>
</tr>
<tr>
<td>Edit Mobiles</td>
<td>12/09/2003</td>
</tr>
<tr>
<td>Edit Objects</td>
<td>12/09/2003</td>
</tr>
<tr>
<td>Emote</td>
<td>06/26/2003</td>
</tr>
<tr>
<td>Equipment</td>
<td>08/28/2002</td>
</tr>
<tr>
<td>Examine</td>
<td>09/27/2002</td>
</tr>
<tr>
<td>Flee</td>
<td>11/15/2002</td>
</tr>
<tr>
<td>Follow</td>
<td>05/30/2002</td>
</tr>
<tr>
<td>Get</td>
<td>07/05/2002</td>
</tr>
<tr>
<td>Give</td>
<td>09/19/2002</td>
</tr>
<tr>
<td>Go</td>
<td>04/29/2002</td>
</tr>
<tr>
<td>GoTo</td>
<td>07/08/2002</td>
</tr>
<tr>
<td>GoToArrive</td>
<td>07/17/2002</td>
</tr>
<tr>
<td>GoToDepart</td>
<td>07/17/2002</td>
</tr>
<tr>
<td>Group</td>
<td>05/09/2002</td>
</tr>
<tr>
<td>Gsay</td>
<td>05/10/2002</td>
</tr>
<tr>
<td>Hail</td>
<td>06/02/2003</td>
</tr>
<tr>
<td>Help</td>
<td>04/08/2002</td>
</tr>
<tr>
<td>Inventory</td>
<td>07/03/2002</td>
</tr>
<tr>
<td>Invisible</td>
<td>10/14/2003</td>
</tr>
<tr>
<td>Kill</td>
<td>11/13/2002</td>
</tr>
<tr>
<td>List</td>
<td>03/11/2003</td>
</tr>
<tr>
<td>Load Mobile</td>
<td>10/11/2002</td>
</tr>
<tr>
<td>Load Object</td>
<td>07/03/2002</td>
</tr>
<tr>
<td>Logon</td>
<td>04/01/2002</td>
</tr>
<tr>
<td>Look</td>
<td>05/01/2002</td>
</tr>
<tr>
<td>Money</td>
<td>04/10/2002</td>
</tr>
<tr>
<td>Motd</td>
<td>03/27/2003</td>
</tr>
<tr>
<td>Password</td>
<td>03/03/2003</td>
</tr>
<tr>
<td>Played</td>
<td>12/06/2002</td>
</tr>
<tr>
<td>Player file</td>
<td>04/18/2002</td>
</tr>
<tr>
<td>Quit</td>
<td>04/01/2002</td>
</tr>
<tr>
<td>Refresh</td>
<td>06/25/2003</td>
</tr>
<tr>
<td>Remove</td>
<td>08/28/2002</td>
</tr>
<tr>
<td>Restore</td>
<td>12/10/2002</td>
</tr>
<tr>
<td>RoomInfo</td>
<td>07/10/2002</td>
</tr>
<tr>
<td>Save</td>
<td>07/17/2002</td>
</tr>
<tr>
<td>Say</td>
<td>05/03/2002</td>
</tr>
<tr>
<td>Sell</td>
<td>03/11/2003</td>
</tr>
<tr>
<td>Show</td>
<td>06/26/2003</td>
</tr>
<tr>
<td>Sit</td>
<td>05/03/2002</td>
</tr>
<tr>
<td>Skills</td>
<td>05/08/2003</td>
</tr>
<tr>
<td>Sleep</td>
<td>07/12/2002</td>
</tr>
<tr>
<td>Socials</td>
<td>04/05/2002</td>
</tr>
<tr>
<td>Spawn</td>
<td>12/19/2002</td>
</tr>
<tr>
<td>Stand</td>
<td>05/03/2002</td>
</tr>
<tr>
<td>Status</td>
<td>04/04/2002</td>
</tr>
<tr>
<td>Stop</td>
<td>04/01/2002</td>
</tr>
<tr>
<td>Tell</td>
<td>04/02/2002</td>
</tr>
<tr>
<td>Time</td>
<td>12/06/2002</td>
</tr>
<tr>
<td>Title</td>
<td>04/30/2002</td>
</tr>
<tr>
<td>Train</td>
<td>04/15/2003</td>
</tr>
<tr>
<td>Wake</td>
<td>07/12/2002</td>
</tr>
<tr>
<td>Wear</td>
<td>08/28/2002</td>
</tr>
<tr>
<td>Where Mobile</td>
<td>12/03/2002</td>
</tr>
<tr>
<td>Who</td>
<td>04/03/2002</td>
</tr>
<tr>
<td>Wield</td>
<td>08/30/2002</td>
</tr>
</tbody>
</table>

# OMugs Source Code and Headers

The source code file names, descriptions, and number of lines are listed
below and are divided into the following sections:

Server 🡪 The game code

Tools 🡪 Supporting code like validation

WinApp 🡪 Windows application code

Osi 🡪 OMugs Script Interpreter.

## Server cpp files

|                   |                                             |           |
|-------------------|---------------------------------------------|----------:|
| **File**          | **Description**                             | **Lines** |
| BigDog.cpp        | Starting point for all OMugs stuff          |       152 |
| Calendar.cpp      | Maintains the game calendar                 |       343 |
| Communication.cpp | Winsock tcp/ip telnet player communications |      7447 |
| Descriptor.cpp    | Maintains a linked list of connection data  |       152 |
| Dnode.cpp         | Connection data                             |        83 |
| Help.cpp          | Displays help file entries to player        |       157 |
| Log.cpp           | Logs messages to disk file                  |       107 |
| Mobile.cpp        | Manages Mobiles                             |      1607 |
| Object.cpp        | Manages Objects                             |      1592 |
| Player.cpp        | Manages Players                             |      1378 |
| Room.cpp          | Manages Rooms                               |       649 |
| Shop.cpp          | Manages Shops                               |       255 |
| Social.cpp        | Displays social messages to players         |       332 |
| Utility.cpp       | General purpose utility stuff               |       637 |
| Violence.cpp      | Fight related routines                      |       623 |
| World.cpp         | Functions to make the world come alive      |      1109 |
| **Total**         |                                             | **16623** |

## Server header files

|                 |                             |           |
|-----------------|-----------------------------|----------:|
| **File**        | **Description**             | **Lines** |
| BigDog.h        | The Big Dog Head            |        31 |
| Calendar.h      | Maintains the game calendar |        78 |
| Color.h         | Defines ascii color codes   |        86 |
| Communication.h | Defines Communication class |       151 |
| Config.h        | Configuration header        |       228 |
| Descriptior.h   | Defines Descriptor class    |        52 |
| Dnode.h         | Defines Dnode class         |        92 |
| Help.h          | Defines Help class          |        46 |
| Log.h           | Defines Log class           |        38 |
| Mobile.h        | Defines Mobile class        |        93 |
| Object.h        | Defines Object class        |        90 |
| Player.h        | Defines Player class        |       137 |
| Room.h          | Defines Room class          |        53 |
| Shop.h          | Defines Shop class          |        40 |
| Social.h        | Defines Social class        |        50 |
| Utility.h       | Defines Utility class       |        51 |
| Violence.h      | Defines Violence class      |        51 |
| World.h         | Defines World class         |        52 |
| **Total**       |                             |  **1419** |

## Osi cpp files

|              |                                               |           |
|--------------|-----------------------------------------------|----------:|
| **File**     | **Description**                               | **Lines** |
| Buffer.cpp   | I/O routines                                  |       125 |
| Executor.cpp | Execute statements                            |       489 |
| Icode.cpp    | Manages the intermediate code                 |       264 |
| Parser.cpp   | Analyze and translate script per syntax rules |       480 |
| RunStack.cpp | Runtime stack for expression evaluation       |        69 |
| Scanner.cpp  | Scans for tokens                              |       148 |
| Symbol.cpp   | Manages the symbol table                      |       171 |
| Token.cpp    | Extracts tokens                               |       307 |
| **Total**    |                                               |  **2053** |

## Osi header files

|            |                        |           |
|------------|------------------------|----------:|
| **File**   | **Description**        | **Lines** |
| Buffer.h   | Defines Buffer class   |        53 |
| Executor.h | Defines Executor class |        68 |
| Icode.h    | Defines Icode class    |        70 |
| Parser.h   | Defines Parser class   |        72 |
| RunStack.h | Defines RunStack class |        44 |
| Scanner.h  | Defines Scanner class  |        50 |
| Symbol.h   | Defines Symbol class   |        64 |
| Token.h    | Defines Token class    |        71 |
| **Total**  |                        |   **492** |

## Tools cpp files

|                   |                                       |           |
|-------------------|---------------------------------------|----------:|
| **File**          | **Description**                       | **Lines** |
| GenerateRooms.cpp | Room generation                       |       565 |
| LineCount.cpp     | Count lines in OMugs source code      |       184 |
| Validate.cpp      | Validate Rooms, Objects, Mobiles, etc |      1543 |
| WhoIsOnline.cpp   | Generate 'who' web content            |       230 |
| **Total**         |                                       |  **2522** |

## Tools header files

|                 |                            |           |
|-----------------|----------------------------|----------:|
| **File**        | **Description**            | **Lines** |
| GenerateRooms.h | Define GenerateRooms class |       104 |
| LineCount.h     | Define LineCount class     |        58 |
| Validate.h      | Define Validate class      |        50 |
| WhoIsOnline.h   | Define WhoIsOnline class   |        59 |
| **Total**       |                            |   **271** |

## WinApp cpp files

|                      |                     |           |
|----------------------|---------------------|----------:|
| **File**             | **Description**     | **Lines** |
| BuildMobiles.cpp     | Edit mobiles        |       796 |
| BuildMobiesList.cpp  | Edit mobiles list   |       184 |
| BuildObjects.cpp     | Edit objects        |       745 |
| BuildObjectsList.cpp | Edit objects list   |       507 |
| ChildFrm.cpp         | Windows application |        66 |
| MainFrm.cpp          | Windows application |       281 |
| OMugs.cpp            | Windows application |       196 |
| OMugsDoc.cpp         | Windows application |        78 |
| OMugsView.cpp        | Windows application |        99 |
| StdAfx.cpp           | Standard headers    |         8 |
| **Total**            |                     |  **2960** |

## WinApp header files

|                    |                               |           |
|--------------------|-------------------------------|----------:|
| **File**           | **Description**               | **Lines** |
| BuildMobiles.h     | Define BuildMobiles class     |       114 |
| BuildMobiesList.h  | Define BuildMobilesList class |        74 |
| BuildObjects.h     | Define BuildObjects class     |        98 |
| BuildObjectsList.h | Define BuildObjectsList class |        90 |
| ChildFrm.h         | Windows application           |        51 |
| MainFrm.h          | Windows application           |        85 |
| OMugs.h            | Windows application           |        56 |
| OMugsDoc.h         | Windows application           |        55 |
| OMugsView.h        | Windows application           |        63 |
| Resource.h         | Windows resources             |        95 |
| StdAfx.h           | Standard headers              |        28 |
| **Total**          |                               |   **809** |

# Summary

This summary includes all source code, documents, spreadsheets, world
files, etc. In other words, the whole OMugs directory. Only the 'debug'
directories are emptied before recording the sizes. File Count, Folder
Count, and Total Bytes are obtained by right-clicking on the OMugs
folder and selecting 'properties' and using the 'Size:' and 'Contains:'
information. These last three numbers might fluxuate significantly due
to testing new features. For example, the 10/27/2003 numbers for File
Count, Folder Count, and Total Bytes are inflated due to work on
scripting.

<table style="width:86%;">
<colgroup>
<col style="width: 14%" />
<col style="width: 9%" />
<col style="width: 10%" />
<col style="width: 9%" />
<col style="width: 10%" />
<col style="width: 9%" />
<col style="width: 9%" />
<col style="width: 12%" />
</colgroup>
<tbody>
<tr>
<td></td>
<td colspan="3" style="text-align: center;"><strong>Line
Counts</strong></td>
<td style="text-align: center;"><strong>WinZip</strong></td>
<td style="text-align: center;"><strong>File</strong></td>
<td style="text-align: center;"><strong>Folder</strong></td>
<td style="text-align: center;"><strong>Total</strong></td>
</tr>
<tr>
<td><strong>Date</strong></td>
<td style="text-align: right;"><strong>CPP</strong></td>
<td style="text-align: right;"><strong>Header</strong></td>
<td style="text-align: right;"><strong>Total</strong></td>
<td style="text-align: center;"><strong>Size</strong></td>
<td style="text-align: center;"><strong>Count</strong></td>
<td style="text-align: center;"><strong>Count</strong></td>
<td style="text-align: center;"><strong>Bytes</strong></td>
</tr>
<tr>
<td>05/09/2002</td>
<td style="text-align: right;">3,803</td>
<td style="text-align: right;">774</td>
<td style="text-align: right;">4,577</td>
<td style="text-align: right;">379KB</td>
<td style="text-align: center;">?</td>
<td style="text-align: center;">?</td>
<td style="text-align: center;">?</td>
</tr>
<tr>
<td>05/30/2002</td>
<td style="text-align: right;">4,194</td>
<td style="text-align: right;">763</td>
<td style="text-align: right;">4,957</td>
<td style="text-align: right;">389KB</td>
<td style="text-align: center;">?</td>
<td style="text-align: center;">?</td>
<td style="text-align: center;">?</td>
</tr>
<tr>
<td>07/03/2002</td>
<td style="text-align: right;">4,449</td>
<td style="text-align: right;">709</td>
<td style="text-align: right;">5,158</td>
<td style="text-align: right;">404KB</td>
<td style="text-align: center;">?</td>
<td style="text-align: center;">?</td>
<td style="text-align: center;">?</td>
</tr>
<tr>
<td>08/28/2002</td>
<td style="text-align: right;">6,154</td>
<td style="text-align: right;">734</td>
<td style="text-align: right;">6,888</td>
<td style="text-align: right;">434KB</td>
<td style="text-align: center;">122</td>
<td style="text-align: center;">23</td>
<td style="text-align: center;">1,167,908</td>
</tr>
<tr>
<td>10/11/2002</td>
<td style="text-align: right;">6,764</td>
<td style="text-align: right;">867</td>
<td style="text-align: right;">7,631</td>
<td style="text-align: right;">464KB</td>
<td style="text-align: center;">133</td>
<td style="text-align: center;">25</td>
<td style="text-align: center;">1,312,219</td>
</tr>
<tr>
<td>12/03/2002</td>
<td style="text-align: right;">8,585</td>
<td style="text-align: right;">949</td>
<td style="text-align: right;">9,534</td>
<td style="text-align: right;">504KB</td>
<td style="text-align: center;">143</td>
<td style="text-align: center;">32</td>
<td style="text-align: center;">1,469,118</td>
</tr>
<tr>
<td>03/31/2003</td>
<td style="text-align: right;">11,924</td>
<td style="text-align: right;">1,115</td>
<td style="text-align: right;">13,039</td>
<td style="text-align: right;">343KB</td>
<td style="text-align: center;">231</td>
<td style="text-align: center;">44</td>
<td style="text-align: center;">1,467,575</td>
</tr>
<tr>
<td>05/20/2003</td>
<td style="text-align: right;">14,328</td>
<td style="text-align: right;">1,190</td>
<td style="text-align: right;">15,518</td>
<td style="text-align: right;">634KB</td>
<td style="text-align: center;">548</td>
<td style="text-align: center;">50</td>
<td style="text-align: center;">3,140,338</td>
</tr>
<tr>
<td>07/01/2003</td>
<td style="text-align: right;">16,802</td>
<td style="text-align: right;">1,598</td>
<td style="text-align: right;">18,400</td>
<td style="text-align: right;">1.17MB</td>
<td style="text-align: center;">519</td>
<td style="text-align: center;">53</td>
<td style="text-align: center;">4,504,937</td>
</tr>
<tr>
<td>10/27/2003</td>
<td style="text-align: right;">18,317</td>
<td style="text-align: right;">1,643</td>
<td style="text-align: right;">19,960</td>
<td style="text-align: center;">2.07MB</td>
<td style="text-align: center;">1004</td>
<td style="text-align: center;">86</td>
<td style="text-align: center;">8,979,168</td>
</tr>
<tr>
<td>12/09/2003</td>
<td style="text-align: right;">21.876</td>
<td style="text-align: right;">2,826</td>
<td style="text-align: right;">24,702</td>
<td style="text-align: right;">1.43MB</td>
<td style="text-align: center;">777</td>
<td style="text-align: center;">63</td>
<td style="text-align: center;">5.260.925</td>
</tr>
<tr>
<td>12/23/2003</td>
<td style="text-align: right;">23,022</td>
<td style="text-align: right;">2,860</td>
<td style="text-align: right;">25,882</td>
<td style="text-align: right;">1.40MB</td>
<td style="text-align: center;">796</td>
<td style="text-align: center;">59</td>
<td style="text-align: center;">4,607,547</td>
</tr>
<tr>
<td>06/17/2004</td>
<td style="text-align: right;">24,158</td>
<td style="text-align: right;">2,991</td>
<td style="text-align: right;">27,149</td>
<td style="text-align: right;">1.95MB</td>
<td style="text-align: center;">997</td>
<td style="text-align: center;">60</td>
<td style="text-align: center;">5,447,962</td>
</tr>
<tr>
<td></td>
<td style="text-align: right;"></td>
<td style="text-align: right;"></td>
<td style="text-align: right;"></td>
<td style="text-align: right;"></td>
<td style="text-align: center;"></td>
<td style="text-align: center;"></td>
<td style="text-align: center;"></td>
</tr>
</tbody>
</table>
