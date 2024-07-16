package constant

// xPay || funPay 代付
const (
	XPayPaymentSuccessCode     = 0    //成功
	XPayPaymentAbnormalCode    = 9999 //异常
	XPayPaymentStateAbnormal   = 0    //订单异常
	XPayPaymentStateProcessing = 1    //代付中(回调)
	XPayPaymentStateSuccess    = 2    //代付成功（回调）
	XPayPaymentStateFailed     = 3    //3-代付下单失败(注意不回调，不回调，不回调)
	XPayPaymentStateRevoke     = 4    //4-撤销
)

// inPay
const (
	InPayPaymentSuccessCode = 100 //成功 code=100即为成功；其他值为失败(若万一出现无返回值的情况，请将您的订单状态转换为处理中，以免引起不必要的损失)
)

const (
	LuckyInPayPaymentSuccessCode = "00000" //成功（注意：成功是5个0）

	LuckyInPayPaymentStateProcessing = 1 //代付中
	LuckyInPayPaymentStateSuccess    = 2 //代付成功
	LuckyInPayPaymentStateFailed     = 3 //代付失败
)
