package gateway

import (
	"net/http"

	"github.com/TykTechnologies/tyk/certs"
	"github.com/TykTechnologies/tyk/internal/crypto"
)

// CertificateCheckMW is used if domain was not detected or multiple APIs bind on the same domain. In this case authentification check happens not on TLS side but on HTTP level using this middleware
type CertificateCheckMW struct {
	BaseMiddleware
}

func (m *CertificateCheckMW) Name() string {
	return "CertificateCheckMW"
}

func (m *CertificateCheckMW) EnabledForSpec() bool {
	return m.Spec.UseMutualTLSAuth
}

func (m *CertificateCheckMW) ProcessRequest(w http.ResponseWriter, r *http.Request, _ interface{}) (error, int) {
	if ctxGetRequestStatus(r) == StatusOkAndIgnore {
		return nil, http.StatusOK
	}

	if m.Spec.UseMutualTLSAuth {
		certIDs := append(m.Spec.ClientCertificates, m.Spec.GlobalConfig.Security.Certificates.API...)
		apiCerts := m.Gw.CertificateManager.List(certIDs, certs.CertificatePublic)
		if err := crypto.ValidateRequestCerts(r, apiCerts); err != nil {
			return err, http.StatusForbidden
		}
	}
	return nil, http.StatusOK
}
