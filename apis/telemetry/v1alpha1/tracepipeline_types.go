/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TracePipelineSpec defines the desired state of TracePipeline
type TracePipelineSpec struct {
	// Defines a destination for shipping trace data. Only one can be defined per pipeline.
	Output TracePipelineOutput `json:"output"`
}

// TracePipelineOutput defines the output configuration section.
type TracePipelineOutput struct {
	// Configures the underlying Otel Collector with an [OTLP exporter](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/otlpexporter/README.md). If you switch `protocol`to `http`, an [OTLP HTTP exporter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/otlphttpexporter) is used.
	Otlp *OtlpOutput `json:"otlp"`
}

type TracePipelineConditionType string

// These are the valid statuses of TracePipeline.
const (
	TracePipelinePending TracePipelineConditionType = "Pending"
	TracePipelineRunning TracePipelineConditionType = "Running"
)

// TracePipelineCondition contains details for the current condition of this LogPipeline.
type TracePipelineCondition struct {
	// Point in time the condition transitioned into a different state.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason of last transition.
	Reason string `json:"reason,omitempty"`
	// The possible transition types are:<br>- `Running`: The instance is ready and usable.<br>- `Pending`: The pipeline is being activated.
	Type TracePipelineConditionType `json:"type,omitempty"`
}

// Defines the observed state of TracePipeline.
type TracePipelineStatus struct {
	// An array of conditions describing the status of the pipeline.
	Conditions []TracePipelineCondition `json:"conditions,omitempty"`
}

func NewTracePipelineCondition(reason string, condType TracePipelineConditionType) *TracePipelineCondition {
	return &TracePipelineCondition{
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Type:               condType,
	}
}

func (tps *TracePipelineStatus) GetCondition(condType TracePipelineConditionType) *TracePipelineCondition {
	for cond := range tps.Conditions {
		if tps.Conditions[cond].Type == condType {
			return &tps.Conditions[cond]
		}
	}
	return nil
}

func (tps *TracePipelineStatus) HasCondition(condition TracePipelineConditionType) bool {
	return tps.GetCondition(condition) != nil
}

func (tps *TracePipelineStatus) SetCondition(cond TracePipelineCondition) {
	currentCond := tps.GetCondition(cond.Type)
	if currentCond != nil && currentCond.Reason == cond.Reason {
		return
	}
	if currentCond != nil {
		cond.LastTransitionTime = currentCond.LastTransitionTime
	}
	newConditions := filterTracePipelineCondition(tps.Conditions, cond.Type)
	tps.Conditions = append(newConditions, cond)
}

func filterTracePipelineCondition(conditions []TracePipelineCondition, condType TracePipelineConditionType) []TracePipelineCondition {
	var newConditions []TracePipelineCondition
	for _, cond := range conditions {
		if cond.Type == condType {
			continue
		}
		newConditions = append(newConditions, cond)
	}
	return newConditions
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.conditions[-1].type`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// TracePipeline is the Schema for the tracepipelines API
type TracePipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Defines the desired state of TracePipeline
	Spec TracePipelineSpec `json:"spec,omitempty"`
	// Shows the observed state of the TracePipeline
	Status TracePipelineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TracePipelineList contains a list of TracePipeline
type TracePipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TracePipeline `json:"items"`
}

//nolint:gochecknoinits // SchemeBuilder's registration is required.
func init() {
	SchemeBuilder.Register(&TracePipeline{}, &TracePipelineList{})
}
