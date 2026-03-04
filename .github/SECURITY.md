# Security Policy

## Supported Versions

The following versions of Holiday API Indonesia are currently being supported with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of Holiday API Indonesia seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to: **ilramdhan@gmail.com**

Please include the following information in your report:

- **Type of vulnerability** (e.g., SQL injection, XSS, authentication bypass)
- **Affected versions** (e.g., v1.0.0)
- **Description of the vulnerability** with as much detail as possible
- **Steps to reproduce** the vulnerability
- **Potential impact** of the vulnerability
- **Suggested fix** (if you have one)

### Response Timeline

We will acknowledge receipt of your vulnerability report within **48 hours** and will send you a more detailed response within **72 hours** indicating the next steps in handling your report.

After the initial reply to your report, the security team will endeavor to keep you informed of the progress towards a fix and full announcement, and may ask for additional information or guidance.

### What to Expect

1. **Acknowledgment**: We will acknowledge your email within 48 hours
2. **Assessment**: We will assess the vulnerability and determine its severity
3. **Fix Development**: We will work on a fix and may reach out for clarification
4. **Release**: Once fixed, we will release a security patch
5. **Disclosure**: We will coordinate with you on public disclosure timing

### Security Best Practices

When deploying Holiday API Indonesia, please ensure you follow these security best practices:

1. **Environment Variables**: Never commit `.env` files or secrets to version control
2. **JWT Secret**: Use a strong, random JWT secret key (minimum 32 characters)
3. **Admin API Key**: Use a secure, random admin API key
4. **HTTPS**: Always use HTTPS in production
5. **Rate Limiting**: Enable rate limiting to prevent abuse
6. **Regular Updates**: Keep dependencies updated
7. **Database Security**: Secure your SQLite database file with proper permissions

### Known Security Features

The following security features are implemented in Holiday API Indonesia:

- ✅ JWT Authentication with secure token handling
- ✅ Password hashing using bcrypt
- ✅ Input sanitization and validation
- ✅ SQL injection protection
- ✅ XSS protection headers
- ✅ Rate limiting per IP/user
- ✅ CORS protection
- ✅ Security headers (HSTS, CSP, X-Frame-Options, etc.)
- ✅ Audit logging for security events

### Security Updates

Security updates will be released as patch versions (e.g., v1.0.1) and will be clearly marked in the release notes. We recommend always running the latest version to ensure you have the latest security fixes.

### Hall of Fame

We would like to thank the following security researchers who have responsibly disclosed vulnerabilities to us:

*No reported vulnerabilities yet - thank you for keeping our community safe!*

---

Thank you for helping keep Holiday API Indonesia and our users safe!