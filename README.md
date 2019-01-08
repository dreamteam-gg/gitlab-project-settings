# GitLab Project Settings

## Usage

```
Usage of ./gitlab-project-settings:
  -config string
    	Path to configuration file (default "./config.yml")
  -dry-run
    	Dry run mode
```

## Configuration

Configuration supports passing top level fields as env variable with `GITLAB_` prefix:
```
# strings
GITLAB_GITLAB_PRIVATE_TOKEN=<sometoken>

# lists joined by spaces
GITLAB_ONLY_PROJECTS="some-project second-project"
```

Complete example:
```
---
group_id: devops # groups name string
gitlab_url: https://gitlab.com/api/v4
gitlab_private_token: asdgfdhgfjhg
stop_on_error: true
create_missing: true # create projects from only_projects if they are missing

exclude_projects: []
only_projects:
  - some-project
  - second-project

# settings for all projects
settings:
  project: # ref: https://docs.gitlab.com/ee/api/projects.html#edit-project
    only_allow_merge_if_pipeline_succeeds: true
    only_allow_merge_if_all_discussions_are_resolved: true
    resolve_outdated_diff_discussions: true
    printing_merge_request_link_enabled: true
    snippets_enabled: false
    wiki_enabled: false
    merge_method: ff
  approvals: # ref: https://docs.gitlab.com/ee/api/merge_request_approvals.html#change-configuration
    approvals_before_merge: 1
    reset_approvals_on_push: true
    disable_overriding_approvers_per_merge_request: true
  approvers: # ref https://docs.gitlab.com/ee/api/merge_request_approvals.html#change-allowed-approvers
    approver_ids: []
    approver_group_ids: ["devops"] # group name as string
  protected_branches: # ref https://docs.gitlab.com/ee/api/protected_branches.html#protect-repository-branches
    master:
      push_access_level: NoAccess # NoAccess, Developer, Maintainer, Admin
      merge_access_level: Maintainer
      unprotect_access_level: NoAccess
  services: # https://docs.gitlab.com/ee/api/services.html#createedit-slack-service
    slack:
      merge_requests_events: true
      job_events: true
      webhook: https://hooks.slack.com/services/123456578
      username: Harold1

# override settings per project
overrides:
  "some-project":
    webhooks:
    "https://somehook.example.com/events":
      merge_requests_events: true
      enable_ssl_verification: true
```
