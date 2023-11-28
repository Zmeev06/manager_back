package handlers

import (
	"bytes"
	"embed"
	"fmt"
	"path"
	"stupidauth/models"
	"stupidauth/repos"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

//go:embed templates
var createTmpl embed.FS

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
	tmpl, err := template.ParseFS(createTmpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "create", &in); err != nil {
		return err
	}
	_, err = conn.DomainDefineXML(buf.String())
	if err != nil {
		return err
	}
	return nil
}
