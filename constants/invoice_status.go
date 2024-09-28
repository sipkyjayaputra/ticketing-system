package constants

type InvoiceStatus int32

const (
	InvoiceStatus_Approval InvoiceStatus = 0
	InvoiceStatus_Revision InvoiceStatus = 1
	InvoiceStatus_Approved InvoiceStatus = 2
	InvoiceStatus_Rejected InvoiceStatus = -1
)
