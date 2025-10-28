package main

import (
	"context"
	"log"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

func main() {
	handler := jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
		switch req.Method {
		case "initialize":
			// Respond with minimal capabilities
			return map[string]interface{}{"capabilities": map[string]interface{}{}}, nil
		case "shutdown":
			return nil, nil
			// Add more LSP methods here
		}
		return nil, nil
	})
	stream := jsonrpc2.NewBufferedStream(os.Stdin, jsonrpc2.VSCodeObjectCodec{})
	conn := jsonrpc2.NewConn(context.Background(), stream, handler)
	<-conn.DisconnectNotify()
	log.Println("Popcorn LSP server exited")
}
