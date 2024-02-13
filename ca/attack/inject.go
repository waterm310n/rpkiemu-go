package attack

import (
	"fmt"
	"log/slog"

	"github.com/waterm310n/rpkiemu-go/ca/krillop"
)

func injectAttack(attack *Attack, caOps map[string]*krillop.KrillK8sCA) {
	switch attack.AttackObject {
	case AttacKObject_RESOURCECERT:
		for _, attackData := range attack.AttackData {
			injectResourceCert(attackData, caOps)
		}
	case AttacKObject_ROAS:
		for _, attackData := range attack.AttackData {
			publishPoint := attackData.PublishPoint
			handle := attackData.HandleName
			for _, binding := range attackData.Bindings {
				caOps[publishPoint].AddAsnIpPair(handle, binding.Ip, fmt.Sprintf("%v", binding.Asn))
			}
		}
	default:
		slog.Error(fmt.Sprintf("AttackObject %v is wrong", attack.AttackObject))
		return
	}
}

func injectResourceCert(attackData *AttackData, caOps map[string]*krillop.KrillK8sCA) {
	handle := attackData.HandleName
	parentHandle := attackData.ParentHandleName
	publishPoint := attackData.PublishPoint
	parentPublishPoint := attackData.ParentPublishPoint
	ipv4 := attackData.Ipv4Resource
	ipv6 := attackData.Ipv6Resource
	var asnes []string
	for _, num := range attackData.ASNes {
		asnes = append(asnes, fmt.Sprintf("%v", num))
	}
	caOps[parentPublishPoint].Inject(handle, parentHandle, publishPoint, parentPublishPoint, ipv4, ipv6, asnes, caOps)
}
