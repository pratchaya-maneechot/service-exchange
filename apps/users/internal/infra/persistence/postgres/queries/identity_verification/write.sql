-- name: CreateIdentityVerification :one
INSERT INTO identity_verifications (
    id, user_id, document_type, document_number, document_urls, status, submitted_at, verified_at, reviewer_id, rejection_reason
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW(), $7, $8, $9
) RETURNING *;

-- name: UpdateIdentityVerificationStatus :one
UPDATE identity_verifications
SET
    status = $2,
    verified_at = $3,
    reviewer_id = $4,
    rejection_reason = $5
WHERE id = $1
RETURNING *;