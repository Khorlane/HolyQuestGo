package server

// Send output to other players in the same room as player
func SendToRoom(TargetRoomId string, MsgText string) {
}

// Placeholder for Communication-related functions
func SockCheckForNewConnections() {
}

func SockRecv() {
}

func SockOpenPort(port int) {
}

func SockClosePort(port int) {
}

// Close a socket handle (placeholder for C++ closesocket)
func CloseSocket(fd int) int {
  // TODO: replace with real socket close when socket layer is added
  return 0
}
