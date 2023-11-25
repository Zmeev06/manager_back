package handlers

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"stupidauth/repos"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/digitalocean/go-libvirt"
	"github.com/elliotchance/sshtunnel"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/ssh"
)

type UUID [16]byte // {{{

func (this UUID) MarshalJSON() ([]byte, error) {
	v := hex.EncodeToString(this[:])
	return json.Marshal(v)
}
func (this *UUID) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	bytes, err := hex.DecodeString(v)
	copy(this[:], bytes)
	return err
} // }}}

type DomainInfo struct {
	Name   string `json:"name"`
	ID     int32  `json:"id"`
	UUID   UUID   `json:"uuid"`
	State  uint8  `json:"state"`
	MaxMem uint64 `json:"mem_max"`
	Mem    uint64 `json:"mem_used"`
	Cpus   uint16 `json:"cpus"`
	// CpuTime   uint64 `json:"cpu_teim"`
}

func VmList(ctx *fiber.Ctx) error {
	type Input struct {
		Host string `json:"host"`
	}
	var in Input
	if err := ctx.BodyParser(&in); err != nil {
		return err
	}
	var conn *libvirt.Libvirt
	if in.Host == "" {
		var err error
		conn, err = getLocalLibvirt()
		if err != nil {
			return err
		}
	} else {
		var (
			err error
			tun *sshtunnel.SSHTunnel
		)
		conn, err, tun = getRemoteLibvirt(ctx, in.Host)
		if err != nil {
			return err
		}
		defer tun.Close()
	}
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
			UUID(v.UUID),
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
func getLocalLibvirt() (*libvirt.Libvirt, error) {
	c, err := net.DialTimeout("unix", "/var/run/libvirt/libvirt-sock", time.Second*2)
	if err != nil {
		return nil, err
	}
	conn := libvirt.New(c)
	if err := conn.Connect(); err != nil {
		return nil, err
	}
	return conn, nil
}
func getRemoteLibvirt(ctx *fiber.Ctx, host string) (conn *libvirt.Libvirt, err error, tun *sshtunnel.SSHTunnel) {
	var buf bytes.Buffer
	username := getUserFromJwt(ctx)
	var key rsa.PrivateKey
	err = repos.Keys.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(username))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return gob.NewDecoder(&buf).Decode(&key)
		})
	})
	if err != nil {
		return
	}
	sshsigner, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return
	}
	tun, err = sshtunnel.NewSSHTunnel(
		fmt.Sprintf("root@%s", host),
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
