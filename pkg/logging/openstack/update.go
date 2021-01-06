package openstack

import (
	"io"

	"github.com/fastly/cli/pkg/common"
	"github.com/fastly/cli/pkg/compute/manifest"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/text"
	"github.com/fastly/go-fastly/v2/fastly"
)

// UpdateCommand calls the Fastly API to update OpenStack logging endpoints.
type UpdateCommand struct {
	common.Base
	manifest manifest.Data

	//required
	EndpointName string
	Version      int

	// optional
	NewName           common.OptionalString
	BucketName        common.OptionalString
	AccessKey         common.OptionalString
	User              common.OptionalString
	URL               common.OptionalString
	Path              common.OptionalString
	Period            common.OptionalUint
	GzipLevel         common.OptionalUint
	Format            common.OptionalString
	FormatVersion     common.OptionalUint
	ResponseCondition common.OptionalString
	MessageType       common.OptionalString
	TimestampFormat   common.OptionalString
	Placement         common.OptionalString
	PublicKey         common.OptionalString
}

// NewUpdateCommand returns a usable command registered under the parent.
func NewUpdateCommand(parent common.Registerer, globals *config.Data) *UpdateCommand {
	var c UpdateCommand
	c.Globals = globals
	c.manifest.File.Read(manifest.Filename)

	c.CmdClause = parent.Command("update", "Update an OpenStack logging endpoint on a Fastly service version")

	c.CmdClause.Flag("service-id", "Service ID").Short('s').StringVar(&c.manifest.Flag.ServiceID)
	c.CmdClause.Flag("version", "Number of service version").Required().IntVar(&c.Version)
	c.CmdClause.Flag("name", "The name of the OpenStack logging object").Short('n').Required().StringVar(&c.EndpointName)

	c.CmdClause.Flag("new-name", "New name of the OpenStack logging object").Action(c.NewName.Set).StringVar(&c.NewName.Value)
	c.CmdClause.Flag("bucket", "The name of the Openstack Space").Action(c.BucketName.Set).StringVar(&c.BucketName.Value)
	c.CmdClause.Flag("access-key", "Your OpenStack account access key").Action(c.AccessKey.Set).StringVar(&c.AccessKey.Value)
	c.CmdClause.Flag("user", "The username for your OpenStack account.").Action(c.User.Set).StringVar(&c.User.Value)
	c.CmdClause.Flag("url", "Your OpenStack auth url.").Action(c.URL.Set).StringVar(&c.URL.Value)
	c.CmdClause.Flag("path", "The path to upload logs to").Action(c.Path.Set).StringVar(&c.Path.Value)
	c.CmdClause.Flag("period", "How frequently log files are finalized so they can be available for reading (in seconds, default 3600)").Action(c.Period.Set).UintVar(&c.Period.Value)
	c.CmdClause.Flag("gzip-level", "What level of GZIP encoding to have when dumping logs (default 0, no compression)").Action(c.GzipLevel.Set).UintVar(&c.GzipLevel.Value)
	c.CmdClause.Flag("format", "Apache style log formatting").Action(c.Format.Set).StringVar(&c.Format.Value)
	c.CmdClause.Flag("format-version", "The version of the custom logging format used for the configured endpoint. Can be either 2 (default) or 1").Action(c.FormatVersion.Set).UintVar(&c.FormatVersion.Value)
	c.CmdClause.Flag("response-condition", "The name of an existing condition in the configured endpoint, or leave blank to always execute").Action(c.ResponseCondition.Set).StringVar(&c.ResponseCondition.Value)
	c.CmdClause.Flag("message-type", "How the message should be formatted. One of: classic (default), loggly, logplex or blank").Action(c.MessageType.Set).StringVar(&c.MessageType.Value)
	c.CmdClause.Flag("timestamp-format", `strftime specified timestamp formatting (default "%Y-%m-%dT%H:%M:%S.000")`).Action(c.TimestampFormat.Set).StringVar(&c.TimestampFormat.Value)
	c.CmdClause.Flag("placement", "Where in the generated VCL the logging call should be placed, overriding any format_version default. Can be none or waf_debug").Action(c.Placement.Set).StringVar(&c.Placement.Value)
	c.CmdClause.Flag("public-key", "A PGP public key that Fastly will use to encrypt your log files before writing them to disk").Action(c.PublicKey.Set).StringVar(&c.PublicKey.Value)

	return &c
}

// createInput transforms values parsed from CLI flags into an object to be used by the API client library.
func (c *UpdateCommand) createInput() (*fastly.UpdateOpenstackInput, error) {
	serviceID, source := c.manifest.ServiceID()
	if source == manifest.SourceUndefined {
		return nil, errors.ErrNoServiceID
	}

	openstack, err := c.Globals.Client.GetOpenstack(&fastly.GetOpenstackInput{
		ServiceID:      serviceID,
		Name:           c.EndpointName,
		ServiceVersion: c.Version,
	})
	if err != nil {
		return nil, err
	}

	input := fastly.UpdateOpenstackInput{
		ServiceID:         openstack.ServiceID,
		ServiceVersion:    openstack.ServiceVersion,
		Name:              openstack.Name,
		NewName:           fastly.String(openstack.Name),
		BucketName:        fastly.String(openstack.BucketName),
		AccessKey:         fastly.String(openstack.AccessKey),
		User:              fastly.String(openstack.User),
		URL:               fastly.String(openstack.URL),
		Path:              fastly.String(openstack.Path),
		Period:            fastly.Uint(openstack.Period),
		GzipLevel:         fastly.Uint(openstack.GzipLevel),
		Format:            fastly.String(openstack.Format),
		FormatVersion:     fastly.Uint(openstack.FormatVersion),
		ResponseCondition: fastly.String(openstack.ResponseCondition),
		MessageType:       fastly.String(openstack.MessageType),
		TimestampFormat:   fastly.String(openstack.TimestampFormat),
		Placement:         fastly.String(openstack.Placement),
		PublicKey:         fastly.String(openstack.PublicKey),
	}

	// Set new values if set by user.
	if c.NewName.WasSet {
		input.NewName = fastly.String(c.NewName.Value)
	}

	if c.BucketName.WasSet {
		input.BucketName = fastly.String(c.BucketName.Value)
	}

	if c.AccessKey.WasSet {
		input.AccessKey = fastly.String(c.AccessKey.Value)
	}

	if c.User.WasSet {
		input.User = fastly.String(c.User.Value)
	}

	if c.URL.WasSet {
		input.URL = fastly.String(c.URL.Value)
	}

	if c.Path.WasSet {
		input.Path = fastly.String(c.Path.Value)
	}

	if c.Period.WasSet {
		input.Period = fastly.Uint(c.Period.Value)
	}

	if c.GzipLevel.WasSet {
		input.GzipLevel = fastly.Uint(c.GzipLevel.Value)
	}

	if c.Format.WasSet {
		input.Format = fastly.String(c.Format.Value)
	}

	if c.FormatVersion.WasSet {
		input.FormatVersion = fastly.Uint(c.FormatVersion.Value)
	}

	if c.ResponseCondition.WasSet {
		input.ResponseCondition = fastly.String(c.ResponseCondition.Value)
	}

	if c.MessageType.WasSet {
		input.MessageType = fastly.String(c.MessageType.Value)
	}

	if c.TimestampFormat.WasSet {
		input.TimestampFormat = fastly.String(c.TimestampFormat.Value)
	}

	if c.Placement.WasSet {
		input.Placement = fastly.String(c.Placement.Value)
	}

	if c.PublicKey.WasSet {
		input.PublicKey = fastly.String(c.PublicKey.Value)
	}

	return &input, nil
}

// Exec invokes the application logic for the command.
func (c *UpdateCommand) Exec(in io.Reader, out io.Writer) error {
	input, err := c.createInput()
	if err != nil {
		return err
	}

	openstack, err := c.Globals.Client.UpdateOpenstack(input)
	if err != nil {
		return err
	}

	text.Success(out, "Updated OpenStack logging endpoint %s (service %s version %d)", openstack.Name, openstack.ServiceID, openstack.ServiceVersion)
	return nil
}