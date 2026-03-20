//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Validate.go                                           *
// Usage:     Validates game data and logs errors                   *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
	"bufio"
	"os"
	"path/filepath"
)

// Log Validation Error
func LogValErr(Message string, FileName string) {
  if StrGetLength(Message) > 50 {
    // Message is too long, chop it off
    Message = StrLeft(Message, 50)
  }
  Message = " - " + Message + " "
  for StrGetLength(Message) < 55 {
    Message += "-"
  }
  Message += "> "
  LogBuf  = "ValErr"
  LogBuf += Message
  LogBuf += FileName
  LogIt(LogBuf)
  ValErr = true
}

// ValidateIt
func ValidateIt(ValidationType string) bool {
  ValErr = false
  ValidationType = StrMakeLower(ValidationType)
  if ValidationType == "all" {
    ValidateAll()
  } else {
    if ValidationType == "mobiles" {
      ValidateLibraryMobiles()
      ValidateLibraryWorldMobiles()
    } else {
      if ValidationType == "objects" {
        ValidateLibraryObjects()
        ValidateLibraryLoot()
        ValidateLibraryShops()
      } else {
        if ValidationType == "rooms" {
          ValidateLibraryRooms()
          ValidateLibraryShops()
          ValidateLibraryWorldMobiles()
        }
      }
    }
  }
  if ValErr {
    LogBuf = "ValErr - Validation failed!!"
    LogIt(LogBuf)
  } else {
    LogBuf = "Validation successful!!"
    LogIt(LogBuf)
  }
  return ValErr
}

// Validate all
func ValidateAll() {
	ValidateLibraryLoot()
	ValidateLibraryMobiles()
	ValidateLibraryObjects()
	ValidateLibraryRooms()
	ValidateLibraryShops()
	ValidateLibraryWorldMobiles()
	ValidateRunningPlayers()
	ValidateRunningPlayersPlayerEqu()
	ValidateRunningPlayersPlayerObj()
	ValidateRunningRoomMob()
	ValidateRunningRoomObj()
}

// Validate LibraryLoot
func ValidateLibraryLoot() {
  var FileName          string
  var LineCount         int
  var Message           string
  var ObjectId          string
  var ObjectIdFileName  string
  var LootFile         *os.File
  var LootFileName      string

  LogBuf = "Begin validation LibraryLoot"
  LogIt(LogBuf)
  // Get list of all LibraryLoot files
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryLoot - Change directory to HomeDir failed")
    os.Exit(1)
  }
  if ChgDir(LOOT_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryLoot - Change directory to LOOT_DIR failed")
    os.Exit(1)
  }
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateLibraryLoot - ReadDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open loot file
    LootFileName = DirEntry.Name()
    LootFile, err = os.Open(LootFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateLibraryLoot - Open loot file failed: " + LootFileName)
      os.Exit(1)
    }
    LineCount = 0
    Scanner := bufio.NewScanner(LootFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = ""
    }
    for Stuff != "" {
      // For all lines
      LineCount++
      ObjectId = StrGetWord(Stuff, 3)
      //************
      //* ObjectId *
      //************
      ObjectIdFileName = filepath.Join(HomeDir, OBJECTS_DIR, ObjectId+".txt")
      if !FileExist(ObjectIdFileName) {
        // ObjectId file not found
        Message  = "Object file"
        Message += " '"
        Message += ObjectId
        Message += "' "
        Message += "not found"
        FileName = LootFileName
        LogValErr(Message, FileName)
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        Stuff = ""
      }
    }
    LootFile.Close()
  }
  LogBuf = "Done  validating LibraryLoot"
  LogIt(LogBuf)
}

// Validate LibraryMobiles
func ValidateLibraryMobiles() {
  var FieldName       string
  var FieldValue      string
  var FileName        string
  var LineCount       int
  var LootFileName    string
  var Message         string
  var MobileFile     *os.File
  var MobileFileName  string
  var MobileId        string

  LogBuf = "Begin validation LibraryMobiles"
  // Get list of all LibraryMobiles files
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryMobiles - Change directory to HomeDir failed")
    os.Exit(1)
  }
  LogIt(LogBuf)
  if ChgDir(MOBILES_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryMobiles - Change directory to MOBILES_DIR failed")
    os.Exit(1)
  }
  // Get list of all LibraryMobiles files
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateLibraryMobiles - ReadDir failed")
    os.Exit(1)
  }
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryMobiles - Change directory to HomeDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open mobile file
    MobileFileName  = DirEntry.Name()
    MobileId        = StrLeft(MobileFileName, StrGetLength(MobileFileName)-4)
    MobileFileName  = MOBILES_DIR + MobileFileName
    MobileFile, err = os.Open(MobileFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateLibraryMobiles - Open mobile file failed")
      os.Exit(1)
    }
    LineCount = 0
    Scanner := bufio.NewScanner(MobileFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = ""
    }
    for Stuff != "" {
      // For all lines
      LineCount++
      FieldName  = StrGetWord(Stuff, 1)
      FieldValue = StrGetWord(Stuff, 2)
      //********************************
      //* MobileId field must be first *
      //********************************
      if LineCount == 1 {
        // First line
        if FieldName != "MobileId:" {
          // Must be MobileId
          Message = "MobileId field must be the first field"
          FileName = MobileFileName
          LogValErr(Message, FileName)
        }
      }
      //*************
      //* MobileId: *
      //*************
      if FieldName == "MobileId:" {
        // MobileId field validation
        if MobileId != FieldValue {
          // MobileId must match file name
          Message = "MobileId must match file name"
          FileName = MobileFileName
          LogValErr(Message, FileName)
        }
      }
      //********
      //* Sex: *
      //********
      if FieldName == "Sex:" {
        // Sex field validation
        if StrIsNotWord(FieldValue, "F M N") {
          // Invalid mobile sex
          Message = "Mobile sex is invalid"
          FileName = MobileFileName
          LogValErr(Message, FileName)
        }
      }
      //***********
      //* Action: *
      //***********
      if FieldName == "Action:" {
        // Sex field validation
        if StrIsNotWord(FieldValue, "None Aggro Faction Destroy Help NoMove Wimpy") {
          // Invalid mobile action
          Message = "Mobile action is invalid"
          FileName = MobileFileName
          LogValErr(Message, FileName)
        }
      }
      //************
      //* Faction: *
      //************
      if FieldName == "Faction:" {
        // Faction field validation
        if StrIsNotWord(FieldValue, "Evil Lawless Neutral Lawful Good") {
          // Invalid mobile faction
          Message  = "Mobile faction is invalid"
          FileName = MobileFileName
          LogValErr(Message, FileName)
        }
      }
      //***********
      //* Attack: *
      //***********
      if FieldName == "Attack:" {
        // Faction field validation
        FieldValue = StrMakeLower(FieldValue)
        if StrIsNotWord(FieldValue, "bites claws crushes hits mauls pierces punches slashes stabs stings thrashes") {
          // Invalid mobile attack
          Message  = "Mobile attack is invalid"
          FileName = MobileFileName
          LogValErr(Message, FileName)
        }
      }
      //*********
      //* Loot: *
      //*********
      if FieldName == "Loot:" {
        // Loot field validation
        LootFileName  = LOOT_DIR
        LootFileName += FieldValue
        LootFileName += ".txt"
        if !FileExist(LootFileName) {
          // Loot file not found
          Message = "Loot file not found"
          FileName = MobileFileName
          LogValErr(Message, FileName)
        }
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        Stuff = ""
      }
    }
    MobileFile.Close()
  }
  LogBuf = "Done  validating LibraryMobiles"
  LogIt(LogBuf)
}

// Validate LibraryObjects
func ValidateLibraryObjects() {
  var FieldName       string
  var FieldValue      string
  var FileName        string
  var LineCount       int
  var Message         string
  var ObjectFile     *os.File
  var ObjectFileName  string
  var ObjectId        string

  LogBuf = "Begin validation LibraryObjects"
  LogIt(LogBuf)
  if ChgDir(OBJECTS_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryObjects - Change directory to OBJECTS_DIR failed")
    os.Exit(1)
  }
  // Get list of all LibarryObjects files
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateLibraryObjects - ReadDir failed")
    os.Exit(1)
  }
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryObjects - Change directory to HomeDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open object file
    ObjectFileName  = DirEntry.Name()
    ObjectId        = StrLeft(ObjectFileName, StrGetLength(ObjectFileName)-4)
    ObjectFileName  = OBJECTS_DIR + ObjectFileName
    ObjectFile, err = os.Open(ObjectFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateLibraryObjects - Open object file failed")
      os.Exit(1)
    }
    LineCount = 0
    Scanner := bufio.NewScanner(ObjectFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = ""
    }
    for Stuff != "" {
      // For all lines
      LineCount++
      FieldName  = StrGetWord(Stuff, 1)
      FieldValue = StrGetWord(Stuff, 2)
      //********************************
      //* ObjectId field must be first *
      //********************************
      if LineCount == 1 {
        // First line
        if FieldName != "ObjectId:" {
          // Must be ObjectId
          Message  = "ObjectId field must be the first field"
          FileName = ObjectFileName
          LogValErr(Message, FileName)
        }
      }
      //*************
      //* ObjectId: *
      //*************
      if FieldName == "ObjectId:" {
        // ObjectId field validation
        if ObjectId != FieldValue {
          // ObjectId must match file name
          Message  = "ObjectId must match file name"
          FileName = ObjectFileName
          LogValErr(Message, FileName)
        }
      }
      //*********
      //* Type: *
      //*********
      if FieldName == "Type:" {
        // Type field validation
        if StrIsNotWord(FieldValue, "Armor Container Drink Food Junk Key Light NoTake Treasure Weapon") {
          // Invalid object type
          Message  = "Object type is invalid"
          FileName = ObjectFileName
          LogValErr(Message, FileName)
        }
        //***************
        //* Type: Armor *
        //***************
        if FieldValue == "armor" {
          //***************
          //* ArmorValue: *
          //***************
          if Scanner.Scan() {
            Stuff = Scanner.Text()
          } else {
            Stuff = ""
          }
          LineCount++
          FieldName = StrGetWord(Stuff, 1)
          FieldValue = StrGetWord(Stuff, 2)
          if FieldName != "ArmorValue:" {
            // ArmorValue must follow 'Type: Armor' specification
            Message  = "ArmorValue must follow 'Type: Armor' specification"
            FileName = ObjectFileName
            LogValErr(Message, FileName)
          }
          //**************
          //* ArmorWear: *
          //**************
          if Scanner.Scan() {
            Stuff = Scanner.Text()
          } else {
            Stuff = ""
          }
          LineCount++
          FieldName = StrGetWord(Stuff, 1)
          FieldValue = StrGetWord(Stuff, 2)
          if FieldName != "ArmorWear:" {
            // ArmorWear must follow 'ArmorValue' specification
            Message  = "ArmorWear must follow 'ArmorValue' specification"
            FileName = ObjectFileName
            LogValErr(Message, FileName)
          } else {
            // Validate 'wear' positions
            if StrIsNotWord(FieldValue, "Head Ear Neck Shoulders Chest Back Arms Wrist Hands Finger Shield Waist Legs Ankle Feet") {
              // Invalid wear position
              Message = "Wear position is invalid"
              FileName = ObjectFileName
              LogValErr(Message, FileName)
            }
          }
        }
        //****************
        //* Type: Weapon *
        //****************
        if FieldValue == "weapon" {
          //***************
          //* WeaponType: *
          //***************
          if Scanner.Scan() {
            Stuff = Scanner.Text()
          } else {
            Stuff = ""
          }
          LineCount++
          FieldName  = StrGetWord(Stuff, 1)
          FieldValue = StrGetWord(Stuff, 2)
          if FieldName != "WeaponType:" {
            // WeaponType must follow 'Type: Weapon' specification
            Message  = "WeaponType must follow 'Type: Weapon' specification"
            FileName = ObjectFileName
            LogValErr(Message, FileName)
          } else {
            // Validate WeaponType
            if StrIsNotWord(FieldValue, "Axe Club Dagger Hammer Spear Staff Sword") {
              // Invalid weapon type
              Message  = "Weapon type is invalid"
              FileName = ObjectFileName
              LogValErr(Message, FileName)
            }
          }
          //*****************
          //* WeaponDamage: *
          //*****************
          if Scanner.Scan() {
            Stuff = Scanner.Text()
          } else {
            Stuff = ""
          }
          LineCount++
          FieldName  = StrGetWord(Stuff, 1)
          FieldValue = StrGetWord(Stuff, 2)
          if FieldName != "WeaponDamage:" {
            // WeaponDamage must follow WeaponType specification
            Message  = "WeaponDamage must follow WeaponType specification"
            FileName = ObjectFileName
            LogValErr(Message, FileName)
          }
        }
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        Stuff = ""
      }
    }
    ObjectFile.Close()
  }
  LogBuf = "Done  validating LibraryObjects"
  LogIt(LogBuf)
}

// Validate LibraryRooms
func ValidateLibraryRooms() {
  var ExitToRoomIdFileName  string
  var FieldName             string
  var FieldValue            string
  var FileName              string
  var i                     int
  var j                     int
  var LineCount             int
  var Message               string
  var RoomFile             *os.File
  var RoomFileName          string
  var RoomId                string
  var RoomTypeError         bool

  LogBuf = "Begin validation LibraryRooms"
  LogIt(LogBuf)
  if ChgDir(ROOMS_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryRooms - Change directory to ROOMS_DIR failed")
    os.Exit(1)
  }
  // Get list of all LibraryRooms files
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateLibraryRooms - ReadDir failed")
    os.Exit(1)
  }
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRooms - Change directory to HomeDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open room file
    RoomFileName  = DirEntry.Name()
    RoomId        = StrLeft(RoomFileName, StrGetLength(RoomFileName)-4)
    RoomFileName  = ROOMS_DIR + RoomFileName
    RoomFile, err = os.Open(RoomFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateLibraryRooms - Open room file failed")
      os.Exit(1)
    }
    LineCount = 0
    Scanner := bufio.NewScanner(RoomFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = "End of Room"
    }
    for Stuff != "End of Room" {
      // For all lines
      LineCount++
      FieldName = StrGetWord(Stuff, 1)
      FieldValue = StrGetWords(Stuff, 2)
      //********************************
      //* RoomId field must be first *
      //********************************
      if LineCount == 1 {
        // First line
        if FieldName != "RoomId:" {
          // Must be RoomId
          Message  = "RoomId field must be the first field"
          FileName = RoomFileName
          LogValErr(Message, FileName)
        }
      }
      //***********
      //* RoomId: *
      //***********
      if FieldName == "RoomId:" {
        // RoomId field validation
        if RoomId != FieldValue {
          // RoomId must match file name
          Message  = "RoomId must match file name"
          FileName = RoomFileName
          LogValErr(Message, FileName)
        }
      }
      //*************
      //* RoomType: *
      //*************
      RoomTypeError = false
      if FieldName == "RoomType:" {
        // RoomType validation
        j = StrCountWords(FieldValue)
        for i = 1; i <= j; i++ {
          // Check each word in FieldValue
          TmpStr = StrGetWord(FieldValue, i)
          if StrIsWord(TmpStr, "None Dark Drink NoFight NoNPC") {
            // Valid RoomType
            if TmpStr != "None" {
              // Exclude RoomType of 'None'
            }
          } else {
            // Invalid RoomType
            Message  = "RoomType has an invalid entry"
            FileName = RoomFileName
            LogValErr(Message, FileName)
            RoomTypeError = true
            break
          }
        }
      }
      //************
      //* Terrain: *
      //************
      if FieldName == "Terrain:" {
        // Terrain validation
        if StrIsWord(FieldValue, "Inside Street Road Field Forest Swamp Desert Hill Mountain") {
          // Valid Terrain
        } else {
          // Invalid Terrain
          Message  = "Terrain is invalid"
          FileName = RoomFileName
          LogValErr(Message, FileName)
        }
      }
      //*****************
      //* ExitToRoomId: *
      //*****************
      if FieldName == "ExitToRoomId:" {
        // ExitToRoomId field validation
        ExitToRoomIdFileName = ROOMS_DIR
        ExitToRoomIdFileName += FieldValue
        ExitToRoomIdFileName += ".txt"
        if !FileExist(ExitToRoomIdFileName) {
          // ExitToRoom file not found
          Message  = "Room file"
          Message += " '"
          Message += FieldValue
          Message += "' "
          Message += "not found"
          FileName = RoomFileName
          LogValErr(Message, FileName)
        }
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        // End of file reached
        Stuff = "End of Room"
      }
      _ = RoomTypeError
    }
    RoomFile.Close()
  }
  LogBuf = "Done  validating LibraryRooms"
  LogIt(LogBuf)
}

// Validate LibraryShops
func ValidateLibraryShops() {
  var FieldName         string
  var FieldValue        string
  var FileName          string
  var LineCount         int
  var Message           string
  var ObjectId          string
  var ObjectIdFileName  string
  var ShopFile         *os.File
  var ShopFileName      string
  var PlayerName        string

  LogBuf = "Begin validation LibraryShops"
  LogIt(LogBuf)
  if ChgDir(SHOPS_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryShops - Change directory to SHOPS_DIR failed")
    os.Exit(1)
  }
  // Get list of all LibraryShops files
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateLibraryShops - ReadDir failed")
    os.Exit(1)
  }
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryShops - Change directory to HomeDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open player file
    ShopFileName  = DirEntry.Name()
    PlayerName    = StrLeft(ShopFileName, StrGetLength(ShopFileName)-4)
    _ = PlayerName
    ShopFileName  = SHOPS_DIR + ShopFileName
    ShopFile, err = os.Open(ShopFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateLibraryShops - Open shop file failed")
      os.Exit(1)
    }
    LineCount = 0
    Scanner := bufio.NewScanner(ShopFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = "End of Shop"
    }
    for Stuff != "End of Shop" {
      // For all lines
      LineCount++
      FieldName  = StrGetWord(Stuff, 1)
      FieldValue = StrGetWord(Stuff, 2)
      if FieldName != "Item:" {
        // Not an item line
        if Scanner.Scan() {
          Stuff = Scanner.Text()
        } else {
          // ObjectId file not found
          if Stuff != "End of Shop" {
            // 'End of Shop' must be last line
            Message = "'End of Shop' must be the last line"
            FileName = ShopFileName
            LogValErr(Message, FileName)
            Stuff = "End of Shop"
          }
        }
        continue
      }
      //********
      //* Item *
      //********
      ObjectId = FieldValue
      ObjectIdFileName = OBJECTS_DIR
      ObjectIdFileName += ObjectId
      ObjectIdFileName += ".txt"
      if !FileExist(ObjectIdFileName) {
        // ObjectId file not found
        Message = "Object file"
        Message += " '"
        Message += ObjectId
        Message += "' "
        Message += "not found"
        FileName = ShopFileName
        LogValErr(Message, FileName)
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        Stuff = "End of Shop"
      }
    }
    ShopFile.Close()
  }
  LogBuf = "Done  validating LibraryShops"
  LogIt(LogBuf)
}

// Validate LibraryWorldMobiles
func ValidateLibraryWorldMobiles() {
  var FieldName            string
  var FieldValue           string
  var FileName             string
  var LineCount            int
  var Message              string
  var MobileId             string
  var MobileIdFileName     string
  var RoomIdFileName       string
  var WorldMobileFile     *os.File
  var WorldMobileFileName  string
  var WorldMobileName      string

  LogBuf = "Begin validation LibraryWorldMobiles"
  LogIt(LogBuf)
  if ChgDir(WORLD_MOBILES_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryWorldMobiles - Change directory to WORLD_MOBILES_DIR failed")
    os.Exit(1)
  }
  // Get list of all LibraryWorldMobiles files
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateLibraryWorldMobiles - ReadDir failed")
    os.Exit(1)
  }
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateLibraryWorldMobiles - Change directory to HomeDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open world mobile file
    WorldMobileFileName = DirEntry.Name()
    if WorldMobileFileName == "ReadMe.txt" {
      // Skip ReadMe files
      continue
    }
    WorldMobileName      = StrLeft(WorldMobileFileName, StrGetLength(WorldMobileFileName)-4)
    MobileId             = WorldMobileName
    WorldMobileFileName  = WORLD_MOBILES_DIR + WorldMobileFileName
    WorldMobileFile, err = os.Open(WorldMobileFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateLibraryWorldMobiles - Open world mobile file failed")
      os.Exit(1)
    }
    //************
    //* MobileId *
    //************
    MobileIdFileName  = MOBILES_DIR
    MobileIdFileName += MobileId
    MobileIdFileName += ".txt"
    if !FileExist(MobileIdFileName) {
      // MobileId file not found
      Message  = "Mobile file"
      Message += " '"
      Message += MobileId
      Message += "' "
      Message += "not found"
      FileName = WorldMobileFileName
      LogValErr(Message, FileName)
    }
    //***********************
    //* Check file contents *
    //***********************
    LineCount = 0
    Scanner := bufio.NewScanner(WorldMobileFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = ""
    }
    for Stuff != "" {
      // For all lines
      LineCount++
      FieldName  = StrGetWord(Stuff, 1)
      FieldValue = StrGetWord(Stuff, 2)
      //**********
      //* RoomId *
      //**********
      if FieldName == "RoomId:" {
        // RoomId field validation
        RoomIdFileName  = ROOMS_DIR
        RoomIdFileName += FieldValue
        RoomIdFileName += ".txt"
        if !FileExist(RoomIdFileName) {
          // RoomId file not found
          Message = "Room file"
          Message += " '"
          Message += FieldValue
          Message += "' "
          Message += "not found"
          FileName = WorldMobileFileName
          LogValErr(Message, FileName)
        }
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        Stuff = ""
      }
    }
    WorldMobileFile.Close()
  }
  LogBuf = "Done  validating LibraryWorldMobiles"
  LogIt(LogBuf)
}

// Validate RunningPlayers
func ValidateRunningPlayers() {
  var FieldName       string
  var FieldValue      string
  var FileName        string
  var LineCount       int
  var Message         string
  var PlayerFile     *os.File
  var PlayerFileName  string
  var PlayerName      string
  var RoomIdFileName  string

  LogBuf = "Begin validation RunningPlayers"
  LogIt(LogBuf)
  if ChgDir(PLAYER_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRunningPlayers - Change directory to PLAYER_DIR failed")
    os.Exit(1)
  }
  // Get list of all RunningPlayers files
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateRunningPlayers - ReadDir failed")
    os.Exit(1)
  }
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRunningPlayers - Change directory to HomeDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open player file
    PlayerFileName  = DirEntry.Name()
    PlayerName      = StrLeft(PlayerFileName, StrGetLength(PlayerFileName)-4)
    PlayerFileName  = PLAYER_DIR + PlayerFileName
    PlayerFile, err = os.Open(PlayerFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateRunningPlayers - Open player file failed")
      os.Exit(1)
    }
    LineCount = 0
    Scanner := bufio.NewScanner(PlayerFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = ""
    }
    for Stuff != "" {
      // For all lines
      LineCount++
      StrReplace(&Stuff, ":", " ")
      FieldName  = StrGetWord(Stuff, 1)
      FieldValue = StrGetWord(Stuff, 2)
      //****************************
      //* Name field must be first *
      //****************************
      if LineCount == 1 {
        // First line
        if FieldName != "Name" {
          // Must be Name
          Message  = "Name field must be the first field"
          FileName = PlayerFileName
          LogValErr(Message, FileName)
        }
      }
      //********
      //* Name *
      //********
      if FieldName == "Name" {
        // Name field validation
        if PlayerName != FieldValue {
          // Name must match file name
          Message  = "Name must match file name"
          FileName = PlayerFileName
          LogValErr(Message, FileName)
        }
      }
      //**********
      //* RoomId *
      //**********
      if FieldName == "RoomId" {
        // RoomId field validation
        RoomIdFileName  = ROOMS_DIR
        RoomIdFileName += FieldValue
        RoomIdFileName += ".txt"
        if !FileExist(RoomIdFileName) {
          // RoomId file not found
          Message  = "Room file"
          Message += " '"
          Message += FieldValue
          Message += "' "
          Message += "not found"
          FileName = PlayerFileName
          LogValErr(Message, FileName)
        }
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        Stuff = ""
      }
    }
    PlayerFile.Close()
  }
  LogBuf = "Done  validating RunningPlayers"
  LogIt(LogBuf)
}

// Validate RunningPlayersPlayerEqu
func ValidateRunningPlayersPlayerEqu() {
  var FileName           string
  var LineCount          int
  var Message            string
  var ObjectId           string
  var ObjectIdFileName   string
  var PlayerEquFile     *os.File
  var PlayerEquFileName  string
  var PlayerName         string
  var WearPosition       string

  LogBuf = "Begin validation RunningPlayersPlayerEqu"
  LogIt(LogBuf)
  if ChgDir(PLAYER_EQU_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRunningPlayersPlayerEqu - Change directory to PLAYER_EQU_DIR failed")
    os.Exit(1)
  }
  // Get list of all RunningPlayersPlayerEqu files
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateRunningPlayersPlayerEqu - ReadDir failed")
    os.Exit(1)
  }
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRunningPlayersPlayerEqu - Change directory to HomeDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open player file
    PlayerEquFileName  = DirEntry.Name()
    PlayerName         = StrLeft(PlayerEquFileName, StrGetLength(PlayerEquFileName)-4)
    _ = PlayerName
    PlayerEquFileName  = PLAYER_EQU_DIR + PlayerEquFileName
    PlayerEquFile, err = os.Open(PlayerEquFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateRunningPlayersPlayerEqu - Open player file failed")
      os.Exit(1)
    }
    LineCount = 0
    Scanner := bufio.NewScanner(PlayerEquFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = ""
    }
    for Stuff != "" {
      // For all lines
      LineCount++
      WearPosition = StrGetWord(Stuff, 1)
      ObjectId     = StrGetWord(Stuff, 2)
      //*****************
      //* Wear position *
      //*****************
      if WearPosition < "01" {
        // Wear position must be 01 - 20
        Message  = "Wear position < 01"
        FileName = PlayerEquFileName
        LogValErr(Message, FileName)
      }
      if WearPosition > "20" {
        // Wear position must be 01 - 20
        Message  = "Wear position > 20"
        FileName = PlayerEquFileName
        LogValErr(Message, FileName)
      }
      //************
      //* ObjectId *
      //************
      ObjectIdFileName  = OBJECTS_DIR
      ObjectIdFileName += ObjectId
      ObjectIdFileName += ".txt"
      if !FileExist(ObjectIdFileName) {
        // ObjectId file not found
        Message  = "Object file"
        Message += " '"
        Message += ObjectId
        Message += "' "
        Message += "not found"
        FileName = PlayerEquFileName
        LogValErr(Message, FileName)
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        Stuff = ""
      }
    }
    PlayerEquFile.Close()
  }
  LogBuf = "Done  validating RunningPlayersPlayerEqu"
  LogIt(LogBuf)
}

// Validate RunningPlayersPlayerObj
func ValidateRunningPlayersPlayerObj() {
  var FileName           string
  var LineCount          int
  var Message            string
  var ObjectId           string
  var ObjectIdFileName   string
  var PlayerObjFile     *os.File
  var PlayerObjFileName  string
  var PlayerName         string

  LogBuf = "Begin validation RunningPlayersPlayerObj"
  LogIt(LogBuf)
  if ChgDir(PLAYER_OBJ_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRunningPlayersPlayerObj - Change directory to PLAYER_OBJ_DIR failed")
    os.Exit(1)
  }
  // Get list of all RunningPlayersPlayerObj files
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateRunningPlayersPlayerObj - ReadDir failed")
    os.Exit(1)
  }
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRunningPlayersPlayerObj - Change directory to HomeDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open player file
    PlayerObjFileName  = DirEntry.Name()
    PlayerName         = StrLeft(PlayerObjFileName, StrGetLength(PlayerObjFileName)-4)
    _ = PlayerName
    PlayerObjFileName  = PLAYER_OBJ_DIR + PlayerObjFileName
    PlayerObjFile, err = os.Open(PlayerObjFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateRunningPlayersPlayerObj - Open player file failed")
      os.Exit(1)
    }
    LineCount = 0
    Scanner := bufio.NewScanner(PlayerObjFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = ""
    }
    for Stuff != "" {
      // For all lines
      LineCount++
      ObjectId = StrGetWord(Stuff, 2)
      //************
      //* ObjectId *
      //************
      ObjectIdFileName  = OBJECTS_DIR
      ObjectIdFileName += ObjectId
      ObjectIdFileName += ".txt"
      if !FileExist(ObjectIdFileName) {
        // ObjectId file not found
        Message  = "Object file"
        Message += " '"
        Message += ObjectId
        Message += "' "
        Message += "not found"
        FileName = PlayerObjFileName
        LogValErr(Message, FileName)
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        Stuff = ""
      }
    }
    PlayerObjFile.Close()
  }
  LogBuf = "Done  validating RunningPlayersPlayerObj"
  LogIt(LogBuf)
}

// Validate ValidateRunningRoomMob
func ValidateRunningRoomMob() {
  var FileName          string
  var LineCount         int
  var Message           string
  var MobileId          string
  var MobileIdFileName  string
  var PositionOfDot     int
  var RoomId            string
  var RoomIdFileName    string
  var RoomMobFile      *os.File
  var RoomMobFileName   string

  LogBuf = "Begin validation RunningRoomMob"
  LogIt(LogBuf)
  if ChgDir(ROOM_MOB_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRunningRoomMob - Change directory to ROOM_MOB_DIR failed")
    os.Exit(1)
  }
  // Get list of all RunningRoomMob files
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateRunningRoomMob - ReadDir failed")
    os.Exit(1)
  }
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRunningRoomMob - Change directory to HomeDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open world mobile file
    RoomMobFileName = DirEntry.Name()
    if RoomMobFileName == "ReadMe.txt" {
      // Skip ReadMe files
      continue
    }
    RoomId           = StrLeft(RoomMobFileName, StrGetLength(RoomMobFileName)-4)
    RoomMobFileName  = ROOM_MOB_DIR + RoomMobFileName
    RoomMobFile, err = os.Open(RoomMobFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateRunningRoomMob - Open world mobile file failed")
      os.Exit(1)
    }
    //**********
    //* RoomId *
    //**********
    RoomIdFileName  = ROOMS_DIR
    RoomIdFileName += RoomId
    RoomIdFileName += ".txt"
    if !FileExist(RoomIdFileName) {
      // RoomId file not found
      Message = "Room file"
      Message += " '"
      Message += RoomId
      Message += "' "
      Message += "not found"
      FileName = RoomMobFileName
      LogValErr(Message, FileName)
    }
    //***********************
    //* Check file contents *
    //***********************
    LineCount = 0
    Scanner := bufio.NewScanner(RoomMobFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = ""
    }
    for Stuff != "" {
      // For all lines
      LineCount++
      MobileId = StrGetWord(Stuff, 2)
      PositionOfDot = StrFindFirstChar(MobileId, '.')
      if PositionOfDot > 1 {
        // Mobile is hurt
        MobileId = StrLeft(MobileId, PositionOfDot)
      }
      //************
      //* MobileId *
      //************
      MobileIdFileName = MOBILES_DIR
      MobileIdFileName += MobileId
      MobileIdFileName += ".txt"
      if !FileExist(MobileIdFileName) {
        // MobileId file not found
        Message  = "Mobile file"
        Message += " '"
        Message += MobileId
        Message += "' "
        Message += "not found"
        FileName = RoomMobFileName
        LogValErr(Message, FileName)
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        Stuff = ""
      }
    }
    RoomMobFile.Close()
  }
  LogBuf = "Done  validating RunningRoomMob"
  LogIt(LogBuf)
}

// Validate ValidateRunningRoomObj
func ValidateRunningRoomObj() {
  var FileName          string
  var LineCount         int
  var Message           string
  var ObjectId          string
  var ObjectIdFileName  string
  var RoomId            string
  var RoomIdFileName    string
  var RoomObjFile      *os.File
  var RoomObjFileName   string

  LogBuf = "Begin validation RunningRoomObj"
  LogIt(LogBuf)
  if ChgDir(ROOM_OBJ_DIR) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRunningRoomObj - Change directory to ROOM_OBJ_DIR failed")
    os.Exit(1)
  }
  // Get list of all RunningRoomObj files
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogIt("Validate::ValidateRunningRoomObj - ReadDir failed")
    os.Exit(1)
  }
  if ChgDir(HomeDir) != nil {
    // Change directory failed
    LogIt("Validate::ValidateRunningRoomObj - Change directory to HomeDir failed")
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    // Open RunningRoomObj file
    RoomObjFileName = DirEntry.Name()
    if RoomObjFileName == "ReadMe.txt" {
      // Skip ReadMe files
      continue
    }
    RoomId           = StrLeft(RoomObjFileName, StrGetLength(RoomObjFileName)-4)
    RoomObjFileName  = ROOM_OBJ_DIR + RoomObjFileName
    RoomObjFile, err = os.Open(RoomObjFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogIt("Validate::ValidateRunningRoomObj - Open world mobile file failed")
      os.Exit(1)
    }
    //**********
    //* RoomId *
    //**********
    RoomIdFileName  = ROOMS_DIR
    RoomIdFileName += RoomId
    RoomIdFileName += ".txt"
    if !FileExist(RoomIdFileName) {
      // RoomId file not found
      Message  = "Room file"
      Message += " '"
      Message += RoomId
      Message += "' "
      Message += "not found"
      FileName = RoomObjFileName
      LogValErr(Message, FileName)
    }
    //***********************
    //* Check file contents *
    //***********************
    LineCount = 0
    Scanner := bufio.NewScanner(RoomObjFile)
    if Scanner.Scan() {
      Stuff = Scanner.Text()
    } else {
      Stuff = ""
    }
    for Stuff != "" {
      // For all lines
      LineCount++
      ObjectId = StrGetWord(Stuff, 2)
      //************
      //* ObjectId *
      //************
      ObjectIdFileName = OBJECTS_DIR
      ObjectIdFileName += ObjectId
      ObjectIdFileName += ".txt"
      if !FileExist(ObjectIdFileName) {
        // ObjectId file not found
        Message  = "Object file"
        Message += " '"
        Message += ObjectId
        Message += "' "
        Message += "not found"
        FileName = RoomObjFileName
        LogValErr(Message, FileName)
      }
      if Scanner.Scan() {
        Stuff = Scanner.Text()
      } else {
        Stuff = ""
      }
    }
    RoomObjFile.Close()
  }
  LogBuf = "Done  validating RunningRoomObj"
  LogIt(LogBuf)
}
