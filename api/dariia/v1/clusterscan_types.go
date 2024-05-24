package v1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    batchv1 "k8s.io/api/batch/v1"
    batchv1beta1 "k8s.io/api/batch/v1beta1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ClusterScan is the Schema for the clusterscans API
type ClusterScan struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   ClusterScanSpec   `json:"spec,omitempty"`
    Status ClusterScanStatus `json:"status,omitempty"`
}

// ClusterScanSpec defines the desired state of ClusterScan
type ClusterScanSpec struct {
    // JobTemplate defines the job to be created for a one-off execution
    JobTemplate batchv1.JobSpec `json:"jobTemplate,omitempty"`

    // CronJobTemplate defines the cron job to be created for recurring executions
    CronJobTemplate batchv1beta1.CronJobSpec `json:"cronJobTemplate,omitempty"`

    // Schedule defines whether the job is one-off or recurring
    Schedule string `json:"schedule,omitempty"` // e.g., "*/5 * * * *" for a recurring job, empty for one-off
}

// ClusterScanStatus defines the observed state of ClusterScan
type ClusterScanStatus struct {
    // JobName is the name of the created Job or CronJob
    JobName string `json:"jobName,omitempty"`

    // LastExecutionTime is the last time the job was executed
    LastExecutionTime *metav1.Time `json:"lastExecutionTime,omitempty"`

    // LastExecutionResult records the result of the last job execution
    LastExecutionResult string `json:"lastExecutionResult,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterScanList contains a list of ClusterScan
type ClusterScanList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []ClusterScan `json:"items"`
}

func init() {
    SchemeBuilder.Register(&ClusterScan{}, &ClusterScanList{})
}

