package info

import (
	"context"
	"fmt"

	"github.com/superfly/flyctl/api"
	"github.com/superfly/flyctl/client"
	"github.com/superfly/flyctl/internal/command/apps"
	"github.com/superfly/flyctl/internal/command/services"
	"github.com/superfly/flyctl/internal/flag"
	"github.com/superfly/flyctl/internal/render"
	"github.com/superfly/flyctl/iostreams"
)

func showMachineInfo(ctx context.Context, appName string) error {
	var (
		client    = client.FromContext(ctx).API()
		jsonOuput = flag.GetBool(ctx, "json")
	)

	if jsonOuput {
		return fmt.Errorf("outputting to json is not yet supported")
	}

	appInfo, err := client.GetAppInfo(ctx, appName)
	if err != nil {
		return err
	}

	app, err := client.GetAppCompact(ctx, appName)
	if err != nil {
		return err
	}

	ctx, err = apps.BuildContext(ctx, app)
	if err != nil {
		return err
	}

	if err := showMachineAppInfo(ctx, app); err != nil {
		return err
	}

	if err := services.ShowMachineServiceInfo(ctx, appInfo); err != nil {
		return err
	}

	if err := showMachineIPInfo(ctx, appName); err != nil {
		return err
	}

	return nil
}

func showMachineAppInfo(ctx context.Context, app *api.AppCompact) error {
	var (
		io = iostreams.FromContext(ctx)
	)
	rows := [][]string{
		{
			app.Name,
			app.Organization.Slug,
			app.PlatformVersion,
			app.Hostname,
		},
	}
	var cols = []string{"Name", "Owner", "Platform", "Hostname"}

	if err := render.VerticalTable(io.Out, "App", rows, cols...); err != nil {
		return err
	}

	return nil
}

func showMachineIPInfo(ctx context.Context, appName string) error {
	var (
		io     = iostreams.FromContext(ctx)
		client = client.FromContext(ctx).API()
	)

	info, err := client.GetAppInfo(ctx, appName)
	if err != nil {
		return err
	}

	ips := [][]string{}

	for _, ip := range info.IPAddresses.Nodes {
		fields := []string{
			ip.Type,
			ip.Address,
			ip.Region,
			formatRelativeTime(ip.CreatedAt),
		}
		ips = append(ips, fields)
	}

	_ = render.Table(io.Out, "IP Addresses", ips, "Type", "Address", "Region", "Created at")

	return nil
}
