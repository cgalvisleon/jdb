package jdb

const (
	OF_SYSTEM  = "de sistema"
	ACTIVE     = "activo"
	ARCHIVED   = "archivado"
	CANCELLED  = "cancelado"
	FOR_DELETE = "eliminado"
	PENDING    = "pendiente"
	APPROVAL   = "aprobado"
	REFUSED    = "rechazado"
	IN_PROCESS = "en proceso"
	STOP       = "detenido"
	FAILED     = "fallido"
)

var LIST_STATES = []string{OF_SYSTEM, ACTIVE, ARCHIVED, CANCELLED, FOR_DELETE, PENDING, APPROVAL, REFUSED, IN_PROCESS, STOP, FAILED}
