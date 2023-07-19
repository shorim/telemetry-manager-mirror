package gateway

import (
	"github.com/kyma-project/telemetry-manager/internal/otelcollector/config/builder/common"
)

func makeProcessorsConfig() ProcessorsConfig {
	k8sAttributes := []string{
		"k8s.pod.name",
		"k8s.node.name",
		"k8s.namespace.name",
		"k8s.deployment.name",
		"k8s.statefulset.name",
		"k8s.daemonset.name",
		"k8s.cronjob.name",
		"k8s.job.name",
	}

	podAssociations := []common.PodAssociations{
		{
			Sources: []common.PodAssociation{
				{
					From: "resource_attribute",
					Name: "k8s.pod.ip",
				},
			},
		},
		{
			Sources: []common.PodAssociation{
				{
					From: "resource_attribute",
					Name: "k8s.pod.uid",
				},
			},
		},
		{
			Sources: []common.PodAssociation{
				{
					From: "connection",
				},
			},
		},
	}
	return ProcessorsConfig{
		BaseProcessorsConfig: common.BaseProcessorsConfig{
			Batch: &common.BatchProcessorConfig{
				SendBatchSize:    512,
				Timeout:          "10s",
				SendBatchMaxSize: 512,
			},
			MemoryLimiter: &common.MemoryLimiterConfig{
				CheckInterval:        "1s",
				LimitPercentage:      75,
				SpikeLimitPercentage: 10,
			},
			K8sAttributes: &common.K8sAttributesProcessorConfig{
				AuthType:    "serviceAccount",
				Passthrough: false,
				Extract: common.ExtractK8sMetadataConfig{
					Metadata: k8sAttributes,
				},
				PodAssociation: podAssociations,
			},
			Resource: &common.ResourceProcessorConfig{
				Attributes: []common.AttributeAction{
					{
						Action: "insert",
						Key:    "k8s.cluster.name",
						Value:  "${KUBERNETES_SERVICE_HOST}",
					},
				},
			},
		},
		SpanFilter: FilterProcessorConfig{
			Traces: TraceConfig{
				Span: makeSpanFilterConfig(),
			},
		},
	}
}

func makeSpanFilterConfig() []string {
	return []string{
		"(attributes[\"http.method\"] == \"GET\") and (attributes[\"component\"] == \"proxy\") and (attributes[\"OperationName\"] == \"Egress\") and (resource.attributes[\"service.name\"] == \"grafana.kyma-system\")",
		"(attributes[\"http.method\"] == \"GET\") and (attributes[\"component\"] == \"proxy\") and (attributes[\"OperationName\"] == \"Ingress\") and (resource.attributes[\"service.name\"] == \"grafana.kyma-system\")",
		"(attributes[\"http.method\"] == \"GET\") and (attributes[\"component\"] == \"proxy\") and (attributes[\"OperationName\"] == \"Ingress\") and (IsMatch(attributes[\"http.url\"], \".+/metrics\") == true) and (resource.attributes[\"k8s.namespace.name\"] == \"kyma-system\")",
		"(attributes[\"http.method\"] == \"GET\") and (attributes[\"component\"] == \"proxy\") and (attributes[\"OperationName\"] == \"Ingress\") and (IsMatch(attributes[\"http.url\"], \".+/healthz(/.*)?\") == true) and (resource.attributes[\"k8s.namespace.name\"] == \"kyma-system\")",
		"(attributes[\"http.method\"] == \"GET\") and (attributes[\"component\"] == \"proxy\") and (attributes[\"OperationName\"] == \"Ingress\") and (attributes[\"user_agent\"] == \"vm_promscrape\")",
		"(attributes[\"http.method\"] == \"POST\") and (attributes[\"component\"] == \"proxy\") and (attributes[\"OperationName\"] == \"Egress\") and (IsMatch(attributes[\"http.url\"], \"http(s)?:\\\\/\\\\/telemetry-otlp-traces\\\\.kyma-system(\\\\..*)?:(4318|4317).*\") == true)",
		"(attributes[\"http.method\"] == \"POST\") and (attributes[\"component\"] == \"proxy\") and (attributes[\"OperationName\"] == \"Egress\") and (IsMatch(attributes[\"http.url\"], \"http(s)?:\\\\/\\\\/telemetry-trace-collector-internal\\\\.kyma-system(\\\\..*)?:(55678).*\") == true)",
		"(attributes[\"http.method\"] == \"POST\") and (attributes[\"component\"] == \"proxy\") and (attributes[\"OperationName\"] == \"Egress\") and (resource.attributes[\"service.name\"] == \"telemetry-fluent-bit.kyma-system\")",
	}
}