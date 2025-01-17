/*
Copyright (C) 2022-2024 ApeCloud Co., Ltd

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

// BackupScheduleSpec defines the desired state of BackupSchedule.
type BackupScheduleSpec struct {
	// Which backupPolicy is applied to perform this backup.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern:=`^[a-z0-9]([a-z0-9\.\-]*[a-z0-9])?$`
	BackupPolicyName string `json:"backupPolicyName"`

	// startingDeadlineMinutes defines the deadline in minutes for starting the
	// backup workload if it misses scheduled time for any reason.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1440
	StartingDeadlineMinutes *int64 `json:"startingDeadlineMinutes,omitempty"`

	// schedules defines the list of backup schedules.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Schedules []SchedulePolicy `json:"schedules"`
}

type SchedulePolicy struct {
	// enabled specifies whether the backup schedule is enabled or not.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// backupMethod specifies the backup method name that is defined in backupPolicy.
	// +kubebuilder:validation:Required
	BackupMethod string `json:"backupMethod"`

	// the cron expression for schedule, the timezone is in UTC.
	// see https://en.wikipedia.org/wiki/Cron.
	// +kubebuilder:validation:Required
	CronExpression string `json:"cronExpression"`

	// retentionPeriod determines a duration up to which the backup should be kept.
	// controller will remove all backups that are older than the RetentionPeriod.
	// For example, RetentionPeriod of `30d` will keep only the backups of last 30 days.
	// Sample duration format:
	//
	// - years: 	2y
	// - months: 	6mo
	// - days: 		30d
	// - hours: 	12h
	// - minutes: 	30m
	//
	// You can also combine the above durations. For example: 30d12h30m
	//
	// +optional
	// +kubebuilder:default="7d"
	RetentionPeriod RetentionPeriod `json:"retentionPeriod,omitempty"`
}

// BackupScheduleStatus defines the observed state of BackupSchedule.
type BackupScheduleStatus struct {
	// phase describes the phase of the BackupSchedule.
	// +optional
	Phase BackupSchedulePhase `json:"phase,omitempty"`

	// observedGeneration is the most recent generation observed for this
	// BackupSchedule. It refers to the BackupSchedule's generation, which is
	// updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// failureReason is an error that caused the backup to fail.
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// schedules describes the status of each schedule.
	// +optional
	Schedules map[string]ScheduleStatus `json:"schedules,omitempty"`
}

// BackupSchedulePhase defines the phase of BackupSchedule
type BackupSchedulePhase string

const (
	// BackupSchedulePhaseAvailable means the backup schedule is available.
	BackupSchedulePhaseAvailable BackupSchedulePhase = "Available"

	// BackupSchedulePhaseFailed means the backup schedule is failed.
	BackupSchedulePhaseFailed BackupSchedulePhase = "Failed"
)

// ScheduleStatus defines the status of each schedule.
type ScheduleStatus struct {
	// phase describes the phase of the schedule.
	// +optional
	Phase SchedulePhase `json:"phase,omitempty"`

	// failureReason is an error that caused the backup to fail.
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// lastScheduleTime records the last time the backup was scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`

	// lastSuccessfulTime records the last time the backup was successfully completed.
	// +optional
	LastSuccessfulTime *metav1.Time `json:"lastSuccessfulTime,omitempty"`
}

// SchedulePhase defines the phase of schedule
type SchedulePhase string

const (
	ScheduleRunning SchedulePhase = "Running"
	ScheduleFailed  SchedulePhase = "Failed"
)

// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kubeblocks},scope=Namespaced,shortName=bs
// +kubebuilder:printcolumn:name="STATUS",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=`.metadata.creationTimestamp`

// BackupSchedule is the Schema for the backupschedules API.
type BackupSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupScheduleSpec   `json:"spec,omitempty"`
	Status BackupScheduleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BackupScheduleList contains a list of BackupSchedule.
type BackupScheduleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupSchedule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BackupSchedule{}, &BackupScheduleList{})
}
