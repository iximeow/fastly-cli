package stats

import (
	"fmt"
	"io"

	"github.com/fastly/cli/pkg/cmd"
	"github.com/fastly/cli/pkg/global"
	"github.com/fastly/cli/pkg/text"
)

// RegionsCommand exposes the Stats Regions API.
type RegionsCommand struct {
	cmd.Base
}

// NewRegionsCommand returns a new command registered under parent.
func NewRegionsCommand(parent cmd.Registerer, g *global.Data) *RegionsCommand {
	var c RegionsCommand
	c.Globals = g
	c.CmdClause = parent.Command("regions", "List stats regions")
	return &c
}

// Exec implements the command interface.
func (c *RegionsCommand) Exec(_ io.Reader, out io.Writer) error {
	resp, err := c.Globals.APIClient.GetRegions()
	if err != nil {
		c.Globals.ErrLog.Add(err)
		return fmt.Errorf("fetching regions: %w", err)
	}

	for _, region := range resp.Data {
		text.Output(out, "%s", region)
	}

	return nil
}
