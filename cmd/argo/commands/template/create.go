package template

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/argoproj/argo-workflows/v3/cmd/argo/commands/client"
	"github.com/argoproj/argo-workflows/v3/cmd/argo/commands/common"
	workflowtemplatepkg "github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflowtemplate"
)

type cliCreateOpts struct {
	output common.EnumFlagValue // --output
	strict bool                 // --strict
}

func NewCreateCommand() *cobra.Command {
	var cliCreateOpts = cliCreateOpts{output: common.NewPrintWorkflowOutputValue("")}
	command := &cobra.Command{
		Use:   "create FILE1 FILE2...",
		Short: "create a workflow template",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return CreateWorkflowTemplates(cmd.Context(), args, &cliCreateOpts)
		},
	}
	command.Flags().VarP(&cliCreateOpts.output, "output", "o", "Output format. "+cliCreateOpts.output.Usage())
	command.Flags().BoolVar(&cliCreateOpts.strict, "strict", true, "perform strict workflow validation")
	return command
}

func CreateWorkflowTemplates(ctx context.Context, filePaths []string, cliOpts *cliCreateOpts) error {
	if cliOpts == nil {
		cliOpts = &cliCreateOpts{}
	}
	ctx, apiClient, err := client.NewAPIClient(ctx)
	if err != nil {
		return err
	}
	serviceClient, err := apiClient.NewWorkflowTemplateServiceClient()
	if err != nil {
		return err
	}

	workflowTemplates := generateWorkflowTemplates(ctx, filePaths, cliOpts.strict)

	for _, wftmpl := range workflowTemplates {
		if wftmpl.Namespace == "" {
			wftmpl.Namespace = client.Namespace(ctx)
		}
		created, err := serviceClient.CreateWorkflowTemplate(ctx, &workflowtemplatepkg.WorkflowTemplateCreateRequest{
			Namespace: wftmpl.Namespace,
			Template:  &wftmpl,
		})
		if err != nil {
			return fmt.Errorf("failed to create workflow template: %v", err)
		}
		printWorkflowTemplate(created, cliOpts.output.String())
	}
	return nil
}
