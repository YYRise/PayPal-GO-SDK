package paypalsdk

import "fmt"

// https://developer.paypal.com/docs/api/webhooks/v1/#definition-event_type
type EventType struct {
	Name             string             `json:"name"`
	Description      string             `json:"description,omitempty"`
	Status           string             `json:"status,omitempty"`
	ResourceVersions []*ResourceVersion `json:"resource_versions,omitempty"`
}

// https://developer.paypal.com/docs/api/webhooks/v1/#definition-resource_version
type ResourceVersion struct {
	ResourceVersion string `json:"resource_version"`
}

// https://developer.paypal.com/docs/api/webhooks/v1/#webhooks_create
type CreateWebhookReq struct {
	Url        string       `json:"url"`
	EventTypes []*EventType `json:"event_types"`
}

type Webhook struct {
	ID string `json:"id"`
	*CreateWebhookReq
	Links *[]LinkDescription `json:"links"`
}

/*
// POST https://api.sandbox.paypal.com/v1/notifications/webhooks \
// Create webhook
*/

func (c *Client) CreateWebhook(q *CreateWebhookReq) (*Webhook, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks"), q)
	rsp := &Webhook{}
	if err != nil {
		return rsp, err
	}
	err = c.SendWithAuth(req, rsp)
	return rsp, err
}

// https://developer.paypal.com/docs/api/webhooks/v1/#event-type_list
type WebhookList struct {
	Webhooks []*Webhook `json:"webhooks,omitempty"`
}

/*
// GET https://api.sandbox.paypal.com/v1/notifications/webhooks
// List webhooks
// anchor_type = APPLICATION or ACCOUNT. Default: APPLICATION.
*/

func (c *Client) ListWebhooks(anchor_type string) (results *WebhookList, err error) {
	var url string
	if anchor_type == "ACCOUNT" {
		url = fmt.Sprintf("%s%s?anchor_type=ACCOUNT", c.APIBase, "/v1/notifications/webhooks")
	} else {
		url = fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks")
	}
	req, err := c.NewRequest("GET", url, nil)
	rsp := &WebhookList{}
	if err != nil {
		return rsp, err
	}
	err = c.SendWithAuth(req, rsp)
	return rsp, err
}

/*
// DELETE https://api.sandbox.paypal.com/v1/notifications/webhooks/{webhook_id}
// Delete webhook
*/
func (c *Client) DeleteWebhook(id string) (err error) {

	url := fmt.Sprintf("%s%s/%s", c.APIBase, "/v1/notifications/webhooks", id)
	req, err := c.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}
