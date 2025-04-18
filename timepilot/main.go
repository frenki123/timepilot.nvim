package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// Your custom event params
type TrackEventParams struct {
	Type      string `json:"type"`
	Filepath  string `json:"filepath"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
    fmt.Println("Starting...")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var req Request
		err := json.Unmarshal(scanner.Bytes(), &req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse: %v\n", err)
			continue
		}

		switch req.Method {
		case "track_event":
			var p TrackEventParams
			json.Unmarshal(req.Params, &p)
            _ = os.WriteFile("test.log", []byte(p.Type), 0644)
            
			//fmt.Fprintf(os.Stderr, "Received event: %+v\n", p)

			res := Response{
				Jsonrpc: "2.0",
				ID:      req.ID,
				Result:  "OK",
			}
			_ = json.NewEncoder(os.Stdout).Encode(res)
            os.Stdout.Sync()

		default:
			res := Response{
				Jsonrpc: "2.0",
				ID:      req.ID,
				Error:   "Unknown method",
			}
			_ = json.NewEncoder(os.Stdout).Encode(res)
		}
	}
}
