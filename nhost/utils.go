package nhost

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
)

// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func appendEnvVars(payload map[interface{}]interface{}, prefix string) []string {
	var response []string
	for key, item := range payload {
		switch item := item.(type) {
		/*
			case map[interface{}]interface{}:
				response = append(response, appendEnvVars(item, prefix)...)
		*/
		case map[interface{}]interface{}:
			for key, value := range item {
				switch value := value.(type) {
				case map[interface{}]interface{}:
					for newkey, newvalue := range value {
						if newvalue != "" {
							response = append(response, fmt.Sprintf("%s_%v_%v=%v", prefix, strings.ToUpper(fmt.Sprint(key)), strings.ToUpper(fmt.Sprint(newkey)), newvalue))
						}
					}
				case interface{}, string:
					if value != "" {
						response = append(response, fmt.Sprintf("%s_%v=%v", prefix, strings.ToUpper(fmt.Sprint(key)), value))
					}
				}
			}
		case interface{}:
			if item != "" {
				response = append(response, fmt.Sprintf("%s_%v=%v", prefix, strings.ToUpper(fmt.Sprint(key)), item))
			}
		}
	}
	return response
}

// generate a random 128 byte key
func generateRandomKey() string {
	key := make([]byte, 128)
	rand.Read(key)
	return hex.EncodeToString(key)
}

func GetPort(low, hi int) int {

	// generate a random port value
	port := strconv.Itoa(low + rand.Intn(hi-low))

	// validate wehther the port is available
	if !portAvaiable(port) {
		return GetPort(low, hi)
	}

	// return the value, if it's available
	response, _ := strconv.Atoi(port)
	return response
}

func portAvaiable(port string) bool {

	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return false
	}

	ln.Close()
	return true
}

func GetContainerName(name string) string {
	return strings.Join([]string{PREFIX, name}, "_")
}

func openbrowser(url string) error {

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

func loadRepository() (*git.Repository, error) {

	log.Debug("Loading local git repository")
	return git.PlainOpen(WORKING_DIR)
}

func getCurrentBranch(repo *git.Repository) string {
	head, err := repo.Head()
	if err != nil {
		return ""
	}

	if head.Name().IsBranch() {
		return head.Name().Short()
	}

	return ""
}
