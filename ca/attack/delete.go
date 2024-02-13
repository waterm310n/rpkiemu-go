package attack

import (
	"fmt"
	"log/slog"

	"github.com/waterm310n/rpkiemu-go/ca/krillop"
)

func deleteAttack(attack *Attack, caOps map[string]*krillop.KrillK8sCA) {
	switch attack.AttackObject {
	case AttacKObject_RESOURCECERT:
		for _, attackData := range attack.AttackData {
			publishPoint := attackData.PublishPoint
			handle := attackData.HandleName
			parentHandle := attackData.ParentHandleName
			parentPublishPoint := attackData.ParentPublishPoint
			caOps[publishPoint].DeleteCert(handle, parentHandle, caOps[parentPublishPoint])
		}
	case AttacKObject_ROAS:
		for _, attackData := range attack.AttackData {
			publishPoint := attackData.PublishPoint
			handle := attackData.HandleName
			asn := attackData.ASN
			caOps[publishPoint].DeleteRoa(handle, asn)
		}
	default:
		slog.Error(fmt.Sprintf("AttackObject %v is wrong", attack.AttackObject))
		return
	}
}
