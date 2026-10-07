# JLF-Payload-Obfuscator

    Designed by jakelo.ai · Coded with AI assistance

Payload mutation engine: feed it one payload, get back a deduplicated
wordlist of encoding and filter-evasion variants. One binary, zero
dependencies, works offline.

## Why this exists

WAF rules, input filters, and sanitizers all fail the same way: they pattern
-match. If your rule blocks `UNION SELECT`, it probably misses `UnIoN/**/
SeLeCt`. If your filter strips `eval(`, it probably misses
`$_=(~hex2bin(...));`.

`jlfpayload` exists to answer one question empirically: **do my filters
actually hold?** Feed your rules every variant this tool generates. If any
variant slips through, the rule is wrong — fix it before someone else runs
the same test against you.

Primary uses:

- **Validating WAF/filter rules** on infrastructure you own or are authorized
  to test — generate the bypass variants, replay them against your own rules,
  measure what lands.
- **Authorized engagements** — pentest and bug bounty work where the target's
  filters need to be tested within scope, with the wordlist generated locally
  before you touch the target.

## What it does NOT do

- No network calls. Ever. Input is a string or a local file; output is a
  local wordlist.
- No exploitation. It transforms text. What the variants are sent to is
  entirely outside this tool.
- No payload database. You bring the payloads; it brings the mutations.

## Install

```bash
git clone https://github.com/jakeloai/jlfpayload.git
cd jlfpayload
go build -o jlfpayload .
./jlfpayload -version
```

Requires Go 1.21+. No CGO, no modules to download — `go build` works
offline after clone.

## Usage

```bash
# Single payload, all mutation classes
./jlfpayload -p "' OR 1=1 -- -"

# A wordlist of your base payloads, SQL mutations only
./jlfpayload -f payloads.txt -m sql -o sql_variants.txt

# PHP filter evasion, version-aware (legacy 5.x / modern 7.x / latest 8.x)
./jlfpayload -p "system('id');" -m php -phpv latest

# XSS context mutations
./jlfpayload -p '<svg onload=alert(1)>' -m xss

# Machine-readable summary (variants count, path, SHA-256) for scripts
./jlfpayload -f payloads.txt -q
```

### Modes

| Mode | Mutation classes |
|------|------------------|
| `general` | URL encode, double URL encode, hex (plain / `\x` / `0x`), HTML entities (dec / hex), unicode escapes, base64 (+ URL-safe combos) |
| `sql` | general + `/**/` and `+` comment substitution, keyword case randomization, mixed comment+case, then re-encoded variants |
| `xss` | general + tag/event case randomization, `javascript:` protocol case mixing, combined mutations, re-encoded variants |
| `php` | general + XOR dynamic strings (`$_=(hex2bin(...)^...)`), NOT dynamic strings (`$_=(~hex2bin(...))`), base64 `eval`/`assert` wrappers, version-specific `preg_replace` / `create_function` variants |
| `all` | every class above |

PHP version flag matters because the useful evasion surface changed between
versions: `preg_replace /e` died in PHP 7, `create_function` died in PHP 8.
Generating dead payloads for a live PHP 8 target is noise; this tool won't.

### Output

- Deduplicated wordlist, one variant per line.
- SHA-256 of the output file printed after generation — embed it in your test
  notes so a report can reference exactly which variant set was used.
- `-q` prints `count / path / sha256` on stdout for pipeline use.

## Legal & responsible use

Generating obfuscated variants is legal; **sending them at systems you are
not authorized to test is not.** This tool is for testing defenses you own,
and for authorized engagements under explicit written scope. That is the
entire intended use. The developer assumes no liability for misuse.

Found a bug in the mutations (a variant that should decode but doesn't)?
See [SECURITY.md](SECURITY.md).

## Development

Designed and reviewed by jakelo.ai; coded with AI assistance under the
JakeLo Framework. Every mutation class in this file was chosen for a
documented reason — see comments in `main.go`. If a variant doesn't
survive a real filter, that's a bug: open an issue.

## License

MIT © 2026 JakeLo
