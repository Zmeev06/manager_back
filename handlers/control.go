package handlers

import (
	"stupidauth/models"

	"github.com/digitalocean/go-libvirt"
	"codeberg.org/shinyzero0/sshtunnel"
	"github.com/gofiber/fiber/v2"
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
