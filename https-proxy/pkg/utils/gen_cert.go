package utils

import (
	"math/rand"
	"os/exec"
	"strconv"
)

const maxRand = 10 ^ 6

func GenHostCert(path, name, host, savePath string) error {
	genCmd := exec.Command(path+name, host, strconv.Itoa(rand.Intn(maxRand)), path, savePath)
	_, err := genCmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
