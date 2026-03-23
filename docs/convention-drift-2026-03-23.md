# Convention Drift Check (2026-03-23)

`CODEX.md` was not present anywhere under `/workspace`, so this review is based on `CLAUDE.md`, `docs/development.md`, and current coverage data from `go test -coverprofile=/tmp/cover.out ./...`.

No fixes were applied as part of this pass.

## stdlib -> core.*

Confirmed drift against the documented `CLAUDE.md` rules: none.

- No uses found of `fmt.Errorf`, `errors.New`, `os.ReadFile`, `os.WriteFile`, or `os.MkdirAll`.

Broader-scope candidates only if the missing `CODEX.md` was meant to forbid more stdlib usage than `CLAUDE.md` currently documents:

- `config.go:22`, `config.go:23` use `os.Getenv(...)`.
- `cmd/forgegen/main.go:16`, `cmd/forgegen/main.go:17`, `cmd/forgegen/main.go:27`, `cmd/forgegen/main.go:28` use `fmt.Fprintf(os.Stderr, ...)` and `os.Exit(1)`.
- `cmd/forgegen/generator_test.go:26`, `cmd/forgegen/generator_test.go:52`, `cmd/forgegen/generator_test.go:88`, `cmd/forgegen/generator_test.go:130` use `os.ReadDir(...)` in tests.

## UK English

Drift found in hand-written docs:

- `README.md:2` badge text uses `License`.
- `README.md:9` heading uses `License`.
- `CONTRIBUTING.md:34` heading uses `License`.

## Missing Tests

Coverage basis: `go test -coverprofile=/tmp/cover.out ./...` on 2026-03-23. Line references below are function start lines from `go tool cover -func=/tmp/cover.out`.

No dedicated service test file:

- `milestones.go:10`, `milestones.go:20`, `milestones.go:26`, `milestones.go:36` (`MilestoneService`; no `milestones_test.go` present).

Zero-coverage functions and methods:

- `actions.go:29` `IterRepoSecrets`; `actions.go:55` `IterRepoVariables`; `actions.go:81` `IterOrgSecrets`; `actions.go:93` `IterOrgVariables`
- `admin.go:27` `IterUsers`; `admin.go:64` `IterOrgs`; `admin.go:80` `IterCron`
- `branches.go:25` `ListBranchProtections`; `branches.go:31` `IterBranchProtections`; `branches.go:37` `GetBranchProtection`; `branches.go:57` `EditBranchProtection`; `branches.go:67` `DeleteBranchProtection`
- `client.go:211` `do`
- `cmd/forgegen/generator.go:158` `enumConstName`
- `cmd/forgegen/main.go:9` `main`
- `commits.go:34` `ListAll`; `commits.go:39` `Iter`
- `issues.go:31` `Unpin`; `issues.go:37` `SetDeadline`; `issues.go:44` `AddReaction`; `issues.go:51` `DeleteReaction`; `issues.go:58` `StartStopwatch`; `issues.go:64` `StopStopwatch`; `issues.go:70` `AddLabels`; `issues.go:77` `RemoveLabel`; `issues.go:83` `ListComments`; `issues.go:89` `IterComments`; `issues.go:106` `toAnySlice`
- `labels.go:28` `IterRepoLabels`; `labels.go:76` `IterOrgLabels`
- `milestones.go:20` `ListAll`; `milestones.go:26` `Get`; `milestones.go:36` `Create`
- `notifications.go:27` `Iter`; `notifications.go:38` `IterRepo`
- `orgs.go:31` `IterMembers`; `orgs.go:37` `AddMember`; `orgs.go:43` `RemoveMember`; `orgs.go:49` `ListUserOrgs`; `orgs.go:55` `IterUserOrgs`; `orgs.go:61` `ListMyOrgs`; `orgs.go:66` `IterMyOrgs`
- `packages.go:28` `Iter`; `packages.go:56` `IterFiles`
- `pulls.go:32` `Update`; `pulls.go:38` `ListReviews`; `pulls.go:44` `IterReviews`; `pulls.go:50` `SubmitReview`; `pulls.go:60` `DismissReview`; `pulls.go:67` `UndismissReview`
- `releases.go:35` `DeleteByTag`; `releases.go:41` `ListAssets`; `releases.go:47` `IterAssets`; `releases.go:53` `GetAsset`; `releases.go:63` `DeleteAsset`
- `repos.go:29` `IterOrgRepos`; `repos.go:34` `ListUserRepos`; `repos.go:39` `IterUserRepos`; `repos.go:58` `Transfer`; `repos.go:63` `AcceptTransfer`; `repos.go:68` `RejectTransfer`; `repos.go:73` `MirrorSync`
- `teams.go:31` `IterMembers`; `teams.go:43` `RemoveMember`; `teams.go:49` `ListRepos`; `teams.go:55` `IterRepos`; `teams.go:61` `AddRepo`; `teams.go:67` `RemoveRepo`; `teams.go:73` `ListOrgTeams`; `teams.go:79` `IterOrgTeams`
- `users.go:40` `IterFollowers`; `users.go:46` `ListFollowing`; `users.go:52` `IterFollowing`; `users.go:58` `Follow`; `users.go:64` `Unfollow`; `users.go:70` `ListStarred`; `users.go:76` `IterStarred`; `users.go:82` `Star`; `users.go:88` `Unstar`
- `webhooks.go:38` `IterOrgHooks`

## SPDX Headers

No tracked `.go`, `.md`, `go.mod`, or `.gitignore` files currently start with an SPDX header.

Repo meta and docs:

- `.gitignore:1`, `go.mod:1`, `CLAUDE.md:1`, `CONTRIBUTING.md:1`, `README.md:1`, `docs/architecture.md:1`, `docs/development.md:1`, `docs/index.md:1`

Hand-written package files:

- `actions.go:1`, `admin.go:1`, `branches.go:1`, `client.go:1`, `commits.go:1`, `config.go:1`, `contents.go:1`, `doc.go:1`, `forge.go:1`, `issues.go:1`, `labels.go:1`, `milestones.go:1`, `misc.go:1`, `notifications.go:1`, `orgs.go:1`, `packages.go:1`, `pagination.go:1`, `params.go:1`, `pulls.go:1`, `releases.go:1`, `repos.go:1`, `resource.go:1`, `teams.go:1`, `users.go:1`, `webhooks.go:1`, `wiki.go:1`

Hand-written test files:

- `actions_test.go:1`, `admin_test.go:1`, `branches_test.go:1`, `client_test.go:1`, `commits_test.go:1`, `config_test.go:1`, `contents_test.go:1`, `forge_test.go:1`, `issues_test.go:1`, `labels_test.go:1`, `misc_test.go:1`, `notifications_test.go:1`, `orgs_test.go:1`, `packages_test.go:1`, `pagination_test.go:1`, `params_test.go:1`, `pulls_test.go:1`, `releases_test.go:1`, `resource_test.go:1`, `teams_test.go:1`, `users_test.go:1`, `webhooks_test.go:1`, `wiki_test.go:1`

Generator and tooling:

- `cmd/forgegen/generator.go:1`, `cmd/forgegen/generator_test.go:1`, `cmd/forgegen/main.go:1`, `cmd/forgegen/parser.go:1`, `cmd/forgegen/parser_test.go:1`

Generated `types/` files (would need generator-owned header handling rather than manual edits):

- `types/action.go:1`, `types/activity.go:1`, `types/admin.go:1`, `types/branch.go:1`, `types/comment.go:1`, `types/commit.go:1`, `types/common.go:1`, `types/content.go:1`, `types/error.go:1`, `types/federation.go:1`, `types/generate.go:1`, `types/git.go:1`, `types/hook.go:1`, `types/issue.go:1`, `types/key.go:1`, `types/label.go:1`, `types/milestone.go:1`, `types/misc.go:1`, `types/notification.go:1`, `types/oauth.go:1`, `types/org.go:1`, `types/package.go:1`, `types/pr.go:1`, `types/quota.go:1`, `types/reaction.go:1`, `types/release.go:1`, `types/repo.go:1`, `types/review.go:1`, `types/settings.go:1`, `types/status.go:1`, `types/tag.go:1`, `types/team.go:1`, `types/time_tracking.go:1`, `types/topic.go:1`, `types/user.go:1`, `types/wiki.go:1`
