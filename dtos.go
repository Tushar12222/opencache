package main

type PostKeyRequest struct {
	Value any `json:"value"`
}

type PatchKeyRequest struct {
	Value any `json:"value"`
}

type PostKeyResponse struct {
	Message any `json:"message"`
}

type PatchKeyResponse struct {
	Message any `json:"message"`
}

type GetKeyResponse struct {
	Message any `json:"message"`
	Value   any `json:"value"`
}

type DeleteKeyResponse struct {
	Message any `json:"message"`
}
