package secretstore

import (
	"io"

	"github.com/fastly/cli/pkg/cmd"
	fsterr "github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/global"
	"github.com/fastly/cli/pkg/manifest"
	"github.com/fastly/cli/pkg/text"
	"github.com/fastly/go-fastly/v7/fastly"
)

// NewCreateCommand returns a usable command registered under the parent.
func NewCreateCommand(parent cmd.Registerer, g *global.Data, m manifest.Data) *CreateCommand {
	c := CreateCommand{
		Base: cmd.Base{
			Globals: g,
		},
		manifest: m,
	}

	c.CmdClause = parent.Command("create", "Create a new secret store")

	// Required.
	c.RegisterFlag(storeNameFlag(&c.Input.Name)) // --name

	// Optional.
	c.RegisterFlagBool(c.JSONFlag()) // --json

	return &c
}

// CreateCommand calls the Fastly API to create an appropriate resource.
type CreateCommand struct {
	cmd.Base
	cmd.JSONOutput

	Input    fastly.CreateSecretStoreInput
	manifest manifest.Data
}

// Exec invokes the application logic for the command.
func (cmd *CreateCommand) Exec(_ io.Reader, out io.Writer) error {
	if cmd.Globals.Verbose() && cmd.JSONOutput.Enabled {
		return fsterr.ErrInvalidVerboseJSONCombo
	}

	o, err := cmd.Globals.APIClient.CreateSecretStore(&cmd.Input)
	if err != nil {
		cmd.Globals.ErrLog.Add(err)
		return err
	}

	if ok, err := cmd.WriteJSON(out, o); ok {
		return err
	}

	text.Success(out, "Created secret store %s (name %s)", o.ID, o.Name)

	return nil
}