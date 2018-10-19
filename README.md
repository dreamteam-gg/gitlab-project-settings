# GitLab Project Settings

## Configuration

```
---
namespace_id: 1
gitlab_url: https://gitlab.com/api/v4
gitlab_private_token: AABBCCDD
stop_on_error: true

exclude_projects: []
only_projects: []

settings:
  approvals_before_merge: 2
  disable_overriding_approvers_per_merge_request: true
  reset_approvals_on_push: true
  merge_requests_author_approval: false
  only_allow_merge_if_pipeline_succeeds: true
  only_allow_merge_if_all_discussions_are_resolved: true
  resolve_outdated_diff_discussions: true
  printing_merge_request_link_enabled: true
```
