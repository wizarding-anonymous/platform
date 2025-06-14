#!/bin/bash
# generate_keys.sh
# Generates RSA private and public keys for JWT signing.
# Ensure OpenSSL is installed.

OUTPUT_DIR="configs/keys" # Relative to backend/auth-service
mkdir -p "$OUTPUT_DIR"

PRIVATE_KEY_FILE="$OUTPUT_DIR/jwtRS256.key"
PUBLIC_KEY_FILE="$OUTPUT_DIR/jwtRS256.key.pub"

if [ -f "$PRIVATE_KEY_FILE" ] || [ -f "$PUBLIC_KEY_FILE" ]; then
  echo "Key files already exist in $OUTPUT_DIR. Remove them first to regenerate."
  exit 1
fi

echo "Generating RSA private key (2048 bit)..."
openssl genrsa -out "$PRIVATE_KEY_FILE" 2048

echo "Extracting public key from private key..."
openssl rsa -in "$PRIVATE_KEY_FILE" -pubout -out "$PUBLIC_KEY_FILE"

echo "Keys generated:"
echo "Private Key: $PRIVATE_KEY_FILE"
echo "Public Key:  $PUBLIC_KEY_FILE"

# Set restrictive permissions for the private key
chmod 600 "$PRIVATE_KEY_FILE"
chmod 644 "$PUBLIC_KEY_FILE" # Public key can be readable

echo "Permissions set."
echo "Remember to add '$OUTPUT_DIR/*.key' to your .gitignore if it's not already there!"
echo "The public key ($PUBLIC_KEY_FILE) can be shared with services that need to validate JWTs."
echo "The private key ($PRIVATE_KEY_FILE) MUST be kept secret."
