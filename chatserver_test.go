package runtelchat

import (
	"net"
	"testing"
)

//TestChatServerConnects ensures that the chat server can be fired up and ran
//and a connected to
func TestChatServerConnects(t *testing.T) {
	go ListenAndServe(defaultConfig)
	addr := defaultConfig.Host + ":" + defaultConfig.Port
	conn, err := net.Dial("tcp", addr)

	defer conn.Close()

	if err != nil {
		t.Fatal(err)
	}

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		t.Fatal(err)
	}
	if n != len([]byte(welcomePrompt)) {
		t.Fatal("Incorrect server response")
	}

}
