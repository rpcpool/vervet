package vervet

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/vault/api"
)

type vaultClient struct {
	apiClient *api.Client
	url       *url.URL
}

func newVaultClient(addr string) (*vaultClient, error) {
	vault := new(vaultClient)

	config := &api.Config{
		Address: addr,
	}

	api, err := api.NewClient(config)
	if err != nil {
		return vault, err
	}

	url, err := url.Parse(addr)
	if err != nil {
		return vault, err
	}

	vault.apiClient = api
	vault.url = url

	return vault, nil
}

// connect to Vault server and execute rekey operation
func (vault *vaultClient) rekey(keys []string) (*api.RekeyStatusResponse, error) {
	resp, err := vault.apiClient.Sys().RekeyStatus()
	if err != nil {
		return nil, err
	}

	if !resp.Started {
		return resp, fmt.Errorf("%s - rekey process has not been started", vault.url.Host)
	}

	complete := false
	nonce := resp.Nonce
	for _, key := range keys {
		resp, err := vault.apiClient.Sys().RekeyUpdate(key, nonce)
		if err != nil {
			return nil, err
		}

		if resp.Complete {
			complete = true
			break
		}
	}

	PrintInfo(fmt.Sprintf("%s - provided %d rekey key share(s) toward rekyey progress", vault.url.Host, len(keys)))

	resp, err = vault.apiClient.Sys().RekeyStatus()
	if err != nil {
		return nil, err
	}

	if complete {
		PrintSuccess(fmt.Sprintf("%s - rekey complete", vault.url.Host))
	}

	return resp, nil
}

// connect to Vault server and execute unseal operation
func (vault *vaultClient) unseal(keys []string) (*api.SealStatusResponse, error) {
	resp, err := vault.apiClient.Sys().SealStatus()
	if err != nil {
		return nil, err
	}

	if !resp.Initialized {
		return resp, fmt.Errorf("%s - Vault server is not initialized", vault.url.Host)
	}

	// if node is already unsealed, skip it
	if !resp.Sealed {
		PrintSuccess(vault.url.Host + " - already unsealed, skipping unseal operation")
		return resp, nil
	}

	for _, key := range keys {
		resp, err = vault.apiClient.Sys().Unseal(key)
		if err != nil {
			return nil, err
		}

		if !resp.Sealed {
			break
		}
	}

	PrintInfo(fmt.Sprintf("%s - provided %d unseal key share(s) toward unseal progress", vault.url.Host, len(keys)))

	resp, err = vault.apiClient.Sys().SealStatus()
	if err != nil {
		return nil, err
	}

	if !resp.Sealed {
		PrintSuccess(fmt.Sprintf("%s - Vault unsealed", vault.url.Host))
	}

	return resp, nil
}

// connect to Vault server and execute unseal operation
func (vault *vaultClient) generateRoot(keys []string) (*api.GenerateRootStatusResponse, error) {
	resp, err := vault.apiClient.Sys().GenerateRootStatus()
	if err != nil {
		return nil, err
	}

	// if node is already unsealed, skip it
	if !resp.Started {
		PrintWarning(vault.url.Host + " - root token generation process has not been started")
		return resp, nil
	}

	nonce := resp.Nonce
	for _, key := range keys {
		resp, err = vault.apiClient.Sys().GenerateRootUpdate(key, nonce)
		if err != nil {
			return nil, err
		}

		msg := fmt.Sprintf("%s - provided unseal key share, root token generation progress: %d of %d key shares",
			vault.url.Host, resp.Progress, resp.Required)
		PrintInfo(msg)

		if resp.Complete {
			msg = fmt.Sprintf("%s - root token generation complete", vault.url.Host)
			PrintSuccess(msg)

			return resp, nil
		}
	}

	return resp, nil
}

func printSealStatus(resp *api.SealStatusResponse) {
	status := "unsealed"
	if resp.Sealed {
		status = "sealed"
	} else {
		PrintKV("Cluster name", resp.ClusterName)
		PrintKV("Cluster ID", resp.ClusterID)
	}

	PrintKV("Seal status", status)
	PrintKV("Key threshold/shares", fmt.Sprintf("%d/%d", resp.T, resp.N))
	PrintKV("Progress", fmt.Sprintf("%d/%d", resp.Progress, resp.T))
	PrintKV("Version", resp.Version)
}

func printGenRootStatus(resp *api.GenerateRootStatusResponse) {
	status := "not started"
	if resp.Started {
		status = "started"

		if resp.Complete {
			status = "complete"
		}
	}

	PrintKV("Root generation", status)

	if resp.Started {
		PrintKV("Nonce", resp.Nonce)
		PrintKV("Progress", fmt.Sprintf("%d/%d", resp.Progress, resp.Required))

		if resp.PGPFingerprint != "" {
			PrintKV("PGP fingerprint", resp.PGPFingerprint)
		}
	}

	if resp.EncodedRootToken != "" {
		PrintKV("Encoded root token", resp.EncodedRootToken)
	}
}

func printRekeyStatus(resp *api.RekeyStatusResponse) {
	status := "not started"
	if resp.Started {
		status = "started"

		if resp.Required <= resp.Progress {
			status = "complete"
		}
	}

	PrintKV("Rekey", status)

	if resp.Started {
		PrintKV("Nonce", resp.Nonce)
		PrintKV("Progress", fmt.Sprintf("%d/%d", resp.Progress, resp.Required))
		PrintKV("New threshold", fmt.Sprintf("%d/%d", resp.T, resp.N))
		PrintKV("Backup", fmt.Sprintf("%t", resp.Backup))
		PrintKV("Verification required", fmt.Sprintf("%t", resp.VerificationRequired))
		PrintKV("Verification nonce", resp.VerificationNonce)

		for _, fingerprint := range resp.PGPFingerprints {
			PrintKV("PGP fingerprint", fingerprint)
		}
	}
}
