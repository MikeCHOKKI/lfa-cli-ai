# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 0.1.x   | ✅        |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability:

1. **Do not** open a public issue
2. Email: [mikechokki5@gmail.com](mailto:mikechokki5@gmail.com)
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will respond within 48 hours and work to address the issue promptly.

## Security Best Practices

- All binary downloads use HTTPS
- Binary extraction uses `filepath.Base()` to prevent path traversal
- `io.LimitReader` (500MB) prevents decompression bomb attacks
- MCP package versions are pinned (npx @0.6.2, uvx @0.1.4)
- No hardcoded secrets — all tokens use environment variables
- `.env` is excluded from version control
