package component

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	// "github.com/openshift/odo/pkg/odo/util/completion"
	"github.com/redhat-developer/odo-fork/pkg/component"
	"github.com/redhat-developer/odo-fork/pkg/config"
	odoutil "github.com/redhat-developer/odo-fork/pkg/kdo/util"

	// appCmd "github.com/redhat-developer/odo-fork/pkg/kdo/cli/application"
	"github.com/redhat-developer/odo-fork/pkg/kclient"
	projectCmd "github.com/redhat-developer/odo-fork/pkg/kdo/cli/project"
	"github.com/redhat-developer/odo-fork/pkg/kdo/cli/ui"
	"github.com/redhat-developer/odo-fork/pkg/kdo/genericclioptions"
	"github.com/redhat-developer/odo-fork/pkg/log"
	"github.com/redhat-developer/odo-fork/pkg/url"

	ktemplates "k8s.io/kubectl/pkg/util/templates"
)

// DeleteRecommendedCommandName is the recommended delete command name
const DeleteRecommendedCommandName = "delete"

var deleteExample = ktemplates.Examples(`  # Delete component named 'frontend'. 
%[1]s frontend
%[1]s frontend --all
  `)

// DeleteOptions is a container to attach complete, validate and run pattern
type DeleteOptions struct {
	componentForceDeleteFlag bool
	componentDeleteAllFlag   bool
	componentContext         string
	*CommonPushOptions
}

// NewDeleteOptions returns new instance of DeleteOptions
func NewDeleteOptions() *DeleteOptions {
	return &DeleteOptions{false, false, "", &CommonPushOptions{}}
}

// Complete completes log args
func (do *DeleteOptions) Complete(name string, cmd *cobra.Command, args []string) (err error) {
	err = do.CommonPushOptions.Complete(name, cmd, args)
	return
}

// Validate validates the list parameters
func (do *DeleteOptions) Validate() (err error) {
	if do.Context.Namespace == "" || do.Application == "" {
		return odoutil.ThrowContextError()
	}
	isExists, err := component.Exists(do.Client, do.componentName, do.Application)
	if err != nil {
		return err
	}
	if !isExists {
		return fmt.Errorf("failed to delete component %s as it doesn't exist", do.componentName)
	}
	return
}

// Run has the logic to perform the required actions as part of command
func (do *DeleteOptions) Run() (err error) {
	glog.V(4).Infof("component delete called")
	glog.V(4).Infof("args: %#v", do)

	err = printDeleteComponentInfo(do.Client, do.componentName, do.Context.Application, do.Context.Namespace)
	if err != nil {
		return err
	}

	if do.componentForceDeleteFlag || ui.Proceed(fmt.Sprintf("Are you sure you want to delete %v from %v?", do.componentName, do.Application)) {
		err := component.Delete(do.Client, do.componentName, do.Application)
		if err != nil {
			return err
		}
		log.Successf("Component %s from application %s has been deleted", do.componentName, do.Application)

	} else {
		return fmt.Errorf("Aborting deletion of component: %v", do.componentName)
	}

	if do.componentDeleteAllFlag {
		if do.componentForceDeleteFlag || ui.Proceed(fmt.Sprintf("Are you sure you want to delete local config for %v?", do.componentName)) {
			cfg, err := config.NewLocalConfigInfo(do.componentContext)
			if err != nil {
				return err
			}

			err = cfg.DeleteConfigDir()
			if err != nil {
				return err
			}

			log.Successf("Config for the Component %s has been deleted", do.componentName)
		} else {
			return fmt.Errorf("Aborting deletion of config for component: %s", do.componentName)
		}
	}

	return
}

// NewCmdDelete implements the delete odo command
func NewCmdDelete(name, fullName string) *cobra.Command {

	do := NewDeleteOptions()

	var componentDeleteCmd = &cobra.Command{
		Use:     fmt.Sprintf("%s <component_name>", name),
		Short:   "Delete an existing component",
		Long:    "Delete an existing component.",
		Example: fmt.Sprintf(deleteExample, fullName),
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			genericclioptions.GenericRun(do, cmd, args)
		},
	}

	componentDeleteCmd.Flags().BoolVarP(&do.componentForceDeleteFlag, "force", "f", false, "Delete component without prompting")
	componentDeleteCmd.Flags().BoolVarP(&do.componentDeleteAllFlag, "all", "a", false, "Delete component and local config")

	// Add a defined annotation in order to appear in the help menu
	componentDeleteCmd.Annotations = map[string]string{"command": "component"}
	componentDeleteCmd.SetUsageTemplate(odoutil.CmdUsageTemplate)
	// completion.RegisterCommandHandler(componentDeleteCmd, completion.ComponentNameCompletionHandler)
	//Adding `--context` flag
	genericclioptions.AddContextFlag(componentDeleteCmd, &do.componentContext)

	//Adding `--project` flag
	projectCmd.AddProjectFlag(componentDeleteCmd)
	//Adding `--application` flag
	// appCmd.AddApplicationFlag(componentDeleteCmd)

	return componentDeleteCmd
}

func printDeleteComponentInfo(client *kclient.Client, componentName string, appName string, projectName string) error {
	componentDesc, err := component.GetComponent(client, componentName, appName, projectName)
	if err != nil {
		return errors.Wrap(err, "unable to get component description")
	}

	if len(componentDesc.Spec.URL) != 0 {
		log.Info("This component has following urls that will be deleted with component")
		ul, err := url.List(client, componentDesc.Name, appName)
		if err != nil {
			return errors.Wrap(err, "Could not get url list")
		}
		for _, u := range ul.Items {
			log.Info("URL named", u.GetName(), "with host", u.Spec.Host, "having protocol", u.Spec.Protocol, "at port", u.Spec.Port)
		}
	}

	// storages, err := storage.List(client, componentDesc.Name, appName)
	// odoutil.LogErrorAndExit(err, "")
	// if len(storages.Items) != 0 {
	// 	log.Info("This component has following storages which will be deleted with the component")
	// 	for _, storageName := range componentDesc.Spec.Storage {
	// 		store := storages.Get(storageName)
	// 		log.Info("Storage", store.GetName(), "of size", store.Spec.Size)
	// 	}
	// }
	return nil
}
