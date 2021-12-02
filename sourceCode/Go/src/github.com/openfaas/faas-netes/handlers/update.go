
package handlers

import (
	"fmt"
	"github.com/openfaas/faas-netes/gpu/repository"
	"github.com/openfaas/faas-netes/k8s"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"net/http"
	"time"

	ptypes "github.com/openfaas/faas-provider/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MakeUpdateHandler update specified function
func MakeUpdateHandler(defaultNamespace string, factory k8s.FunctionFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*
		if r.Body != nil {
			defer r.Body.Close()
		}

		body, _ := ioutil.ReadAll(r.Body)

		request := ptypes.FunctionDeployment{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			wrappedErr := fmt.Errorf("update: unable to unmarshal request: %s", err.Error())
			http.Error(w, wrappedErr.Error(), http.StatusBadRequest)
			return
		}

		lookupNamespace := defaultNamespace
		if len(request.Namespace) > 0 {
			lookupNamespace = request.Namespace
		}

		if lookupNamespace == "kube-system" {
			http.Error(w, "update: unable to list within the kube-system namespace", http.StatusUnauthorized)
			return
		}

		annotations := buildAnnotations(request)
		if err, status := updatePodSpec(lookupNamespace, factory, request, annotations); err != nil {
			if !k8s.IsNotFound(err) {
				log.Printf("update: error updating deployment= %s.%s, error: %s \n", request.Service, lookupNamespace, err)
				return
			}
			wrappedErr := fmt.Errorf("update: unable update Deployment= %s.%s, error: %s \n", request.Service, lookupNamespace, err.Error())
			http.Error(w, wrappedErr.Error(), status)
			return
		}

		if err, status := updateService(lookupNamespace, factory, request, annotations); err != nil {
			if !k8s.IsNotFound(err) {
				log.Printf("update: error updating service= %s.%s, error: %s \n", request.Service, lookupNamespace, err)
			}
			wrappedErr := fmt.Errorf("update: unable update Service= %s.%s, error: %s \n", request.Service, request.Namespace, err.Error())
			http.Error(w, wrappedErr.Error(), status)
			return
		}
*/
		w.WriteHeader(http.StatusAccepted)
	}
}

func updatePodSpec(
	functionNamespace string,
	factory k8s.FunctionFactory,
	request ptypes.FunctionDeployment,
	annotations map[string]string) (err error, httpStatus int) {

	labelPod := labels.SelectorFromSet(map[string]string{"faas_function": request.Service})
	listPodOptions := metav1.ListOptions {
		LabelSelector: labelPod.String(),
	}
	// This makes sure we don't delete non-labeled deployments
	podList, findPodsErr := factory.Client.CoreV1().Pods(functionNamespace).List(listPodOptions)
	if findPodsErr != nil {
		return findPodsErr, http.StatusNotFound
	}

	for i := 0; i < len(podList.Items); i++ {
		if len(podList.Items[i].Spec.Containers) > 0 {
			podList.Items[i].Spec.Containers[0].Image = request.Image

			// Disabling update support to prevent unexpected mutations of deployed functions,
			// since imagePullPolicy is now configurable. This could be reconsidered later depending
			// on desired behavior, but will need to be updated to take config.
			//deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy = v1.PullAlways

			podList.Items[i].Spec.Containers[0].Env = buildEnvVars(&request)

			factory.ConfigureReadOnlyRootFilesystem(request, &podList.Items[i])
			factory.ConfigureContainerUserID(&podList.Items[i])

			podList.Items[i].Spec.NodeSelector = createSelector(request.Constraints)

			label := map[string]string {
				"faas_function": request.Service,
				"uid": fmt.Sprintf("%d", time.Now().Nanosecond()),
			}

			if request.Labels != nil {
				for k, v := range *request.Labels {
					label[k] = v
				}
			}

			// deployment.Labels = labels
			podList.Items[i].ObjectMeta.Labels = label
			podList.Items[i].ObjectMeta.Annotations = annotations

			resources, resourceErr := createResources(request)
			if resourceErr != nil {
				return resourceErr, http.StatusBadRequest
			}

			podList.Items[i].Spec.Containers[0].Resources = *resources

			var serviceAccount string

			if request.Annotations != nil {
				annotations := *request.Annotations
				if val, ok := annotations["com.openfaas.serviceaccount"]; ok && len(val) > 0 {
					serviceAccount = val
				}
			}
			podList.Items[i].Spec.ServiceAccountName = serviceAccount

			secrets := k8s.NewSecretsClient(factory.Client)
			existingSecrets, err := secrets.GetSecrets(functionNamespace, request.Secrets)
			if err != nil {
				return err, http.StatusBadRequest
			}

			err = factory.ConfigureSecrets(request, &podList.Items[i], existingSecrets)
			if err != nil {
				log.Println(err)
				return err, http.StatusBadRequest
			}

			probes, err := factory.MakeProbes(request)
			if err != nil {
				return err, http.StatusBadRequest
			}

			podList.Items[i].Spec.Containers[0].LivenessProbe = probes.Liveness
			podList.Items[i].Spec.Containers[0].ReadinessProbe = probes.Readiness
		}

		if _, updateErr := factory.Client.CoreV1().
			Pods(functionNamespace).
			Update(&podList.Items[i]); updateErr != nil {

			return updateErr, http.StatusInternalServerError
		}
	}
	if len(podList.Items) > 0 {
		repository.UpdateFuncSpec(request.Service, &podList.Items[0],nil)
		repository.UpdateFuncConstrains(request.Service, request.Constraints)
		repository.UpdateFuncRequestResources(request.Service, request.Requests)
	}
	return nil, http.StatusAccepted
}

func updateService(functionNamespace string, factory k8s.FunctionFactory, request ptypes.FunctionDeployment, annotations map[string]string) (err error, httpStatus int) {
	getOpts := metav1.GetOptions{}

	service, findServiceErr := factory.Client.CoreV1().Services(functionNamespace).Get(request.Service, getOpts)
	if findServiceErr != nil {
		return findServiceErr, http.StatusNotFound
	}
	if service.Name == "" {
		return nil, http.StatusAccepted
	}
	service.Annotations = annotations

	if _, updateErr := factory.Client.CoreV1().Services(functionNamespace).Update(service);
	updateErr != nil {
		return updateErr, http.StatusInternalServerError
	}

	repository.UpdateFuncSpec(request.Service, nil, service)
	return nil, http.StatusAccepted
}
