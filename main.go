package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

var usage = `usage: vault-mv [options] <source> <destination>
Options:
-p  Preserve data path. (default=true)
`

func main() {

	var sourcePath, destPath string
	var preserveData bool

	flag.BoolVar(&preserveData, "p", true, "")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "%s", usage)
		os.Exit(0)
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
	}

	sourcePath = flag.Args()[0]
	destPath = flag.Args()[1]

	config := api.DefaultConfig()
	vaultClient, err := api.NewClient(config)
	if err != nil {
		logrus.Println(err)
		os.Exit(1)
	}

	logical := vaultClient.Logical()
	secret, err := logical.List(sourcePath)
	if err != nil {
		logrus.Fatal(err)
	}

	if secret == nil {
		secret, err := logical.Read(sourcePath)
		if err != nil {
			log.Fatal(err)
		}

		if !preserveData {
			_, err := logical.Delete(sourcePath)
			if err != nil {
				log.Fatal(err)
			}
		}

		// If the sourePath does not exist this returns a <nil> secret
		// and it would cause a panic due nil pointer dereference.
		if secret != nil {
			secret, err = logical.Write(destPath, secret.Data)
			if err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}

	walkVaultPath(sourcePath, destPath, preserveData, logical)
}

func walkVaultPath(vaultPath, destPath string, preserveDataSource bool, logical *api.Logical) {
	secret, _ := logical.List(vaultPath)
	if secret == nil {

		switch preserveDataSource {
		case false:
			fmt.Printf("moving %s to %s\n", vaultPath, destPath)
		default:
			fmt.Printf("copying %s to %s\n", vaultPath, destPath)
		}
		secret, err := logical.Read(vaultPath)
		if err != nil {
			logrus.Fatal(err)
		}

		secret, err = logical.Write(destPath, secret.Data)
		if err != nil {
			logrus.Fatal(err)
		}

		if !preserveDataSource {
			_, err := logical.Delete(vaultPath)
			if err != nil {
				log.Fatal(err)
			}
		}
		return
	}

	n := secret.Data["keys"]
	for _, key := range n.([]interface{}) {
		newPath := filepath.Join(vaultPath, key.(string))
		dPath := filepath.Join(destPath, key.(string))
		walkVaultPath(newPath, dPath, preserveDataSource, logical)
	}
}
