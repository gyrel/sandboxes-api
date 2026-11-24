# Changelog

### Added (contract 0.0.38)

- Add MCP gateway lifecycle, inspection and credential management contracts.

## Unreleased

### Added (contract 0.0.37)

- `KitPolicyService` (`CheckKit`, `GetKitPolicyConfig`) exposes the daemon's kit
  source policy for display and fail-closed checks. `CheckKit` takes
  `source_kind` (directory, zip, git, oci), `source_ref` (a URL/path identifier
  only), an optional `sandbox_id`, and `signature_present`; it returns `allowed`,
  a stable `reason` (`OK`, `SOURCE_NOT_ALLOWED`, `LOCAL_DISABLED`,
  `SIGNATURE_REQUIRED`, `POLICY_UNAVAILABLE`, `INVALID_REFERENCE`), a
  human-readable `message`, and `checked_settings` (setting names only). A daemon
  without a policy engine returns `POLICY_UNAVAILABLE` with `allowed=false`;
  org/managed deny is sticky; unknown enum values are `INVALID_ARGUMENT`. No
  field carries secret material.
- `GetKitPolicyConfig` returns the effective `allowed_sources`,
  `allow_local_kits`, `require_signature` and the read-only `managed` flag (org
  policy override). Local settings may only tighten a managed policy.

### Compatibility

- Additive: no existing field, message, service or RPC was changed or removed.
  Readers predating 0.0.37 ignore the new service.

### Added (contract 0.0.36)

- `ApprovalService` (`List`, `Inspect`, `Respond`) exposes host-side pending
  approvals for `policy approval ls/inspect/respond`. An `Approval` carries a
  stable ID, kind (host command, MCP registration, MCP call, other), requesting
  sandbox, action summary, normalized HTTPS destination host and path (no
  query, userinfo or fragment), header names with static-or-vault class and
  vault reference IDs, bounded environment variable names, a redacted cwd
  summary and a creation timestamp. No field carries a header, environment or
  secret value. `Inspect` and `List` filter by sandbox.
- `Respond` takes an APPROVE or DENY decision and an optional ONCE or EXACT
  persist scope. Persisting an EXACT rule requires a trusted OS-backed approval
  provider. A daemon without one fails closed: it refuses APPROVE and EXACT
  with `FAILED_PRECONDITION` rather than downgrading them.
- `SecretService.SetCustom` creates or updates a custom secret that a sandbox
  uses through a placeholder. Sources are an inline write-only value, a
  `reference` (`op://` or AWS ARN) or a `command`, reusing `SecretSource` for
  refresh (on-demand or duration) and `no_verify`. Reference and command
  sources require host approval and return `approval_required` plus an
  `approval_id` until approved. Placeholders are explicit or generated (with
  `{rand}`), injected through `env`, and substituted by the proxy for `hosts`
  patterns (exact, IP, `*`, `**`). Responses never carry the value or resolver
  text; the write response echoes the resolved placeholder in full so the CLI
  can print the generated value (owner decision 2026-11-24; the same string is
  injected into the guest environment, so this adds no exposure), while
  inventories keep only the masked form.
- `SecretSummary.env_name`, `hosts` and `placeholder_masked` expose custom-secret
  placement. The placeholder is masked to at most four trailing characters.

### Compatibility

- Additive: no existing field, message or RPC was changed or removed. Readers
  predating 0.0.36 ignore the new RPCs and summary fields.
- Local-surface only. Cloud `--cloud` flags, `--header` and `--format` for
  literal header injection are not modeled; Rheic is vault-only.

### Added (contract 0.0.35)

- `TemplateService.Export` RPC (`TemplateServiceExportRequest{name,
  tar_path}` → `{size_bytes}`): docker-save export of a recorded template;
  existing destination files are refused. (Contract 0.0.34 shipped without
  this RPC due to a changelog-only merge; 0.0.35 carries the intended
  additive surface.)

## Unreleased

### Added (contract 0.0.33)

- `TemplateService` (Save/Load/Remove/List/Inspect) with `Template`
  records: save optionally writes a docker-save compatible tar to a host
  path and echoes a non-secret notice (guest-resident secrets are not
  exported); load consumes OCI or docker-save archives with optional name
  override; inspect returns name-only env entries plus working dir and
  user (presence-only contract, no secret material).

### Added (contract 0.0.32)

- `CheckPolicyRequest.url` and `.method` let `policy check network` evaluate
  an HTTP(S) request target against the effective policy. Credentials and
  fragments are rejected; the query string is never logged.

### Added

- `NetworkRule.methods`, `NetworkRule.path` and `NetworkRule.query`
  (repeated `HTTPQueryCondition`) carry HTTP L7 policy conditions, matching
  the gateway policy engine's semantics: methods are upper-cased and `ANY`
  marks a rule HTTP-conditioned without relaxing non-HTTP traffic; `path` is
  a segment-aware glob (`**` spans segments); every query condition must
  match with an exact name and a glob value. Local RPC contract advances to
  0.0.31.

### Compatibility

- Additive proto fields; readers predating this release ignore them.
  Daemons before contract 0.0.31 must reject rules carrying HTTP conditions
  rather than silently treating them as plain domain allows (fail closed).
- The module tag `v0.0.30` (contract 0.0.30, account `username`, and
  contract 0.0.29, `Operation.progress`) shipped without changelog entries;
  this section covers the contract line from 0.0.29 forward.

## v0.0.29 — 2026-11-19

### Added

- Add optional `SkillsOptions` to `CreateRequest` and `Sandbox` with explicit
  `off`, `readonly` and `readwrite` modes plus a host store path. Omitted options
  keep shared skills disabled.
- Advance the local RPC contract version to 0.0.28. Downstream modules must bump
  `github.com/rheic/sandboxes-api` to this release before using the field.

### Compatibility

- Adding two fields is wire-compatible for readers and writers that predate
  shared skills; older daemons ignore the request field. Clients must gate the
  option on API contract 0.0.28 so an older daemon is never asked to downgrade a
  requested mode silently.


## v0.0.28 — 2026-11-16

### Changed

- Remove the mandatory external release-review gate while retaining candidate validation, CI, and immutable tag protections.
- Accept an `Unreleased` heading before the newest version during release validation.
- Preserve API contract version 0.0.27; this release only changes release tooling.

## Local RPC contract 0.0.27

- Add non-secret metadata to `SecretSummary`: creation and update times,
  source kind, refresh policy and capability, coarse value type, and origin.
  The summary never carries a secret value, reference, command or environment
  variable value.
- Add creation and update times and origin to `RegistryCredentialSummary`.
- Advance the local RPC contract version to 0.0.27. Downstream modules must
  bump `github.com/rheic/sandboxes-api` to the release that carries this
  contract.

## Local RPC contract 0.0.26

- Add `SandboxService.GetOrgPolicy`, a read-only view of the locally stored,
  signed organization policy bundle. It reports the organization, bundle
  schema version and revision, a payload fingerprint, enforcement state, load
  time, and parsed network rules with an `enforced` flag. It never includes
  signatures, tokens, or secret material.
- Add optional `NetworkRule.origin` (`local` or `org`) so policy listings can
  show provenance. Older daemons leave it empty; clients must treat an empty
  value as `local`.
- Advance the local RPC contract version to 0.0.26. Downstream modules must
  bump `github.com/rheic/sandboxes-api` to the release that carries this
  contract.

## v0.0.25
- Add optional `Workspace.guest_path` for mounting a validated host path at an
  explicit guest destination. Empty values retain same-path mount behavior.
- Advance the local RPC contract version to 0.0.25 so clients can reject older
  daemons that would ignore a non-empty guest destination.
- Include the typed clone request and sandbox clone metadata in the unpublished
  0.0.25 contract candidate; the endpoint is optional and omitted while stopped.

## v0.0.24 — 2026-11-09

- Deprecate static remote MCP request headers; supported daemons reject them
  because sensitivity cannot be inferred from custom header names.
- Require host-vault-backed `secret_headers` templates for all remote MCP headers.
- Clarify that pre-1.0 breaking contract changes follow the documented alpha migration policy.

## v0.0.23 — 2026-11-08

- Bind M7 pre-tag approval to the exact reviewed commit.
- Add the `secret_headers` MCP registration field for host-vault placeholder
  templates, keeping credential values out of the registration RPC and state.

## v0.0.20 — 2026-11-07

- Add optional `url` and non-secret `headers` to `MCPRegisterServerRequest`
  for the remote streamable-HTTP MCP transport. Exactly one of `command` or
  `url` must be set; remote hosts are constrained by the daemon's fail-closed
  allowlist.

## v0.0.19 — 2026-11-07

- Add presence-only refresh metadata (`refresh_present`, `refresh_coordinated`)
  to `AccountService.Status`. No token material, expiry values or validity
  claims are exposed.

## v0.0.18 — 2026-11-07

- Add the local `SandboxService.MCPCall` contract for invoking a named MCP tool
  in a sandbox. Daemon-side registry lookup, policy enforcement, bounded
  payload validation, and approval remain required before execution.
- Add a redacted `SandboxService.ListMCPServers` inventory contract for local
  `mcp ls`; it exposes gateway/transport and usability flags, never credentials
  or OAuth identity details, and must not start servers or auth flows.
- Add explicit stdio registration/removal contracts and one-use registration approval challenges.
- Add one-use MCP call approval challenges. The approval RPC carries only the challenge ID; it is not proof of consent. The daemon must obtain approval through a trusted host approval provider before execution.
- Document that persisted stdio command arguments are non-secret configuration; credentials require a host credential broker.

## v0.0.16 — 2026-11-06

### Added

- Add an optional host-side dynamic secret source to service-secret set/create RPCs, including a 1Password/AWS reference or shell command, refresh policy, verification controls, and captured provider account/profile selectors.

### Compatibility

- Existing value-based secret requests are unchanged. Older daemons ignore the additive field and reject the empty value; SDK clients must check the daemon API version before sending a dynamic source.

## v0.0.15 — 2026-11-05

### Added

- Add sandbox-scoped OAuth metadata to local Create using nested messages aligned with the Docker Kits v2 declaration: token endpoint, resource hosts, sentinels, credential file, response fields, passthrough, and `skipIfEnv`.

### Compatibility

- This contract carries kit metadata only. Daemons must fail closed when custom OAuth interception is unavailable; the contract alone does not imply runtime support.

## v0.0.12 — 2026-11-03

### Added

- Add the account device-login stream and credential-presence status RPCs.
- Publish a CycloneDX module SBOM as a CI artifact.

### Changed

- Normalize the Go module and generated protobuf package paths to the lowercase `github.com/quoiner/sandboxes-api` namespace.
- Update `golang.org/x/net` to v0.62.0, including fixes for GO-2026-6860 and GO-2026-6866.

