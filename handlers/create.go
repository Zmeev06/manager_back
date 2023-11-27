package handlers

import (
	"bytes"
	"fmt"
	"path"
	"stupidauth/models"
	"stupidauth/repos"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Create(ctx *fiber.Ctx) error {
	type Input struct {
		models.VmControlInput
		Name   string `json:"name"`
		Cpus   uint   `json:"cpus"`
		UUID   string
		Memory uint64
		Image  string `json:"image"`
		Drive  string
	}
	var in Input
	if err := ctx.BodyParser(&in); err != nil {
		return err
	}
	conn, err, tun := getRemoteLibvirt(ctx, "root@"+in.Host)
	if err != nil {
		return err
	}
	in.Memory = 1024 * 512
	in.UUID = uuid.New().String()
	defer tun.Close()
	key, err := repos.GetKey(getUserFromJwt(ctx))
	if err != nil {
		return err
	}
	sess, err := sshSess(in.Host, key)
	if err != nil {
		return err
	}
	imagepath := path.Join("/var/images", fmt.Sprintf("%s.qcow2", in.UUID))
	out, err := sess.CombinedOutput(fmt.Sprintf("cp %s %s", "/var/images/imageimage.qcow2", imagepath))
	if err != nil {
		ctx.Context().SetBody(out)
		return err
	}
	in.Drive = imagepath
	tmpl, err := template.New("domain").Parse(
		`<domain type='qemu'>
  <name>{{.Name}}</name>
  <uuid>{{.UUID}}</uuid>
  <memory>{{.Memory}}</memory>
  <vcpu>{{.Cpus}}</vcpu>
  <os>
    <type arch='x86_64' machine='pc'>hvm</type>
    <boot dev='cdrom'/>
  </os>
  <devices>
    <emulator>/usr/bin/qemu-system-x86_64</emulator>
    <disk type='file' device='cdrom'>
      <source file='/root/downloads/{{.Image}}'/>
      <target dev='hdc'/>
      <readonly/>
    </disk>
    <disk type='file' device='disk'>
      <source file='{{.Drive}}'/>
      <target dev='hda'/>
    </disk>
    <interface type='network'>
      <source network='default'/>
    </interface>
    <graphics type='vnc' port='-1'>
		<listen type='address' address='0.0.0.0'/>
	</graphics>
  </devices>
</domain>`)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, &in); err != nil {
		return err
	}
	_, err = conn.DomainDefineXML(buf.String())
	if err != nil {
		return err
	}
	return nil
}
