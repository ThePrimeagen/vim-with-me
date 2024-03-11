package main

import (
	"fmt"
	"log"

	"chat.theprimeagen.com/pkg/chat"
)

//
// func readFromStdin() chan string {
//     buffer := make([]byte, 1024)
//     out := make(chan string)
//
//     go func() {
//         for {
//             count, err := os.Stdin.Read(buffer)
//             if err == io.EOF {
//                 break
//             }
//             out <- string(buffer[:count])
//         }
//     }()
//
//     return out
// }
//
// func createTCPServer() chan string {
//     out := make(chan string)
//
//     go func() {
//         listener, err := net.Listen("tcp", ":42069")
//         if err != nil {
//             log.Fatal("You are a horrible human being", err)
//         }
//         defer listener.Close()
//
//         for {
//             conn, err := listener.Accept()
//             if err != nil {
//                 log.Fatal("You like amouranth", err)
//             }
//             go func(c net.Conn) {
//                 defer c.Close()
//                 for {
//                     str := <-out
//                     str = fmt.Sprintf("%d:%s", len(str), str)
//                     _, err := c.Write([]byte(str))
//                     if err != nil {
//                         fmt.Printf("Error writing to client: %s\n", err)
//                         break
//                     }
//                 }
//             }(conn)
//         }
//     }()
//
//     return out
// }
//
// //var allowableChars = []string{"<dot>", "<backspace>", "<space>", "<esc>", "<cr>", "<tab>"};
// //func contains(arr []string, str string) bool {
// //    for _, a := range arr {
// //        if a == str {
// //            return true
// //        }
// //    }
// //    return false
// //}
// //
// func main() {
//     // read from standard in line by line
//     // stdin := readFromStdin()
//     tcpOut := createTCPServer()
//
//     /*
//     processor := processors.NewTDProcessor(5)
//
//     for {
//         select {
//         case s := <-stdin:
//             processor.Process(strings.TrimSpace(s))
//
//         case point := <-processor.Out():
//             fmt.Printf("Got a point: %s\n", point)
//             tcpOut <- point
//         }
//     }
//     */
//
//     ticker := time.NewTicker(16 * time.Millisecond)
//
//     // create an 80x24 grid string arrays
//     length := 80 * 24
//     grid := make([]byte, length)
//     count := 0
//
//     for {
//         select {
//         case <-ticker.C:
//             for i := 0; i < length; i++ {
//                 if (i + count) % 4 == 0 {
//                     grid[i] = 'X'
//                 } else {
//                     grid[i] = ' '
//                 }
//             }
//
//             str := fmt.Sprintf("r:%s", string(grid))
//             tcpOut <- str
//             count++
//         }
//     }
//
// }
//
func main() {
    c, err := chat.FromChatProgram("./chat.js")
    if err != nil {
        log.Fatal("Error creating chat program", err)
    }

    for msg := range c.Chat {
        fmt.Printf("Got a message: %v\n", msg)
    }
}

