package v3public

import (
	"context"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	AuthProviderGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "AuthProvider",
	}
	AuthProviderResource = metav1.APIResource{
		Name:         "authproviders",
		SingularName: "authprovider",
		Namespaced:   false,
		Kind:         AuthProviderGroupVersionKind.Kind,
	}
)

type AuthProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuthProvider
}

type AuthProviderHandlerFunc func(key string, obj *AuthProvider) (*AuthProvider, error)

type AuthProviderLister interface {
	List(namespace string, selector labels.Selector) (ret []*AuthProvider, err error)
	Get(namespace, name string) (*AuthProvider, error)
}

type AuthProviderController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() AuthProviderLister
	AddHandler(ctx context.Context, name string, handler AuthProviderHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler AuthProviderHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type AuthProviderInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*AuthProvider) (*AuthProvider, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*AuthProvider, error)
	Get(name string, opts metav1.GetOptions) (*AuthProvider, error)
	Update(*AuthProvider) (*AuthProvider, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*AuthProviderList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() AuthProviderController
	AddHandler(ctx context.Context, name string, sync AuthProviderHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle AuthProviderLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync AuthProviderHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle AuthProviderLifecycle)
}

type authProviderLister struct {
	controller *authProviderController
}

func (l *authProviderLister) List(namespace string, selector labels.Selector) (ret []*AuthProvider, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*AuthProvider))
	})
	return
}

func (l *authProviderLister) Get(namespace, name string) (*AuthProvider, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    AuthProviderGroupVersionKind.Group,
			Resource: "authProvider",
		}, key)
	}
	return obj.(*AuthProvider), nil
}

type authProviderController struct {
	controller.GenericController
}

func (c *authProviderController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *authProviderController) Lister() AuthProviderLister {
	return &authProviderLister{
		controller: c,
	}
}

func (c *authProviderController) AddHandler(ctx context.Context, name string, handler AuthProviderHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*AuthProvider); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *authProviderController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler AuthProviderHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*AuthProvider); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type authProviderFactory struct {
}

func (c authProviderFactory) Object() runtime.Object {
	return &AuthProvider{}
}

func (c authProviderFactory) List() runtime.Object {
	return &AuthProviderList{}
}

func (s *authProviderClient) Controller() AuthProviderController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.authProviderControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(AuthProviderGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &authProviderController{
		GenericController: genericController,
	}

	s.client.authProviderControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type authProviderClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   AuthProviderController
}

func (s *authProviderClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *authProviderClient) Create(o *AuthProvider) (*AuthProvider, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*AuthProvider), err
}

func (s *authProviderClient) Get(name string, opts metav1.GetOptions) (*AuthProvider, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*AuthProvider), err
}

func (s *authProviderClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*AuthProvider, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*AuthProvider), err
}

func (s *authProviderClient) Update(o *AuthProvider) (*AuthProvider, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*AuthProvider), err
}

func (s *authProviderClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *authProviderClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *authProviderClient) List(opts metav1.ListOptions) (*AuthProviderList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*AuthProviderList), err
}

func (s *authProviderClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *authProviderClient) Patch(o *AuthProvider, data []byte, subresources ...string) (*AuthProvider, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*AuthProvider), err
}

func (s *authProviderClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *authProviderClient) AddHandler(ctx context.Context, name string, sync AuthProviderHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *authProviderClient) AddLifecycle(ctx context.Context, name string, lifecycle AuthProviderLifecycle) {
	sync := NewAuthProviderLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *authProviderClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync AuthProviderHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *authProviderClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle AuthProviderLifecycle) {
	sync := NewAuthProviderLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}
