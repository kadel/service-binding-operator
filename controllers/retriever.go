package controllers

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/redhat-developer/service-binding-operator/api/v1alpha1"
	"github.com/redhat-developer/service-binding-operator/pkg/envvars"
	"github.com/redhat-developer/service-binding-operator/pkg/log"
)

// retriever reads all data referred in plan instance, and store in a secret.
type retriever struct {
	logger *log.Log          // logger instance
	client dynamic.Interface // Kubernetes API client
}

// createServiceIndexPath returns a string slice with fields representing a path to a resource in the
// environment variable context. This function cleans fields that might contain invalid characters to
// be used in Go template; for example, a Group might contain the "." character, which makes it
// harder to refer using Go template direct accessors and is substituted by an underbar "_".
func createServiceIndexPath(name string, gvk schema.GroupVersionKind) []string {
	return []string{
		gvk.Version,
		strings.ReplaceAll(gvk.Group, ".", "_"),
		gvk.Kind,
		strings.ReplaceAll(name, "-", "_"),
	}

}

func buildServiceEnvVars(svcCtx *serviceContext, globalNamePrefix string) (map[string]string, error) {
	prefixes := []string{}
	if len(globalNamePrefix) > 0 {
		prefixes = append(prefixes, globalNamePrefix)
	}
	if svcCtx.namePrefix != nil && len(*svcCtx.namePrefix) > 0 {
		prefixes = append(prefixes, *svcCtx.namePrefix)
	}
	if svcCtx.namePrefix == nil {
		prefixes = append(prefixes, svcCtx.service.GroupVersionKind().Kind)
	}

	return envvars.Build(svcCtx.envVars, prefixes...)
}

func (r *retriever) processServiceContext(
	svcCtx *serviceContext,
	mappingsCtx map[string]interface{},
	globalNamePrefix string,
) (map[string][]byte, error) {
	svcEnvVars, err := buildServiceEnvVars(svcCtx, globalNamePrefix)
	if err != nil {
		return nil, err
	}

	// contribute the entire resource to the context shared with the custom env parser
	gvk := svcCtx.service.GetObjectKind().GroupVersionKind()

	// add an entry in the custom environment variable context, allowing the user to use the
	// following expression:
	//
	// `{{ index "v1alpha1" "postgresql.baiju.dev" "Database", "db-testing", "status", "connectionUrl" }}`
	err = unstructured.SetNestedField(
		mappingsCtx, svcCtx.service.Object, gvk.Version, gvk.Group, gvk.Kind,
		svcCtx.service.GetName())
	if err != nil {
		return nil, err
	}

	// add an entry in the custom environment variable context with modified key names (group
	// names have the "." separator changed to underbar and "-" in the resource name is changed
	// to underbar "_" as well).
	//
	// `{{ .v1alpha1.postgresql_baiju_dev.Database.db_testing.status.connectionUrl }}`
	err = unstructured.SetNestedField(
		mappingsCtx,
		svcCtx.service.Object,
		createServiceIndexPath(svcCtx.service.GetName(), svcCtx.service.GroupVersionKind())...,
	)
	if err != nil {
		return nil, err
	}

	// add an entry in the custom environment variable context with the informed 'id'.
	//
	// `{{ .db_testing.status.connectionUrl }}`
	if svcCtx.id != nil {
		err = unstructured.SetNestedField(
			mappingsCtx,
			svcCtx.service.Object,
			*svcCtx.id,
		)
		if err != nil {
			return nil, err
		}
	}

	envVars := make(map[string][]byte, len(svcEnvVars))
	for k, v := range svcEnvVars {
		envVars[k] = []byte(v)
	}

	return envVars, nil
}

// ProcessServiceContexts returns environment variables and volume keys from a ServiceContext slice.
func (r *retriever) ProcessServiceContexts(
	globalNamePrefix string,
	svcCtxs serviceContextList,
	envVarTemplates []v1alpha1.Mapping,
) (map[string][]byte, error) {
	mappingsCtx := make(map[string]interface{})
	envVars := make(map[string][]byte)

	for _, svcCtx := range svcCtxs {
		s, err := r.processServiceContext(svcCtx, mappingsCtx, globalNamePrefix)
		if err != nil {
			return nil, err
		}
		for k, v := range s {
			envVars[k] = []byte(v)
		}
	}

	envParser := newMappingsParser(envVarTemplates, mappingsCtx)
	mappingsList, err := envParser.Parse()
	if err != nil {
		r.logger.Error(
			err, "Creating envVars", "Templates", envVarTemplates, "TemplateContext", mappingsCtx)
		return nil, err
	}

	for k, v := range mappingsList {
		envVars[k] = []byte(v.(string))
	}

	return envVars, nil
}

// NewRetriever instantiate a new retriever instance.
func NewRetriever(
	client dynamic.Interface,
) *retriever {
	return &retriever{
		logger: log.NewLog("retriever"),
		client: client,
	}
}
