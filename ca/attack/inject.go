package attack

import (
	"fmt"
	"log/slog"

	"github.com/waterm310n/rpkiemu-go/ca/krillop"
)

func injectAttack(attack *Attack, caOps map[string]*krillop.KrillK8sCA) {
	switch attack.AttackObject {
	case AttacKObject_RESOURCECERT:
		//Todo AttacKObject_RESOURCECERT
		fmt.Print("123")
	case AttacKObject_ROAS:
		
	default:
		slog.Error(fmt.Sprintf("AttackObject %v is wrong", attack.AttackObject))
		return
	}
}
