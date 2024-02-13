package attack

import (
	"fmt"
	"log/slog"

	"github.com/waterm310n/rpkiemu-go/ca/krillop"
)

func modificateAttack(attack *Attack, caOps map[string]*krillop.KrillK8sCA) {
	switch attack.AttackObject {
	case AttacKObject_RESOURCECERT:
		for _, attackData := range attack.AttackData {
			handle := attackData.HandleName
			parentHandle := attackData.ParentHandleName
			parentPublishPoint := attackData.ParentPublishPoint
			//考虑是恶意修改行为，因此只在上级所在的发布点对下级的资源修改，而下级是否知道事情发生应该是无关紧要的。
			ipv4 := attackData.Ipv4Resource
			ipv6 := attackData.Ipv6Resource
			var asnes []string
			for _, num := range attackData.ASNes {
				asnes = append(asnes, fmt.Sprintf("%v", num))
			}
			caOps[parentPublishPoint].Modificate(handle, parentHandle, ipv4, ipv6, asnes)

		}
	case AttacKObject_ROAS:
		for _, attackData := range attack.AttackData {
			publishPoint := attackData.PublishPoint
			handle := attackData.HandleName
			for _, binding := range attackData.PreBindings {
				caOps[publishPoint].RemoveAsnIpPair(handle, binding.Ip, fmt.Sprintf("%v", binding.Asn))
			}
			for _, binding := range attackData.AfterBindings {
				caOps[publishPoint].AddAsnIpPair(handle, binding.Ip, fmt.Sprintf("%v", binding.Asn))
			}
		}
	default:
		slog.Error(fmt.Sprintf("AttackObject %v is wrong", attack.AttackObject))
		return
	}
}
