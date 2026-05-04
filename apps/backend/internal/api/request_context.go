package api

import "context"

type requestMetaKey struct{}

type requestMeta struct {
	ownerUserID string
	apiKey      string
}

func withRequestMeta(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestMetaKey{}, &requestMeta{})
}

func setOwnerUserIDFromVerifiedAuth(ctx context.Context, ownerUserID string) {
	if ownerUserID == "" {
		return
	}
	meta, _ := ctx.Value(requestMetaKey{}).(*requestMeta)
	if meta == nil {
		return
	}
	meta.ownerUserID = ownerUserID
}

func getOwnerUserID(ctx context.Context) string {
	meta, _ := ctx.Value(requestMetaKey{}).(*requestMeta)
	if meta == nil {
		return ""
	}
	return meta.ownerUserID
}

func setAPIKeyFromCheckRequest(ctx context.Context, apiKey string) {
	if apiKey == "" {
		return
	}
	meta, _ := ctx.Value(requestMetaKey{}).(*requestMeta)
	if meta == nil {
		return
	}
	meta.apiKey = apiKey
}

func getAPIKeyFromCheckRequest(ctx context.Context) string {
	meta, _ := ctx.Value(requestMetaKey{}).(*requestMeta)
	if meta == nil {
		return ""
	}
	return meta.apiKey
}
