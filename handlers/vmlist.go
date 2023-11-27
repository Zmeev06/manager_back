package handlers

import (
	"math"
	"stupidauth/models"

	"github.com/digitalocean/go-libvirt"
	"github.com/gofiber/fiber/v2"
)

type DomainInfo struct {
	Name   string      `json:"name"`
	ID     int32       `json:"id"`
	UUID   models.UUID `json:"uuid"`
	State  uint8       `json:"state"`
	MaxMem uint64      `json:"mem_max"`
	Mem    uint64      `json:"mem_used"`
	Cpus   uint16      `json:"cpus"`
	// CpuTime   uint64 `json:"cpu_teim"`
}

func VmList(ctx *fiber.Ctx) error {
	var in models.VmControlInput
	if err := ctx.BodyParser(&in); err != nil {
		return err
	}
	var conn *libvirt.Libvirt
	conn, err, tun := getRemoteLibvirt(ctx, in.Host)
	if err != nil {
		return err
	}
	defer tun.Close()
	dms, n, err := conn.ConnectListAllDomains(math.MaxInt32, 0)
	if err != nil {
		return err
	}
	dminfos := make([]DomainInfo, n)
	for k, v := range dms {
		state, mmem, mem, cpus, _, err := conn.DomainGetInfo(v)
		if err != nil {
			continue
		}
		dminfos[k] = DomainInfo{
			v.Name,
			v.ID,
			models.UUID(v.UUID),
			state,
			mmem,
			mem,
			cpus,
		}
	}
	return ctx.JSON(dminfos)
}
