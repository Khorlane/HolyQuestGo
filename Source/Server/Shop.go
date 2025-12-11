//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Shop.go                                               *
// Usage:     Manages shop entities and their interactions          *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
  "bufio"
  "fmt"
  "os"
)

// Is a room a valid shop?
func IsShop(RoomId string) bool {
  ShopFileName := SHOPS_DIR + RoomId + ".txt"
  return FileExist(ShopFileName)
}

// IsShopObject checks if a shop buys and sells a specific object.
func IsShopObj(RoomId string, ObjectName string) bool {
  ShopFileName := SHOPS_DIR + RoomId + ".txt"
  ShopFile, err := os.Open(ShopFileName)
  if err != nil {
    // No such file, but there should be one. This is bad!
    LogIt("Shop::IsShopObject - Shop does not exist")
    return false
  }
  defer ShopFile.Close()
  scanner := bufio.NewScanner(ShopFile)
  for scanner.Scan() {
    Stuff := scanner.Text()
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
    Stuff = StrMakeLower(Stuff)
    if StrGetWord(Stuff, 1) == "item:" {
      ObjectId := StrGetWord(Stuff, 2)
      ObjectName = StrMakeLower(ObjectName)
      if ObjectName == ObjectId {
        // Found a match
        return true
      }
    }
  }
  if err := scanner.Err(); err != nil {
    LogIt("Shop::IsShopObject - Error reading shop file")
    return false
  }
  // Object not found in shop item list
  ShopFile.Close()
  //***************************************************
  //* No match found, try getting match using 'names' *
  //***************************************************
  ShopFile, err = os.Open(ShopFileName)
  if err != nil {
    // No such file, but there should be one. This is bad!
    LogIt("Shop::IsShopObject - Shop does not exist")
    return false
  }
  defer ShopFile.Close()
  scanner = bufio.NewScanner(ShopFile)
  for scanner.Scan() {
    Stuff := scanner.Text()
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
    Stuff = StrMakeLower(Stuff)
    if StrGetWord(Stuff, 1) == "item:" {
      ObjectId := StrGetWord(Stuff, 2)
      NamesCheck := pObject.Names
      NamesCheck = StrMakeLower(NamesCheck)
      if StrFind(NamesCheck, ObjectName) != -1 {
        // Match, Object found in this shop
        return true
      } else {
        LogBuf = ObjectId + " is an invalid shop item - Shop::IsShopObject"
        LogIt(LogBuf)
      }
    }
  }
  if err := scanner.Err(); err != nil {
    LogIt("Shop::IsShopObject - Error reading shop file")
    return false
  }
  // No match found, Object is not buyable from this shop
  return false
}

// List the items available for buying and selling in the shop
func ListObjects(RoomId string, PlayerOut *string) {
  var (
    i             int
    j             int
    ObjectId      string
    ShopFileName  string
    ShopFile     *os.File
    ShopText      string
  )

  ShopFileName = SHOPS_DIR + RoomId + ".txt"
  ShopFile, err := os.Open(ShopFileName)
  if err != nil {
    // No such file???, But there should be, This is bad!
    LogIt("Shop::ListObjects - Shop does not exist")
    return
  }
  defer ShopFile.Close()
  scanner := bufio.NewScanner(ShopFile)
  // Shop welcome message
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "&W"
  pDnodeActor.PlayerOut += Stuff
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  // Headings
  pDnodeActor.PlayerOut += "\r\n"
  // Line one
  Buf = fmt.Sprintf("%-45s", "Items you may buy and sell")
  ShopText = Buf
  pDnodeActor.PlayerOut += ShopText
  pDnodeActor.PlayerOut += " "
  Buf = fmt.Sprintf("%-6s", "Amount")
  ShopText = Buf
  pDnodeActor.PlayerOut += ShopText
  pDnodeActor.PlayerOut += "\r\n"
  // Line two
  Buf = fmt.Sprintf("%-45s", " ")
  ShopText = Buf
  StrReplace(&ShopText, " ", "-")
  pDnodeActor.PlayerOut += ShopText
  pDnodeActor.PlayerOut += " "
  Buf = fmt.Sprintf("%-6s", " ")
  ShopText = Buf
  StrReplace(&ShopText, " ", "-")
  pDnodeActor.PlayerOut += ShopText
  pDnodeActor.PlayerOut += "\r\n"
  // List items for trade
  for scanner.Scan() {
    Stuff = scanner.Text()
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
    Stuff = StrMakeLower(Stuff)
    if Stuff == "End of Items" {
      break
    }
    if StrGetWord(Stuff, 1) == "item:" {
      ObjectId = StrGetWord(Stuff, 2)
      // Assuming IsObject returns (*Object, bool)
      IsObject(ObjectId)
      if pObject != nil {
        ShopText = fmt.Sprintf("%-45s", pObject.Desc1)
        TmpStr := ShopText
        i = StrCountChar(TmpStr, '&')
        i = i * 2
        for j = 1; j <= i; j++ {
          ShopText += " "
        }
        pDnodeActor.PlayerOut += ShopText
        pDnodeActor.PlayerOut += " "
        ShopText = fmt.Sprintf("%6d", pObject.Cost)
        pDnodeActor.PlayerOut += ShopText
        pDnodeActor.PlayerOut += "\r\n"
        // No need to delete in Go, GC handles it
      } else {
        LogBuf := ObjectId + " is an invalid shop item - Shop::ListObjects"
        LogIt(LogBuf)
      }
    }
  }
}