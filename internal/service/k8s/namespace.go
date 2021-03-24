package k8ss

import (
	"context"
	"go-xops/assets/k8s"
	"go-xops/pkg/common"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNameSpaces(ctx context.Context) ([]k8s.NameSpace, error) {
	result := []k8s.NameSpace{}
	data, err := common.ClientSet.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range data.Items {
		result = append(result, k8s.NameSpace{
			Name:   namespace.ObjectMeta.Name,
			UID:    string(namespace.ObjectMeta.UID),
			Status: strings.ToLower(string(namespace.Status.Phase)),
		})
	}
	return result, nil
}

func GetNamespace(ctx context.Context, name string) (k8s.NameSpace, error) {
	result := k8s.NameSpace{}
	data, err := common.ClientSet.CoreV1().Namespaces().Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return result, nil
	}

	result.Name = data.ObjectMeta.Name
	result.Status = strings.ToLower(string(data.Status.Phase))
	result.UID = string(data.ObjectMeta.UID)
	return result, nil
}
