package handlers

import (
	"fmt"
	"math"
	"net"
	"stupidauth/models"
	"stupidauth/repos"
	"time"

	"github.com/digitalocean/go-libvirt"
	"codeberg.org/shinyzero0/sshtunnel"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/ssh"
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
func getUserFromJwt(ctx *fiber.Ctx) string {
	token := ctx.Locals("user").(*jwt.Token)
	return token.Claims.(jwt.MapClaims)["user"].(string)
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
