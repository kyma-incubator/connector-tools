package controller

import (
	"context"
	"fmt"
	"github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/compass"
	errorWrap "github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/error"
	"github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/open_connectors"
	log "github.com/sirupsen/logrus"
	"sync"
)



//Controller controls the synchronization between the Open Connectors Tenant and The Compass Registry
type Controller struct {
	compass        compass.Connector
	openConnectors open_connectors.Connector
	tags           []string
	namePrefix     string
}

func New(context context.Context, compass compass.Connector, openConnectors open_connectors.Connector,
	tags []string, namePrefix string) (*Controller, error) {
	log.Debugf("creating new controller instance with tags %v", tags)

	return &Controller{
		compass:        compass,
		openConnectors: openConnectors,
		tags:           tags,
		namePrefix:     namePrefix,
	}, nil
}

func (c *Controller) DetermineStatus(ctx context.Context) (instancesToAdd []open_connectors.Instance,
	existingApps []compass.Application, appsToDelete []compass.Application, err error) {

	log.Debug("comparing Open Connectors connector instances to compass applications")

	openConnectorsInstances, err := c.openConnectors.GetConnectorInstances(ctx, c.tags)
	if err != nil {
		log.Errorf("error receiving open connectors instances: %s", err.Error())

		//check if error is properly wrapped
		if _, ok := err.(errorWrap.Error); !ok {
			return nil, nil,nil, errorWrap.WrapError(err, "error receiving open "+
				"connectors instances")
		}
		return nil, nil, nil, err
	}

	log.Tracef("received %d instances from open connectors", len(openConnectorsInstances))

	openConnectorsContext, err := c.openConnectors.GetOpenConnectorsContext(ctx)
	if err != nil {
		log.Errorf("error reading open connectors context: %s", err.Error())

		if _, ok := err.(errorWrap.Error); !ok {
			return nil, nil,nil, errorWrap.WrapError(err, "error reading open connectors "+
				"context")
		}
		return nil, nil,nil, err
	}

	compassApplications, err := c.compass.GetApplications(ctx, openConnectorsContext)

	if err != nil {
		log.Errorf("error reading compass applications: %s", err.Error())

		if _, ok := err.(errorWrap.Error); !ok {
			return nil, nil,nil, errorWrap.WrapError(err, "error reading compass "+
				"applications")
		}
		return nil, nil, nil, err
	}

	//classifying instances, already existing ones get removed, ultimately leaving the ones to be deleted
	for i := 0; i < len(openConnectorsInstances); {
		exists := false
		for j := 0; j < len(compassApplications); j++ {

			//Application exists
			if openConnectorsInstances[i].ID == compassApplications[j].ConnectorInstanceID {
				existingApps = append(existingApps,compassApplications[j])
				openConnectorsInstances = append(openConnectorsInstances[:i], openConnectorsInstances[i+1:]...)
				compassApplications = append(compassApplications[:j], compassApplications[j+1:]...)
				exists = true
				break
			}
		}

		//if instance does not exist we increment (otherwise instance is already removed)
		if !exists {
			instancesToAdd = append(instancesToAdd, openConnectorsInstances[i])
			i++
		}
	}

	return instancesToAdd, existingApps, compassApplications, nil
}

func (c *Controller) createNewApplications(ctx context.Context, instancesToAdd []open_connectors.Instance) error {

	log.Debugf("creating %d new compass applications", len(instancesToAdd))

	var wg sync.WaitGroup

	openConnectorsContext, err := c.openConnectors.GetOpenConnectorsContext(ctx)
	if err != nil {
		log.Errorf("error reading open connectors context: %s", err.Error())

		if _, ok := err.(errorWrap.Error); !ok {
			return errorWrap.WrapError(err, "error reading open connectors "+
				"context")
		}
		return err
	}

	errors := make(chan error, len(instancesToAdd))
	defer close(errors)
	wg.Add(len(instancesToAdd))

	for i := range instancesToAdd {
		log.Debugf("adding connector instance %q (context %q)", instancesToAdd[i].ID, openConnectorsContext)

		go func(instance open_connectors.Instance) {
			defer wg.Done()
			log.Debugf("creating new compass application %q for instance (context %q)", instance.ID,
				openConnectorsContext)

			spec, err := c.openConnectors.GetOpenAPISpec(ctx, instance.ID, "")
			if err != nil {
				errors <- err
				return
			}

			appID, err := c.compass.CreateApplication(ctx,
				createApplicationName(c.namePrefix, &instance),
				createApplicationDescription(&instance),
				openConnectorsContext,
				instance.ID)
			if err != nil {
				errors <- err
				return
			}

			_, err = c.compass.CreateAPIForApplication(ctx,
				appID,
				fmt.Sprintf("%s: %s", instance.ConnectorName, instance.Name),
				instancesToAdd[i].ConnectorKey,
				" ",
				c.openConnectors.GetOpenConnectorsAPIURL(ctx),
				c.openConnectors.CreateAPIAuthorizationHeader(ctx, &instance),
				spec)

			if err != nil {
				errors <- err
				return
			}
			return
		}(instancesToAdd[i])

	}

	wg.Wait()
	for range instancesToAdd {
		select {
		case err = <-errors:
			log.Errorf("error creating compass application: %s", err.Error())
		default:
		}
	}

	return err
}

func (c *Controller) deleteApplications(ctx context.Context, appsToDelete []compass.Application) error {
	log.Debugf("deleting %d compass applications", len(appsToDelete))

	var wg sync.WaitGroup

	errors := make(chan error, len(appsToDelete))
	defer close(errors)
	wg.Add(len(appsToDelete))

	for i := range appsToDelete {
		log.Debugf("deleting compass application %q", appsToDelete[i].ID)

		go func(application compass.Application) {
			defer wg.Done()

			_, err := c.compass.DeleteApplication(ctx, application.ID)
			if err != nil {
				errors <- err
			}
			return

		}(appsToDelete[i])
	}

	wg.Wait()

	var err error
	for range appsToDelete {
		select {
		case err = <-errors:
			log.Errorf("error deleting compass application: %s", err.Error())
		default:
		}
	}

	return err
}

func (c *Controller) Synchronize(ctx context.Context) error {

	instancesToAdd, _, appsToDelete, err := c.DetermineStatus(ctx)
	if err != nil {
		log.Errorf("error comparing state between SAP Cloud Platform Open Connectors and Compass : %s",
			err.Error())
		return err
	}

	errors := make(chan error, 2)
	defer close(errors)

	var wg sync.WaitGroup

	//Add new applications if applicable
	if len(instancesToAdd) > 0 {
		wg.Add(1)
		go func(ctx context.Context, error <-chan error, instanceToAdd []open_connectors.Instance) {
			defer wg.Done()
			if err := c.createNewApplications(ctx, instancesToAdd); err != nil {
				errors <- err
			}
			return
		}(ctx, errors, instancesToAdd)
	}

	//Delete unused applications
	if len(appsToDelete) > 0 {
		wg.Add(1)
		go func(ctx context.Context, error <-chan error, appsToDelete []compass.Application) {
			defer wg.Done()
			if err := c.deleteApplications(ctx, appsToDelete); err != nil {
				errors <- err
			}
			return
		}(ctx, errors, appsToDelete)
	}
	wg.Wait()

	for i := 0; i < 2; i++ {
		select {
		case err = <-errors:
			log.Errorf("error adding or deleting application: %s", err.Error())
		default:
		}
	}

	return err
}

func createApplicationName(prefix string, instance *open_connectors.Instance) string {
	return fmt.Sprintf("%s-%s", prefix, instance.ID)
}

func createApplicationDescription(instance *open_connectors.Instance) string {
	return fmt.Sprintf("SAP Cloud Platform Open Connectors %s: %s", instance.ConnectorName, instance.Name)
}
