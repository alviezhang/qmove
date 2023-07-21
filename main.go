package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/alecthomas/kingpin/v2"
)

var (
	source      = kingpin.Flag("source", "Source download path").Short('s').Required().String()
	destination = kingpin.Flag("destination", "Destination path").Short('d').Required().String()
	category    = kingpin.Flag("category", "Category").Short('c').Required().String()

	permission = kingpin.Flag("perm", "Permission in number format").Short('p').String()
)

func getTargetDirectory(destination string, category string) string {
	return filepath.Join(destination, category)
}

func createDirectory(path string, uid int, gid int, permission int) error {
	// check if directory exists
	var fmode os.FileMode
	if permission == -1 {
		fmode = os.ModePerm
	} else {
		fmode = os.FileMode(permission)
	}

	err := os.MkdirAll(path, fmode)
	if err != nil {
		return err
	}
	// change owner and group
	if uid != -1 || gid != -1 {
		err := os.Chown(path, uid, gid)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	targetDirectory := getTargetDirectory(*destination, *category)

	uid := -1
	gid := -1
	perm := -1

	owner := ""
	group := ""

	if owner != "" {
		// Parse uid and gid
		_uid, err := strconv.ParseInt(owner, 10, 32)
		if err != nil {
			user, err := user.Lookup(owner)
			if err != nil {
				fmt.Printf("Error getting user: %s\n", owner)
				os.Exit(1)
			}
			_uid, err = strconv.ParseInt(user.Uid, 10, 32)
			if err != nil {
				fmt.Printf("Error parsing uid: %s\n", user.Uid)
				os.Exit(1)
			}
		}
		uid = int(_uid)
	}

	if group != "" {
		_gid, err := strconv.ParseInt(group, 10, 32)
		if err != nil {
			group, err := user.LookupGroup(group)
			if err != nil {
				fmt.Printf("Error getting group: %s\n", group)
				os.Exit(1)
			}
			_gid, err = strconv.ParseInt(group.Gid, 10, 32)
			if err != nil {
				fmt.Printf("Error parsing gid: %s\n", group.Gid)
				os.Exit(1)
			}
		}
		gid = int(_gid)
	}

	if *permission != "" {
		// parse permission
		_perm, err := strconv.ParseInt(*permission, 8, 32)

		if err != nil {
			fmt.Printf("Error parsing permission: %s\n", *permission)
			os.Exit(1)
		}
		perm = int(_perm)
	}

	err := createDirectory(targetDirectory, uid, gid, perm)

	if err != nil {
		fmt.Printf("Error creating directory: %s\n", targetDirectory)
		os.Exit(1)
	}

	targetPath := filepath.Join(targetDirectory, filepath.Base(*source))

	err = os.Rename(*source, targetPath)

	if err != nil {
		fmt.Printf("Error when mv file:\n%s", err)
		os.Exit(2)
	}

	switch {
	case uid != -1 || gid != -1:
		err := os.Chown(targetPath, uid, gid)

		if err != nil {
			fmt.Printf("Error when setting user and group:\n%s", err)
			os.Exit(3)
		}

	case perm != -1:
		err := os.Chmod(targetPath, os.FileMode(perm))

		if err != nil {
			fmt.Printf("Error when change permission:\n%s", err)
			os.Exit(4)
		}
	}
}
