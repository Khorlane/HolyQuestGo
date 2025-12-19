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
  }
  return false
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
    LogIt("Shop::IsShopObj - Shop does not exist")
    os.Exit(1) // _endthread()
  }
  scanner := bufio.NewScanner(ShopFile)
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  for Stuff != "End of Items" {
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
    Stuff = StrMakeLower(Stuff)
    if StrGetWord(Stuff, 1) == "item:" {
      ObjectId = StrGetWord(Stuff, 2)
      ObjectName = StrMakeLower(ObjectName)
      if ObjectName == ObjectId {
        IsObject(ObjectId)
        if pObject != nil {
          ShopFile.Close()
          return
        } else {
          LogBuf := ObjectId + " is an invalid shop item - Shop::IsShopObj"
          LogIt(LogBuf)
          pObject = nil
          ShopFile.Close()
          return
        }
      }
    }
    if !scanner.Scan() {
      break
    }
    Stuff = scanner.Text()
  }
  ShopFile.Close()
  ShopFile, err = os.Open(ShopFileName)
  if err != nil {
    LogIt("Shop::IsShopObj - Shop does not exist")
    os.Exit(1)
  }
  scanner = bufio.NewScanner(ShopFile)
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  for Stuff != "End of Items" {
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
    Stuff = StrMakeLower(Stuff)
    if StrGetWord(Stuff, 1) == "item:" {
      ObjectId = StrGetWord(Stuff, 2)
      IsObject(ObjectId)
      if pObject != nil {
        NamesCheck = pObject.Names
        NamesCheck = StrMakeLower(NamesCheck)
        Result = StrFind(NamesCheck, ObjectName)
        if Result != -1 {
          ShopFile.Close()
          return
        } else {
          pObject = nil
        }
      } else {
        LogBuf := ObjectId + " is an invalid shop item - Shop::IsShopObj"
        LogIt(LogBuf)
        pObject = nil
      }
    }
    if !scanner.Scan() {
      break
    }
    Stuff = scanner.Text()
  }
  ShopFile.Close()
}

// List objects that can be bought and sold
func ListObjects(pDnodeActor *Dnode) {
  var i int
  var j int
  var ObjectId string
  var ShopFileName string
  var ShopFile *os.File
  var ShopText string

  ShopFileName = SHOPS_DIR
  ShopFileName += pDnodeActor.pPlayer.RoomId
  ShopFileName += ".txt"
  f, err := os.Open(ShopFileName)
  if err != nil {
    LogIt("Shop::ListObjects - Shop does not exist")
    os.Exit(1)
  }
  ShopFile = f
  Scanner := bufio.NewScanner(ShopFile)
  if Scanner.Scan() {
    Stuff = Scanner.Text()
  }
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "&W"
  pDnodeActor.PlayerOut += Stuff
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "\r\n"
  Buf = fmt.Sprintf("%-45s", "Items you may buy and sell")
  ShopText = Buf
  pDnodeActor.PlayerOut += ShopText
  pDnodeActor.PlayerOut += " "
  Buf = fmt.Sprintf("%-6s", "Amount")
  ShopText = Buf
  pDnodeActor.PlayerOut += ShopText
  pDnodeActor.PlayerOut += "\r\n"
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
  if Scanner.Scan() {
    Stuff = Scanner.Text()
  }
  for Stuff != "End of Items" {
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
    Stuff = StrMakeLower(Stuff)
    if StrGetWord(Stuff, 1) == "item:" {
      ObjectId = StrGetWord(Stuff, 2)
      pObject = nil
      IsObject(ObjectId)
      if pObject != nil {
        Buf = fmt.Sprintf("%-45s", pObject.Desc1)
        ShopText = Buf
        TmpStr = ShopText
        i = StrCountChar(TmpStr, '&')
        i = i * 2
        for j = 1; j <= i; j++ {
          ShopText += " "
        }
        pDnodeActor.PlayerOut += ShopText
        pDnodeActor.PlayerOut += " "
        Buf = fmt.Sprintf("%6d", pObject.Cost)
        ShopText = Buf
        pDnodeActor.PlayerOut += ShopText
        pDnodeActor.PlayerOut += "\r\n"
        pObject = nil
      } else {
        LogBuf = ObjectId + " is an invalid shop item - Shop::ListObjects"
        LogIt(LogBuf)
      }
    }
    if !Scanner.Scan() {
      break
    }
    Stuff = Scanner.Text()
  }
  ShopFile.Close()
}