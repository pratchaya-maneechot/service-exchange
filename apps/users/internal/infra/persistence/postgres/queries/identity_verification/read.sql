-- name: FindIdentityVerificationsByUserID :many
SELECT
    id, user_id, document_type, document_number, document_urls, status, submitted_at, verified_at, reviewer_id, rejection_reason
FROM identity_verifications
WHERE user_id = $1
ORDER BY submitted_at DESC;