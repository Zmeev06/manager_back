package handlers

import (
	"fmt"
	"net"
	"stupidauth/models"
	"stupidauth/repos"
	"time"

	"codeberg.org/shinyzero0/sshtunnel"
	"github.com/digitalocean/go-libvirt"
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
func Delete(ctx *fiber.Ctx) error{
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
	sshsigner, err := ssh.NewSignerFromKey(&key)
	if err != nil {
		return
	}
	tun, err = sshtunnel.NewSSHTunnel(
		fmt.Sprintf("%s", host),
		ssh.PublicKeys(sshsigner),
		"127.0.0.1:16509",
		"0")
	if err != nil {
		return
	}
	go tun.Start()
	time.Sleep(1 * time.Second)
	c, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("localhost:%d", tun.Local.Port),
		time.Second*5)
	if err != nil {
		return
	}
	conn = libvirt.New(c)
	if err = conn.Connect(); err != nil {
		return
	}
	return
}
