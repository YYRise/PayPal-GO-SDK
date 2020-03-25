package paypalsdk

import "time"

type E_EventResourceType string

const (
	E_EVENT_RESOURCE_TYPE_SUBCRIPTION E_EventResourceType = "subscription"
	E_EVENT_RESOURCE_TYPE_SALE        E_EventResourceType = "sale"
)

const (
	E_EVENT_TYPE_PAYMENT_SALE_COMPLETED = "PAYMENT.SALE.COMPLETED"
	E_EVENT_TYPE_PAYMENT_SALE_REFUNDED  = "PAYMENT.SALE.REFUNDED"
	E_EVENT_TYPE_PAYMENT_SALE_DENIED    = "PAYMENT.SALE.DENIED"
	E_EVENT_TYPE_PAYMENT_SALE_PENDING   = "PAYMENT.SALE.PENDING"

	E_EVENT_TYPE_BILLING_PLAN_CREATED = "BILLING.PLAN.CREATED"

	E_EVENT_TYPE_BILLING_SUBSCRIPTION_CREATED        = "BILLING.SUBSCRIPTION.CREATED"
	E_EVENT_TYPE_BILLING_SUBSCRIPTION_SUSPENDED      = "BILLING.SUBSCRIPTION.SUSPENDED"
	E_EVENT_TYPE_BILLING_SUBSCRIPTION_UPDATED        = "BILLING.SUBSCRIPTION.UPDATED"
	E_EVENT_TYPE_BILLING_SUBSCRIPTION_CANCELLED      = "BILLING.SUBSCRIPTION.CANCELLED"
	E_EVENT_TYPE_BILLING_SUBSCRIPTION_ACTIVATED      = "BILLING.SUBSCRIPTION.ACTIVATED"
	E_EVENT_TYPE_BILLING_SUBSCRIPTION_PAYMENT_FAILED = "BILLING.SUBSCRIPTION.PAYMENT.FAILED"
	E_EVENT_TYPE_BILLING_SUBSCRIPTION_RENEWED        = "BILLING.SUBSCRIPTION.RENEWED"
)

type Event struct {
	Id           string              `json:"id"`
	CreateTime   time.Time           `json:"create_time,omitempty"`
	ResourceType E_EventResourceType `json:"resource_type,omitempty"`
	EventVersion string              `json:"event_version,omitempty"`
	EventType    string              `json:"event_type,omitempty"`
	Summary      string              `json:"summary,omitempty"`
	Resource     interface{}         `json:"resource,omitempty"`
	Status       string              `json:"status,omitempty"`
	Links        []*Link             `json:"links,omitempty"`
}

func (e *Event) Sale() *Sale {
	if s, ok := e.Resource.(*Sale); ok {
		return s
	}
	return nil
}

func (e *Event) Subscription() *Subscription {
	if s, ok := e.Resource.(*Subscription); ok {
		return s
	}
	return nil
}
