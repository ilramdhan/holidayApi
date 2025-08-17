# Secure Deployment Guide - Holiday API Indonesia

## 🔒 Security First Deployment

**IMPORTANT**: This guide contains security-critical steps that MUST be followed for production deployment.

## 📋 Pre-Deployment Checklist

### 1. Environment Security
- [ ] Generate strong JWT secret key (minimum 32 characters)
- [ ] Set secure database path with proper permissions
- [ ] Configure HTTPS/TLS certificates
- [ ] Set up firewall rules
- [ ] Configure reverse proxy (nginx/apache)

### 2. Database Security
- [ ] Change default database location
- [ ] Set proper file permissions (600)
- [ ] Enable database encryption if needed
- [ ] Configure backup strategy

### 3. Application Security
- [ ] Create secure admin credentials
- [ ] Remove default/example credentials
- [ ] Configure rate limiting appropriately
- [ ] Set up monitoring and alerting

## 🚀 Deployment Steps

### Step 1: Environment Configuration

Create production `.env` file:
```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=120s

# Database Configuration
DATABASE_PATH=/secure/path/to/holidays.db
MIGRATIONS_PATH=./migrations

# JWT Configuration (CRITICAL - CHANGE THESE)
JWT_SECRET_KEY=YOUR_SUPER_SECURE_JWT_SECRET_KEY_MIN_32_CHARS_RANDOM
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=168h

# Rate Limiting
RATE_LIMIT_RPM=100
RATE_LIMIT_BURST=20
```

### Step 2: Generate Secure JWT Secret

```bash
# Generate secure random key (Linux/Mac)
openssl rand -base64 32

# Or use online generator (ensure HTTPS)
# https://generate-random.org/api-key-generator
```

### Step 3: Create Admin User Securely

**Option A: Via Database (Recommended)**
```sql
-- Generate password hash first
-- Use bcrypt online tool or create via API

INSERT INTO users (username, email, password, role, is_active) VALUES 
('your-admin-username', 'admin@yourcompany.com', 'YOUR_BCRYPT_HASH', 'super_admin', TRUE);
```

**Option B: Via API (After deployment)**
```bash
# First, temporarily enable registration endpoint
# Then create admin user via API
# Finally, disable public registration
```

### Step 4: Database Permissions

```bash
# Set secure permissions
chmod 600 /path/to/holidays.db
chown app-user:app-group /path/to/holidays.db

# Secure directory
chmod 700 /path/to/database/directory
```

### Step 5: Reverse Proxy Configuration

**Nginx Example:**
```nginx
server {
    listen 443 ssl http2;
    server_name api.yourcompany.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 🔐 Security Hardening

### 1. System Level
```bash
# Create dedicated user
useradd -r -s /bin/false holidayapi

# Set up systemd service
sudo systemctl enable holidayapi
sudo systemctl start holidayapi
```

### 2. Application Level
- Enable all security middleware
- Configure CORS for specific domains
- Set up request logging
- Configure audit log retention

### 3. Network Level
- Use HTTPS only
- Configure firewall (UFW/iptables)
- Set up fail2ban for brute force protection
- Use VPN for admin access if needed

## 📊 Monitoring & Alerting

### 1. Application Monitoring
- Monitor failed login attempts
- Track API response times
- Monitor database performance
- Set up health check endpoints

### 2. Security Monitoring
- Monitor audit logs
- Set up alerts for suspicious activity
- Track rate limit violations
- Monitor certificate expiration

### 3. Log Management
```bash
# Configure log rotation
/var/log/holidayapi/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 holidayapi holidayapi
}
```

## 🔄 Maintenance

### Regular Tasks
- [ ] Update dependencies monthly
- [ ] Rotate JWT secrets quarterly
- [ ] Review audit logs weekly
- [ ] Update SSL certificates before expiration
- [ ] Backup database daily

### Security Updates
- [ ] Monitor security advisories
- [ ] Test updates in staging first
- [ ] Have rollback plan ready
- [ ] Document all changes

## 🚨 Incident Response

### If Credentials Compromised
1. Immediately rotate JWT secret
2. Force logout all users
3. Review audit logs
4. Change admin passwords
5. Notify stakeholders

### If Database Compromised
1. Take system offline
2. Assess damage scope
3. Restore from clean backup
4. Investigate breach vector
5. Implement additional security

## 📞 Support

For security issues:
- Create private GitHub issue
- Contact security team directly
- Follow responsible disclosure

---

**Remember**: Security is not a one-time setup but an ongoing process. Regular reviews and updates are essential for maintaining a secure system.
