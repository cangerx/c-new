package operation_setting

import "github.com/QuantumNous/new-api/setting/config"

type QuotaSetting struct {
	EnableFreeModelPreConsume bool `json:"enable_free_model_pre_consume"` // 是否对免费模型启用预消耗

	// 默认任务计费模式：per_call（按次）或 per_second（按秒）。
	// 任务/视频模型未单独配置价格时，使用该模式与价格兜底计费。
	DefaultTaskBillingMode string  `json:"default_task_billing_mode"` // "per_call" | "per_second"
	DefaultTaskPrice       float64 `json:"default_task_price"`        // 默认价格（per_call 为每次价格，per_second 为每秒价格）
}

// 默认配置
var quotaSetting = QuotaSetting{
	EnableFreeModelPreConsume: true,
	DefaultTaskBillingMode:    "per_second",
	DefaultTaskPrice:          0,
}

func init() {
	// 注册到全局配置管理器
	config.GlobalConfig.Register("quota_setting", &quotaSetting)
}

func GetQuotaSetting() *QuotaSetting {
	return &quotaSetting
}
