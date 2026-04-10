package client

type OncallSchedule struct {
	ID          string  `json:"id"`
	OrgID       string  `json:"orgId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Timezone    string  `json:"timezone"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type OncallScheduleWithLayers struct {
	OncallSchedule
	Layers []OncallLayer `json:"layers"`
}

type OncallLayer struct {
	ID                 string   `json:"id"`
	ScheduleID         string   `json:"scheduleId"`
	Name               string   `json:"name"`
	RotationType       string   `json:"rotationType"`
	RotationLengthDays int      `json:"rotationLengthDays"`
	HandoffTime        string   `json:"handoffTime"`
	EffectiveFrom      string   `json:"effectiveFrom"`
	EffectiveUntil     *string  `json:"effectiveUntil"`
	ParticipantIDs     []string `json:"participantIds"`
	CreatedAt          string   `json:"createdAt"`
	UpdatedAt          string   `json:"updatedAt"`
}

type CreateScheduleRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Timezone    string  `json:"timezone"`
}

type UpdateScheduleRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Timezone    string  `json:"timezone"`
}

type CreateLayerRequest struct {
	Name               string   `json:"name"`
	RotationType       string   `json:"rotationType"`
	RotationLengthDays int      `json:"rotationLengthDays"`
	HandoffTime        string   `json:"handoffTime"`
	EffectiveFrom      string   `json:"effectiveFrom"`
	EffectiveUntil     *string  `json:"effectiveUntil,omitempty"`
	ParticipantIDs     []string `json:"participantIds"`
}

type UpdateLayerRequest struct {
	Name               string   `json:"name"`
	RotationType       string   `json:"rotationType"`
	RotationLengthDays int      `json:"rotationLengthDays"`
	HandoffTime        string   `json:"handoffTime"`
	EffectiveFrom      string   `json:"effectiveFrom"`
	EffectiveUntil     *string  `json:"effectiveUntil"`
	ParticipantIDs     []string `json:"participantIds"`
}

type EscalationPolicy struct {
	ID          string  `json:"id"`
	OrgID       string  `json:"orgId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type EscalationPolicyWithSteps struct {
	EscalationPolicy
	Steps []EscalationStep `json:"steps"`
}

type EscalationStep struct {
	ID                   string `json:"id"`
	Position             int    `json:"position"`
	TargetType           string `json:"targetType"`
	TargetID             string `json:"targetId"`
	EscalateAfterSeconds int    `json:"escalateAfterSeconds"`
	CreatedAt            string `json:"createdAt"`
	UpdatedAt            string `json:"updatedAt"`
}

type EscalationStepInput struct {
	Position             int    `json:"position"`
	TargetType           string `json:"targetType"`
	TargetID             string `json:"targetId"`
	EscalateAfterSeconds int    `json:"escalateAfterSeconds"`
}

type CreateEscalationPolicyRequest struct {
	Name        string                `json:"name"`
	Description *string               `json:"description,omitempty"`
	Steps       []EscalationStepInput `json:"steps"`
}

type UpdateEscalationPolicyRequest struct {
	Name        string                `json:"name"`
	Description *string               `json:"description"`
	Steps       []EscalationStepInput `json:"steps"`
}
