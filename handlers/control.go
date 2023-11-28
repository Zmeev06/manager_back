package handlers

import (
	"fmt"
	"net"
	"stupidauth/models"
	"stupidauth/repos"
	"time"

	"codeberg.org/shinyzero0/sshtunnel"
	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/ssh"
)

func Start(ctx *fiber.Ctx) error {
	var in models.VmControlByUUIDInput
	if err := ctx.BodyParser(&in); err != nil {
		return err
	}
	var conn *libvirt.Libvirt
	var (
		err error
		tun *sshtunnel.SSHTunnel
	)
	conn, err, tun = getRemoteLibvirt(ctx, in.Host)
	if err != nil {
		return err
	}
	defer tun.Close()
	defer conn.ConnectClose()
	dm, err := conn.DomainLookupByUUID(libvirt.UUID(in.UUID))
	if err != nil {
		return err
	}
	return conn.DomainCreate(dm)
}

func Stop(ctx *fiber.Ctx) error {
	var in models.VmControlByUUIDInput
	if err := ctx.BodyParser(&in); err != nil {
		return err
	}
	var conn *libvirt.Libvirt
	var (
		err error
		tun *sshtunnel.SSHTunnel
	)
	conn, err, tun = getRemoteLibvirt(ctx, in.Host)
	if err != nil {
		return err
	}
	defer tun.Close()
	defer conn.ConnectClose()
	dm, err := conn.DomainLookupByUUID(libvirt.UUID(in.UUID))
	if err != nil {
		return err
	}
	return conn.DomainShutdown(dm)
}
func Delete(ctx *fiber.Ctx) error {
	var in models.VmControlByUUIDInput
	if err := ctx.BodyParser(&in); err != nil {
		return err
	}
	var conn *libvirt.Libvirt
	var (
		err error
		tun *sshtunnel.SSHTunnel
	)
	conn, err, tun = getRemoteLibvirt(ctx, in.Host)
	if err != nil {
		return err
	}
	defer tun.Close()
	defer conn.ConnectClose()
	dom, err := conn.DomainLookupByUUID(libvirt.UUID(in.UUID))
	if err != nil {
		return err
	}
	return conn.DomainUndefine(dom)
}
func getRemoteLibvirt(ctx *fiber.Ctx, host string) (conn *libvirt.Libvirt, err error, tun *sshtunnel.SSHTunnel) {
	username := getUserFromJwt(ctx)
	key, err := repos.GetKey(username)
	if err != nil {
		return
	}
	h, p, err := net.SplitHostPort(host)
	if err.Error() != "missing port in address" {
		return
	} else {
		p = "22"
	}
	newhost := fmt.Sprintf("%s:%s", h, p)
	sshsigner, err := ssh.NewSignerFromKey(&key)
	if err != nil {
		return
	}
	tun, err = sshtunnel.NewSSHTunnel(
		fmt.Sprintf("%s:22", host),
		ssh.PublicKeys(sshsigner),
		"/var/run/libvirt/libvirt-sock",
		"0")
	if err != nil {
		return
	}
	tun.Log = sshtunnel.Logger{
		LogFunc: func(f string, items ...interface{}) {
			fmt.Printf(f, items)
			fmt.Println()
		},
	}
	// tun.Local.Proto = "tcp"
	tun.Remote.Proto = "unix"
	go tun.Start()
	time.Sleep(1 * time.Second)
	c, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("localhost:%d", tun.LocalAddr.(*net.TCPAddr).Port),
		time.Second*5)
	if err != nil {
		return
	}
	conn = libvirt.NewWithDialer(dialers.NewAlreadyConnected(c))
	if err = conn.Connect(); err != nil {
		return
	}
	return
}
