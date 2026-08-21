package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addTemplateContentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadGetCommand(ctx, root, options, readGetSpec[client.TemplateDetail]{
		Use:   "content TEMPLATE_ID",
		Short: "Get exact template data and scope",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.TemplateDetail, error) {
			return client.GetTemplateDetail(ctx, runtime.graphqlClient, id)
		},
		Write: writeTemplateDetail,
	})
}

func addTemplateCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	flags := templateWriteFlags{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.TemplateDetail]{
		Use:   "create",
		Short: "Create an issue template in the pinned team",
		Args:  cobra.NoArgs,
		Configure: func(command *cobra.Command) {
			bindTemplateCreateFlags(command, &flags)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, _ []string,
		) (client.TemplateDetail, error) {
			data, err := readJSONObjectFile(flags.DataFile)
			if err != nil {
				return client.TemplateDetail{}, err
			}

			return client.CreateTemplate(ctx, runtime.graphqlClient, runtime.config.Target,
				client.TemplateCreateRequest{
					ID:   flags.ID,
					Name: flags.Name,
					Type: flags.Type,
					Data: data,
				})
		},
		Write: writeTemplateDetail,
	})
}

func addTemplateUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	flags := templateWriteFlags{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.TemplateDetail]{
		Use:   "update TEMPLATE_ID",
		Short: "Update an issue template after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			bindTemplateUpdateFlags(command, &flags)
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.TemplateDetail, error) {
			request, err := templateUpdateRequest(command, flags, args[0])
			if err != nil {
				return client.TemplateDetail{}, err
			}

			return client.UpdateTemplate(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeTemplateDetail,
	})
}

type templateWriteFlags struct {
	ID       string
	Name     string
	Type     string
	DataFile string
}

func bindTemplateCreateFlags(command *cobra.Command, flags *templateWriteFlags) {
	command.Flags().StringVar(&flags.ID, "id", "", "caller-supplied template UUID v4")
	command.Flags().StringVar(&flags.Name, "name", "", "template name")
	command.Flags().StringVar(&flags.Type, "type", "", "template type; only issue is accepted")
	command.Flags().StringVar(&flags.DataFile, "data-file", "", "local JSON object file for template data")
}

func bindTemplateUpdateFlags(command *cobra.Command, flags *templateWriteFlags) {
	command.Flags().StringVar(&flags.Name, "name", "", "new template name")
	command.Flags().StringVar(&flags.DataFile, "data-file", "", "local JSON object file for template data")
}

func templateUpdateRequest(
	command *cobra.Command,
	flags templateWriteFlags,
	id string,
) (client.TemplateUpdateRequest, error) {
	request := client.TemplateUpdateRequest{ID: id}
	if command.Flags().Changed("name") {
		request.Name = &flags.Name
	}
	if !command.Flags().Changed("data-file") {
		return request, nil
	}
	data, err := readJSONObjectFile(flags.DataFile)
	if err != nil {
		return client.TemplateUpdateRequest{}, err
	}
	request.Data = data

	return request, nil
}

func writeTemplateDetail(command *cobra.Command, options *rootOptions, template client.TemplateDetail) error {
	return writeItemLine(
		command, options, template, template.ID,
		"%s %s [%s] %s",
		template.ID, template.Name, template.Type, templateScopeLabel(template.TeamKey),
	)
}
