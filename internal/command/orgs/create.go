package orgs

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/superfly/flyctl/api"
	"github.com/superfly/flyctl/iostreams"

	"github.com/superfly/flyctl/client"
	"github.com/superfly/flyctl/internal/command"
	"github.com/superfly/flyctl/internal/config"
	"github.com/superfly/flyctl/internal/flag"
	"github.com/superfly/flyctl/internal/prompt"
	"github.com/superfly/flyctl/internal/render"
)

func newCreate() *cobra.Command {
	const (
		long = `Create a new organization. Other users can be invited to join the
organization later.
`
		short = "Create an organization"
		usage = "create [name]"
	)

	cmd := command.New(usage, short, long, runCreate,
		command.RequireSession)

	flag.Add(cmd,
		flag.Bool{
			Name:        "apps-v2-default-on",
			Description: "Configure this org to use apps v2 by default for new apps",
			Default:     false,
		},
	)
	cmd.Args = cobra.MaximumNArgs(1)

	flag.Add(cmd, flag.JSONOutput())
	return cmd
}

func runCreate(ctx context.Context) error {
	name, err := nameFromFirstArgOrPrompt(ctx)
	if err != nil {
		return err
	}

	client := client.FromContext(ctx).API()

	var org *api.Organization
	if flag.GetBool(ctx, "apps-v2-default-on") {
		org, err = client.CreateOrganizationWithAppsV2DefaultOn(ctx, name)
	} else {
		org, err = client.CreateOrganization(ctx, name)
	}
	if err != nil {
		return fmt.Errorf("failed creating organization: %w", err)
	}

	if io := iostreams.FromContext(ctx); config.FromContext(ctx).JSONOutput {
		_ = render.JSON(io.Out, org)
	} else {
		printOrg(io.Out, org, true)
	}

	return nil
}

func nameFromFirstArgOrPrompt(ctx context.Context) (name string, err error) {
	if name = flag.FirstArg(ctx); name != "" {
		return
	}

	const msg = "Enter Organization Name:"

	if err = prompt.String(ctx, &name, msg, "", true); prompt.IsNonInteractive(err) {
		err = prompt.NonInteractiveError("name argument must be specified when not running interactively")
	}

	return
}
