package paypalsdk

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-patch
type E_PatchOp string

const (
	E_PATCH_OP_ADD     E_PatchOp = "add"
	E_PATCH_OP_REMOVE  E_PatchOp = "remove"
	E_PATCH_OP_REPLACE E_PatchOp = "replace"
	E_PATCH_OP_MOVE    E_PatchOp = "move"
	E_PATCH_OP_COPY    E_PatchOp = "copy"
	E_PATCH_OP_TEST    E_PatchOp = "test"
)

type Patch struct {
	Op    E_PatchOp   `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
	From  string      `json:"from,omitempty"`
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-money
type Money struct {
	CurrencyCode string `json:"currency_code,omitempty"` // len=3, eg: USD ……
	Value        string `json:"value,omitempty"`         // len<=32, 必须数字，eg：123.45
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-subscriber
type Subscriber struct {
	Name            *Name           `json:"name,omitempty"` // 只支持 given_name 和 surname
	EmailAddress    string          `json:"email_address,omitempty"`
	PayerId         string          `json:"payer_id ,omitempty"` // PayPal为付款人分配的ID, 只读
	ShippingAddress *ShippingDetail `json:"shipping_address,omitempty"`
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-name
type Name struct {
	Prefix            string `json:"prefix,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	Surname           string `json:"surname,omitempty"`
	MiddleName        string `json:"middle_name,omitempty"`
	Suffix            string `json:"suffix,omitempty"`
	AlternateFullName string `json:"alternate_full_name,omitempty"`
	FullName          string `json:"full_name,omitempty"`
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-shipping_detail
type ShippingDetail struct {
	Name    *ShippingDetailName            `json:"name,omitempty"`
	Address *ShippingDetailAddressPortable `json:"address,omitempty"` // 只支持 address_line_1, address_line_2, admin_area_1, admin_area_2, postal_code, and country_code properties.
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-shipping_detail.name
type ShippingDetailName struct {
	FullName string `json:"full_name,omitempty"`
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-shipping_detail.address_portable
type ShippingDetailAddressPortable struct {
	AddressLine_1 string `json:"address_line_1,omitempty"`
	AddressLine_2 string `json:"address_line_2,omitempty"`
	AdminArea_2   string `json:"admin_area_2,omitempty"` // 城市、城镇或村庄。小于admin_area_level_1。
	AdminArea_1   string `json:"admin_area_1,omitempty"` // 国家/地区最高级别的分区，通常是省、州
	PostalCode    string `json:"postal_code,omitempty"`  // 邮编
	CountryCode   string `json:"country_code"`           // eg: GB
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-application_context
type E_ShippingPreference string

const (
	E_SHIPPING_PREFERENCE_GET_FROM_FILE        E_ShippingPreference = "GET_FROM_FILE"        // 在贝宝网站上获得客户提供的送货地址。
	E_SHIPPING_PREFERENCE_NO_SHIPPING                               = "NO_SHIPPING"          // 从PayPal网站编辑送货地址。 推荐用于数字商品。
	E_SHIPPING_PREFERENCE_SET_PROVIDED_ADDRESS                      = "SET_PROVIDED_ADDRESS" // 获取商家提供的地址。 客户无法在PayPal网站上更改此地址。 如果商家未通过地址，则客户可以在PayPal页面上选择地址。
)

type E_UserAction string

const (
	E_USER_ACTION_CONTINUE      E_UserAction = "CONTINUE" /* 将客户重定向到PayPal订阅同意页面后，将出现“继续”按钮。当您要控制订阅的激活并且不希望PayPal激活订阅时，请使用此选项。*/
	E_USER_ACTION_SUBSCRIBE_NOW E_UserAction = "SUBSCRIBE_NOW"
)

type ApplicationContext struct {
	BrandName          string               `json:"brand_name,omitempty"`          // 1<=len<=127
	Locale             string               `json:"locale,omitempty"`              // 2<=len<=10 eg: da-DK, he-IL, id-ID, ja-JP, no-NO, pt-BR, ru-RU, sv-SE, th-TH, zh-CN, zh-HK, or zh-TW
	ShippingPreference E_ShippingPreference `json:"shipping_preference,omitempty"` // Default: GET_FROM_FILE.
	UserAction         E_UserAction         `json:"user_action,omitempty"`         // Default: SUBSCRIBE_NOW.
	PaymentMethod      PaymentMethod        `json:"payment_method,omitempty"`
	ReturnUrl          string               `json:"return_url"`
	CancelUrl          string               `json:"cancel_url"`
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-payment_method
type E_PayeePreferred string

const (
	E_PAYEE_PREFERRED_UNRESTRICTED               E_PayeePreferred = "UNRESTRICTED"               // 接受客户的任何类型的付款。
	E_PAYEE_PREFERRED_IMMEDIATE_PAYMENT_REQUIRED E_PayeePreferred = "IMMEDIATE_PAYMENT_REQUIRED" // 只接受客户的即时付款。例如，信用卡、PayPal 余额或即时 ACH。确保在捕获时付款没有"挂起"状态。
)

type PaymentMethod struct {
	PayerSelected  string           `json:"payer_selected,omitempty"`  //Default: PAYPAL.
	PayeePreferred E_PayeePreferred `json:"payee_preferred,omitempty"` // Default: UNRESTRICTED.
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-subscription_billing_info
type BillingInfo struct {
	OutstandingBalance  *Money                `json:"outstanding_balance"`
	CycleExecutions     []*CycleExecutions    `json:"cycle_executions,omitempty"`
	LastPayment         *LastPaymentDetails   `json:"last_payment,omitempty"`          //只读
	NextBillingTime     string                `json:"next_billing_time,omitempty"`     // 只读
	FinalPaymentTime    string                `json:"final_payment_time,omitempty"`    // 只读
	FailedPaymentsCount int                   `json:"failed_payments_count,omitempty"` //[0, 999] // 连续付款失败数。成功付款后重置为 0。如果达到payment_failure_threshold值，订阅将更新为"暂停"状态。
	LastFailedPayment   *FailedPaymentDetails `json:"last_failed_payment,omitempty"`   //只读
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-cycle_execution
type E_TenureType string

const (
	E_TENURE_TYPE_REGULAR E_TenureType = "REGULAR" // 常规计费周期。
	E_TENURE_TYPE_TRIAL                = "TRIAL"   // 试用计费周期。
)

type CycleExecutions struct {
	TenureType                  E_TenureType `json:"tenure_type"`                              // 只读
	Sequence                    int          `json:"sequence"`                                 //[0, 99] 在其他计费周期中运行此周期的顺序。
	CyclesCompleted             int          `json:"cycles_completed"`                         //[0, 9999] 已完成的计费周期数。
	CyclesRemaining             int          `json:"cycles_remaining,omitempty"`               //只读, [0, 9999] 对于有限的计费周期，cycles_remaining是剩余周期的数量。对于无限计费周期，cycles_remaining设置为 0。
	CurrentPricingSchemeVersion int          `json:"current_pricing_scheme_version,omitempty"` //只读, [0, 99] 计费周期的活动定价方案版本。
	TotalCycles                 int          `json:"total_cycles,omitempty"`                   //只读, [0, 999] 此计费周期运行的次数。试用计费周期对于total_cycles只能有 1 。常规计费周期可以具有无限周期（total_cycles值为 0）或有限数量的周期（total_cycles值介于 1 和 999 之间）。
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-last_payment_details
type LastPaymentDetails struct {
	Amount *Money `json:"amount"`
	Time   string `json:"time"` // 只读
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-failed_payment_details
type E_ReasonCode string

const (
	E_REASPN_CODE_PAYMENT_DENIED                 E_ReasonCode = "PAYMENT_DENIED"                 //由于一个或多个客户问题，PayPal 拒绝付款。
	E_REASPN_CODE_COMPLIANCE_VIOLATION           E_ReasonCode = "COMPLIANCE_VIOLATION"           //由于违反合规易被拒绝。
	E_REASPN_CODE_PAYEE_ACCOUNT_LOCKED_OR_CLOSED E_ReasonCode = "PAYEE_ACCOUNT_LOCKED_OR_CLOSED" //收件人帐户已锁定或关闭，无法接收付款。
)

type FailedPaymentDetails struct {
	Amount               *Money `json:"amount"`
	Time                 string `json:"time"`                              // 只读
	reason_code          string `json:"reason_code,omitempty"`             // 只读, [1, 120]
	NextPaymentRetryTime string `json:"next_payment_retry_time,omitempty"` // 只读
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-link_description
type LinkDescription struct {
	Href   string  `json:"href"`
	Rel    LinkRel `json:"rel"`
	Method string  `json:"method,omitempty"` // GET, POST, PUT, DELETE, HEAD, CONNECT, OPTIONS, PATCH
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-billing_cycle
type BillingCycle struct {
	PricingScheme *PricingScheme `json:"pricing_scheme"`
	Frequency     *Frequency     `json:"frequency"` //计费周期的频率详细信息
	TenureType    E_TenureType   `json:"tenure_type"`
	Sequence      int            `json:"sequence"`               // [1, 99] 在其他计费周期中，此周期的运行顺序。例如，试用计费周期的序列为1，而常规计费周期的序列为2，因此试用周期在常规周期之前运行。
	TotalCycles   int            `json:"total_cycles,omitempty"` //[0, 999] 此计费周期运行的次数。试用计费周期对于total_cycles只能有 1 。常规计费周期可以具有无限周期（total_cycles值为 0）或有限数量的周期（total_cycles值介于 1 和 999 之间）。
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-pricing_scheme
type PricingScheme struct {
	Version    int     `json:"version,omitempty"` // [0, 999]
	FixedPrice *Money  `json:"fixed_price,omitempty"`
	CreateTime *string `json:"create_time,omitempty"`
	UpdateTime *string `json:"update_time,omitempty"`
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-frequency
type E_FrequencyInterval string

const (
	E_FREQUENCY_INTERVAL_DAY   E_FrequencyInterval = "DAY"
	E_FREQUENCY_INTERVAL_WEEK  E_FrequencyInterval = "WEEK"
	E_FREQUENCY_INTERVAL_MONTH E_FrequencyInterval = "MONTH"
	E_FREQUENCY_INTERVAL_YEAR  E_FrequencyInterval = "YEAR"
)

type Frequency struct {
	IntervalUnit  E_FrequencyInterval `json:"interval_unit"`
	IntervalCount int                 `json:"interval_count,omitempty"` // Default: 1.
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-taxes

type Taxes struct {
	Percentage string `json:"percentage"` //帐单金额的百分比。
	Inclusive  bool   `json:"inclusive"`  // 指示税金是否已包含在计费金额中。Default: true.
}

// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-payment_preferences
type PaymentPreferences struct {
	AutoBillOutstanding     bool   `json:"auto_bill_outstanding"` // Default: true. 是否\在下一个计费周期中自动计费未结金额。
	SetupFee                *Money `json:"setup_fee"`
	SetupFeeFailureAction   string `json:"setup_fee_failure_action"`  //[CONTINUE, CANCEL] Default: CANCEL. 初始付款失败，则对订阅执行的操作。
	PaymentFailureThreshold int    `json:"payment_failure_threshold"` // Default: 0. 暂停订阅之前的最大付款失败数。
}

//type ApplicationContext struct {
//	BrandName          string `json:"brand_name"` // 1<=len<=127
//	Locale             string `json:"locale"`     // 2<=len<=10 eg: da-DK, he-IL, id-ID, ja-JP, no-NO, pt-BR, ru-RU, sv-SE, th-TH, zh-CN, zh-HK, or zh-TW
//	LandingPage        string `json:"landing_page"`
//	ShippingPreference string `json:"shipping_preference"` // GET_FROM_FILE, NO_SHIPPING, SET_PROVIDED_ADDRESS
//	UserAction         string `json:"user_action"`
//}
