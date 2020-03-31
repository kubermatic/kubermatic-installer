package kubermatic

import (
	"errors"
	"fmt"

	"github.com/kubermatic/kubermatic-installer/pkg/shared/operatorv1alpha1"
	"github.com/kubermatic/kubermatic-installer/pkg/yamled"
	"github.com/sirupsen/logrus"
)

func ValidateConfiguration(config *operatorv1alpha1.KubermaticConfiguration, helmValues *yamled.Document, logger logrus.FieldLogger) (*operatorv1alpha1.KubermaticConfiguration, *yamled.Document, []error) {
	kubermaticFailures := validateKubermaticConfiguration(config)
	for idx, e := range kubermaticFailures {
		kubermaticFailures[idx] = prefixError("KubermaticConfiguration: ", e)
	}

	helmFailures := validateHelmValues(config, helmValues, logger)
	for idx, e := range helmFailures {
		helmFailures[idx] = prefixError("Helm values: ", e)
	}

	return config, helmValues, append(kubermaticFailures, helmFailures...)
}

func validateKubermaticConfiguration(config *operatorv1alpha1.KubermaticConfiguration) []error {
	failures := []error{}

	if config.Spec.Ingress.Domain == "" {
		failures = append(failures, errors.New("spec.ingress.domain cannot be left empty"))
	}

	return failures
}

func validateHelmValues(config *operatorv1alpha1.KubermaticConfiguration, helmValues *yamled.Document, logger logrus.FieldLogger) []error {
	failures := []error{}

	if domain, _ := helmValues.GetString(yamled.Path{"dex", "ingress", "host"}); domain == "" {
		logger.WithField("domain", config.Spec.Ingress.Domain).Warn("dex.ingress.host is empty, setting to spec.ingress.domain from KubermaticConfiguration")
		helmValues.Set(yamled.Path{"dex", "ingress", "host"}, config.Spec.Ingress.Domain)
	}

	return failures
}

func prefixError(prefix string, e error) error {
	return fmt.Errorf("%s%v", prefix, e)
}
