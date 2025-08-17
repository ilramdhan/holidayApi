#!/bin/bash

# Holiday API - Secure Admin Setup Script
# This script helps generate secure admin credentials

echo "🔐 Holiday API - Secure Admin Setup"
echo "===================================="
echo ""

# Check if required tools are available
if ! command -v openssl &> /dev/null; then
    echo "❌ Error: openssl is required but not installed"
    echo "Please install openssl and try again"
    exit 1
fi

# Function to validate password strength
validate_password() {
    local password="$1"
    
    if [[ ${#password} -lt 8 ]]; then
        echo "❌ Password must be at least 8 characters long"
        return 1
    fi
    
    if [[ ! "$password" =~ [A-Z] ]]; then
        echo "❌ Password must contain at least one uppercase letter"
        return 1
    fi
    
    if [[ ! "$password" =~ [a-z] ]]; then
        echo "❌ Password must contain at least one lowercase letter"
        return 1
    fi
    
    if [[ ! "$password" =~ [0-9] ]]; then
        echo "❌ Password must contain at least one digit"
        return 1
    fi
    
    if [[ ! "$password" =~ [^a-zA-Z0-9] ]]; then
        echo "❌ Password must contain at least one special character"
        return 1
    fi
    
    return 0
}

# Get username
read -p "Enter admin username: " username
if [[ -z "$username" ]]; then
    echo "❌ Username cannot be empty"
    exit 1
fi

# Get email
read -p "Enter admin email: " email
if [[ -z "$email" ]]; then
    echo "❌ Email cannot be empty"
    exit 1
fi

# Get password securely
echo "Enter admin password (input will be hidden):"
read -s password

if ! validate_password "$password"; then
    exit 1
fi

# Confirm password
echo "Confirm admin password:"
read -s confirm_password

if [[ "$password" != "$confirm_password" ]]; then
    echo "❌ Passwords do not match"
    exit 1
fi

echo ""
echo "⏳ Generating secure password hash..."

# Generate bcrypt hash using openssl and python/node if available
if command -v python3 &> /dev/null; then
    # Use Python to generate bcrypt hash
    hash=$(python3 -c "
import bcrypt
import sys
password = sys.argv[1].encode('utf-8')
hash = bcrypt.hashpw(password, bcrypt.gensalt())
print(hash.decode('utf-8'))
" "$password" 2>/dev/null)
    
    if [[ -z "$hash" ]]; then
        echo "❌ Failed to generate hash with Python. Please install python3-bcrypt"
        echo "Alternative: Use online bcrypt generator (ensure HTTPS)"
        echo "Password to hash: $password"
        exit 1
    fi
elif command -v node &> /dev/null; then
    # Use Node.js to generate bcrypt hash
    hash=$(node -e "
const bcrypt = require('bcrypt');
const password = process.argv[1];
const hash = bcrypt.hashSync(password, 10);
console.log(hash);
" "$password" 2>/dev/null)
    
    if [[ -z "$hash" ]]; then
        echo "❌ Failed to generate hash with Node.js. Please install bcrypt package"
        echo "Run: npm install -g bcrypt"
        echo "Alternative: Use online bcrypt generator (ensure HTTPS)"
        echo "Password to hash: $password"
        exit 1
    fi
else
    echo "❌ Neither Python3 nor Node.js with bcrypt is available"
    echo ""
    echo "Please install one of the following:"
    echo "1. Python3 with bcrypt: pip3 install bcrypt"
    echo "2. Node.js with bcrypt: npm install -g bcrypt"
    echo ""
    echo "Alternative: Use online bcrypt generator (ensure HTTPS)"
    echo "Password to hash: $password"
    exit 1
fi

# Generate SQL
echo ""
echo "✅ Admin user configuration generated!"
echo ""
echo "SQL to insert admin user:"
echo "========================"
echo "INSERT INTO users (username, email, password, role, is_active) VALUES"
echo "('$username', '$email', '$hash', 'super_admin', TRUE);"
echo ""
echo "🔒 Security Notes:"
echo "- Store this SQL securely"
echo "- Run it directly on your database"
echo "- Delete this output after use"
echo "- Never commit credentials to version control"
echo ""
echo "💾 Save to file? (y/n)"
read -r save_choice

if [[ "$save_choice" =~ ^[Yy]$ ]]; then
    filename="admin-setup-$(date +%Y%m%d-%H%M%S).sql"
    echo "INSERT INTO users (username, email, password, role, is_active) VALUES" > "$filename"
    echo "('$username', '$email', '$hash', 'super_admin', TRUE);" >> "$filename"
    echo "✅ SQL saved to: $filename"
    echo "⚠️  Remember to delete this file after use!"
fi

echo ""
echo "🚀 Next steps:"
echo "1. Run the SQL on your database"
echo "2. Test login with your credentials"
echo "3. Delete any temporary files"
echo "4. Update your deployment documentation"
