package profile

import (
	"fmt"
	"io"

	"github.com/fastly/cli/pkg/cmd"
	fsterr "github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/global"
	"github.com/fastly/cli/pkg/profile"
	"github.com/fastly/cli/pkg/text"
)

// TokenCommand represents a Kingpin command.
type TokenCommand struct {
	cmd.Base
	profile string
}

// NewTokenCommand returns a new command registered in the parent.
func NewTokenCommand(parent cmd.Registerer, g *global.Data) *TokenCommand {
	var c TokenCommand
	c.Globals = g
	c.CmdClause = parent.Command("token", "Print access token")
	c.CmdClause.Flag("name", "Print access token for the named profile").Short('n').StringVar(&c.profile)
	return &c
}

// Exec implements the command interface.
func (c *TokenCommand) Exec(_ io.Reader, out io.Writer) (err error) {
	if c.profile == "" {
		if name, p := profile.Default(c.Globals.Config.Profiles); name != "" {
			text.Output(out, p.Token)
			return nil
		}
		return fsterr.RemediationError{
			Inner:       fmt.Errorf("no profiles available"),
			Remediation: fsterr.ProfileRemediation,
		}
	}

	if name, p := profile.Get(c.profile, c.Globals.Config.Profiles); name != "" {
		text.Output(out, p.Token)
		return nil
	}
	msg := fmt.Sprintf(profile.DoesNotExist, c.profile)
	return fsterr.RemediationError{
		Inner:       fmt.Errorf(msg),
		Remediation: fsterr.ProfileRemediation,
	}
}
