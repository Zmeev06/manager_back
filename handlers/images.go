package handlers

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"stupidauth/models"
	"stupidauth/repos"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/ssh"
)

func Images(ctx *fiber.Ctx) error {
	var in models.VmControlInput
	if err := ctx.BodyParser(&in); err != nil {
		return err
	}
	username := getUserFromJwt(ctx)
	key, err := repos.GetKey(username)
	if err != nil {
		return err
	}
	sess, err := sshSess(in.Host, &key)
	out, err := sess.Output("ls /root/downloads")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	// _, err = buf.Write(out)
	if err != nil {
		return err
	}
	var lines = []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return ctx.JSON(lines)
}
func sshSess(host string, key *rsa.PrivateKey) (sess *ssh.Session, err error) {
	sshsigner, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return
	}
	config := ssh.ClientConfig{
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(sshsigner)},
	}
	cli, err := ssh.Dial("tcp", host+":22", &config)
	if err != nil {
		return
	}
	return cli.NewSession()
}
