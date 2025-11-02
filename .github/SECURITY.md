# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Currently supported versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.x.x   | :white_check_mark: |

**Note**: This project is currently in beta (0.x.x versions). Once we reach 1.0.0, we will establish a formal support policy for older versions.

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to:

📧 **gerry@gerrymiller.com**

### What to Include

Please include the following information in your report:

- **Description**: Clear description of the vulnerability
- **Impact**: What could an attacker accomplish with this vulnerability
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Affected Versions**: Which versions are affected
- **Proof of Concept**: Code or commands demonstrating the issue (if applicable)
- **Suggested Fix**: Your ideas on how to fix it (if any)
- **Disclosure Timeline**: When you plan to publicly disclose (we request 90 days)

### Response Timeline

- **Initial Response**: Within 48 hours of receipt
- **Status Update**: Within 7 days with assessment and action plan
- **Resolution**: Depends on severity and complexity
  - Critical: Immediate priority
  - High: Within 30 days
  - Medium: Within 60 days
  - Low: Within 90 days

### What to Expect

1. **Acknowledgment**: We'll confirm receipt of your report
2. **Assessment**: We'll investigate and assess the severity
3. **Updates**: Regular updates on our progress
4. **Resolution**: We'll develop and test a fix
5. **Disclosure**: Coordinated disclosure once fixed
6. **Credit**: Public acknowledgment in release notes (if desired)

## Security Practices

This project follows security best practices:

### Automated Security

- **Dependency Scanning**: Snyk monitors dependencies for known vulnerabilities
- **Code Scanning**: CodeQL and Gosec analyze code for security issues
- **Secret Scanning**: GitHub scans for accidentally committed secrets
- **Dependabot**: Automated security updates for dependencies
- **CI/CD Security**: All PRs undergo security checks before merge

### Development Security

- **Input Validation**: All user inputs are validated and sanitized
- **API Key Protection**: Secrets stored in environment variables, never committed
- **Least Privilege**: Components operate with minimum necessary permissions
- **Error Handling**: Errors are handled securely without leaking sensitive information
- **Dependencies**: Regular updates and security audits

### Deployment Security

- **HTTPS Only**: All API communications use TLS
- **Authentication**: OpenAI API keys required and validated
- **Rate Limiting**: Protection against abuse (for API deployments)
- **Container Security**: Docker images scanned for vulnerabilities
- **Access Control**: Principle of least privilege for all services

## Known Security Considerations

### API Keys

This project requires OpenAI API keys. Users must:
- Store keys securely in environment variables or config files (not in code)
- Use separate keys for development and production
- Rotate keys periodically
- Set spending limits on OpenAI accounts
- Never commit keys to version control

### Document Processing

- PDF and HTML parsers may be vulnerable to malicious documents
- Use caution when processing untrusted documents
- Consider sandboxing document processing in production
- Validate document sizes and complexity before processing

### Vector Database

- Qdrant should be deployed with authentication in production
- Network access should be restricted
- Data encryption at rest recommended for sensitive documents

### LLM Prompt Injection

- User queries are sent to LLMs and could contain prompt injection attempts
- System prompts are designed to be robust but not 100% injection-proof
- Validate and sanitize user inputs before sending to LLMs
- Monitor LLM outputs for unexpected behavior

## Security Updates

Security updates are released as:
- **Patch versions** (0.1.x) for minor security fixes
- **Minor versions** (0.x.0) for moderate security updates
- **Immediate releases** for critical vulnerabilities

Subscribe to:
- **GitHub Security Advisories**: Watch this repo for security alerts
- **Release Notes**: Check releases for security-related changes
- **Dependabot PRs**: Review and merge security updates promptly

## Bug Bounty Program

We currently do not have a bug bounty program. However, we deeply appreciate security researchers who help improve our security posture and will acknowledge contributions in release notes.

## Contact

For security-related questions or concerns:

📧 Email: gerry@gerrymiller.com  
🔒 PGP Key: Available upon request  
🐛 GitHub: [@gerrymiller](https://github.com/gerrymiller)

---

**Thank you for helping keep Deep Thinking Agent and its users safe!**
