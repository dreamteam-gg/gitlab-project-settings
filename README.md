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

### Top level configuration variables

Configuration supports passing top level fields as env variable with `GITLAB_` prefix:
```shell
# strings
GITLAB_GITLAB_PRIVATE_TOKEN=<sometoken>

# lists joined by spaces
GITLAB_ONLY_PROJECTS="some-project second-project"
```

### Environment variable substitution

Any key in configuration can be set as env variable. For example set Jira service password via `JIRA_PASSWORD` env variable:

```yaml
...
project_settings:
  services:
    # https://docs.gitlab.com/ee/api/services.html#createedit-jira-service
    jira:
      url: https://jira.example.com
      username: admin
      password: ${JIRA_PASSWORD}
...
```

### Masking output

To hide field values from diff output add them to `mask` list in configuration:
```yaml
...
mask:
  - password
...
```

### Complete example:

```yaml
---
group_id: devops # groups name string, does not work with personal accounts
gitlab_url: https://gitlab.com/api/v4
gitlab_private_token: asdgfdhgfjhg
create_missing: true # create projects from only_projects if they are missing

exclude_projects: []
only_projects:
  - some-project
  - second-project

group_settings:
  members:
    some_user_name: Maintainer # Guest, Reporter,Developer, Maintainer, Owner

# mask diff for any fields containing this words
mask:
  - password
  - webhook

# settings for all projects
project_settings:
  project: # ref: https://docs.gitlab.com/ee/api/projects.html#edit-project
    only_allow_merge_if_pipeline_succeeds: true
    only_allow_merge_if_all_discussions_are_resolved: true
    resolve_outdated_diff_discussions: true
    printing_merge_request_link_enabled: true
    snippets_enabled: false
    wiki_enabled: false
    merge_method: ff
    shared_runners_enabled: false
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
      allowed_to_merge:
        - user_id: user_name
        - group_id: devops
      allowed_to_push:
        - user_id: user_name
  services: # https://docs.gitlab.com/ee/api/services.html#createedit-slack-service
    slack:
      merge_requests_events: true
      notify_only_broken_pipelines: false
      notify_only_default_branch: false
      push_events: false
      issues_events: false
      confidential_issues_events: false
      tag_push_events: false
      note_events: false
      pipeline_events: false
      wiki_page_events: false
      webhook: https://hooks.slack.com/services/123456578
      username: Harold1
  webhooks: # https://docs.gitlab.com/ee/api/projects.html#add-project-hook
    "https://somehook.example.com/events":
      merge_requests_events: true
      enable_ssl_verification: true
  deploy_keys:
    - deploy-key-title
    - second-deploy-key

# override settings per project
overrides:
  "some-project": # must be in quotes
    webhooks:
      "https://somehook.example.com/events":
        merge_requests_events: true
        enable_ssl_verification: true
```
