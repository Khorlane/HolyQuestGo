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

// Is a valid shop?
func IsShop(RoomId string) bool {
  var ShopFileName string

  ShopFileName = SHOPS_DIR
  ShopFileName += RoomId
  ShopFileName += ".txt"
  if FileExist(ShopFileName) {
    return true
  } else {
    return false
  }
}

// Is this shop buying and selling this object?
func IsShopObj(RoomId string, ObjectName string) {
  var NamesCheck string
  var ObjectId string
  var Result int
  var ShopFileName string
  var ShopFile *os.File

  ShopFileName = SHOPS_DIR
  ShopFileName += RoomId
  ShopFileName += ".txt"
  ShopFile, err := os.Open(ShopFileName)
  if err != nil {
    // No such file???, But there should be, This is bad!
    LogIt("Shop::IsShopObj - Shop does not exist")
    os.Exit(1) // _endthread()
  }
  scanner := bufio.NewScanner(ShopFile)
  scanner.Scan()
  Stuff = scanner.Text()
  for Stuff != "End of Items" {
    // Read 'item' lines in ShopFile
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
    Stuff = StrMakeLower(Stuff)
    if StrGetWord(Stuff, 1) == "item:" {
      // Found an item
      ObjectId = StrGetWord(Stuff, 2)
      ObjectName = StrMakeLower(ObjectName)
      if ObjectName == ObjectId {
        // Found a match
        IsObject(ObjectId)
        if pObject != nil {
          // Object exists
          ShopFile.Close()
          return
        } else {
          // Object does not exist, Log it
          LogBuf := ObjectId + " is an invalid shop item - Shop::IsShopObj"
          LogIt(LogBuf)
          pObject = nil
          ShopFile.Close()
          return
        }
      }
    }
    scanner.Scan()
    Stuff = scanner.Text()
  }
  // Object not found in shop item list
  ShopFile.Close()
  //***************************************************
  //* No match found, try getting match using 'names' *
  //***************************************************
  ShopFile, err = os.Open(ShopFileName)
  if err != nil {
    // No such file???, But there should be, This is bad!
    LogIt("Shop::IsShopObj - Shop does not exist")
    os.Exit(1)
  }
  scanner = bufio.NewScanner(ShopFile)
  scanner.Scan()
  Stuff = scanner.Text()
  for Stuff != "End of Items" {
    // Read 'item' lines in ShopFile
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
    Stuff = StrMakeLower(Stuff)
    if StrGetWord(Stuff, 1) == "item:" {
      // Found an item
      ObjectId = StrGetWord(Stuff, 2)
      IsObject(ObjectId)
      if pObject != nil {
        // Check for a match
        NamesCheck = pObject.Names
        NamesCheck = StrMakeLower(NamesCheck)
        Result = StrFind(NamesCheck, ObjectName)
        if Result != -1 {
          // Match, Object found in this shop
          ShopFile.Close()
          return
        } else {
          pObject = nil
        }
      } else {
        // Object does not exist, Log it
        LogBuf := ObjectId + " is an invalid shop item - Shop::IsShopObj"
        LogIt(LogBuf)
        pObject = nil
      }
    }
    scanner.Scan()
    Stuff = scanner.Text()
  }
  ShopFile.Close()
  // No match found, Object is not buyable from this shop
}

// List objects that can be bought and sold
func ListObjects() {
  var i             int
  var j             int
  var ObjectId      string
  var ShopFileName  string
  var ShopFile     *os.File
  var ShopText      string

  ShopFileName = SHOPS_DIR
  ShopFileName += pDnodeActor.pPlayer.RoomId
  ShopFileName += ".txt"
  f, err := os.Open(ShopFileName)
  if err != nil {
    // No such file???, But there should be, This is bad!
    LogIt("Shop::ListObjects - Shop does not exist")
    os.Exit(1)
  }
  ShopFile = f
  Scanner := bufio.NewScanner(ShopFile)
  // Shop welcome message
  Scanner.Scan()
  Stuff = Scanner.Text()
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
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "End of Items" {
    // Read 'item' lines in ShopFile
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
    Stuff = StrMakeLower(Stuff)
    if StrGetWord(Stuff, 1) == "item:" {
      // Found an item
      ObjectId = StrGetWord(Stuff, 2)
      pObject = nil
      IsObject(ObjectId) // Sets pObject
      if pObject != nil {
        // Format shop item text
        Buf = fmt.Sprintf("%-45s", pObject.Desc1)
        ShopText = Buf
        TmpStr = ShopText
        i = StrCountChar(TmpStr, '&')
        i = i * 2
        for j = 1; j <= i; j++ {
          // Color codes will be removed, so adjust length
          ShopText += " "
        }
        pDnodeActor.PlayerOut += ShopText
        pDnodeActor.PlayerOut += " "
        Buf = fmt.Sprintf("%6d", pObject.Cost)
        ShopText = Buf
        pDnodeActor.PlayerOut += ShopText
        pDnodeActor.PlayerOut += "\r\n"
        // Done with object
        pObject = nil
      } else {
        // Object does not exist, Log it
        LogBuf = ObjectId + " is an invalid shop item - Shop::ListObjects"
        LogIt(LogBuf)
      }
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  ShopFile.Close()
}
