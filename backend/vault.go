package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"io/ioutil"
	"k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

type Vault struct {
	kubeClient *kubernetes.Clientset
}

type secret struct {
	Name string            `json:"name"`
	Data map[string][]byte `json:"data"`
}
type secretListItem struct {
	Name    string `json:"name"`
	Objects int    `json:"objects"`
}

// Init - Init function for initializing connection with kubernetes
func (m *Vault) Init(kubeconfig string, debug bool) error {
	var config *rest.Config

	if debug {
		fileInfo, err := os.Stat(kubeconfig)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			kubeconfig = filepath.Join(kubeconfig, ".kube", "config")
		}

		conf, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return err
		}
		config = conf
	} else {
		conf, err := rest.InClusterConfig()
		if err != nil {
			return errors.New(fmt.Sprintf("cannot reach kubernetes API: %v", err))
		}
		config = conf
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot create k8s client: %v", err))
	}

	m.kubeClient = kubeClient

	return nil
}

// NamespacesList - NamespacesList handler which return a list namespaces
func (m *Vault) NamespacesList(w http.ResponseWriter, r *http.Request) {
	n := m.kubeClient.CoreV1().Namespaces()
	list, err := n.List(metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var namespaces = []string{}
	for _, item := range list.Items {
		namespaces = append(namespaces, item.Name)
	}

	b, err := json.Marshal(namespaces)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// Get - Get handler which return a list of secrets in selected namespaces
func (m *Vault) Get(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	_, err := m.kubeClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if k8sErrors.IsNotFound(err) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	n := m.kubeClient.CoreV1().Secrets(namespace)
	list, err := n.List(metav1.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var secretsList = []secretListItem{}
	for _, item := range list.Items {
		secretsList = append(secretsList, secretListItem{
			Name:    item.Name,
			Objects: len(item.Data),
		})
	}

	b, err := json.Marshal(secretsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// Find - Find handler which find and return selected secret from selected namespace
func (m *Vault) Find(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	_, err := m.kubeClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if k8sErrors.IsNotFound(err) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secretName := chi.URLParam(r, "secret")

	s, err := m.kubeClient.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	if k8sErrors.IsNotFound(err) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var secret = secret{
		Name: s.Name,
		Data: s.Data,
	}

	b, err := json.Marshal(secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// Add - Add handler which add new secret to selected namespace
func (m *Vault) Add(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	_, err := m.kubeClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if k8sErrors.IsNotFound(err) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	secretName := chi.URLParam(r, "secret")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var secretData = map[string]string{}
	err = json.Unmarshal(body, &secretData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data = map[string][]byte{}

	const configMapKeyFmt = `[-._a-zA-Z0-9]+`
	var configMapKeyRegexp = regexp.MustCompile("^" + configMapKeyFmt + "$")

	for key, value := range secretData {
		if !configMapKeyRegexp.MatchString(key) {
			http.Error(w, "secret name must consist of alphanumeric characters, '-', '_' or '.'", http.StatusBadRequest)
			return
		}

		var stringValue = []byte(value)
		v, err := base64.StdEncoding.DecodeString(value)
		if err == nil {
			stringValue = v
		}

		data[key] = []byte(stringValue)
	}

	newSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: data,
	}

	var s *v1.Secret
	s, err = m.kubeClient.CoreV1().Secrets(namespace).Create(newSecret)
	if k8sErrors.IsAlreadyExists(err) {
		s, err = m.kubeClient.CoreV1().Secrets(namespace).Update(newSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(secret{
		Name: s.Name,
		Data: s.Data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

// Delete - Delete handler which find and delete selected secret from selected namespace
func (m *Vault) Delete(w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	_, err := m.kubeClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if k8sErrors.IsNotFound(err) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	secretName := chi.URLParam(r, "secret")

	err = m.kubeClient.CoreV1().Secrets(namespace).Delete(secretName, &metav1.DeleteOptions{})

	if k8sErrors.IsNotFound(err) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
