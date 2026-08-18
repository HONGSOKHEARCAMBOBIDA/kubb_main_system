package model

type AdmissionState string

const (
	AdmissionStateCreated   AdmissionState = "created"
	AdmissionStateSubmitted AdmissionState = "submitted"
	AdmissionStateApproved  AdmissionState = "approved"
	AdmissionStateRejected  AdmissionState = "rejected"
	AdmissionStateCancelled AdmissionState = "cancelled"
)
