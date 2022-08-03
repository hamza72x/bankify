package util

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	col "github.com/hamza72x/go-color"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func LogFatal(err error, cause string) {
	if err != nil {
		log.Println("----------------------------")
		log.Println(col.Purple(cause))
		log.Println(col.Red(err))
		log.Println("+++++++++++++++++++++++++++++")
		os.Exit(1)
	}
}

// PrettyPrint prints a interface with all fields
func PrettyPrint(data interface{}) {
	var p []byte
	p, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%s \n", p)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
