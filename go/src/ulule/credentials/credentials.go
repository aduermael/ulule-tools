package credentials

import (
	"github.com/Sirupsen/logrus"
	"os"
)

type getStringFunc func(string) (string, error)

func Get(f getStringFunc) (username, apikey string) {
	if len(os.Args) > 1 {
		username = os.Args[1]
	}
	if len(os.Args) > 2 {
		apikey = os.Args[2]
	}
	var err error
	for username == "" {
		username, err = f("username> ")
		if err != nil {
			logrus.Fatal(err)
		}
	}
	for apikey == "" {
		apikey, err = f("apikey> ")
		if err != nil {
			logrus.Fatal(err)
		}
	}
	return
}
