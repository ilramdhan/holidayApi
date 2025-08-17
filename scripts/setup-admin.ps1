# Holiday API - Secure Admin Setup Script (PowerShell)
# This script helps generate secure admin credentials

Write-Host "🔐 Holiday API - Secure Admin Setup" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""

# Function to validate password strength
function Test-PasswordStrength {
    param([string]$Password)
    
    if ($Password.Length -lt 8) {
        Write-Host "❌ Password must be at least 8 characters long" -ForegroundColor Red
        return $false
    }
    
    if ($Password -cnotmatch '[A-Z]') {
        Write-Host "❌ Password must contain at least one uppercase letter" -ForegroundColor Red
        return $false
    }
    
    if ($Password -cnotmatch '[a-z]') {
        Write-Host "❌ Password must contain at least one lowercase letter" -ForegroundColor Red
        return $false
    }
    
    if ($Password -notmatch '[0-9]') {
        Write-Host "❌ Password must contain at least one digit" -ForegroundColor Red
        return $false
    }
    
    if ($Password -notmatch '[^a-zA-Z0-9]') {
        Write-Host "❌ Password must contain at least one special character" -ForegroundColor Red
        return $false
    }
    
    return $true
}

# Function to generate bcrypt hash using .NET
function New-BcryptHash {
    param([string]$Password)
    
    try {
        # Try to use BCrypt.Net if available
        Add-Type -Path "BCrypt.Net.dll" -ErrorAction Stop
        return [BCrypt.Net.BCrypt]::HashPassword($Password)
    }
    catch {
        # Fallback: Use a simple hash (NOT RECOMMENDED for production)
        Write-Host "⚠️  Warning: BCrypt.Net not available, using fallback method" -ForegroundColor Yellow
        Write-Host "For production, please install BCrypt.Net or use online generator" -ForegroundColor Yellow
        
        # Generate a placeholder hash
        $bytes = [System.Text.Encoding]::UTF8.GetBytes($Password + "salt")
        $hash = [System.Security.Cryptography.SHA256]::Create().ComputeHash($bytes)
        $hashString = [System.Convert]::ToBase64String($hash)
        return "`$2a`$10`$" + $hashString.Substring(0, 53)
    }
}

# Get username
$username = Read-Host "Enter admin username"
if ([string]::IsNullOrWhiteSpace($username)) {
    Write-Host "❌ Username cannot be empty" -ForegroundColor Red
    exit 1
}

# Get email
$email = Read-Host "Enter admin email"
if ([string]::IsNullOrWhiteSpace($email)) {
    Write-Host "❌ Email cannot be empty" -ForegroundColor Red
    exit 1
}

# Get password securely
$password = Read-Host "Enter admin password" -AsSecureString
$passwordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($password))

if (-not (Test-PasswordStrength $passwordPlain)) {
    exit 1
}

# Confirm password
$confirmPassword = Read-Host "Confirm admin password" -AsSecureString
$confirmPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($confirmPassword))

if ($passwordPlain -ne $confirmPlain) {
    Write-Host "❌ Passwords do not match" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "⏳ Generating secure password hash..." -ForegroundColor Yellow

# Generate hash
$hash = New-BcryptHash $passwordPlain

# Clear password variables
$passwordPlain = $null
$confirmPlain = $null

# Generate SQL
Write-Host ""
Write-Host "✅ Admin user configuration generated!" -ForegroundColor Green
Write-Host ""
Write-Host "SQL to insert admin user:" -ForegroundColor Cyan
Write-Host "========================" -ForegroundColor Cyan
Write-Host "INSERT INTO users (username, email, password, role, is_active) VALUES" -ForegroundColor White
Write-Host "('$username', '$email', '$hash', 'super_admin', TRUE);" -ForegroundColor White
Write-Host ""
Write-Host "🔒 Security Notes:" -ForegroundColor Yellow
Write-Host "- Store this SQL securely" -ForegroundColor White
Write-Host "- Run it directly on your database" -ForegroundColor White
Write-Host "- Delete this output after use" -ForegroundColor White
Write-Host "- Never commit credentials to version control" -ForegroundColor White
Write-Host ""

$saveChoice = Read-Host "💾 Save to file? (y/n)"
if ($saveChoice -match '^[Yy]$') {
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
    $filename = "admin-setup-$timestamp.sql"
    
    $sqlContent = @"
INSERT INTO users (username, email, password, role, is_active) VALUES
('$username', '$email', '$hash', 'super_admin', TRUE);
"@
    
    $sqlContent | Out-File -FilePath $filename -Encoding UTF8
    Write-Host "✅ SQL saved to: $filename" -ForegroundColor Green
    Write-Host "⚠️  Remember to delete this file after use!" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "🚀 Next steps:" -ForegroundColor Cyan
Write-Host "1. Run the SQL on your database" -ForegroundColor White
Write-Host "2. Test login with your credentials" -ForegroundColor White
Write-Host "3. Delete any temporary files" -ForegroundColor White
Write-Host "4. Update your deployment documentation" -ForegroundColor White

# Clear sensitive variables
$hash = $null
[System.GC]::Collect()
