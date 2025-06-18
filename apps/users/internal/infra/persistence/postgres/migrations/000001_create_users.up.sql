-- Create Users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    line_user_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255), -- Can be NULL
    password_hash VARCHAR(255), -- Can be NULL
    status VARCHAR(50) NOT NULL, -- Enum-like string (e.g., 'ACTIVE', 'INACTIVE', 'PENDING_VERIFICATION', 'SUSPENDED')
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE -- Can be NULL
);

-- Create Profiles table (one-to-one with users)
CREATE TABLE profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE, -- ON DELETE CASCADE: if user is deleted, profile is also deleted
    display_name VARCHAR(255) NOT NULL,
    first_name VARCHAR(255), -- Can be NULL
    last_name VARCHAR(255),  -- Can be NULL
    bio TEXT,                -- Can be NULL
    avatar_url TEXT,         -- Can be NULL
    phone_number VARCHAR(50),-- Can be NULL
    address TEXT,            -- Can be NULL
    preferences JSONB NOT NULL DEFAULT '{}' -- Store map[string]any as JSONB, default to empty object
);

-- Create Roles table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY, -- SERIAL for auto-incrementing integer ID
    name VARCHAR(50) UNIQUE NOT NULL, -- Enum-like string (e.g., 'POSTER', 'TASKER', 'ADMIN')
    description TEXT
);

-- Create user_roles join table (many-to-many relationship)
CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id) -- Composite primary key to prevent duplicate assignments
);

-- Create IdentityVerifications table
CREATE TABLE identity_verifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL, -- Enum-like string (e.g., 'NATIONAL_ID', 'PASSPORT')
    document_number VARCHAR(255) NOT NULL, -- Encrypt/hash this in application layer before storing
    document_urls TEXT[] NOT NULL DEFAULT '{}', -- Array of text for URLs, default to empty array
    status VARCHAR(50) NOT NULL, -- Enum-like string (e.g., 'PENDING', 'APPROVED', 'REJECTED')
    submitted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    verified_at TIMESTAMP WITH TIME ZONE, -- Can be NULL
    reviewer_id UUID, -- Can be NULL, references users(id) if reviewer is also a user
    rejection_reason TEXT -- Can be NULL
);

-- Add indexes for common lookups
CREATE INDEX idx_users_line_user_id ON users (line_user_id);
CREATE INDEX idx_profiles_user_id ON profiles (user_id);
CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);
CREATE INDEX idx_identity_verifications_user_id ON identity_verifications (user_id);
CREATE INDEX idx_identity_verifications_status ON identity_verifications (status);

-- Optionally, add initial roles
INSERT INTO roles (name, description) VALUES
('POSTER', 'Can post tasks'),
('TASKER', 'Can take on tasks'),
('ADMIN', 'Administrator with full access');