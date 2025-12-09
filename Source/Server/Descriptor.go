//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Descriptor.go                                         *
// Usage:     Manages player connections and their descriptors      *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import "os"

// Append new connection to connection list
func AppendIt() {
  pDnodeActor.pDnodePrev = pDnodeHead.pDnodePrev
  pDnodeHead.pDnodePrev.pDnodeNext = pDnodeActor
  pDnodeHead.pDnodePrev = pDnodeActor
  pDnodeActor.pDnodeNext = pDnodeHead
}

// Clear descriptor linked list
func ClearDescriptor() {
  DnodeDestructor(pDnodeHead)
  pDnodeHead = nil
}

// Delete connection from connection list and close socket
func DeleteNode() bool {
  if pDnodeCursor.DnodeFd == 0 {
    return false
  }
  Result := CloseSocket(pDnodeCursor.DnodeFd)
  if Result != 0 {
    PrintIt("Descriptor::DeleteNode - Error: Closesocket")
    os.Exit(1)
  }
  pDnodeActor = nil
  pDnode := pDnodeCursor.pDnodePrev
  DnodeDestructor(pDnodeCursor)
  pDnodeCursor = pDnode
  return true
}

// End of Dnode list?
func EndOfDnodeList() bool {
  return pDnodeCursor.DnodeFd == 0
}

// Get Dnode pointer
func GetDnode() *Dnode {
  return pDnodeCursor
}

// Initialize descriptor pointers
func InitDescriptor() {
  pDnodeHead = DnodeConstructor(0, "")
  pDnodeHead.pDnodeNext = pDnodeHead
  pDnodeHead.pDnodePrev = pDnodeHead
  pDnodeCursor = pDnodeHead
}

// Set Dnode pointer to the first Dnode in the list
func SetpDnodeCursorFirst() {
  pDnodeCursor = pDnodeHead.pDnodeNext
}

// Set Dnode pointer to the next Dnode in the list
func SetpDnodeCursorNext() {
  pDnodeCursor = pDnodeCursor.pDnodeNext
}

// Set Dnode pointer to the previous Dnode in the list
func SetpDnodeCursorPrev() {
  pDnodeCursor = pDnodeCursor.pDnodePrev
}