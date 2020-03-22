package paypalsdk

import (
	"fmt"
	"time"
)

const (
	K_SUBSCRIPTION_API = "/v1/billing/subscriptions"
)

/*
1. 创建plan以固定的时间间隔或根据用户订购的数量向用户收取固定金额的费用。
2. 为您的订阅者提供免费或打折的试用版，以获取更多订阅签约。
3. 轻松修改计划价格。
4. 使订户能够升级和降级其plan或更改其订购数量。
5. 自动为失败的付款恢复付款。
*/

type CreateSubscriptionReq struct {
	PlanID             string              `json:"plan_id"`
	StartTime          string              `json:"start_time,omitempty"` // Default: Current time.
	Quantity           string              `json:"quantity,omitempty"`   //  1<=len<=32,数字
	ShippingAmount     *Money              `json:"shipping_amount"`
	Subscriber         *Subscriber         `json:"subscriber"`
	AutoRenewal        bool                `json:"auto_renewal,omitempty"` // 订阅在计费周期完成后是否自动续订。
	ApplicationContext *ApplicationContext `json:"application_context,omitempty"`
}
type E_SubscriptionStatus string

const (
	E_SUBSCRIPTION_STATUS_APPROVAL_PENDING E_SubscriptionStatus = "APPROVAL_PENDING" // 订阅已创建但尚未获得买方批准。
	E_SUBSCRIPTION_STATUS_APPROVAL         E_SubscriptionStatus = "APPROVED"         // 买方已批准。
	E_SUBSCRIPTION_STATUS_ACTIVE           E_SubscriptionStatus = "ACTIVE"
	E_SUBSCRIPTION_STATUS_SUSPENDED        E_SubscriptionStatus = "SUSPENDED"
	E_SUBSCRIPTION_STATUS_CANCELLED        E_SubscriptionStatus = "CANCELLED"
	E_SUBSCRIPTION_STATUS_EXPIRED          E_SubscriptionStatus = "EXPIRED"
)

func (s E_SubscriptionStatus) Int() int {
	if s == E_SUBSCRIPTION_STATUS_ACTIVE {
		return 1
	} else {
		return 0
	}
}

type Subscription struct {
	Status           E_SubscriptionStatus `json:"status,omitempty"`
	StatusChangeNote string               `json:"status_change_note,omitempty"` //订阅状态的原因或注释 1<=len<=128
	StatusUpdateTime time.Time            `json:"status_update_time,omitempty"` // eg: 2020-03-09T12:00:01
	ID               string               `json:"id,omitempty"`                 // paypal生成的订阅 ID。
	PlanID           string               `json:"plan_id,omitempty"`
	StartTime        time.Time            `json:"start_time,omitempty"` // eg: 2020-03-09T12:00:01
	Quantity         string               `json:"quantity,omitempty"`   //  1<=len<=32,数字
	ShippingAmount   *Money               `json:"shipping_amount,omitempty"`
	Subscriber       *Subscriber          `json:"subscriber,omitempty"`
	BillingInfo      *BillingInfo         `json:"billing_info,omitempty"`
	CreateTime       time.Time            `json:"create_time,omitempty"` // 只读
	UpdateTime       time.Time            `json:"update_time,omitempty"` // 只读
	Links            []*LinkDescription   `json:"links,omitempty"`
	AutoRenewal      bool                 `json:"auto_renewal,omitempty"`
}

/*
// POST https://api.sandbox.paypal.com/v1/billing/subscriptions
// 创建成功触发webhook： BILLING.SUBSCRIPTION.CREATED
// 自动支付后触发webhook： PAYMENT.SALE.COMPLETED
// 自动支付失败触发webhook： BILLING.SUBSCRIPTION.PAYMENT.FAILED
// 创建订阅
*/

func (c *Client) CreateSubscription(q *CreateSubscriptionReq) (*Subscription, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions"), q)
	req.Header.Add("Prefer", "return=representation")
	rsp := &Subscription{}
	if err != nil {
		return rsp, err
	}
	err = c.SendWithAuth(req, rsp)
	return rsp, err
}

/*
// GET https://api.sandbox.paypal.com/v1/billing/subscriptions/I-BW452GLLEP1G \
// Show subscription details
*/

/*
// PATCH https://api.sandbox.paypal.com/v1/billing/subscriptions/I-BW452GLLEP1G \
// Update subscription

update the following fields:
	subscriber.shipping_address
	shipping_amount
	billing_info.outstanding_balance

// 触发webhook： BILLING.SUBSCRIPTION.UPDATED
// 更新
*/

func (c *Client) UpdateSubscription(subId string, op E_PatchOp) (*Subscription, error) {
	patchs := []Patch{
		Patch{
			Op:    op,
			Path:  "/status_change_note",
			Value: "/Item out of stock",
		},
	}

	req, err := c.NewRequest("PATCH", fmt.Sprintf("%s%s/%s", c.APIBase, K_SUBSCRIPTION_API, subId), patchs)
	rsp := &Subscription{}
	if err != nil {
		return rsp, err
	}
	err = c.SendWithAuth(req, rsp)
	return rsp, err
}

/*
// POST https://api.sandbox.paypal.com/v1/billing/subscriptions/I-BW452GLLEP1G/suspend \
// Suspend subscription
// 触发webhook： BILLING.SUBSCRIPTION.SUSPENDED
// 暂停
*/

/*
// POST https://api.sandbox.paypal.com/v1/billing/subscriptions/I-BW452GLLEP1G/activate \
// returns 204 No Content
// Activate subscription
// 触发webhook： BILLING.SUBSCRIPTION.ACTIVATED
// 激活
*/
type UpdateSubscriptionReq struct {
	Reason string `json:"reason"` // 1<=len<=128
}

func (c *Client) ActivateSubscription(subId, reason string) error {
	as := &UpdateSubscriptionReq{Reason: reason}
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s/%s/activate", c.APIBase, K_SUBSCRIPTION_API, subId), as)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

/*
// POST https://api.sandbox.paypal.com/v1/billing/subscriptions/I-BW452GLLEP1G/cancel \
// returns 204 No Content
// Cancel subscription
// 触发webhook： BILLING.SUBSCRIPTION.CANCELLED
// 取消
*/
func (c *Client) CancelSubscription(subID, reason string) error {
	as := &UpdateSubscriptionReq{Reason: reason}
	req, err := c.NewRequest("POST", fmt.Sprintf("%s%s/%s/cancel", c.APIBase, K_SUBSCRIPTION_API, subID), as)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

/*
// GET https://api.sandbox.paypal.com/v1/billing/subscriptions/I-BW452GLLEP1G
// Show subscription details
*/

func (c *Client) ShowSubscriptionDetails(subID string) (*Subscription, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s%s/%s", c.APIBase, K_SUBSCRIPTION_API, subID), nil)
	if err != nil {
		return nil, err
	}
	rsp := &Subscription{}
	err = c.SendWithAuth(req, rsp)
	return rsp, err
}

/*
// POST https://api.sandbox.paypal.com/v1/billing/subscriptions/I-BW452GLLEP1G/capture \
// Capture amount on subscription
// 触发webhook： PAYMENT.SALE.COMPLETED
// 订阅时获取用户授权付款。
*/

/*
// GET https://api.sandbox.paypal.com/v1/billing/subscriptions/I-BW452GLLEP1G/transactions?start_time=2018-01-21T07:50:20.940Z&end_time=2018-08-21T07:50:20.940Z" \
// List transactions for a subscription
// 查询订阅的交易列表, 默认当天
*/
type ListTransactionRsp struct {
	Transactions []*SubTransaction  `json:"transactions"`
	TotalItems   int                `json:"total_items"`
	TotalPages   int                `json:"total_pages"`
	Links        []*LinkDescription `json:"links,omitempty"`
}

func (c *Client) ListTransactionsForSubscription(subID, startTime, endTime string) (*ListTransactionRsp, error) {
	url := fmt.Sprintf("%s%s/%s%sstart_time=%s&end_time=%s", c.APIBase, K_SUBSCRIPTION_API, subID, "/transactions?", startTime, endTime)
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	rsp := &ListTransactionRsp{}
	err = c.SendWithAuth(req, rsp)
	return rsp, err
}
