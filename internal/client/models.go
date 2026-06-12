package client

// User is a member of the org as returned by GET /v1/users. firstName,
// lastName, and email are nullable in the API.
type User struct {
	ID        string  `json:"id"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     *string `json:"email"`
}

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

// WorkingHoursCondition gates an escalation step to a Working Hours set.
// When is "during" or "outside". A nil *WorkingHoursCondition means the step
// is unconditional; on a full-replacement update it is sent as an explicit
// null so the server clears any prior condition (ADR-0040).
type WorkingHoursCondition struct {
	WorkingHoursID string `json:"workingHoursId"`
	When           string `json:"when"`
}

type EscalationStep struct {
	ID                   string                 `json:"id"`
	Position             int                    `json:"position"`
	TargetType           string                 `json:"targetType"`
	TargetID             string                 `json:"targetId"`
	EscalateAfterSeconds int                    `json:"escalateAfterSeconds"`
	Condition            *WorkingHoursCondition `json:"condition"`
	CreatedAt            string                 `json:"createdAt"`
	UpdatedAt            string                 `json:"updatedAt"`
}

type EscalationStepInput struct {
	Position             int                    `json:"position"`
	TargetType           string                 `json:"targetType"`
	TargetID             string                 `json:"targetId"`
	EscalateAfterSeconds int                    `json:"escalateAfterSeconds"`
	Condition            *WorkingHoursCondition `json:"condition"`
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

type WorkingHoursWindow struct {
	Days  []int  `json:"days"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type WorkingHours struct {
	ID        string               `json:"id"`
	OrgID     string               `json:"orgId"`
	Name      string               `json:"name"`
	Timezone  string               `json:"timezone"`
	Windows   []WorkingHoursWindow `json:"windows"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt string               `json:"updatedAt"`
}

type CreateWorkingHoursRequest struct {
	Name     string               `json:"name"`
	Timezone string               `json:"timezone"`
	Windows  []WorkingHoursWindow `json:"windows"`
}

type UpdateWorkingHoursRequest struct {
	Name     string               `json:"name"`
	Timezone string               `json:"timezone"`
	Windows  []WorkingHoursWindow `json:"windows"`
}

type RoutingPolicy struct {
	ID          string  `json:"id"`
	OrgID       string  `json:"orgId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsDefault   bool    `json:"isDefault"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type RoutingPolicyWithRules struct {
	RoutingPolicy
	Rules []RoutingRule `json:"rules"`
}

type RoutingRule struct {
	Severity           string  `json:"severity"`
	Outcome            string  `json:"outcome"`
	EscalationPolicyID *string `json:"escalationPolicyId"`
}

type RoutingRuleInput struct {
	Severity           string  `json:"severity"`
	Outcome            string  `json:"outcome"`
	EscalationPolicyID *string `json:"escalationPolicyId"`
}

type CreateRoutingPolicyRequest struct {
	Name        string             `json:"name"`
	Description *string            `json:"description,omitempty"`
	Rules       []RoutingRuleInput `json:"rules"`
}

type UpdateRoutingPolicyRequest struct {
	Name        string             `json:"name"`
	Description *string            `json:"description"`
	Rules       []RoutingRuleInput `json:"rules"`
}

// MaintenanceWindow is the active maintenance interval for a component, or nil
// when the component is not in maintenance.
type MaintenanceWindow struct {
	StartedAt string  `json:"startedAt"`
	EndsAt    *string `json:"endsAt"`
}

type Component struct {
	ID                string             `json:"id"`
	OrgID             string             `json:"orgId"`
	Name              string             `json:"name"`
	Description       *string            `json:"description"`
	Aliases           []string           `json:"aliases"`
	OperationalStatus string             `json:"operationalStatus"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow"`
	CreatedAt         string             `json:"createdAt"`
	UpdatedAt         string             `json:"updatedAt"`
}

type CreateComponentRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Aliases     []string `json:"aliases"`
}

// UpdateComponentRequest carries the full desired state. Aliases is always sent
// (never omitted) so the server replaces the alias set wholesale, matching the
// IaC full-replacement posture (ADR-0003).
type UpdateComponentRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Aliases     []string `json:"aliases"`
}

// RoutingSelector is one ordered selector row on an inbound integration
// (ADR-0039): its matchers are ANDed together and, on the first matching row,
// RoutingPolicyID becomes the policy that applies.
type RoutingSelector struct {
	Matchers        []SelectorMatcher `json:"matchers"`
	RoutingPolicyID string            `json:"routingPolicyId"`
}

type SelectorMatcher struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type InboundIntegration struct {
	ID              string            `json:"id"`
	OrgID           string            `json:"orgId"`
	Name            string            `json:"name"`
	ComponentID     string            `json:"componentId"`
	RoutingPolicyID *string           `json:"routingPolicyId"`
	Selectors       []RoutingSelector `json:"selectors"`
	CreatedAt       string            `json:"createdAt"`
	UpdatedAt       string            `json:"updatedAt"`
}

// SetSelectorsRequest replaces an inbound integration's ordered Routing
// Selectors wholesale (ADR-0039 / full-replacement, ADR-0003).
type SetSelectorsRequest struct {
	Selectors []RoutingSelector `json:"selectors"`
}
